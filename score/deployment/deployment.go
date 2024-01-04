package deployment

import (
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/score/internal"
	"github.com/zegl/kube-score/scorecard"
	v1 "k8s.io/api/apps/v1"
)

func Register(allChecks *checks.Checks, all ks.AllTypes) {
	allChecks.RegisterDeploymentCheck("Deployment Strategy", `Makes sure that all Deploymtes targeted by service use RollingUpdate strategy`, deploymentRolloutStrategy(all.Services()))
}

// deploymentRolloutStrategy checks if a Deployment has the update strategy on RollingUpdate if targeted by a service
func deploymentRolloutStrategy(svcs []ks.Service) func(deployment v1.Deployment) (scorecard.TestScore, error) {
	svcsInNamespace := make(map[string][]map[string]string)
	for _, s := range svcs {
		svc := s.Service()
		if _, ok := svcsInNamespace[svc.Namespace]; !ok {
			svcsInNamespace[svc.Namespace] = []map[string]string{}
		}
		svcsInNamespace[svc.Namespace] = append(svcsInNamespace[svc.Namespace], svc.Spec.Selector)
	}
	return func(deployment v1.Deployment) (score scorecard.TestScore, err error) {
		referencedByService := false

		for _, svcSelector := range svcsInNamespace[deployment.Namespace] {
			if internal.LabelSelectorMatchesLabels(svcSelector, deployment.Spec.Template.Labels) {
				referencedByService = true
				break
			}
		}

		if referencedByService {
			if deployment.Spec.Strategy.Type == v1.RollingUpdateDeploymentStrategyType {
				score.Grade = scorecard.GradeAllOK
			} else {
				score.Grade = scorecard.GradeWarning
				score.AddCommentWithURL("", "Deployment update strategy", "The deployment is used by a service but not using rolling update strategy which can cause interruptions", "https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy")
			}
		} else {
			score.Skipped = true
			score.AddComment("", "Skipped as the deployment strategy does not matter if not targeted by a service", "")
		}

		return
	}
}
