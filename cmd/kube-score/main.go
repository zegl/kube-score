package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	flag "github.com/spf13/pflag"
	"io"
	"os"

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
	okThreshold := fs.Int("threshold-ok", 10, "The score threshold for treating an score as OK. Must be between 1 and 10 (inclusive). Scores graded below this threshold are WARNING or CRITICAL.")
	warningThreshold := fs.Int("threshold-warning", 5, "The score threshold for treating a score as WARNING. Grades below this threshold are CRITICAL. Must be between 1 and 10 (inclusive).")
	verboseOutput := fs.Bool("v", false, "Verbose output")
	printHelp := fs.Bool("help", false, "Print help")
	outputFormat := fs.String("output-format", "human", "Set to 'human' or 'ci'. If set to ci, kube-score will output the program in a format that is easier to parse by other programs.")
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

	if *okThreshold < 1 || *okThreshold > 10 ||
		*warningThreshold < 1 || *warningThreshold > 10 {
		fs.Usage()
		return fmt.Errorf("Error: --threshold-ok and --threshold-warning must be set to a value between 1 and 10 inclusive.")
	}

	if *outputFormat != "human" && *outputFormat != "ci" {
		fs.Usage()
		return fmt.Errorf("Error: --output-format must be set to: 'human' or 'ci'")
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

	cnf := config.Configuration{
		AllFiles:                              allFilePointers,
		VerboseOutput:                         *verboseOutput,
		IgnoreContainerCpuLimitRequirement:    *ignoreContainerCpuLimit,
		IgnoreContainerMemoryLimitRequirement: *ignoreContainerMemoryLimit,
	}

	parsedFiles, err := parser.ParseFiles(cnf)
	if err != nil {
		panic(err)
	}

	scoreCard, err := score.Score(parsedFiles, cnf)
	if err != nil {
		panic(err)
	}

	ignoredTests := make(map[string]struct{})
	for _, testID := range *ignoreTests {
		ignoredTests[testID] = struct{}{}
	}

	hasWarning := false
	hasCritical := false

	// Detect which output format we should use
	humanOutput := *outputFormat == "human"

	for _, scoredObject := range scoreCard {
		// Headers for each object
		if humanOutput {
			color.New(color.FgMagenta).Printf("%s/%s %s", scoredObject.TypeMeta.APIVersion, scoredObject.TypeMeta.Kind, scoredObject.ObjectMeta.Name)
			if scoredObject.ObjectMeta.Namespace != "" {
				color.New(color.FgMagenta).Printf(" in %s\n", scoredObject.ObjectMeta.Namespace)
			} else {
				fmt.Println()
			}
		}

		for _, card := range scoredObject.Checks {
			if _, ok := ignoredTests[card.Check.ID]; ok {
				continue
			}

			var col color.Attribute
			var status string

			if card.Grade >= scorecard.Grade(*okThreshold) {
				// Higher than or equal to --threshold-ok
				col = color.FgGreen
				status = "OK"
			} else if card.Grade >= scorecard.Grade(*warningThreshold) {
				// Higher than or equal to --threshold-warning
				col = color.FgYellow
				status = "WARNING"
				hasWarning = true
			} else {
				// All lower than both --threshold-ok and --threshold-warning are critical
				col = color.FgRed
				status = "CRITICAL"
				hasCritical = true
			}

			if humanOutput {
				color.New(col).Printf("    [%s] %s\n", status, card.Check.Name)

				for _, comment := range card.Comments {
					fmt.Printf("        * ")

					if len(comment.Path) > 0 {
						fmt.Printf("%s -> ", comment.Path)
					}

					fmt.Print(comment.Summary)

					if len(comment.Description) > 0 {
						fmt.Printf("\n             %s", comment.Description)
					}

					fmt.Println()
				}
			} else {
				// "Machine" / CI friendly output
				for _, comment := range card.Comments {
					message := comment.Summary
					if comment.Path != "" {
						message = "(" + comment.Path + ") " + comment.Summary
					}

					fmt.Printf("[%s] %s: %s\n",
						status,
						scoredObject.HumanFriendlyRef(),
						message,
					)
				}

			}
		}
	}

	if hasCritical {
		os.Exit(1)
	} else if hasWarning && *exitOneOnWarning {
		os.Exit(1)
	} else {
		os.Exit(0)
	}

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
		output.Write([]string{c.ID, c.TargetType, c.Comment})
	}
	output.Flush()
}
