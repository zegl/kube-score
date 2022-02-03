package score

import (
	"testing"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

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
