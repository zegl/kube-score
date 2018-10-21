package score

import "testing"

func TestStatefulSetPodDisruptionBudgetMatches(t *testing.T) {
	testExpectedScore(t, "statefulset-poddisruptionbudget-matches.yaml", "StatefulSet has PodDisruptionBudget", 10)
}

func TestStatefulSetPodDisruptionBudgetExpressionMatches(t *testing.T) {
	testExpectedScore(t, "statefulset-poddisruptionbudget-expression-matches.yaml", "StatefulSet has PodDisruptionBudget", 10)
}

func TestStatefulSetPodDisruptionBudgetExpressionNoMatch(t *testing.T) {
	testExpectedScore(t, "statefulset-poddisruptionbudget-expression-no-match.yaml", "StatefulSet has PodDisruptionBudget", 1)
}

func TestStatefulSetPodDisruptionBudgetNoMatch(t *testing.T) {
	testExpectedScore(t, "statefulset-poddisruptionbudget-no-match.yaml", "StatefulSet has PodDisruptionBudget", 1)
}

func TestDeploymentPodDisruptionBudgetMatches(t *testing.T) {
	testExpectedScore(t, "deployment-poddisruptionbudget-matches.yaml", "Deployment has PodDisruptionBudget", 10)
}

func TestDeploymentPodDisruptionBudgetExpressionMatches(t *testing.T) {
	testExpectedScore(t, "deployment-poddisruptionbudget-expression-matches.yaml", "Deployment has PodDisruptionBudget", 10)
}

func TestDeploymentPodDisruptionBudgetExpressionNoMatch(t *testing.T) {
	testExpectedScore(t, "deployment-poddisruptionbudget-expression-no-match.yaml", "Deployment has PodDisruptionBudget", 1)
}

func TestDeploymentPodDisruptionBudgetNoMatch(t *testing.T) {
	testExpectedScore(t, "deployment-poddisruptionbudget-no-match.yaml", "Deployment has PodDisruptionBudget", 1)
}
