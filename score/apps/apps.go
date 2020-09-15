package apps

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"

	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/score/internal"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks, allHPAs []ks.HpaTargeter) {
	allChecks.RegisterDeploymentCheck("Deployment has host PodAntiAffinity", "Makes sure that a podAntiAffinity has been set that prevents multiple pods from being scheduled on the same node. https://kubernetes.io/docs/concepts/configuration/assign-pod-node/", deploymentHasAntiAffinity)
	allChecks.RegisterStatefulSetCheck("StatefulSet has host PodAntiAffinity", "Makes sure that a podAntiAffinity has been set that prevents multiple pods from being scheduled on the same node. https://kubernetes.io/docs/concepts/configuration/assign-pod-node/", statefulsetHasAntiAffinity)
	allChecks.RegisterDeploymentCheck("Deployment targeted by HPA does not have replicas configured", "Makes sure that Deployments using a HorizontalPodAutoscaler doesn't have a statically configured replica count set", hpaDeploymentNoReplicas(allHPAs))
}

func hpaDeploymentNoReplicas(allHPAs []ks.HpaTargeter) func(deployment appsv1.Deployment) (scorecard.TestScore, error) {
	return func(deployment appsv1.Deployment) (score scorecard.TestScore, err error) {
		// If is targeted by a HPA
		for _, hpa := range allHPAs {
			target := hpa.HpaTarget()

			if hpa.GetObjectMeta().Namespace == deployment.Namespace &&
				strings.ToLower(target.Kind) == strings.ToLower(deployment.Kind) &&
				target.Name == deployment.Name {

				if deployment.Spec.Replicas == nil {
					score.Grade = scorecard.GradeAllOK
					return
				}

				score.Grade = scorecard.GradeCritical
				score.AddComment("", "The deployment is targeted by a HPA, but a static replica count is configured in the DeploymentSpec", "When replicas is both statically set and managed by the HPA, the replicas will be changed to the statically configured count when the spec is applied, even if the HPA wants the replica count to be higher.")
				return
			}
		}

		score.Grade = scorecard.GradeAllOK
		score.Skipped = true
		score.AddComment("", "Skipped because the deployment is not targeted by a HorizontalPodAutoscaler", "")
		return
	}
}

func deploymentHasAntiAffinity(deployment appsv1.Deployment) (score scorecard.TestScore, err error) {
	// Ignore if the deployment only has a single replica
	// If replicas is not explicitly set, we'll still warn if the anti affinity is missing
	// as that might indicate use of a Horizontal Pod Autoscaler
	if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas < 2 {
		score.Skipped = true
		score.AddComment("", "Skipped because the deployment has less than 2 replicas", "")
		return
	}

	warn := func() {
		score.Grade = scorecard.GradeWarning
		score.AddComment("", "Deployment does not have a host podAntiAffinity set", "It's recommended to set a podAntiAffinity that stops multiple pods from a deployment from being scheduled on the same node. This increases availability in case the node becomes unavailable.")
	}

	affinity := deployment.Spec.Template.Spec.Affinity
	if affinity == nil || affinity.PodAntiAffinity == nil {
		warn()
		return
	}

	lables := internal.MapLables(deployment.Spec.Template.GetObjectMeta().GetLabels())

	if hasPodAntiAffinity(lables, affinity) {
		score.Grade = scorecard.GradeAllOK
		return
	}

	warn()
	return
}

func statefulsetHasAntiAffinity(statefulset appsv1.StatefulSet) (score scorecard.TestScore, err error) {
	// Ignore if the statefulset only has a single replica
	// If replicas is not explicitly set, we'll still warn if the anti affinity is missing
	// as that might indicate use of a Horizontal Pod Autoscaler
	if statefulset.Spec.Replicas != nil && *statefulset.Spec.Replicas < 2 {
		score.Skipped = true
		score.AddComment("", "Skipped because the statefulset has less than 2 replicas", "")
		return
	}

	warn := func() {
		score.Grade = scorecard.GradeWarning
		score.AddComment("", "StatefulSet does not have a host podAntiAffinity set", "It's recommended to set a podAntiAffinity that stops multiple pods from a statefulset from being scheduled on the same node. This increases availability in case the node becomes unavailable.")
	}

	affinity := statefulset.Spec.Template.Spec.Affinity
	if affinity == nil || affinity.PodAntiAffinity == nil {
		warn()
		return
	}

	lables := internal.MapLables(statefulset.Spec.Template.GetObjectMeta().GetLabels())

	if hasPodAntiAffinity(lables, affinity) {
		score.Grade = scorecard.GradeAllOK
		return
	}

	warn()
	return
}

func hasPodAntiAffinity(selfLables internal.MapLables, affinity *corev1.Affinity) bool {
	for _, pref := range affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
		if pref.PodAffinityTerm.TopologyKey == "kubernetes.io/hostname" {
			if selector, err := metav1.LabelSelectorAsSelector(pref.PodAffinityTerm.LabelSelector); err == nil {
				if selector.Matches(internal.MapLables(selfLables)) {
					return true
				}
			}
		}
	}

	for _, req := range affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution {
		if req.TopologyKey == "kubernetes.io/hostname" {
			if selector, err := metav1.LabelSelectorAsSelector(req.LabelSelector); err == nil {
				if selector.Matches(internal.MapLables(selfLables)) {
					return true
				}
			}
		}
	}

	return false
}
