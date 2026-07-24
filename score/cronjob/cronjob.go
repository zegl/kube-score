package cronjob

import (
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks) {
	allChecks.RegisterCronJobCheck("CronJob has deadline", `Makes sure that all CronJobs has a configured deadline`, cronJobHasDeadline)
	allChecks.RegisterCronJobCheck("CronJob RestartPolicy", `Makes sure CronJobs have a valid RestartPolicy`, cronJobHasRestartPolicy)
	// Optional: many clusters intentionally rely on the controller-manager time zone.
	// Enable when schedules must not silently follow the control-plane TZ.
	allChecks.RegisterOptionalCronJobCheck("CronJob has time zone", `Makes sure that CronJobs set an explicit timeZone for the schedule`, cronJobHasTimeZone)
}

func cronJobHasDeadline(job ks.CronJob) (score scorecard.TestScore, err error) {
	if job.StartingDeadlineSeconds() == nil {
		score.Grade = scorecard.GradeCritical
		score.AddComment("", "The CronJob should have startingDeadlineSeconds configured",
			"This makes sure that jobs are automatically cancelled if they can not be scheduled")
		return
	}

	score.Grade = scorecard.GradeAllOK
	return
}

// CronJob restartPolicy must be "OnFailure" or "Never". It cannot be empty (unspecified)
func cronJobHasRestartPolicy(job ks.CronJob) (score scorecard.TestScore, err error) {
	podTmpl := job.GetPodTemplateSpec()
	restartPolicy := podTmpl.Spec.RestartPolicy

	if len(restartPolicy) > 0 {
		if restartPolicy == "Never" || restartPolicy == "OnFailure" {
			score.Grade = scorecard.GradeAllOK
		} else {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "The CronJob must have a valid RestartPolicy configured",
				"Valid CronJob RestartPolicy settings are Never or OnFailure")
		}
	} else {
		score.Grade = scorecard.GradeCritical
		score.AddComment("", "The CronJob is missing a valid RestartPolicy",
			"Valid CronJob RestartPolicy settings are Never or OnFailure")
	}

	return
}

// cronJobHasTimeZone warns when the schedule time zone is left implicit.
// Kubernetes defaults missing timeZone to the kube-controller-manager process
// time zone, which can shift schedules after control-plane moves or upgrades.
// https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#time-zones
func cronJobHasTimeZone(job ks.CronJob) (score scorecard.TestScore, err error) {
	tz := job.TimeZone()
	if tz == nil || len(*tz) == 0 {
		score.Grade = scorecard.GradeWarning
		score.AddComment("", "The CronJob should set timeZone",
			"Without timeZone the schedule uses the kube-controller-manager process time zone, which can cause incorrect run times")
		return
	}

	score.Grade = scorecard.GradeAllOK
	return
}
