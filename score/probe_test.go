package score

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPodProbesAllMissing(t *testing.T) {
	testExpectedScore(t, "pod-probes-all-missing.yaml", "Pod Probes", 0)
}

func TestPodProbesMissingReady(t *testing.T) {
	testExpectedScore(t, "pod-probes-missing-ready.yaml", "Pod Probes", 10)
}

func TestPodProbesIdenticalHTTP(t *testing.T) {
	testExpectedScore(t, "pod-probes-identical-http.yaml", "Pod Probes", 7)
}

func TestPodProbesIdenticalTCP(t *testing.T) {
	testExpectedScore(t, "pod-probes-identical-tcp.yaml", "Pod Probes", 7)
}

func TestPodProbesIdenticalExec(t *testing.T) {
	testExpectedScore(t, "pod-probes-identical-exec.yaml", "Pod Probes", 7)
}

func TestProbesTargetedByService(t *testing.T) {
	comments := testGetComments(t, "pod-probes-targeted-by-service.yaml", "Pod Probes")
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container is missing a readinessProbe", comments[0].Summary)

	testExpectedScore(t, "pod-probes-targeted-by-service.yaml", "Pod Probes", 0)
}

func TestProbesTargetedByServiceSameNamespace(t *testing.T) {
	comments := testGetComments(t, "pod-probes-targeted-by-service-same-namespace.yaml", "Pod Probes")
	assert.Len(t, comments, 1)
	assert.Equal(t, "Container is missing a readinessProbe", comments[0].Summary)

	testExpectedScore(t, "pod-probes-targeted-by-service-same-namespace.yaml", "Pod Probes", 0)
}

func TestProbesTargetedByServiceDifferentNamespace(t *testing.T) {
	comments := testGetComments(t, "pod-probes-targeted-by-service-different-namespace.yaml", "Pod Probes")
	assert.Len(t, comments, 0)
	testExpectedScore(t, "pod-probes-targeted-by-service-different-namespace.yaml", "Pod Probes", 10)
}

func TestProbesTargetedByServiceNotTargeted(t *testing.T) {
	testExpectedScore(t, "pod-probes-not-targeted-by-service.yaml", "Pod Probes", 10)
}

