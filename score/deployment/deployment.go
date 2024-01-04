package deployment

import (
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/score/internal"
	"github.com/zegl/kube-score/scorecard"
	v1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/utils/ptr"
)

func Register(allChecks *checks.Checks, all ks.AllTypes) {
	allChecks.RegisterDeploymentCheck("Deployment Strategy", `Makes sure that all Deployments targeted by service use RollingUpdate strategy`, deploymentRolloutStrategy(all.Services()))
	allChecks.RegisterDeploymentCheck("Deployment Replicas", `Makes sure that Deployment has multiple replicas`, deploymentReplicas(all.Services(), all.HorizontalPodAutoscalers()))
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
			if deployment.Spec.Strategy.Type == v1.RollingUpdateDeploymentStrategyType || deployment.Spec.Strategy.Type == "" {
				score.Grade = scorecard.GradeAllOK
			} else {
				score.Grade = scorecard.GradeWarning
				score.AddCommentWithURL("", "Deployment update strategy", "The deployment is used by a service but not using the RollingUpdate strategy which can cause interruptions. Set .spec.strategy.type to RollingUpdate.", "https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy")
			}
		} else {
			score.Skipped = true
			score.AddComment("", "Skipped as the Deployment is not targeted by a service", "")
		}

		return
	}
}

// deploymentReplicas checks if a Deployment has >= 2 replicas if not (targeted by service || has HorizontalPodAutoscaler)
func deploymentReplicas(svcs []ks.Service, hpas []ks.HpaTargeter) func(deployment v1.Deployment) (scorecard.TestScore, error) {
	svcsInNamespace := make(map[string][]map[string]string)
	for _, s := range svcs {
		svc := s.Service()
		if _, ok := svcsInNamespace[svc.Namespace]; !ok {
			svcsInNamespace[svc.Namespace] = []map[string]string{}
		}
		svcsInNamespace[svc.Namespace] = append(svcsInNamespace[svc.Namespace], svc.Spec.Selector)
	}
	hpasInNamespace := make(map[string][]autoscalingv1.CrossVersionObjectReference)
	for _, hpa := range hpas {
		hpaTarget := hpa.HpaTarget()
		hpaMeta := hpa.GetObjectMeta()
		if _, ok := hpasInNamespace[hpaMeta.Namespace]; !ok {
			hpasInNamespace[hpaMeta.Namespace] = []autoscalingv1.CrossVersionObjectReference{}
		}
		hpasInNamespace[hpaMeta.Namespace] = append(hpasInNamespace[hpaMeta.Namespace], hpaTarget)
	}

	return func(deployment v1.Deployment) (score scorecard.TestScore, err error) {
		referencedByService := false
		hasHPA := false

		for _, svcSelector := range svcsInNamespace[deployment.Namespace] {
			if internal.LabelSelectorMatchesLabels(svcSelector, deployment.Spec.Template.Labels) {
				referencedByService = true
				break
			}
		}

		for _, hpaTarget := range hpasInNamespace[deployment.Namespace] {
			if deployment.TypeMeta.APIVersion == hpaTarget.APIVersion &&
				deployment.TypeMeta.Kind == hpaTarget.Kind &&
				deployment.ObjectMeta.Name == hpaTarget.Name {
				hasHPA = true
				break
			}
		}

		if !referencedByService || hasHPA {
			score.Skipped = true
			score.AddComment("", "Skipped as the Deployment is not targeted by service or is controlled by a HorizontalPodAutoscaler", "")
		} else {
			if ptr.Deref(deployment.Spec.Replicas, 1) >= 2 {
				score.Grade = scorecard.GradeAllOK
			} else {
				score.Grade = scorecard.GradeWarning
				score.AddComment("", "Deployment few replicas", "Deployments targeted by Services are recommended to have at least 2 replicas to prevent unwanted downtime.")
			}
		}

		return
	}
}
