package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestPodHasNoMatchingNetworkPolicy(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-not-matching.yaml", "Pod NetworkPolicy", 1)
}

func TestPodHasMatchingNetworkPolicy(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-matching.yaml", "Pod NetworkPolicy", 10)
}

func TestPodHasMatchingIngressNetworkPolicy(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-matching-only-ingress.yaml", "Pod NetworkPolicy", 5)
}

func TestPodHasMatchingEgressNetworkPolicy(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-matching-only-egress.yaml", "Pod NetworkPolicy", 5)
}

func TestNetworkPolicyTargetsPod(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-targets-pod.yaml", "NetworkPolicy targets Pod", 10)
}

func TestNetworkPolicyTargetsPodInDeployment(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-targets-pod-deployment.yaml", "NetworkPolicy targets Pod", 10)
}

func TestNetworkPolicyTargetsPodNotMatching(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-targets-pod-not-matching.yaml", "NetworkPolicy targets Pod", 1)
}

func TestNetworkPolicyDeploymentNamespaceMatching(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-deployment-matching.yaml", "NetworkPolicy targets Pod", scorecard.GradeAllOK)
	testExpectedScore(t, "networkpolicy-deployment-matching.yaml", "Pod NetworkPolicy", scorecard.GradeAllOK)
}

func TestNetworkPolicyStatefulSetNamespaceMatching(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-statefulset-matching.yaml", "NetworkPolicy targets Pod", scorecard.GradeAllOK)
	testExpectedScore(t, "networkpolicy-statefulset-matching.yaml", "Pod NetworkPolicy", scorecard.GradeAllOK)
}

func TestNetworkPolicyCronJobNamespaceMatching(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-cronjob-matching.yaml", "NetworkPolicy targets Pod", scorecard.GradeAllOK)
	testExpectedScore(t, "networkpolicy-cronjob-matching.yaml", "Pod NetworkPolicy", scorecard.GradeAllOK)
}

func TestNetworkPolicyDeploymentNamespaceNotMatchingSelector(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-deployment-not-matching-selector.yaml", "NetworkPolicy targets Pod", scorecard.GradeCritical)
	testExpectedScore(t, "networkpolicy-deployment-not-matching-selector.yaml", "Pod NetworkPolicy", scorecard.GradeCritical)
}

func TestNetworkPolicyStatefulSetNamespaceNotMatchingSelector(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-statefulset-not-matching-selector.yaml", "NetworkPolicy targets Pod", scorecard.GradeCritical)
	testExpectedScore(t, "networkpolicy-statefulset-not-matching-selector.yaml", "Pod NetworkPolicy", scorecard.GradeCritical)
}

func TestNetworkPolicyCronJobNamespaceNotMatchingSelector(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "networkpolicy-cronjob-not-matching-selector.yaml", "NetworkPolicy targets Pod", scorecard.GradeCritical)
	testExpectedScore(t, "networkpolicy-cronjob-not-matching-selector.yaml", "Pod NetworkPolicy", scorecard.GradeCritical)
}
