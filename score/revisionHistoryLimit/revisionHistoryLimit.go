package revisionHistoryLimit

import (
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
        appsv1 "k8s.io/api/apps/v1"
)

func Register( allChecks *checks.Checks ) {
        allChecks.RegisterDeploymentCheck( "Deployment sets revisionHistoryLimit", `Makes sure that all Deployments set a limit`, deploymentHas() )
}

func deploymentHas() func( appsv1.Deployment ) ( scorecard.TestScore, error ) {
        return func( deployment appsv1.Deployment ) ( score scorecard.TestScore, err error ) {
		limit := deployment.Spec.RevisionHistoryLimit
                if limit != nil {
			if *limit > 5 {
                                score.Grade = scorecard.GradeAlmostOK
				score.AddComment( "", "revisionHistoryLimit greater than 5", "It's recommended to define a revisionHistoryLimit below 5 to avoid data storage impact" )
			} else {
				score.Grade = scorecard.GradeAllOK
			}
			return
                } else {
                        score.Grade = scorecard.GradeCritical
                        score.AddComment( "", "No revisionHistoryLimit Configuration was found", "It's recommended to define a revisionHistoryLimit below 5 to avoid data storage impact" )
                        return
                }
		return
        }
}
