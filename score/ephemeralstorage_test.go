package score

import (
	"testing"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func TestPodContainerStorageEphemeralNoLimit(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-ephemeral-storage-missing-limit.yaml", "Container Resources", scorecard.GradeCritical)
}

func TestPodContainerStorageEphemeralNoRequest(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-ephemeral-storage-missing-request.yaml", "Container Resources", scorecard.GradeWarning)
}
func TestPodContainerStorageEphemeralRequestEqualsLimit(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-ephemeral-storage-request-matches-limit.yaml", "Container Resources", scorecard.GradeAllOK)
}

func TestPodContainerStorageEphemeralRequestNoMatchLimit(t *testing.T) {
	t.Parallel()

	structMap := make(map[string]struct{})
	structMap["container-ephemeral-storage-request-matches-limit"] = struct{}{}

	testExpectedScoreWithConfig(t, config.Configuration{
		AllFiles:             []ks.NamedReader{testFile("pod-ephemeral-storage-request-nomatch-limit.yaml")},
		EnabledOptionalTests: structMap,
	}, "Container Ephemeral Storage Requests and Limits", scorecard.GradeCritical)
}
