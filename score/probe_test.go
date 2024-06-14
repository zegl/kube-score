package score

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zegl/kube-score/scorecard"
)

func TestProbesPodAllMissing(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-all-missing.yaml", "Pod Probes", scorecard.GradeCritical)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container is missing a readinessProbe", comments[0].Summary)
}

func TestProbesServiceAccountName(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-service-account-name.yaml", "Pod Probes", scorecard.GradeCritical)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container is missing a readinessProbe", comments[0].Summary)
}

func TestProbesPodMissingReady(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-missing-ready.yaml", "Pod Probes", scorecard.GradeCritical)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container is missing a readinessProbe", comments[0].Summary)
}

func TestProbesPodIdenticalHTTP(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-identical-http.yaml", "Pod Probes", scorecard.GradeCritical)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container has the same readiness and liveness probe", comments[0].Summary)
}

func TestProbesPodIdenticalTCP(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-identical-tcp.yaml", "Pod Probes", scorecard.GradeCritical)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container has the same readiness and liveness probe", comments[0].Summary)
}

func TestProbesPodIdenticalExec(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-identical-exec.yaml", "Pod Probes", scorecard.GradeCritical)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container has the same readiness and liveness probe", comments[0].Summary)
}

func TestProbesTargetedByService(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-targeted-by-service.yaml", "Pod Probes", scorecard.GradeCritical)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container is missing a readinessProbe", comments[0].Summary)
}

func TestProbesTargetedByServiceSameNamespace(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-targeted-by-service-same-namespace.yaml", "Pod Probes", scorecard.GradeCritical)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container is missing a readinessProbe", comments[0].Summary)
}

func TestProbesTargetedByServiceSameNamespaceMultiLabels(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-targeted-by-service-same-namespace-multi-labels.yaml", "Pod Probes", scorecard.GradeCritical)
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container is missing a readinessProbe", comments[0].Summary)
}

func TestProbesTargetedByServiceDifferentNamespace(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-targeted-by-service-different-namespace.yaml", "Pod Probes", scorecard.GradeAllOK)
	assert.Len(t, comments, 1)
	assert.Equal(t, "The pod is not targeted by a service, skipping probe checks.", comments[0].Summary)
}

func TestProbesTargetedByServiceNotTargeted(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-not-targeted-by-service.yaml", "Pod Probes", scorecard.GradeAllOK)
	assert.Len(t, comments, 1)
	assert.Equal(t, "The pod is not targeted by a service, skipping probe checks.", comments[0].Summary)
}

func TestProbesTargetedByServiceNotTargetedMultiLabels(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-not-targeted-by-service-multi-labels.yaml", "Pod Probes", scorecard.GradeAllOK)
	assert.Len(t, comments, 1)
	assert.Equal(t, "The pod is not targeted by a service, skipping probe checks.", comments[0].Summary)
}

func TestProbesMultipleContainers(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-on-different-containers.yaml", "Pod Probes", scorecard.GradeAllOK)
	assert.Len(t, comments, 0)
}

func TestProbesMultipleContainersInit(t *testing.T) {
	t.Parallel()
	comments := testExpectedScore(t, "pod-probes-on-different-containers-init.yaml", "Pod Probes", scorecard.GradeAllOK)
	assert.Len(t, comments, 0)
}
