package score

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/scorecard"
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

func wasSkipped(t *testing.T, config config.Configuration, testcase string) bool {
	sc, err := testScore(config)
	assert.NoError(t, err)
	for _, objectScore := range sc {
		for _, s := range objectScore.Checks {
			if s.Check.Name == testcase {
				return s.Skipped
			}
		}
	}

	assert.Fail(t, "test was not run")
	return false
}

func testScore(config config.Configuration) (scorecard.Scorecard, error) {
	p, err := parser.New()
	if err != nil {
		return nil, err
	}

	parsed, err := p.ParseFiles(config)
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
		AllFiles:          []ks.NamedReader{testFile(filename)},
		KubernetesVersion: config.Semver{Major: 1, Minor: 18},
	}, testcase, expectedScore)
}

func TestPodContainerNoResources(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-test-resources-none.yaml", "Container Resources", scorecard.GradeCritical)
}

func TestPodContainerResourceLimits(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-test-resources-only-limits.yaml", "Container Resources", scorecard.GradeWarning)
}

func TestPodContainerResourceLimitsAndRequests(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-test-resources-limits-and-requests.yaml", "Container Resources", scorecard.GradeAllOK)
}

func TestPodContainerResourceLimitCpuNotRequired(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, config.Configuration{
		IgnoreContainerCpuLimitRequirement: true,
		AllFiles:                           []ks.NamedReader{testFile("pod-test-resources-limits-and-requests-no-cpu-limit.yaml")},
	}, "Container Resources", scorecard.GradeAllOK)
}

func TestPodContainerResourceLimitCpuRequired(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, config.Configuration{
		IgnoreContainerCpuLimitRequirement: false,
		AllFiles:                           []ks.NamedReader{testFile("pod-test-resources-limits-and-requests-no-cpu-limit.yaml")},
	}, "Container Resources", scorecard.GradeCritical)
}

func TestPodContainerResourceNoLimitRequired(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, config.Configuration{
		IgnoreContainerCpuLimitRequirement:    true,
		IgnoreContainerMemoryLimitRequirement: true,
		AllFiles:                              []ks.NamedReader{testFile("pod-test-resources-no-limits.yaml")},
	}, "Container Resources", scorecard.GradeAllOK)
}

func TestPodContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-resource-requests-equal-limits"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-test-resources-limits-and-requests.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Resource Requests Equal Limits", scorecard.GradeAllOK)
}

func TestPodContainerMemoryRequestsEqualLimits(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-memory-requests-equal-limits"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-test-resources-limits-and-requests.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Memory Requests Equal Limits", scorecard.GradeAllOK)
}

func TestPodContainerCPURequestsEqualLimits(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-cpu-requests-equal-limits"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-test-resources-limits-and-requests.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container CPU Requests Equal Limits", scorecard.GradeAllOK)
}

func TestPodContainerResourceRequestsEqualLimitsNoLimits(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-resource-requests-equal-limits"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-test-resources-no-limits.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Resource Requests Equal Limits", scorecard.GradeCritical)
}

func TestPodContainerResourceRequestsEqualLimitsNoLimitsAnnotation(t *testing.T) {
	t.Parallel()

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:                    []ks.NamedReader{testFile("pod-test-resources-no-limits-annotation-optional.yaml")},
		UseOptionalChecksAnnotation: true,
	}, "Container Resource Requests Equal Limits", scorecard.GradeCritical)
}

func TestPodContainerMemoryRequestsEqualLimitsNoLimits(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-memory-requests-equal-limits"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-test-resources-no-limits.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Memory Requests Equal Limits", scorecard.GradeCritical)
}

func TestPodContainerMemoryRequestsEqualLimitsNoLimitsAnnotation(t *testing.T) {
	t.Parallel()

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:                    []ks.NamedReader{testFile("pod-test-resources-no-limits-annotation-optional.yaml")},
		UseOptionalChecksAnnotation: true,
	}, "Container Memory Requests Equal Limits", scorecard.GradeCritical)
}

