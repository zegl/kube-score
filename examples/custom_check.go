package examples

import (
	"bytes"
	"io"
	"strings"

	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/score"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"

	v1 "k8s.io/api/apps/v1"
)

type namedReader struct {
	io.Reader
	name string
}

func (n namedReader) Name() string {
	return n.name
}

// ExampleCheckObject shows how kube-score can be extended with a custom check function
//
// In this example, raw is a YAML encoded Kubernetes object
func ExampleCheckObject(raw []byte) (*scorecard.Scorecard, error) {
	parser, err := parser.New(nil)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(raw)

	// Parse all objects to read
	allObjects, err := parser.ParseFiles(
		[]domain.NamedReader{
			namedReader{
				Reader: reader,
				name:   "input",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	// Register check functions to run
	checks := checks.New(nil)
	checks.RegisterDeploymentCheck("custom-deployment-check", "A custom kube-score check function", customDeploymentCheck)

	return score.Score(allObjects, checks, &config.RunConfiguration{})
}

func customDeploymentCheck(d v1.Deployment) (scorecard.TestScore, error) {
	if strings.Contains(d.Name, "foo") {
		return scorecard.TestScore{
			Grade: scorecard.GradeCritical,
			Comments: []scorecard.TestScoreComment{{
				Summary: "Deployments names can not contian 'foo'",
			}}}, nil
	}

	return scorecard.TestScore{Grade: scorecard.GradeAllOK}, nil
}
