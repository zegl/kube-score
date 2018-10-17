package score

import (
	"testing"
)

func TestPodHasNoMatchingNetworkPolicy(t *testing.T) {
	testExpectedScore(t, "networkpolicy-not-matching.yaml", "Pod NetworkPolicy", 1)
}

func TestPodHasMatchingNetworkPolicy(t *testing.T) {
	testExpectedScore(t, "networkpolicy-matching.yaml", "Pod NetworkPolicy", 10)
}

func TestPodHasMatchingIngressNetworkPolicy(t *testing.T) {
	testExpectedScore(t, "networkpolicy-matching-only-ingress.yaml", "Pod NetworkPolicy", 5)
}

func TestPodHasMatchingEgressNetworkPolicy(t *testing.T) {
	testExpectedScore(t, "networkpolicy-matching-only-egress.yaml", "Pod NetworkPolicy", 5)
}

func TestNetworkPolicyTargetsPod(t *testing.T) {
	testExpectedScore(t, "networkpolicy-targets-pod.yaml", "NetworkPolicy targets Pod", 10)
}

func TestNetworkPolicyTargetsPodInDeployment(t *testing.T) {
	testExpectedScore(t, "networkpolicy-targets-pod-deployment.yaml", "NetworkPolicy targets Pod", 10)
}

func TestNetworkPolicyTargetsPodNotMatching(t *testing.T) {
	testExpectedScore(t, "networkpolicy-targets-pod-not-matching.yaml", "NetworkPolicy targets Pod", 1)
}
