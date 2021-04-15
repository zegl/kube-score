package cronjob

import (
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks) {
	allChecks.RegisterCronJobCheck("CronJob has deadline", `Makes sure that all CronJobs has a configured deadline`, cronJobHasDeadline)
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
