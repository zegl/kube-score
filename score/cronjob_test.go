package score

import (
	"testing"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
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

func TestCronJobHasRestartPolicyMissing(t *testing.T) {
	t.Parallel()

	for _, v := range []string{"batchv1beta1", "batchv1"} {
		t.Run(v, func(t *testing.T) {
			testExpectedScore(t, "cronjob-"+v+"-restartpolicy-not-set.yaml", "CronJob RestartPolicy", scorecard.GradeCritical)
		})
	}
}

func TestCronJobHasRestartPolicyInvalid(t *testing.T) {
	t.Parallel()

	for _, v := range []string{"batchv1beta1", "batchv1"} {
		t.Run(v, func(t *testing.T) {
			testExpectedScore(t, "cronjob-"+v+"-restartpolicy-invalid.yaml", "CronJob RestartPolicy", scorecard.GradeCritical)
		})
	}
}

func TestCronJobHasRestartPolicyValid(t *testing.T) {
	t.Parallel()

	for _, v := range []string{"batchv1beta1", "batchv1"} {
		t.Run(v, func(t *testing.T) {
			testExpectedScore(t, "cronjob-"+v+"-restartpolicy-valid.yaml", "CronJob RestartPolicy", scorecard.GradeAllOK)
		})
	}
}

func TestCronJobHasTimeZone(t *testing.T) {
	t.Parallel()

	structMap := map[string]struct{}{
		"cronjob-has-time-zone": {},
	}

	for _, v := range []string{"batchv1beta1", "batchv1"} {
		t.Run(v, func(t *testing.T) {
			testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("cronjob-"+v+"-timezone-set.yaml")}, nil, &config.RunConfiguration{
				EnabledOptionalTests: structMap,
			}, "CronJob has time zone", scorecard.GradeAllOK)
		})
	}
}

func TestCronJobNotHasTimeZone(t *testing.T) {
	t.Parallel()

	structMap := map[string]struct{}{
		"cronjob-has-time-zone": {},
	}

	for _, v := range []string{"batchv1beta1", "batchv1"} {
		t.Run(v, func(t *testing.T) {
			testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("cronjob-"+v+"-timezone-not-set.yaml")}, nil, &config.RunConfiguration{
				EnabledOptionalTests: structMap,
			}, "CronJob has time zone", scorecard.GradeWarning)
		})
	}
}
