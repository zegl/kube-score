package sarif

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/sarif"
	"github.com/zegl/kube-score/scorecard"
)

func Output(input *scorecard.Scorecard) io.Reader {
	var results []sarif.Results
	var rules []sarif.Rules

	addRule := func(check domain.Check) {
		for _, r := range rules {
			if r.ID == check.ID {
				return
			}
		}

		rules = append(rules, sarif.Rules{
			ID:   check.ID,
			Name: check.Name,
		})
	}

	for _, v := range *input {
		for _, check := range v.Checks {
			if check.Skipped {
				continue
			}

			var level string
			switch check.Grade {
			case scorecard.GradeCritical:
				level = "error"
			case scorecard.GradeWarning:
				level = "warning"
			default:
				continue
			}

			addRule(check.Check)

			for _, comment := range check.Comments {
				results = append(results, sarif.Results{
					Message: sarif.Message{
						Text: comment.Summary,
					},
					RuleID: check.Check.ID,
					Level:  level,
					Properties: sarif.ResultsProperties{
						IssueConfidence: "HIGH",
						IssueSeverity:   "HIGH",
					},
					Locations: []sarif.Locations{
						{
							PhysicalLocation: sarif.PhysicalLocation{
								ArtifactLocation: sarif.ArtifactLocation{
									URI: "file://" + v.FileLocation.Name,
								},
								ContextRegion: sarif.ContextRegion{
									StartLine: v.FileLocation.Line,
								},
							},
						},
					},
				})
			}
		}
	}

	run := sarif.Run{
		Tool: sarif.Tool{
			Driver: sarif.Driver{
				Name:  "kube-score",
				Rules: rules,
			},
		},
		Results: results,
	}
	res := sarif.Sarif{
		Runs:    []sarif.Run{run},
		Version: "2.1.0",
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
	}

	j, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(j)
}
