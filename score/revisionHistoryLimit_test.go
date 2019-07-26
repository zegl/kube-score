package score

import "testing"

func TestRevisionHistoryLimit(t *testing.T) {
	testExpectedScore(t, "deployment-sets-revisionhistorylimit.yaml", "Deployment sets revisionHistoryLimit", 10)
}
func TestNoRevisionHistoryLimit(t *testing.T) {
        testExpectedScore(t, "deployment-sets-no-revisionhistorylimit.yaml", "Deployment sets revisionHistoryLimit", 1)
}
func TestHighRevisionHistoryLimit(t *testing.T) {
        testExpectedScore(t, "deployment-sets-high-revisionhistorylimit.yaml", "Deployment sets revisionHistoryLimit", 7)
}
