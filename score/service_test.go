package score

import "testing"

func TestServiceTargetsPodDeployment(t *testing.T) {
	testExpectedScore(t, "service-target-deployment.yaml", "Service Targets Pod", 10)
}

func TestServiceNotTargetsPodDeployment(t *testing.T) {
	testExpectedScore(t, "service-not-target-deployment.yaml", "Service Targets Pod", 0)
}

func TestServiceTargetsPodRaw(t *testing.T) {
	testExpectedScore(t, "service-target-pod.yaml", "Service Targets Pod", 10)
}

func TestServiceNotTargetsPodRaw(t *testing.T) {
	testExpectedScore(t, "service-not-target-pod.yaml", "Service Targets Pod", 0)
}

func TestServiceTargetsPodRawMultiLabel(t *testing.T) {
	testExpectedScore(t, "service-target-pod-multi-label.yaml", "Service Targets Pod", 10)
}

func TestServiceNotTargetsPodRawMultiLabel(t *testing.T) {
	testExpectedScore(t, "service-not-target-pod-multi-label.yaml", "Service Targets Pod", 0)
}

func TestServiceTargetsPodRawSameNamespace(t *testing.T) {
	testExpectedScore(t, "service-target-pod-same-namespace.yaml", "Service Targets Pod", 10)
}

func TestServiceTargetsPodRawDifferentNamespace(t *testing.T) {
	testExpectedScore(t, "service-target-pod-different-namespace.yaml", "Service Targets Pod", 0)
}