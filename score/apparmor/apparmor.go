package apparmor

import (
	"strings"

	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
        appsv1 "k8s.io/api/apps/v1"
)

func Register( allChecks *checks.Checks ) {
        allChecks.RegisterDeploymentCheck( "Deployment sets apparmor annotation", `Makes sure that all Deployments set apparmor annotation`, deploymentHas() )
}

func deploymentHas() func( appsv1.Deployment ) ( scorecard.TestScore, error ) {
        return func( deployment appsv1.Deployment ) ( score scorecard.TestScore, err error ) {
		if armor, found := deployment.Spec.Template.GetObjectMeta().GetAnnotations()[ "container.apparmor.security.beta.kubernetes.io" ]; found {
			if strings.Index( armor, "localhost/docker-default" ) != -1 {
				score.Grade = scorecard.GradeAlmostOK
				score.AddComment( "", "apparmor annotation is set:", "It is recommended to not use docker-default and instead customize a profile" )
				return
			}
			score.Grade = scorecard.GradeAllOK
			score.AddComment( "", "apparmor annotation is set:", armor )
			return
		}
		score.Grade = scorecard.GradeCritical
		score.AddComment( "", "apparmor annotation is not set", "It is recommended to set apparmor annotation and customize a profile" )
		return
	}
}
