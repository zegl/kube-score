package score

import (
	"io"
	"os"
	"testing"

	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/scorecard"

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
func testExpectedScoreWithConfig(t *testing.T, config config.Configuration, testcase string, expectedScore scorecard.Grade) []scorecard.TestScoreComment {
	sc, err := testScore(config)
	assert.NoError(t, err)

	for _, objectScore := range sc {
		for _, s := range objectScore.Checks {
			if s.Check.Name == testcase {
				assert.Equal(t, expectedScore, s.Grade)
				return s.Comments
			}
		}
	}

	t.Error("Was not tested")
	return nil
}

func testScore(config config.Configuration) (scorecard.Scorecard, error) {
	parsed, err := parser.ParseFiles(config)
	if err != nil {
		return nil, err
	}

	card, err := Score(parsed, config)
	if err != nil {
		return nil, err
	}

	return *card, err
}

func testExpectedScore(t *testing.T, filename string, testcase string, expectedScore scorecard.Grade) []scorecard.TestScoreComment {
	return testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles: []io.Reader{testFile(filename)},
	}, testcase, expectedScore)
}

func testExpectedScoreReader(t *testing.T, content io.Reader, testcase string, expectedScore scorecard.Grade) []scorecard.TestScoreComment {
	return testExpectedScoreWithConfig(
		t, config.Configuration{
			AllFiles: []io.Reader{content},
		},
		testcase,
		expectedScore,
	)
}

func TestPodContainerNoResources(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-test-resources-none.yaml", "Container Resources", 1)
}

func TestPodContainerResourceLimits(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-test-resources-only-limits.yaml", "Container Resources", 5)
}

func TestPodContainerResourceLimitsAndRequests(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-test-resources-limits-and-requests.yaml", "Container Resources", 10)
}

func TestPodContainerResourceLimitCpuNotRequired(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, config.Configuration{
		IgnoreContainerCpuLimitRequirement: true,
		AllFiles:                           []io.Reader{testFile("pod-test-resources-limits-and-requests-no-cpu-limit.yaml")},
	}, "Container Resources", 10)
}

func TestPodContainerResourceLimitCpuRequired(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, config.Configuration{
		IgnoreContainerCpuLimitRequirement: false,
		AllFiles:                           []io.Reader{testFile("pod-test-resources-limits-and-requests-no-cpu-limit.yaml")},
	}, "Container Resources", 1)
}

func TestPodContainerResourceRequestsEqualLimitsCpuNotRequired(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-resource-requests-equal-limits"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		IgnoreContainerCpuLimitRequirement: true,
		AllFiles:                           []io.Reader{testFile("pod-test-resources-limits-and-requests-no-cpu-limit.yaml")},
		EnabledOptionalTests:               structMap,
	}, "Container Resource Requests Equal Limits", 10)
}

func TestPodContainerResourceNoLimitRequired(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, config.Configuration{
		IgnoreContainerCpuLimitRequirement:    true,
		IgnoreContainerMemoryLimitRequirement: true,
		AllFiles:                              []io.Reader{testFile("pod-test-resources-no-limits.yaml")},
	}, "Container Resources", 10)
}

func TestDeploymentResources(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-test-resources.yaml", "Container Resources", 5)
}

func TestStatefulSetResources(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-test-resources.yaml", "Container Resources", 5)
}

func TestPodContainerTagLatest(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-tag-latest.yaml", "Container Image Tag", 1)
}

func TestPodContainerTagFixed(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-tag-fixed.yaml", "Container Image Tag", 10)
}

func TestPodContainerPullPolicyUndefined(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-pullpolicy-undefined.yaml", "Container Image Pull Policy", 1)
}

func TestPodContainerPullPolicyUndefinedLatestTag(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-pullpolicy-undefined-latest-tag.yaml", "Container Image Pull Policy", 10)
}

func TestPodContainerPullPolicyUndefinedNoTag(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-pullpolicy-undefined-no-tag.yaml", "Container Image Pull Policy", 10)
}

func TestPodContainerPullPolicyNever(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-pullpolicy-never.yaml", "Container Image Pull Policy", 1)
}

func TestPodContainerPullPolicyAlways(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-pullpolicy-always.yaml", "Container Image Pull Policy", 10)
}

func TestConfigMapMultiDash(t *testing.T) {
	t.Parallel()
	_, err := testScore(config.Configuration{
		AllFiles: []io.Reader{testFile("configmap-multi-dash.yaml")},
	})
	assert.Nil(t, err)
}

func TestAnnotationIgnore(t *testing.T) {
	t.Parallel()
	s, err := testScore(config.Configuration{
		VerboseOutput:             0,
		AllFiles:                  []io.Reader{testFile("ignore-annotation-service.yaml")},
		UseIgnoreChecksAnnotation: true,
	})
	assert.Nil(t, err)
	assert.Len(t, s, 1)

	tested := false

	for _, o := range s {
		for _, c := range o.Checks {
			if c.Check.ID == "service-type" {
				assert.True(t, c.Skipped)
				tested = true
			}
		}
		assert.Equal(t, "node-port-service-with-ignore", o.ObjectMeta.Name)
	}
	assert.True(t, tested)
}

func TestAnnotationIgnoreDisabled(t *testing.T) {
	t.Parallel()
	s, err := testScore(config.Configuration{
		VerboseOutput:             0,
		AllFiles:                  []io.Reader{testFile("ignore-annotation-service.yaml")},
		UseIgnoreChecksAnnotation: false,
	})
	assert.Nil(t, err)
	assert.Len(t, s, 1)

	tested := false

	for _, o := range s {
		for _, c := range o.Checks {
			if c.Check.ID == "service-type" {
				assert.False(t, c.Skipped)
				assert.Equal(t, scorecard.GradeWarning, c.Grade)
				tested = true
			}
		}
		assert.Equal(t, "node-port-service-with-ignore", o.ObjectMeta.Name)
	}
	assert.True(t, tested)
}

func TestList(t *testing.T) {
	t.Parallel()
	s, err := testScore(config.Configuration{
		AllFiles: []io.Reader{testFile("list.yaml")},
	})
	assert.Nil(t, err)
	assert.Len(t, s, 2)

	hasService := false
	hasDeployment := false

	for _, obj := range s {
		if obj.ObjectMeta.Name == "list-service-test" {
			hasService = true
		}
		if obj.ObjectMeta.Name == "list-deployment-test" {
			hasDeployment = true
		}
		assert.Condition(t, func() bool { return len(obj.Checks) > 2 })
	}

	assert.True(t, hasService)
	assert.True(t, hasDeployment)
}
