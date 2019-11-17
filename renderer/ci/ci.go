// Package ci is currently considered to be in alpha status, and is not covered
// by the API stability guarantees
package ci

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/zegl/kube-score/scorecard"
)

// "Machine" / CI friendly output
func CI(scoreCard *scorecard.Scorecard) io.Reader {
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
				if card.Skipped {
					fmt.Fprintf(w, "[SKIPPED] %s\n",
						scoredObject.HumanFriendlyRef(),
					)
				} else {
					fmt.Fprintf(w, "[%s] %s\n",
						card.Grade.String(),
						scoredObject.HumanFriendlyRef(),
					)
				}
			}

			for _, comment := range card.Comments {
				message := comment.Summary
				if comment.Path != "" {
					message = "(" + comment.Path + ") " + comment.Summary
				}

				if card.Skipped {
					fmt.Fprintf(w, "[SKIPPED] %s: %s\n",
						scoredObject.HumanFriendlyRef(),
						message,
					)
				} else {
					fmt.Fprintf(w, "[%s] %s: %s\n",
						card.Grade.String(),
						scoredObject.HumanFriendlyRef(),
						message,
					)
				}
			}
		}
	}

	return w
}
