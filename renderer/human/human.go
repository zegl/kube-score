// Package human is currently considered to be in alpha status, and is not covered
// by the API stability guarantees
package human

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/eidolon/wordwrap"
	"github.com/fatih/color"

	"github.com/zegl/kube-score/scorecard"
)

func Human(scoreCard *scorecard.Scorecard, verboseOutput int, termWidth int) io.Reader {
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
		var writtenHeaderChars int
		writtenHeaderChars, _ = color.New(color.FgMagenta).Fprintf(w, "%s/%s %s", scoredObject.TypeMeta.APIVersion, scoredObject.TypeMeta.Kind, scoredObject.ObjectMeta.Name)
		if scoredObject.ObjectMeta.Namespace != "" {
			written2, _ := color.New(color.FgMagenta).Fprintf(w, " in %s", scoredObject.ObjectMeta.Namespace)
			writtenHeaderChars += written2
		}

		// Adjust to termsize
		fmt.Fprint(w, safeRepeat(" ", min(80, termWidth)-writtenHeaderChars-2))

		switch {
		case scoredObject.AnyBelowOrEqualToGrade(scorecard.GradeCritical):
			fmt.Fprintf(w, "💥\n")
		case scoredObject.AnyBelowOrEqualToGrade(scorecard.GradeWarning):
			fmt.Fprintf(w, "🤔\n")
		default:
			fmt.Fprintf(w, "✅\n")
		}

		for _, card := range scoredObject.Checks {
			r := outputHumanStep(card, verboseOutput, termWidth)
			_, err := io.Copy(w, r)
			if err != nil {
				panic(err)
			}

		}
	}

	return w
}

func outputHumanStep(card scorecard.TestScore, verboseOutput int, termWidth int) io.Reader {
	w := bytes.NewBufferString("")

	// Only print skipped items if verbosity is at least 2
	if card.Skipped && verboseOutput < 2 {
		return w
	}

	var col color.Attribute

	switch {
	case card.Skipped || card.Grade >= scorecard.GradeAllOK:
		// Higher than or equal to --threshold-ok
		col = color.FgGreen

		// If verbose output is disabled, skip OK items in the output
		if verboseOutput == 0 {
			return w
		}

	case card.Grade >= scorecard.GradeWarning:
		// Higher than or equal to --threshold-warning
		col = color.FgYellow
	default:
		// All lower than both --threshold-ok and --threshold-warning are critical
		col = color.FgRed
	}

	if card.Skipped {
		color.New(col).Fprintf(w, "    [SKIPPED] %s\n", card.Check.Name)
	} else {
		color.New(col).Fprintf(w, "    [%s] %s\n", card.Grade.String(), card.Check.Name)
	}

	for _, comment := range card.Comments {
		fmt.Fprintf(w, "        · ")

		if len(comment.Path) > 0 {
			fmt.Fprintf(w, "%s -> ", comment.Path)
		}

		fmt.Fprint(w, comment.Summary)

		if len(comment.Description) > 0 {
			wrapWidth := termWidth - 12
			if wrapWidth < 40 {
				wrapWidth = 40
			}
			wrapper := wordwrap.Wrapper(wrapWidth, false)
			wrapped := wrapper(comment.Description)
			fmt.Fprintln(w)
			fmt.Fprint(w, wordwrap.Indent(wrapped, strings.Repeat(" ", 12), false))
		}

		if len(comment.DocumentationURL) > 0 {
			fmt.Fprintln(w)
			fmt.Fprintf(w, "%sMore information: %s", strings.Repeat(" ", 12), comment.DocumentationURL)
		}

		fmt.Fprintln(w)
	}

	return w
}

func safeRepeat(s string, count int) string {
	if count < 0 {
		return ""
	}
	return strings.Repeat(s, count)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
