package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

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
	helpName := execName(os.Args[0])

	fs := flag.NewFlagSet(helpName, flag.ExitOnError)
	setDefault(fs, helpName, "", true)

	command := ""
	cmdArgsOffset := 2

	// When executing kube-score as a kubectl plugin, default to the "score" sub-command to avoid stuttering
	// "kubectl score" is equivalent to "kubectl score score"
	if isKubectlPlugin(helpName) {
		command = "score"
		cmdArgsOffset = 1
	}

	// No command, flag, or file has been specified
	if len(os.Args) <= cmdArgsOffset {
		fs.Usage()
		return
	}

	type cmdFunc func(string, []string)

	cmds := map[string]cmdFunc{
		"score": func(helpName string, args []string) {
			if err := scoreFiles(helpName, args); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to score files: %v", err)
				os.Exit(1)
			}
		},

		"list": func(helpName string, args []string) {
			listChecks(helpName, args)
		},

		"version": func(helpName string, args []string) {
			cmdVersion()
		},

		"help": func(helpName string, args []string) {
			fs.Usage()
			os.Exit(1)
		},
	}

	// If arg 1 is set and is a valid command, always use it as the command to execute, instead of the default
	if _, ok := cmds[os.Args[1]]; ok {
		command = os.Args[1]
		cmdArgsOffset = 2
	}

	// Execute the command, or the help command if no matching command is found
	if ex, ok := cmds[command]; ok {
		ex(helpName, os.Args[cmdArgsOffset:])
	} else {
		cmds["help"](helpName, os.Args[cmdArgsOffset:])
	}
}

func execName(args0 string) string {
	// Detect name of the binary
	binName := path.Base(args0)

	// If executed as a kubectl plugin, replace dash with a space
	// "kubectl-score" -> "kubectl score"
	if strings.HasPrefix(binName, "kubectl-") {
		binName = strings.Replace(binName, "kubectl-", "kubectl ", 1)
	}

	return binName
}

func isKubectlPlugin(helpName string) bool {
	return execName(helpName) == "kubectl score"
}

func setDefault(fs *flag.FlagSet, binName, actionName string, displayForMoreInfo bool) {
	fs.Usage = func() {
		usage := fmt.Sprintf(`Usage of %s:
%s [action] --flags

Actions:
	score	Checks all files in the input, and gives them a score and recommendations
	list	Prints a CSV list of all available score checks
	version	Print the version of kube-score
	help	Print this message`+"\n\n", binName, binName)

		if displayForMoreInfo {
			usage += fmt.Sprintf(`Run "%s [action] --help" for more information about a particular command`, binName)
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

func scoreFiles(binName string, args []string) error {
	fs := flag.NewFlagSet(binName, flag.ExitOnError)
	exitOneOnWarning := fs.Bool("exit-one-on-warning", false, "Exit with code 1 in case of warnings")
	ignoreContainerCpuLimit := fs.Bool("ignore-container-cpu-limit", false, "Disables the requirement of setting a container CPU limit")
	ignoreContainerMemoryLimit := fs.Bool("ignore-container-memory-limit", false, "Disables the requirement of setting a container memory limit")
	verboseOutput := fs.CountP("verbose", "v", "Enable verbose output, can be set multiple times for increased verbosity.")
	printHelp := fs.Bool("help", false, "Print help")
	outputFormat := fs.StringP("output-format", "o", "human", "Set to 'human', 'json' or 'ci'. If set to ci, kube-score will output the program in a format that is easier to parse by other programs.")
	outputVersion := fs.String("output-version", "", "Changes the version of the --output-format. The 'json' format has version 'v2' (default) and 'v1' (deprecated, will be removed in v1.7.0). The 'human' and 'ci' formats has only version 'v1' (default). If not explicitly set, the default version for that particular output format will be used.")
	optionalTests := fs.StringSlice("enable-optional-test", []string{}, "Enable an optional test, can be set multiple times")
	ignoreTests := fs.StringSlice("ignore-test", []string{}, "Disable a test, can be set multiple times")
	disableIgnoreChecksAnnotation := fs.Bool("disable-ignore-checks-annotations", false, "Set to true to disable the effect of the 'kube-score/ignore' annotations")
	setDefault(fs, binName, "score", false)

	err := fs.Parse(args)
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

Usage: %s score [--flag1 --flag2] file1 file2 ...

Use "-" as filename to read from STDIN.`, execName(binName))
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
		return "v2"
	default:
		return "v1"
	}
}

func listChecks(binName string, args []string) {
	fs := flag.NewFlagSet(binName, flag.ExitOnError)
	printHelp := fs.Bool("help", false, "Print help")
	setDefault(fs, binName, "list", false)
	fs.Parse(args)

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
