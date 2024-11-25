package score

import (
	"testing"

	"github.com/zegl/kube-score/config"

	"github.com/stretchr/testify/assert"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func TestServiceTargetsDeploymentStrategyRolling(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment.yaml", "Deployment Strategy", scorecard.GradeAllOK)
}

func TestServiceNotTargetsDeploymentStrategyNotRelevant(t *testing.T) {
	t.Parallel()
	skipped := wasSkipped(t,
		[]ks.NamedReader{testFile("service-not-target-deployment.yaml")}, nil, nil,
		"Deployment Strategy")
	assert.True(t, skipped)
}

func TestServiceTargetsDeploymentStrategyNotRolling(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment-not-rolling.yaml", "Deployment Strategy", scorecard.GradeWarning)
}

func TestServiceTargetsDeploymentStrategyNotSet(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment-strategy-not-set.yaml", "Deployment Strategy", scorecard.GradeAllOK)
}

func TestServiceTargetsDeploymentReplicasOk(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment.yaml", "Deployment Replicas", scorecard.GradeAllOK)
}

func TestServiceNotTargetsDeploymentReplicasNotRelevant(t *testing.T) {
	t.Parallel()
	assert.True(t, wasSkipped(t,
		[]ks.NamedReader{testFile("service-not-target-deployment.yaml")}, nil, nil,
		"Deployment Replicas"))

	summaries := getSummaries(t, []ks.NamedReader{testFile("service-not-target-deployment.yaml")}, nil, nil,
		"Deployment Replicas")
	assert.Contains(t, summaries, "Skipped as the Deployment is not targeted by service")
}

func TestServiceTargetsDeploymentReplicasNok(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("service-target-deployment-replica-1.yaml")}, nil, &config.RunConfiguration{
		MinReplicasDeployment: 2,
	}, "Deployment Replicas", scorecard.GradeWarning)
}

func TestServiceTargetsDeploymentReplicasNok1(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("service-target-deployment-replica-1.yaml")}, nil, &config.RunConfiguration{
		MinReplicasDeployment: 1,
	}, "Deployment Replicas", scorecard.GradeAllOK)
}

func TestServiceTargetsDeploymentReplicasNok5(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("service-target-deployment-replica-1.yaml")}, nil, &config.RunConfiguration{
		MinReplicasDeployment: 5,
	}, "Deployment Replicas", scorecard.GradeWarning)
}

func TestHPATargetsDeployment(t *testing.T) {
	t.Parallel()
	assert.True(t, wasSkipped(t,
		[]ks.NamedReader{testFile("hpa-target-deployment.yaml")}, nil, nil,
		"Deployment Replicas"))

	summaries := getSummaries(t, []ks.NamedReader{testFile("hpa-target-deployment.yaml")}, nil, nil,
		"Deployment Replicas")
	assert.Contains(t, summaries, "Skipped as the Deployment is controlled by a HorizontalPodAutoscaler")
}
