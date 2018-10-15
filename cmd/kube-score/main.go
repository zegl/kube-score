package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/zegl/kube-score/score"
	"github.com/zegl/kube-score/scorecard"
	"io"
	"os"
)

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	exitOneOnWarning := fs.Bool("exit-one-on-warning", false, "Exit with code 1 in case of warnings")
	ignoreContainerCpuLimit := fs.Bool("ignore-container-cpu-limit", false, "Disables the requirement of setting a container CPU limit")
	okThreshold := fs.Int("threshold-ok", 10, "The score threshold for treating an score as OK. Must be between 1 and 10 (inclusive). Scores graded below this threshold are WARNING or CRITICAL.")
	warningThreshold := fs.Int("threshold-warning", 5, "The score threshold for treating a score as WARNING. Grades below this threshold are CRITICAL. Must be between 1 and 10 (inclusive).")
	verboseOutput := fs.Bool("v", false, "Verbose output")
	printHelp := fs.Bool("help", false, "Print help")
	fs.Parse(os.Args[1:])

	if *printHelp {
		fs.Usage()
		return
	}

	if *okThreshold < 1 || *okThreshold > 10 ||
		*warningThreshold < 1 || *warningThreshold > 10 {
		fmt.Println("Error: --threshold-ok and --threshold-warning must be set to a value between 1 and 10 inclusive.")
		fs.Usage()
		os.Exit(1)
	}

	filesToRead := fs.Args()
	if len(filesToRead) == 0 {
		fmt.Println(`Error: No files given as arguments.

Usage: kube-score [--flag1 --flag2] file1 file2 ...

Use "-" as filename to read from STDIN.`)
		fmt.Println()
		fs.Usage()
		os.Exit(1)
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
				panic(err)
			}
		}

		allFilePointers = append(allFilePointers, fp)
	}

	scoreCard, err := score.Score(score.Configuration{
		AllFiles:                           allFilePointers,
		VerboseOutput:                      *verboseOutput,
		IgnoreContainerCpuLimitRequirement: *ignoreContainerCpuLimit,
	})
	if err != nil {
		panic(err)
	}

	hasWarning := false
	hasCritical := false

	for _, resourceScores := range scoreCard.Scores {
		firstCard := resourceScores[0]

		p := color.New(color.FgMagenta)

		p.Printf("%s/%s %s", firstCard.ResourceRef.Version, firstCard.ResourceRef.Kind, firstCard.ResourceRef.Name)

		if firstCard.ResourceRef.Namespace != "" {
			p.Printf(" in %s\n", firstCard.ResourceRef.Namespace)
		} else {
			p.Println()
		}

		for _, card := range resourceScores {

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

			color.New(col).Printf("    [%s] %s\n", status, card.Name)

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
		}
	}

	if hasCritical {
		os.Exit(1)
	} else if hasWarning && *exitOneOnWarning {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
