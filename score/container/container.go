package container

import (
	"strings"

	"github.com/zegl/kube-score/scorecard"

	corev1 "k8s.io/api/core/v1"
)

// ScoreContainerLimit makes sure that the container has resource requests and limits set
// The check for a CPU limit requirement can be enabled via the requireCpuLimit flag parameter
func ScoreContainerLimits(requireCpuLimit bool) func(corev1.PodTemplateSpec) scorecard.TestScore {
	return func(podTemplate corev1.PodTemplateSpec) (score scorecard.TestScore) {
		score.Name = "Container Resources"

		pod := podTemplate.Spec

		allContainers := pod.InitContainers
		allContainers = append(allContainers, pod.Containers...)

		hasMissingLimit := false
		hasMissingRequest := false

		for _, container := range allContainers {
			if container.Resources.Limits.Cpu().IsZero() && requireCpuLimit {
				score.AddComment(container.Name, "CPU limit is not set", "Resource limits are recommended to avoid resource DDOS. Set resources.limits.cpu")
				hasMissingLimit = true
			}
			if container.Resources.Limits.Memory().IsZero() {
				score.AddComment(container.Name, "Memory limit is not set", "Resource limits are recommended to avoid resource DDOS. Set resources.limits.memory")
				hasMissingLimit = true
			}
			if container.Resources.Requests.Cpu().IsZero() {
				score.AddComment(container.Name, "CPU request is not set", "Resource requests are recommended to make sure that the application can start and run without crashing. Set resources.requests.cpu")
				hasMissingRequest = true
			}
			if container.Resources.Requests.Memory().IsZero() {
				score.AddComment(container.Name, "Memory request is not set", "Resource requests are recommended to make sure that the application can start and run without crashing. Set resources.requests.memory")
				hasMissingRequest = true
			}
		}

		if len(allContainers) == 0 {
			score.Grade = 0
			score.AddComment("", "No containers defined", "")
		} else if hasMissingLimit {
			score.Grade = 0
		} else if hasMissingRequest {
			score.Grade = 5
		} else {
			score.Grade = 10
		}

		return
	}
}

// ScoreContainerImageTag checks that no container is using the ":latest" tag
func ScoreContainerImageTag(podTemplate corev1.PodTemplateSpec) (score scorecard.TestScore) {
	score.Name = "Container Image Tag"

	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	hasTagLatest := false

	for _, container := range allContainers {
		tag := containerTag(container.Image)
		if tag == "" || tag == "latest" {
			score.AddComment(container.Name, "Image with latest tag", "Using a fixed tag is recommended to avoid accidental upgrades")
			hasTagLatest = true
		}
	}

	if hasTagLatest {
		score.Grade = 0
	} else {
		score.Grade = 10
	}

	return
}

// ScoreContainerImagePullPolicy checks if the containers ImagePullPolicy is set to PullAlways
func ScoreContainerImagePullPolicy(podTemplate corev1.PodTemplateSpec) (score scorecard.TestScore) {
	score.Name = "Container Image Pull Policy"

	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	hasNonAlways := false

	for _, container := range allContainers {

		// No defined pull policy
		if container.ImagePullPolicy == corev1.PullPolicy("") {
			tag := containerTag(container.Image)
			if tag != "" && tag != "latest" {
				hasNonAlways = true
			}
		} else {
			if container.ImagePullPolicy != corev1.PullAlways {
				score.AddComment(container.Name, "ImagePullPolicy is not set to Always", "It's recommended to always set the ImagePullPolicy to Always, to make sure that the imagePullSecrets are always correct, and to always get the image you want.")
				hasNonAlways = true
			}
		}
	}

	if hasNonAlways {
		score.Grade = 0
	} else {
		score.Grade = 10
	}

	return
}

// containerTag returns the image tag
// An empty string is returned if the image has no tag
func containerTag(image string) string {
	imageParts := strings.Split(image, ":")
	if len(imageParts) > 1 {
		imageVersion := imageParts[len(imageParts)-1]
		return imageVersion
	}
	return ""
}
