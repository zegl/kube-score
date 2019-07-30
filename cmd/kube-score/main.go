package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	flag "github.com/spf13/pflag"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/parser"
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
	verboseOutput := fs.Bool("v", false, "Verbose output")
	printHelp := fs.Bool("help", false, "Print help")
	outputFormat := fs.String("output-format", "human", "Set to 'human', 'json' or 'ci'. If set to ci, kube-score will output the program in a format that is easier to parse by other programs.")
	optionalTests := fs.StringSlice("enable-optional-test", []string{}, "Enable an optional test, can be set multiple times")
	ignoreTests := fs.StringSlice("ignore-test", []string{}, "Disable a test, can be set multiple times")
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

	if *outputFormat == "json" {
		d, _ := json.MarshalIndent(scoreCard, "", "    ")
		w := bytes.NewBufferString("")
		w.WriteString(string(d))
		r = w
	} else if *outputFormat == "human" {
		r = outputHuman(scoreCard)
	} else {
		r = outputCi(scoreCard)
	}

	output, _ := ioutil.ReadAll(r)
	fmt.Print(string(output))
	os.Exit(exitCode)
	return nil
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

func outputHuman(scoreCard *scorecard.Scorecard) io.Reader {
	// Print the items sorted by scorecard key
	var keys []string
	for k := range *scoreCard {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := bytes.NewBufferString("")

	for _, key := range keys {
		scoredObject := (*scoreCard)[key]

		// Headers for each object
		color.New(color.FgMagenta).Fprintf(w, "%s/%s %s", scoredObject.TypeMeta.APIVersion, scoredObject.TypeMeta.Kind, scoredObject.ObjectMeta.Name)
		if scoredObject.ObjectMeta.Namespace != "" {
			color.New(color.FgMagenta).Fprintf(w, " in %s\n", scoredObject.ObjectMeta.Namespace)
		} else {
			fmt.Fprintln(w)
		}

		for _, card := range scoredObject.Checks {
			r := outputHumanStep(card)
			io.Copy(w, r)
		}

	}

	return w
}

func outputHumanStep(card scorecard.TestScore) io.Reader {
	var col color.Attribute

	if card.Grade >= scorecard.GradeAllOK {
		// Higher than or equal to --threshold-ok
		col = color.FgGreen
	} else if card.Grade >= scorecard.GradeWarning {
		// Higher than or equal to --threshold-warning
		col = color.FgYellow
	} else {
		// All lower than both --threshold-ok and --threshold-warning are critical
		col = color.FgRed
	}

	w := bytes.NewBufferString("")

	color.New(col).Fprintf(w, "    [%s] %s\n", card.Grade.String(), card.Check.Name)

	for _, comment := range card.Comments {
		fmt.Fprintf(w, "        * ")

		if len(comment.Path) > 0 {
			fmt.Fprintf(w, "%s -> ", comment.Path)
		}

		fmt.Fprint(w, comment.Summary)

		if len(comment.Description) > 0 {
			fmt.Fprintf(w, "\n             %s", comment.Description)
		}

		fmt.Fprintln(w)
	}

	return w
}

// "Machine" / CI friendly output
func outputCi(scoreCard *scorecard.Scorecard) io.Reader {
	w := bytes.NewBufferString("")

	// Print the items sorted by scorecard key
	var keys []string
	for k := range *scoreCard {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		scoredObject := (*scoreCard)[key]

		for _, card := range scoredObject.Checks {
			if len(card.Comments) == 0 {
				fmt.Fprintf(w, "[%s] %s\n",
					card.Grade.String(),
					scoredObject.HumanFriendlyRef(),
				)
			}

			for _, comment := range card.Comments {
				message := comment.Summary
				if comment.Path != "" {
					message = "(" + comment.Path + ") " + comment.Summary
				}

				fmt.Fprintf(w, "[%s] %s: %s\n",
					card.Grade.String(),
					scoredObject.HumanFriendlyRef(),
					message,
				)
			}
		}
	}

	return w
}

func listToStructMap(items *[]string) map[string]struct{} {
	structMap := make(map[string]struct{})
	for _, testID := range *items {
		structMap[testID] = struct{}{}
	}
	return structMap
}
