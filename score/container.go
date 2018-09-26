package score

import (
	"strings"

	"github.com/zegl/kube-score/scorecard"

	corev1 "k8s.io/api/core/v1"
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

func scoreContainerProbes(podTemplate corev1.PodTemplateSpec) (score scorecard.TestScore) {
	score.Name = "Pod Probes"

	allContainers := podTemplate.Spec.InitContainers
	allContainers = append(allContainers, podTemplate.Spec.Containers...)

	hasReadinessProbe := true
	hasLivenessProbe := true

	probesAreIdentical := false

	for _, container := range allContainers {
		if container.ReadinessProbe == nil  {
			hasReadinessProbe = false
			score.Comments = append(score.Comments, "Container is missing readinessProbe")
		}

		if container.LivenessProbe == nil {
			hasLivenessProbe = false
			score.Comments = append(score.Comments, "Container is missing livenessProbe")
		}

		if container.ReadinessProbe != nil && container.LivenessProbe != nil {

			r := container.ReadinessProbe
			l := container.LivenessProbe

			if r.HTTPGet != nil && l.HTTPGet != nil {
				if r.HTTPGet.Path == l.HTTPGet.Path &&
					r.HTTPGet.Port.IntValue() == l.HTTPGet.Port.IntValue() {
						probesAreIdentical = true
					score.Comments = append(score.Comments, "Container has the same readiness and liveness probe")
				}
			}

			if r.TCPSocket != nil && l.TCPSocket != nil {
				if r.TCPSocket.Port == l.TCPSocket.Port {
					probesAreIdentical = true
					score.Comments = append(score.Comments, "Container has the same readiness and liveness probe")
				}
			}

			if r.Exec != nil && l.Exec != nil {
				if len(r.Exec.Command) == len(l.Exec.Command) {
					hasDifferent := false
					for i, v := range r.Exec.Command {
						if l.Exec.Command[i] != v {
							hasDifferent = true
							break
						}
					}

					if !hasDifferent {
						probesAreIdentical = true
						score.Comments = append(score.Comments, "Container has the same readiness and liveness probe")
					}
				}
			}

		}
	}

	if hasReadinessProbe && hasLivenessProbe {
		if !probesAreIdentical {
			score.Grade = 10
		} else {
			score.Grade = 7
		}
	} else if hasLivenessProbe || hasReadinessProbe {
		score.Grade = 5
	} else {
		score.Grade = 0
	}

	return
}

func scoreContainerSecurityContext(podTemplate corev1.PodTemplateSpec) (score scorecard.TestScore) {
	score.Name = "Container Security Context"

	allContainers := podTemplate.Spec.InitContainers
	allContainers = append(allContainers, podTemplate.Spec.Containers...)

	hasPrivileged := false
	hasWritableRootFS := false
	hasLowUserID := false
	hasLowGroupID := false

	for _, container := range allContainers {

		if container.SecurityContext == nil {
			continue
		}

		sec := container.SecurityContext

		if sec.Privileged != nil && *sec.Privileged {
			hasPrivileged = true
			score.Comments = append(score.Comments, "The pod has a privileged container")
		}

		if sec.ReadOnlyRootFilesystem != nil && *sec.ReadOnlyRootFilesystem == false {
			hasWritableRootFS = true
			score.Comments = append(score.Comments, "The pod has a container with a writable root filesystem")
		}

		if sec.RunAsUser != nil && *sec.RunAsUser < 10000 {
			hasLowUserID = true
			score.Comments = append(score.Comments, "The pod has a container running with a low user ID")
		}

		if sec.RunAsGroup != nil && *sec.RunAsGroup < 10000 {
			hasLowGroupID = true
			score.Comments = append(score.Comments, "The pod has a container running with a low group ID")
		}
	}

	if hasPrivileged || hasWritableRootFS || hasLowUserID || hasLowGroupID {
		score.Grade = 0
	} else {
		score.Grade = 10
	}

	return
}