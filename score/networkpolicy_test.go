package score

import (
	"testing"
)

func TestPodHasNoMatchingNetworkPolicy(t *testing.T) {
	testExpectedScore(t, "networkpolicy-not-matching.yaml", "Pod NetworkPolicy", 0)
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
