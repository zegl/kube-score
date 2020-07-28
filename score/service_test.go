package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestServiceTargetsPodDeployment(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment.yaml", "Service Targets Pod", scorecard.GradeAllOK)
}

func TestServiceNotTargetsPodDeployment(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-not-target-deployment.yaml", "Service Targets Pod", scorecard.GradeCritical)
}

func TestServiceTargetsPodRaw(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-pod.yaml", "Service Targets Pod", scorecard.GradeAllOK)
}

func TestServiceNotTargetsPodRaw(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-not-target-pod.yaml", "Service Targets Pod", scorecard.GradeCritical)
}

func TestServiceTargetsPodRawMultiLabel(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-pod-multi-label.yaml", "Service Targets Pod", scorecard.GradeAllOK)
}

func TestServiceNotTargetsPodRawMultiLabel(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-not-target-pod-multi-label.yaml", "Service Targets Pod", scorecard.GradeCritical)
}

func TestServiceTargetsPodRawSameNamespace(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-pod-same-namespace.yaml", "Service Targets Pod", scorecard.GradeAllOK)
}

func TestServiceTargetsPodRawDifferentNamespace(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-pod-different-namespace.yaml", "Service Targets Pod", scorecard.GradeCritical)
}

func TestServiceTargetsPodDeploymentSameNamespace(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment-same-namespace.yaml", "Service Targets Pod", scorecard.GradeAllOK)
}

func TestServiceTargetsPodDeploymentDifferentNamespace(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-target-deployment-different-namespace.yaml", "Service Targets Pod", scorecard.GradeCritical)
}

func TestServiceExternalName(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-externalname.yaml", "Service Targets Pod", scorecard.GradeAllOK)
}

func TestServiceTypeNodePort(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-type-nodeport.yaml", "Service Type", scorecard.GradeWarning)
}

func TestServiceTypeClusterIP(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-type-clusterip.yaml", "Service Type", scorecard.GradeAllOK)
}

func TestServiceTypeDefault(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "service-type-default.yaml", "Service Type", scorecard.GradeAllOK)
}
