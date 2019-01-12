package score

import (
	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/scorecard"
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

// testExpectedScoreWithConfig runs all tests, but makes sure that the test for "testcase" was executed, and that
// the grade is set to expectedScore. The function returns the comments of "testcase".
func testExpectedScoreWithConfig(t *testing.T, config config.Configuration, filename string, testcase string, expectedScore scorecard.Grade) []scorecard.TestScoreComment {
	sc, err := testScore(config, filename)
	assert.NoError(t, err)

	for _, objectScore := range sc.Objects {
		for _, s := range objectScore.Checks {
			if s.Name == testcase {
				assert.Equal(t, expectedScore, s.Grade)
				return s.Comments
			}
		}
	}

	t.Error("Was not tested")
	return nil
}

func testScore(config config.Configuration, filename string) (*scorecard.Scorecard, error) {
	config.AllFiles = []io.Reader{testFile(filename)}

	parsed, err := parser.ParseFiles(config)
	if err != nil {
		return nil, err
	}

	return Score(parsed, config)
}

func testExpectedScore(t *testing.T, filename string, testcase string, expectedScore scorecard.Grade) []scorecard.TestScoreComment {
	return testExpectedScoreWithConfig(t, config.Configuration{}, filename, testcase, expectedScore)
}

func TestPodContainerNoResources(t *testing.T) {
	testExpectedScore(t, "pod-test-resources-none.yaml", "Container Resources", 1)
}

func TestPodContainerResourceLimits(t *testing.T) {
	testExpectedScore(t, "pod-test-resources-only-limits.yaml", "Container Resources", 5)
}

func TestPodContainerResourceLimitsAndRequests(t *testing.T) {
	testExpectedScore(t, "pod-test-resources-limits-and-requests.yaml", "Container Resources", 10)
}

func TestPodContainerResourceLimitCpuNotRequired(t *testing.T) {
	testExpectedScoreWithConfig(t, config.Configuration{IgnoreContainerCpuLimitRequirement: true}, "pod-test-resources-limits-and-requests-no-cpu-limit.yaml", "Container Resources", 10)
}

func TestPodContainerResourceLimitCpuRequired(t *testing.T) {
	testExpectedScoreWithConfig(t, config.Configuration{IgnoreContainerCpuLimitRequirement: false}, "pod-test-resources-limits-and-requests-no-cpu-limit.yaml", "Container Resources", 1)
}

func TestDeploymentResources(t *testing.T) {
	testExpectedScore(t, "deployment-test-resources.yaml", "Container Resources", 5)
}

func TestStatefulSetResources(t *testing.T) {
	testExpectedScore(t, "statefulset-test-resources.yaml", "Container Resources", 5)
}

func TestPodContainerTagLatest(t *testing.T) {
	testExpectedScore(t, "pod-image-tag-latest.yaml", "Container Image Tag", 1)
}

func TestPodContainerTagFixed(t *testing.T) {
	testExpectedScore(t, "pod-image-tag-fixed.yaml", "Container Image Tag", 10)
}

func TestPodContainerPullPolicyUndefined(t *testing.T) {
	testExpectedScore(t, "pod-image-pullpolicy-undefined.yaml", "Container Image Pull Policy", 1)
}

func TestPodContainerPullPolicyUndefinedLatestTag(t *testing.T) {
	testExpectedScore(t, "pod-image-pullpolicy-undefined-latest-tag.yaml", "Container Image Pull Policy", 10)
}

func TestPodContainerPullPolicyUndefinedNoTag(t *testing.T) {
	testExpectedScore(t, "pod-image-pullpolicy-undefined-no-tag.yaml", "Container Image Pull Policy", 10)
}

func TestPodContainerPullPolicyNever(t *testing.T) {
	testExpectedScore(t, "pod-image-pullpolicy-never.yaml", "Container Image Pull Policy", 1)
}

func TestPodContainerPullPolicyAlways(t *testing.T) {
	testExpectedScore(t, "pod-image-pullpolicy-always.yaml", "Container Image Pull Policy", 10)
}

func TestContainerSecurityContextPrivilegied(t *testing.T) {
	testExpectedScore(t, "pod-security-context-privilegied.yaml", "Container Security Context", 1)
}

func TestContainerSecurityContextNonPrivilegied(t *testing.T) {
	testExpectedScore(t, "pod-security-context-non-privilegied.yaml", "Container Security Context", 10)
}

func TestContainerSecurityContextLowUser(t *testing.T) {
	testExpectedScore(t, "pod-security-context-low-user-id.yaml", "Container Security Context", 1)
}

func TestContainerSecurityContextLowGroup(t *testing.T) {
	testExpectedScore(t, "pod-security-context-low-group-id.yaml", "Container Security Context", 1)
}

func TestContainerSecurityContextHighIds(t *testing.T) {
	testExpectedScore(t, "pod-security-context-high-ids.yaml", "Container Security Context", 10)
}

func TestConfigMapMultiDash(t *testing.T) {
	_, err := testScore(config.Configuration{}, "configmap-multi-dash.yaml")
	assert.Nil(t, err)
}
