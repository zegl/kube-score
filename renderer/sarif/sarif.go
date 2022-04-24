package sarif

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"

	"github.com/owenrumney/go-sarif/sarif"
)

func Output(input *scorecard.Scorecard) (error, io.Reader) {
	report, err := sarif.New(sarif.Version210)
	if err != nil {
		return err, nil
	}

	run := sarif.NewRun("kube-score", "https://kube-score.com/")

	addRule := func(check domain.Check) {
		run.AddRule(check.ID).WithDescription(check.Name)
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

			pb := sarif.NewPropertyBag()
			pb.Add("confidence", "High")
			pb.Add("severity", "High")

			for _, comment := range check.Comments {
				run.AddResult(check.Check.ID).
					WithLevel(level).
					WithMessage(sarif.NewTextMessage(comment.Summary)).
					WithProperties(pb.Properties).
					WithLocation(
						sarif.NewLocationWithPhysicalLocation(
							sarif.NewPhysicalLocation().
								WithArtifactLocation(
									sarif.NewSimpleArtifactLocation("file://" + v.FileLocation.Name),
								).WithRegion(
								sarif.NewSimpleRegion(
									v.FileLocation.Line,
									v.FileLocation.Line,
								),
							),
						),
					)
			}
		}
	}

	report.AddRun(run)

	j, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		return err, nil
	}
	return nil, bytes.NewBuffer(j)
}
