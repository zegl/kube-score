package score

import (
	"github.com/stretchr/testify/assert"
	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestServiceTargetsDeploymentStrategyRolling(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment.yaml", "Deployment Strategy", scorecard.GradeAllOK)
}

func TestServiceNotTargetsDeploymentStrategyNotRelevant(t *testing.T) {
	t.Parallel()
	// Expecting score 0 as it didn't cot rated
	skipped := wasSkipped(t, config.Configuration{
		AllFiles: []ks.NamedReader{testFile("service-not-target-deployment.yaml")},
	}, "Deployment Strategy")
	assert.True(t, skipped)
}

func TestServiceTargetsDeploymentStrategyNotRolling(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment-not-rolling.yaml", "Deployment Strategy", scorecard.GradeWarning)
}

func TestServiceTargetsDeploymentStrategyNotSet(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment-strategy-not-set.yaml", "Deployment Strategy", scorecard.GradeWarning)
}