func TestPodContainerCPURequestsEqualLimitsNoLimits(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-cpu-requests-equal-limits"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-test-resources-no-limits.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container CPU Requests Equal Limits", scorecard.GradeCritical)
}

func TestPodContainerCPURequestsEqualLimitsNoLimitsAnnotation(t *testing.T) {
	t.Parallel()

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:                    []ks.NamedReader{testFile("pod-test-resources-no-limits-annotation-optional.yaml")},
		UseOptionalChecksAnnotation: true,
	}, "Container CPU Requests Equal Limits", scorecard.GradeCritical)
}

func TestDeploymentResources(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-test-resources.yaml", "Container Resources", scorecard.GradeWarning)
}

func TestStatefulSetResources(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-test-resources.yaml", "Container Resources", scorecard.GradeWarning)
}

func TestPodContainerTagLatest(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-tag-latest.yaml", "Container Image Tag", scorecard.GradeCritical)
}

func TestPodContainerTagFixed(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-tag-fixed.yaml", "Container Image Tag", scorecard.GradeAllOK)
}

func TestPodContainerPullPolicyUndefined(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-pullpolicy-undefined.yaml", "Container Image Pull Policy", scorecard.GradeCritical)
}

func TestPodContainerPullPolicyUndefinedLatestTag(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-pullpolicy-undefined-latest-tag.yaml", "Container Image Pull Policy", scorecard.GradeAllOK)
}

func TestPodContainerPullPolicyUndefinedNoTag(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-pullpolicy-undefined-no-tag.yaml", "Container Image Pull Policy", scorecard.GradeAllOK)
}

func TestPodContainerPullPolicyNever(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-pullpolicy-never.yaml", "Container Image Pull Policy", scorecard.GradeCritical)
}

func TestPodContainerPullPolicyAlways(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-image-pullpolicy-always.yaml", "Container Image Pull Policy", scorecard.GradeAllOK)
}

func TestConfigMapMultiDash(t *testing.T) {
	t.Parallel()
	_, err := testScore(config.Configuration{
		AllFiles: []ks.NamedReader{testFile("configmap-multi-dash.yaml")},
	})
	assert.Nil(t, err)
}

func TestAnnotationIgnore(t *testing.T) {
	t.Parallel()
	s, err := testScore(config.Configuration{
		VerboseOutput:             0,
		AllFiles:                  []ks.NamedReader{testFile("ignore-annotation-service.yaml")},
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
		AllFiles:                  []ks.NamedReader{testFile("ignore-annotation-service.yaml")},
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
		AllFiles: []ks.NamedReader{testFile("list.yaml")},
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

// Note the input file specifies a condition that would fail the optional matching request and limit test, but returns GradeAllOK
// when only the default case is evaluated
func TestPodContainerStorageEphemeralRequestAndLimitOK(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-ephemeral-storage-request-nomatch-limit.yaml", "Container Ephemeral Storage Request and Limit", scorecard.GradeAllOK)
}

func TestPodContainerStorageEphemeralNoLimit(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-ephemeral-storage-missing-limit.yaml", "Container Ephemeral Storage Request and Limit", scorecard.GradeCritical)
}

func TestPodContainerStorageEphemeralNoRequest(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-ephemeral-storage-missing-request.yaml", "Container Ephemeral Storage Request and Limit", scorecard.GradeWarning)
}

func TestPodContainerStorageEphemeralRequestEqualsLimit(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-ephemeral-storage-request-equals-limit"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-ephemeral-storage-request-matches-limit.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Ephemeral Storage Request Equals Limit", scorecard.GradeAllOK)
}

func TestPodContainerStorageEphemeralRequestNoMatchLimit(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-ephemeral-storage-request-equals-limit"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-ephemeral-storage-request-nomatch-limit.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Ephemeral Storage Request Equals Limit", scorecard.GradeCritical)
}

func TestPodContainerStorageEphemeralRequestNoMatchLimitAnnotation(t *testing.T) {
	t.Parallel()

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:                    []ks.NamedReader{testFile("pod-ephemeral-storage-request-nomatch-limit-annotation-optional.yaml")},
		UseOptionalChecksAnnotation: true,
	}, "Container Ephemeral Storage Request Equals Limit", scorecard.GradeCritical)
}

