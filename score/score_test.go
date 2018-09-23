package score

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testFile(name string) *os.File {
	fp, err := os.Open("testdata/" + name)
	if err != nil {
		panic(err)
	}
	return fp
}

func testExpectedScore(t *testing.T, filename string, testcase string, expectedScore int) {
	sc, err := Score([]io.Reader{testFile(filename)})
	assert.NoError(t, err)
	tested := false
	for _, objectScore := range sc.Scores {
		for _, s := range objectScore {
			if s.Name == testcase {
				assert.Equal(t, expectedScore, s.Grade)
				tested = true
			}
		}
	}
	assert.True(t, tested, "Was not tested")
}

func TestPodContainerNoResources(t *testing.T) {
	testExpectedScore(t, "pod-test-resources-none.yaml", "Container Resources", 0)
}

func TestPodContainerResourceLimits(t *testing.T) {
	testExpectedScore(t, "pod-test-resources-only-limits.yaml", "Container Resources", 5)
}

func TestPodContainerResourceLimitsAndRequests(t *testing.T) {
	testExpectedScore(t, "pod-test-resources-limits-and-requests.yaml", "Container Resources", 10)
}

func TestDeploymentResources(t *testing.T) {
	testExpectedScore(t, "deployment-test-resources.yaml", "Container Resources", 5)
}

func TestStatefulSetResources(t *testing.T) {
	testExpectedScore(t, "statefulset-test-resources.yaml", "Container Resources", 5)
}

func TestStatefulSetAppsv1beta1(t *testing.T) {
	testExpectedScore(t, "statefulset-appsv1beta1.yaml", "Unstable version", 5)
}

func TestStatefulSetAppsv1beta2(t *testing.T) {
	testExpectedScore(t, "statefulset-appsv1beta2.yaml", "Unstable version", 5)
}

func TestPodContainerTagLatest(t *testing.T) {
	testExpectedScore(t, "pod-image-tag-latest.yaml", "Container Image Tag", 0)
}

func TestPodContainerTagFixed(t *testing.T) {
	testExpectedScore(t, "pod-image-tag-fixed.yaml", "Container Image Tag", 10)
}

func TestPodContainerPullPolicyUndefined(t *testing.T) {
	testExpectedScore(t, "pod-image-pullpolicy-undefined.yaml", "Container Image Pull Policy", 0)
}

func TestPodContainerPullPolicyNever(t *testing.T) {
	testExpectedScore(t, "pod-image-pullpolicy-never.yaml", "Container Image Pull Policy", 0)
}

func TestPodContainerPullPolicyAlways(t *testing.T) {
	testExpectedScore(t, "pod-image-pullpolicy-always.yaml", "Container Image Pull Policy", 10)
}

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

func TestPodProbesAllMissing(t *testing.T) {
	testExpectedScore(t, "pod-probes-all-missing.yaml", "Pod Probes", 0)
}

func TestPodProbesMissingReady(t *testing.T) {
	testExpectedScore(t, "pod-probes-missing-ready.yaml", "Pod Probes", 5)
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

func TestDeploymentExtensionsv1beta1(t *testing.T) {
	testExpectedScore(t, "deployment-extensions-v1beta1.yaml", "Unstable version", 5)
}

func TestDeploymentAppsv1beta1(t *testing.T) {
	testExpectedScore(t, "deployment-appsv1beta1.yaml", "Unstable version", 5)
}

func TestDeploymentAppsv1beta2(t *testing.T) {
	testExpectedScore(t, "deployment-appsv1beta2.yaml", "Unstable version", 5)
}

func TestDaemonSetAppsv1(t *testing.T) {
	testExpectedScore(t, "daemonset-appsv1.yaml", "Unstable version", 10)
}

func TestDaemonSetAppsv1beta2(t *testing.T) {
	testExpectedScore(t, "daemonset-appsv1beta2.yaml", "Unstable version", 5)
}

func TestDaemonSetExtensionsv1beta1(t *testing.T) {
	testExpectedScore(t, "daemonset-extensionsv1beta1.yaml", "Unstable version", 5)
}

func TestCronJobBatchv1beta1(t *testing.T) {
	testExpectedScore(t, "cronjob-batchv1beta1.yaml", "Unstable version", 10)
}

func TestJobBatchv1(t *testing.T) {
	testExpectedScore(t, "job-batchv1.yaml", "Unstable version", 10)
}