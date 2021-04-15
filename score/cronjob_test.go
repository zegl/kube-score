package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestCronJobHasDeadline(t *testing.T) {
	t.Parallel()

	for _, v := range []string{"batchv1beta1", "batchv1"} {
		t.Run(v, func(t *testing.T) {
			testExpectedScore(t, "cronjob-"+v+"-deadline-set.yaml", "CronJob has deadline", scorecard.GradeAllOK)
		})
	}
}

func TestCronJobNotHasDeadline(t *testing.T) {
	t.Parallel()

	for _, v := range []string{"batchv1beta1", "batchv1"} {
		t.Run(v, func(t *testing.T) {
			testExpectedScore(t, "cronjob-"+v+"-deadline-not-set.yaml", "CronJob has deadline", scorecard.GradeCritical)
		})
	}
}

func TestProbesPodCronMissingReady(t *testing.T) {
	t.Parallel()

	for _, v := range []string{"batchv1beta1", "batchv1"} {
		t.Run(v, func(t *testing.T) {
			testExpectedScore(t, "cronjob-"+v+"-deadline-not-set.yaml", "Pod Probes", scorecard.GradeAllOK)
		})
	}
}
