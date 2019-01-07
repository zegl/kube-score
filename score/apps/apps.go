package apps

import (
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func Register(allChecks *checks.Checks) {
	allChecks.RegisterDeploymentCheck("Deployment has host PodAntiAffinity", "Makes sure that a podAntiAffinity has been set that prevents multiple pods from being scheduled on the same node", deploymentHasAntiAffinity)
	allChecks.RegisterStatefulSetCheck("StatefulSet has host PodAntiAffinity", "Makes sure that a podAntiAffinity has been set that prevents multiple pods from being scheduled on the same node", statefulsetHasAntiAffinity)
}

func deploymentHasAntiAffinity(deployment appsv1.Deployment) (score scorecard.TestScore) {
	// Ignore if the deployment only has a single replica
	if deployment.Spec.Replicas == nil || *deployment.Spec.Replicas < 2 {
		score.Grade = scorecard.GradeAllOK
		score.AddComment("", "Skipped", "Skipped because the deployment has less than 2 replicas")
		return
	}

	affinity := deployment.Spec.Template.Spec.Affinity
	if affinity == nil || affinity.PodAntiAffinity == nil {
		score.Grade = scorecard.GradeWarning
		return
	}

	if hasPodAntiAffinity(affinity) {
		score.Grade = scorecard.GradeAllOK
		return
	}

	score.Grade = scorecard.GradeWarning
	score.AddComment("", "Deployment does not have a host podAntiAffinity set", "It's recommended to set a podAntiAffinity that stops multiple pods from a deployment from beeing scheduled on the same node. This increases availability in case the node becomes unavailable.")
	return
}

func statefulsetHasAntiAffinity(statefulset appsv1.StatefulSet) (score scorecard.TestScore) {
	// Ignore if the statefulset only has a single replica
	if statefulset.Spec.Replicas == nil || *statefulset.Spec.Replicas < 2 {
		score.Grade = scorecard.GradeAllOK
		score.AddComment("", "Skipped", "Skipped because the statefulset has less than 2 replicas")
		return
	}

	affinity := statefulset.Spec.Template.Spec.Affinity
	if affinity == nil || affinity.PodAntiAffinity == nil {
		score.Grade = scorecard.GradeWarning
		return
	}

	if hasPodAntiAffinity(affinity) {
		score.Grade = scorecard.GradeAllOK
		return
	}

	score.Grade = scorecard.GradeWarning
	score.AddComment("", "StatefulSet does not have a host podAntiAffinity set", "It's recommended to set a podAntiAffinity that stops multiple pods from a statefulset from beeing scheduled on the same node. This increases availability in case the node becomes unavailable.")
	return
}

func hasPodAntiAffinity(affinity *corev1.Affinity) bool {
	for _, pref := range affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
		if pref.PodAffinityTerm.TopologyKey == "kubernetes.io/hostname" {
			return true
		}
	}

	for _, pref := range affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution {
		if pref.TopologyKey == "kubernetes.io/hostname" {
			return true
		}
	}

	return false
}
