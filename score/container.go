package score

import (
	corev1 "k8s.io/api/core/v1"
	"github.com/zegl/kube-score/scorecard"
	"strings"
)

func scoreContainerLimits(podTemplate corev1.PodTemplateSpec) (score scorecard.TestScore) {
	score.Name = "Container Resources"

	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	hasMissingLimit := false
	hasMissingRequest := false

	for _, container := range allContainers {
		if container.Resources.Limits.Cpu().IsZero() {
			score.Comments = append(score.Comments, "CPU limit is not set")
			hasMissingLimit = true
		}
		if container.Resources.Limits.Memory().IsZero() {
			score.Comments = append(score.Comments, "Memory limit is not set")
			hasMissingLimit = true
		}
		if container.Resources.Requests.Cpu().IsZero() {
			score.Comments = append(score.Comments, "CPU request is not set")
			hasMissingRequest = true
		}
		if container.Resources.Requests.Memory().IsZero() {
			score.Comments = append(score.Comments, "Memory request is not set")
			hasMissingRequest = true
		}
	}

	if len(allContainers) == 0 {
		score.Grade = 0
		score.Comments = append(score.Comments, "No containers defined")
	} else if hasMissingLimit {
		score.Grade = 0
	} else if hasMissingRequest {
		score.Grade = 5
	} else {
		score.Grade = 10
	}

	return
}

func scoreContainerImageTag(podTemplate corev1.PodTemplateSpec) (score scorecard.TestScore) {
	score.Name = "Container Image Tag"

	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	hasTagLatest := false

	for _, container := range allContainers{
		imageParts := strings.Split(container.Image, ":")
		imageVersion := imageParts[len(imageParts)-1]

		if imageVersion == "latest" {
			score.Comments = append(score.Comments, "Image with latest tag")
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

func scoreContainerImagePullPolicy(podTemplate corev1.PodTemplateSpec) (score scorecard.TestScore) {
	score.Name = "Container Image Pull Policy"

	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	hasNonAlways := false

	for _, container := range allContainers{
		if container.ImagePullPolicy != corev1.PullAlways {
			score.Comments = append(score.Comments, "ImagePullPolicy is not set to PullAlways")
			hasNonAlways = true
		}
	}

	if hasNonAlways {
		score.Grade = 0
	} else {
		score.Grade = 10
	}

	return
}
