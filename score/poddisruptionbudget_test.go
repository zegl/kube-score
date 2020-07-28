package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestStatefulSetPodDisruptionBudgetMatches(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-poddisruptionbudget-matches.yaml", "StatefulSet has PodDisruptionBudget", scorecard.GradeAllOK)
}

func TestStatefulSetPodDisruptionBudgetExpressionMatches(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-poddisruptionbudget-expression-matches.yaml", "StatefulSet has PodDisruptionBudget", scorecard.GradeAllOK)
}

func TestStatefulSetPodDisruptionBudgetExpressionNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-poddisruptionbudget-expression-no-match.yaml", "StatefulSet has PodDisruptionBudget", scorecard.GradeCritical)
}

func TestStatefulSetPodDisruptionBudgetNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-poddisruptionbudget-no-match.yaml", "StatefulSet has PodDisruptionBudget", scorecard.GradeCritical)
}

func TestDeploymentPodDisruptionBudgetMatches(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-matches.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeAllOK)
}

func TestDeploymentPodDisruptionBudgetExpressionMatches(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-expression-matches.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeAllOK)
}

func TestDeploymentPodDisruptionBudgetExpressionNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-expression-no-match.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeCritical)
}

func TestDeploymentPodDisruptionBudgetNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-no-match.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeCritical)
}
