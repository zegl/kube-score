package revisionHistoryLimit

import (
//	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
//	"github.com/zegl/kube-score/score/internal"
	"github.com/zegl/kube-score/scorecard"
        appsv1 "k8s.io/api/apps/v1"
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

/*
                        } else {
                                score.Grade = scorecard.GradeAlmostOK
                                score.AddComment("", "Pod has the same readiness and liveness probe", "It's recommended to have different probes for the two different purposes.")
                        }
*/

/*
                if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas < 2 {
                        score.Grade = scorecard.GradeAllOK
                        score.AddComment("", "Skipped", "Skipped because the deployment has less than 2 replicas")
                        return
                }

                match, matchErr := "hello", nil
                if matchErr != nil {
                        err = matchErr
                        return
                }

                if match {
                        score.Grade = scorecard.GradeAllOK
                } else {
                        score.Grade = scorecard.GradeCritical
                        score.AddComment("", "No matching PodDisruptionBudget was found", "It's recommended to define a PodDisruptionBudget to avoid unexpected downtime during Kubernetes maintenance operations, such as when draining a node.")
                }

                return
*/
        }
}
