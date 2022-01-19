package ephemeralstorage

import (
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Register(allChecks *checks.Checks) {
	allChecks.RegisterPodCheck("Container Ephemeral Storage Requests and Limits", `Makes sure that all pods have ephemeral-storage requests and limits set.`, containerStorageEphemeralRequestsAndLimits)
	allChecks.RegisterPodCheck("Container Ephemeral Storage Requests Equal Limits", `Makes sure that all pods have the same ephemeral-storage requests as limits set.`, containerStorageEphemeralRequestsEqualLimits)
}

func containerStorageEphemeralRequestsAndLimits(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {

	allContainers := podTemplate.Spec.InitContainers
	allContainers = append(allContainers, podTemplate.Spec.Containers...)

	hasMissingLimit := false
	hasMissingRequest := false

	for _, container := range allContainers {
		if container.Resources.Limits.StorageEphemeral().IsZero() {
			score.AddComment(container.Name, "Ephemeral Storage limit is not set", "Resource limits are recommended to avoid resource DDOS. Set resources.limits.ephemeral-storage")
			hasMissingLimit = true
		}
		if container.Resources.Requests.StorageEphemeral().IsZero() {
			score.AddComment(container.Name, "Ephemeral Storage request is not set", "Resource requests are recommended to make sure the application can start and run without crashing. Set resource.requests.ephemeral-storage")
			hasMissingRequest = true
		}
	}

	if hasMissingLimit {
		score.Grade = scorecard.GradeCritical
	}
	if hasMissingRequest {
		score.Grade = scorecard.GradeWarning
	}

	return
}

func containerStorageEphemeralRequestsEqualLimits(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {

	allContainers := podTemplate.Spec.InitContainers
	allContainers = append(allContainers, podTemplate.Spec.Containers...)

	resourcesDoNotMatch := false

	for _, container := range allContainers {
		requests := &container.Resources.Requests
		limits := &container.Resources.Limits
		if !requests.StorageEphemeral().Equal(*limits.StorageEphemeral()) {
			score.AddComment(container.Name, "Ephemeral Storage requests does not match limits", "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.ephemeral-storage == resources.limits.ephemeral-storage")
			resourcesDoNotMatch = true
		}
	}

	if resourcesDoNotMatch {
		score.Grade = scorecard.GradeCritical
	} else {
		score.Grade = scorecard.GradeAllOK
	}

	return
}
