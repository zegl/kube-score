package score

import (
	"testing"
)

func TestCronJobHasDeadline(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "cronjob-deadline-set.yaml", "CronJob has deadline", 10)
}

func TestCronJobNotHasDeadline(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "cronjob-deadline-not-set.yaml", "CronJob has deadline", 1)
}

func TestProbesPodCronMissingReady(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "cronjob-deadline-not-set.yaml", "Pod Probes", 10)
}
