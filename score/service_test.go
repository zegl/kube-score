package score

import "testing"

func TestServiceTargetsPodDeployment(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment.yaml", "Service Targets Pod", 10)
}

func TestServiceNotTargetsPodDeployment(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-not-target-deployment.yaml", "Service Targets Pod", 1)
}

func TestServiceTargetsPodRaw(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-pod.yaml", "Service Targets Pod", 10)
}

func TestServiceNotTargetsPodRaw(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-not-target-pod.yaml", "Service Targets Pod", 1)
}

func TestServiceTargetsPodRawMultiLabel(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-pod-multi-label.yaml", "Service Targets Pod", 10)
}

func TestServiceNotTargetsPodRawMultiLabel(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-not-target-pod-multi-label.yaml", "Service Targets Pod", 1)
}

func TestServiceTargetsPodRawSameNamespace(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-pod-same-namespace.yaml", "Service Targets Pod", 10)
}

func TestServiceTargetsPodRawDifferentNamespace(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-pod-different-namespace.yaml", "Service Targets Pod", 1)
}

func TestServiceTargetsPodDeploymentSameNamespace(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment-same-namespace.yaml", "Service Targets Pod", 10)
}

func TestServiceTargetsPodDeploymentDifferentNamespace(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment-different-namespace.yaml", "Service Targets Pod", 1)
}

func TestServiceExternalName(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-externalname.yaml", "Service Targets Pod", 10)
}

func TestServiceTypeNodePort(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-type-nodeport.yaml", "Service Type", 5)
}

func TestServiceTypeClusterIP(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-type-clusterip.yaml", "Service Type", 10)
}

func TestServiceTypeDefault(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-type-default.yaml", "Service Type", 10)
}
