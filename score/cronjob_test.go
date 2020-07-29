package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestCronJobHasDeadline(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "cronjob-deadline-set.yaml", "CronJob has deadline", scorecard.GradeAllOK)
}

func TestCronJobNotHasDeadline(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "cronjob-deadline-not-set.yaml", "CronJob has deadline", scorecard.GradeCritical)
}

func TestProbesPodCronMissingReady(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "cronjob-deadline-not-set.yaml", "Pod Probes", scorecard.GradeAllOK)
}
