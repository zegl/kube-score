package score

import (
	"testing"
)

func TestCronJobHasDeadline(t *testing.T) {
	testExpectedScore(t, "cronjob-deadline-set.yaml", "CronJob has deadline", 10)
}

func TestCronJobNotHasDeadline(t *testing.T) {
	testExpectedScore(t, "cronjob-deadline-not-set.yaml", "CronJob has deadline", 1)
}
