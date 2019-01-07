package score

import (
	"testing"
)

func TestDeploymentHasPodAntiAffinityPreffered(t *testing.T) {
	testExpectedScore(t, "deployment-host-antiaffinity-preffered.yaml", "Deployment has host PodAntiAffinity", 10)
}

func TestDeploymentHasPodAntiAffinityRequired(t *testing.T) {
	testExpectedScore(t, "deployment-host-antiaffinity-required.yaml", "Deployment has host PodAntiAffinity", 10)
}

func TestDeploymentHasPodAntiAffinityNotSet(t *testing.T) {
	testExpectedScore(t, "deployment-host-antiaffinity-not-set.yaml", "Deployment has host PodAntiAffinity", 5)
}

func TestDeploymentHasPodAntiAffinityOneReplica(t *testing.T) {
	testExpectedScore(t, "deployment-host-antiaffinity-1-replica.yaml", "Deployment has host PodAntiAffinity", 10)
}
