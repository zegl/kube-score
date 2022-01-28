package ephemeralstorage

import (
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Register(allChecks *checks.Checks) {
	allChecks.RegisterPodCheck("Container Ephemeral Storage Requests and Limits",
		"Makes sure all pods have ephemeral-storage requests and limits set", containerStorageEphemeralRequestAndLimit)
}

func containerStorageEphemeralRequestAndLimit(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {

	allContainers := podTemplate.Spec.InitContainers
	allContainers = append(allContainers, podTemplate.Spec.Containers...)

	score.Grade = scorecard.GradeAllOK

	for _, container := range allContainers {
		if container.Resources.Limits.StorageEphemeral().IsZero() {
			score.AddComment(container.Name, "Ephemeral Storage limit is not set",
				"Resource limits are recommended to avoid resource DDOS. Set resources.limits.ephemeral-storage")
			score.Grade = scorecard.GradeCritical
		} else if container.Resources.Requests.StorageEphemeral().IsZero() {
			score.AddComment(container.Name, "Ephemeral Storage request is not set",
				"Resource requests are recommended to make sure the application can start and run without crashing. Set resource.requests.ephemeral-storage")
			score.Grade = scorecard.GradeWarning
		} else if !container.Resources.Limits.StorageEphemeral().IsZero() && !container.Resources.Requests.StorageEphemeral().IsZero() {
			requests := &container.Resources.Requests
			limits := &container.Resources.Limits
			if !requests.StorageEphemeral().Equal(*limits.StorageEphemeral()) {
				score.AddComment(container.Name, "Ephemeral Storage request does not match limit", "Having equal requests and limits is recommended to avoid node resource DDOS during spikes")
				score.Grade = scorecard.GradeCritical
			}
		}
	}

	return
}
