package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestStatefulSetPodDisruptionBudgetMatches(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-poddisruptionbudget-v1beta1-matches.yaml", "StatefulSet has PodDisruptionBudget", scorecard.GradeAllOK)
}

func TestStatefulSetPodDisruptionBudgetExpressionMatches(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-poddisruptionbudget-v1beta1-expression-matches.yaml", "StatefulSet has PodDisruptionBudget", scorecard.GradeAllOK)
}

func TestStatefulSetPodDisruptionBudgetExpressionNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-poddisruptionbudget-v1beta1-expression-no-match.yaml", "StatefulSet has PodDisruptionBudget", scorecard.GradeCritical)
}

func TestStatefulSetPodDisruptionBudgetNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-poddisruptionbudget-v1beta1-no-match.yaml", "StatefulSet has PodDisruptionBudget", scorecard.GradeCritical)
}

func TestDeploymentPodDisruptionBudgetMatches(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-v1beta1-matches.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeAllOK)
}

func TestDeploymentPodDisruptionBudgetExpressionMatches(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-v1beta1-expression-matches.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeAllOK)
}

func TestDeploymentPodDisruptionBudgetExpressionNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-v1beta1-expression-no-match.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeCritical)
}

func TestDeploymentPodDisruptionBudgetNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-v1beta1-no-match.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeCritical)
}

func TestDeploymentPodDisruptionBudgetV1Matches(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-v1-matches.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeAllOK)
}

func TestDeploymentPodDisruptionBudgetV1NoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-v1-no-match.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeCritical)
}
