package junit

import (
	"bytes"
	"io"

	"github.com/jstemmer/go-junit-report/v2/junit"
	"github.com/zegl/kube-score/scorecard"
)

// JUnit XML output
func JUnit(scoreCard *scorecard.Scorecard) io.Reader {
	testSuites := junit.Testsuites{
		Name: "kube-score",
	}

	for _, scoredObject := range *scoreCard {
		testsuite := junit.Testsuite{
			Name: scoredObject.HumanFriendlyRef(),
		}

		for _, testScore := range scoredObject.Checks {
			if len(testScore.Comments) == 0 {
				if testScore.Skipped {
					testsuite.AddTestcase(junit.Testcase{
						Name:      testScore.Check.Name,
						Classname: scoredObject.HumanFriendlyRef(),
						Skipped:   &junit.Result{},
					})
				} else {
					if testScore.Grade == scorecard.GradeAlmostOK || testScore.Grade == scorecard.GradeAllOK {
						testsuite.AddTestcase(junit.Testcase{
							Name:      testScore.Check.Name,
							Classname: scoredObject.HumanFriendlyRef(),
						})
					} else {
						testsuite.AddTestcase(junit.Testcase{
							Name:      testScore.Check.Name,
							Classname: scoredObject.HumanFriendlyRef(),
							Failure:   &junit.Result{},
						})
					}
				}
			} else {
				for _, comment := range testScore.Comments {
					message := comment.Summary
					if comment.Path != "" {
						message = "(" + comment.Path + ") " + comment.Summary
					}

					if testScore.Skipped {
						testsuite.AddTestcase(junit.Testcase{
							Name:      testScore.Check.Name,
							Classname: scoredObject.HumanFriendlyRef(),
							Skipped: &junit.Result{
								Message: message,
							},
						})
					} else {
						if testScore.Grade == scorecard.GradeAlmostOK || testScore.Grade == scorecard.GradeAllOK {
							testsuite.AddTestcase(junit.Testcase{
								Name:      testScore.Check.Name,
								Classname: scoredObject.HumanFriendlyRef(),
							})
						} else {
							testsuite.AddTestcase(junit.Testcase{
								Name:      testScore.Check.Name,
								Classname: scoredObject.HumanFriendlyRef(),
								Failure: &junit.Result{
									Message: message,
								},
							})
						}

					}
				}
			}
		}
		testSuites.AddSuite(testsuite)
	}

	buffer := &bytes.Buffer{}
	err := testSuites.WriteXML(buffer)
	if err != nil {
		panic(err)
	}
	return buffer
}
