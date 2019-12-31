package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/renderer/ci"
	"github.com/zegl/kube-score/renderer/human"
	"github.com/zegl/kube-score/renderer/json_v2"
	"github.com/zegl/kube-score/score"
	"github.com/zegl/kube-score/scorecard"
)

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	setDefault(fs, "", true)

	// No command, flag, or file has been specified
	if len(os.Args) == 1 {
		fs.Usage()
		return
	}

	switch os.Args[1] {
	case "score":
		if err := scoreFiles(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to score files: %v", err)
			os.Exit(1)
		}
	case "list":
		listChecks()
	case "version":
		cmdVersion()
	case "help":
		fallthrough
	default:
		fs.Usage()
		os.Exit(1)
	}
}

func setDefault(fs *flag.FlagSet, actionName string, displayForMoreInfo bool) {
	fs.Usage = func() {
		usage := `Usage of kube-score:
kube-score [action] --flags

Actions:
	score	Checks all files in the input, and gives them a score and recommendations
	list	Prints a CSV list of all available score checks
	version	Print the version of kube-score
	help	Print this message` + "\n\n"

		if displayForMoreInfo {
			usage += `Run "kube-score [action] --help" for more information about a particular command`
		}

		if len(actionName) > 0 {
			usage += "Flags for " + actionName + ":"
		}

		fmt.Println(usage)

		if len(actionName) > 0 {
			fs.PrintDefaults()
		}
	}
}

func scoreFiles() error {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	exitOneOnWarning := fs.Bool("exit-one-on-warning", false, "Exit with code 1 in case of warnings")
	ignoreContainerCpuLimit := fs.Bool("ignore-container-cpu-limit", false, "Disables the requirement of setting a container CPU limit")
	ignoreContainerMemoryLimit := fs.Bool("ignore-container-memory-limit", false, "Disables the requirement of setting a container memory limit")
	verboseOutput := fs.CountP("verbose", "v", "Enable verbose output, can be set multiple times for increased verbosity.")
	printHelp := fs.Bool("help", false, "Print help")
	outputFormat := fs.StringP("output-format", "o", "human", "Set to 'human', 'json' or 'ci'. If set to ci, kube-score will output the program in a format that is easier to parse by other programs.")
	outputVersion := fs.String("output-version", "", "Changes the version of the --output-format. The 'json' format has version 'v1' (default) and 'v2'. The 'human' and 'ci' formats has only version 'v1' (default). If not explicitly set, the default version for that particular output format will be used.")
	optionalTests := fs.StringSlice("enable-optional-test", []string{}, "Enable an optional test, can be set multiple times")
	ignoreTests := fs.StringSlice("ignore-test", []string{}, "Disable a test, can be set multiple times")
	disableIgnoreChecksAnnotation := fs.Bool("disable-ignore-checks-annotations", false, "Set to true to disable the effect of the 'kube-score/ignore' annotations")
	setDefault(fs, "score", false)

	err := fs.Parse(os.Args[2:])
	if err != nil {
		return fmt.Errorf("failed to parse files: %s", err)
	}

	if *printHelp {
		fs.Usage()
		return nil
	}

	if *outputFormat != "human" && *outputFormat != "ci" && *outputFormat != "json" {
		fs.Usage()
		return fmt.Errorf("Error: --output-format must be set to: 'human', 'json' or 'ci'")
	}

	filesToRead := fs.Args()
	if len(filesToRead) == 0 {
		return fmt.Errorf(`Error: No files given as arguments.

Usage: kube-score check [--flag1 --flag2] file1 file2 ...

Use "-" as filename to read from STDIN.`)
	}

	var allFilePointers []io.Reader

	for _, file := range filesToRead {
		var fp io.Reader

		if file == "-" {
			fp = os.Stdin
		} else {
			var err error
			fp, err = os.Open(file)
			if err != nil {
				return err
			}
		}

		allFilePointers = append(allFilePointers, fp)
	}

	ignoredTests := listToStructMap(ignoreTests)
	enabledOptionalTests := listToStructMap(optionalTests)

	cnf := config.Configuration{
		AllFiles:                              allFilePointers,
		VerboseOutput:                         *verboseOutput,
		IgnoreContainerCpuLimitRequirement:    *ignoreContainerCpuLimit,
		IgnoreContainerMemoryLimitRequirement: *ignoreContainerMemoryLimit,
		IgnoredTests:                          ignoredTests,
		EnabledOptionalTests:                  enabledOptionalTests,
		UseIgnoreChecksAnnotation:             !*disableIgnoreChecksAnnotation,
	}

	parsedFiles, err := parser.ParseFiles(cnf)
	if err != nil {
		return err
	}

	scoreCard, err := score.Score(parsedFiles, cnf)
	if err != nil {
		return err
	}

	var exitCode int
	if scoreCard.AnyBelowOrEqualToGrade(scorecard.GradeCritical) {
		exitCode = 1
	} else if *exitOneOnWarning && scoreCard.AnyBelowOrEqualToGrade(scorecard.GradeWarning) {
		exitCode = 1
	} else {
		exitCode = 0
	}

	var r io.Reader

	version := getOutputVersion(*outputVersion, *outputFormat)

	if *outputFormat == "json" && version == "v1" {
		d, _ := json.MarshalIndent(scoreCard, "", "    ")
		w := bytes.NewBufferString("")
		w.WriteString(string(d))
		r = w
	} else if *outputFormat == "json" && version == "v2" {
		r = json_v2.Output(scoreCard)
	} else if *outputFormat == "human" && version == "v1" {
		termWidth, _, err := terminal.GetSize(int(os.Stdin.Fd()))
		// Assume a width of 80 if it can't be detected
		if err != nil {
			termWidth = 80
		}
		r = human.Human(scoreCard, *verboseOutput, termWidth)
	} else if *outputFormat == "ci" && version == "v1" {
		r = ci.CI(scoreCard)
	} else {
		return fmt.Errorf("error: Unknown --output-format or --output-version")
	}

	output, _ := ioutil.ReadAll(r)
	fmt.Print(string(output))
	os.Exit(exitCode)
	return nil
}

func getOutputVersion(flagValue, format string) string {
	if len(flagValue) > 0 {
		return flagValue
	}

	switch format {
	case "json":
		return "v1" // TODO: Switch this to v2 in v1.6.0
	default:
		return "v1"
	}
}

func listChecks() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	printHelp := fs.Bool("help", false, "Print help")
	setDefault(fs, "list", false)
	fs.Parse(os.Args[2:])

	if *printHelp {
		fs.Usage()
		return
	}

	allChecks := score.RegisterAllChecks(parser.Empty(), config.Configuration{})

	output := csv.NewWriter(os.Stdout)
	for _, c := range allChecks.All() {
		optionalString := "default"
		if c.Optional {
			optionalString = "optional"
		}
		output.Write([]string{c.ID, c.TargetType, c.Comment, optionalString})
	}
	output.Flush()
}

func listToStructMap(items *[]string) map[string]struct{} {
	structMap := make(map[string]struct{})
	for _, testID := range *items {
		structMap[testID] = struct{}{}
	}
	return structMap
}
