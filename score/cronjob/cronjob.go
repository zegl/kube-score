package cronjob

import (
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks) {
	allChecks.RegisterCronJobCheck("CronJob has deadline", `Makes sure that all CronJobs has a configured deadline`, cronJobHasDeadline)
	allChecks.RegisterCronJobCheck("CronJob RestartPolicy", `Makes sure CronJobs have a valid RestartPolicy`, cronJobHasRestartPolicy)
}

func cronJobHasDeadline(job ks.CronJob) (score scorecard.TestScore) {
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
func cronJobHasRestartPolicy(job ks.CronJob) (score scorecard.TestScore) {

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