func TestPodContainerStorageEphemeralIgnoreAnnotation(t *testing.T) {

	t.Parallel()
	s, err := testScore(config.Configuration{
		VerboseOutput:             0,
		AllFiles:                  []ks.NamedReader{testFile("pod-ephemeral-storage-annotation-ignore.yaml")},
		UseIgnoreChecksAnnotation: true,
	})
	assert.Nil(t, err)
	assert.Len(t, s, 1)

	tested := false

	for _, o := range s {
		for _, c := range o.Checks {
			if c.Check.ID == "container-resources" {
				assert.True(t, c.Skipped)
				tested = true
			}
		}
		assert.Equal(t, "pod-ephemeral-storage-annotation-ignore", o.ObjectMeta.Name)
	}
	assert.True(t, tested)
}

func TestPodContainerPortsContainerPortMissing(t *testing.T) {
	t.Parallel()
	structMap := make(map[string]struct{})
	structMap["container-ports-check"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-container-ports-missing-containerport.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Ports Check", scorecard.GradeCritical)
}

func TestPodContainerPortsContainerPortMissingAnnotation(t *testing.T) {
	t.Parallel()

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:                    []ks.NamedReader{testFile("pod-container-ports-missing-containerport-annotation-optional.yaml")},
		UseOptionalChecksAnnotation: true,
	}, "Container Ports Check", scorecard.GradeCritical)
}

func TestPodContainerPortsDuplicateNames(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-ports-check"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-container-ports-duplicate-names.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Ports Check", scorecard.GradeCritical)
}

func TestPodContainerPortsNameLength(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-ports-check"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-container-ports-name-too-long.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Ports Check", scorecard.GradeCritical)
}

func TestPodContainerPortsOK(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-ports-check"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-container-ports-ok.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Ports Check", scorecard.GradeAllOK)
}

func TestPodEnvOK(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["environment-variable-key-duplication"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-env-ok.yaml")},
		EnabledOptionalTests: structMap,
	}, "Environment Variable Key Duplication", scorecard.GradeAllOK)
}

func TestPodEnvDuplicated(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["environment-variable-key-duplication"] = struct{}{}

	actual := testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-env-duplicated.yaml")},
		EnabledOptionalTests: structMap,
	}, "Environment Variable Key Duplication", scorecard.GradeCritical)

	expected := []scorecard.TestScoreComment{
		{
			Path:        "foobar",
			Summary:     "Environment Variable Key Duplication",
			Description: "Container environment variable key 'bar' is duplicated",
		},
		{
			Path:        "foobar",
			Summary:     "Environment Variable Key Duplication",
			Description: "Container environment variable key 'baz' is duplicated",
		},
	}
	diff := cmp.Diff(expected, actual)
	assert.Empty(t, diff)
}

func TestMultipleIgnoreAnnotations(t *testing.T) {
	t.Parallel()
	s, err := testScore(config.Configuration{
		VerboseOutput:             0,
		AllFiles:                  []ks.NamedReader{testFile("kube-score-ignore-annotations.yaml")},
		UseIgnoreChecksAnnotation: true,
	})
	assert.Nil(t, err)
	assert.Len(t, s, 1)

	tested := false
	skipped := false

	for _, o := range s {
		for _, c := range o.Checks {
			// implied by the ignore container-resources annotation
			if c.Check.ID == "container-ephemeral-storage-request-and-limit" {
				assert.True(t, c.Skipped)
				skipped = true
			}
			// a default check
			if c.Check.ID == "container-image-pull-policy" {
				assert.False(t, c.Skipped)
				assert.Equal(t, scorecard.GradeCritical, c.Grade)
				tested = true
			}
			assert.Equal(t, "kube-score-ignore-annotations", o.ObjectMeta.Name)
		}
	}
	assert.True(t, tested)
	assert.True(t, skipped)
}
