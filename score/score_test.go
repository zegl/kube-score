package score

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func testFile(name string) *os.File {
	fp, err := os.Open("testdata/" + name)
	if err != nil {
		panic(err)
	}
	return fp
}

func testExpectedScore(t *testing.T, filename string, testcase string, expectedScore int) {
	sc, err := Score(testFile(filename))
	assert.NoError(t, err)
	tested := false
	for _, s := range sc.Scores {
		if s.Name == testcase {
			assert.Equal(t, expectedScore, s.Grade)
			tested = true
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
