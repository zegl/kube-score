package score

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
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

func TestDeploymentPodDisruptionBudgetNoPolicy(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-v1beta1-no-policy.yaml", "PodDisruptionBudget has policy", scorecard.GradeCritical)
}

func TestDeploymentPodDisruptionBudgetV1NoPolicy(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-v1-no-policy.yaml", "PodDisruptionBudget has policy", scorecard.GradeCritical)
}

func TestDeploymentPodDisruptionBudgetV1Matches(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-poddisruptionbudget-v1-matches.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeAllOK)
}

func TestDeploymentPodDisruptionBudgetV1NoMatch(t *testing.T) {
	t.Parallel()
	actual := testExpectedScore(t, "deployment-poddisruptionbudget-v1-no-match.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeCritical)

	expected := []scorecard.TestScoreComment{
		{
			Path:        "",
			Summary:     "No matching PodDisruptionBudget was found",
			Description: "It's recommended to define a PodDisruptionBudget to avoid unexpected downtime during Kubernetes maintenance operations, such as when draining a node. ",
		},
	}

	diff := cmp.Diff(expected, actual)
	assert.Empty(t, diff)
}

func TestDeploymentPodDisruptionBudgetV1NoMatchMatchInOtherNamespace(t *testing.T) {
	t.Parallel()
	actual := testExpectedScore(t, "deployment-poddisruptionbudget-v1-different-namespace.yaml", "Deployment has PodDisruptionBudget", scorecard.GradeCritical)

	expected := []scorecard.TestScoreComment{
		{
			Path:        "",
			Summary:     "No matching PodDisruptionBudget was found",
			Description: "It's recommended to define a PodDisruptionBudget to avoid unexpected downtime during Kubernetes maintenance operations, such as when draining a node. A matching budget was found, but in a different namespace. expected='foo' got='[not-foo bar]'",
		},
	}

	diff := cmp.Diff(expected, actual)
	assert.Empty(t, diff)
}
