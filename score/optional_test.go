package score

import (
	"testing"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func TestOptionalSkippedByDefault(t *testing.T) {
	t.Parallel()
	enabledOptionalTests := make(map[string]struct{})
	wasSkipped(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-container-memory-requests.yaml")},
		EnabledOptionalTests: enabledOptionalTests,
	}, "Container Memory Requests Equal Limits")
}

func TestOptionalIgnoredAndEnabled(t *testing.T) {
	t.Parallel()

	enabledOptionalTests := make(map[string]struct{})
	enabledOptionalTests["container-resource-requests-equal-limits"] = struct{}{}

	ignoredTests := make(map[string]struct{})
	ignoredTests["container-resource-requests-equal-limits"] = struct{}{}

	wasSkipped(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-container-memory-requests.yaml")},
		EnabledOptionalTests: enabledOptionalTests,
		IgnoredTests:         ignoredTests,
	}, "Container Memory Requests Equal Limits")
}

func TestOptionalRunCliFlagEnabledDefault(t *testing.T) {
	t.Parallel()

	enabledOptionalTests := make(map[string]struct{})
	enabledOptionalTests["container-resource-requests-equal-limits"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-container-memory-requests.yaml")},
		EnabledOptionalTests: enabledOptionalTests,
	}, "Container Memory Requests Equal Limits", scorecard.GradeCritical)
}

func TestOptionalRunAnnotationEnabled(t *testing.T) {
	t.Parallel()

	enabledOptionalTests := make(map[string]struct{})

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-container-memory-requests-annotation-optional.yaml")},
		EnabledOptionalTests: enabledOptionalTests,
	}, "Container Memory Requests Equal Limits", scorecard.GradeCritical)
}
