package container

import (
	"fmt"
	"strings"

	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Register(allChecks *checks.Checks, cnf config.Configuration) {
	allChecks.RegisterPodCheck("Container Resources", `Makes sure that all pods have resource limits and requests set. The --ignore-container-cpu-limit flag can be used to disable the requirement of having a CPU limit`, containerResources(!cnf.IgnoreContainerCpuLimitRequirement, !cnf.IgnoreContainerMemoryLimitRequirement))
	allChecks.RegisterOptionalPodCheck("Container Resource Requests Equal Limits", `Makes sure that all pods have the same requests as limits on resources set.`, containerResourceRequestsEqualLimits)
	allChecks.RegisterOptionalPodCheck("Container CPU Requests Equal Limits", `Makes sure that all pods have the same CPU requests as limits set.`, containerCPURequestsEqualLimits)
	allChecks.RegisterOptionalPodCheck("Container Memory Requests Equal Limits", `Makes sure that all pods have the same memory requests as limits set.`, containerMemoryRequestsEqualLimits)
	allChecks.RegisterPodCheck("Container Image Tag", `Makes sure that a explicit non-latest tag is used`, containerImageTag)
	allChecks.RegisterPodCheck("Container Image Pull Policy", `Makes sure that the pullPolicy is set to Always. This makes sure that imagePullSecrets are always validated.`, containerImagePullPolicy)
	allChecks.RegisterPodCheck("Container Ephemeral Storage Request and Limit", "Makes sure all pods have ephemeral-storage requests and limits set", containerStorageEphemeralRequestAndLimit)
	allChecks.RegisterOptionalPodCheck("Container Ephemeral Storage Request Equals Limit", "Make sure all pods have matching ephemeral-storage requests and limits", containerStorageEphemeralRequestEqualsLimit)
	allChecks.RegisterOptionalPodCheck("Container Ports Check", "Container Ports Checks", containerPortsCheck)
	allChecks.RegisterPodCheck("Environment Variable Key Duplication", "Makes sure that duplicated environment variable keys are not duplicated", environmentVariableKeyDuplication)
}

// containerResources makes sure that the container has resource requests and limits set
// The check for a CPU limit requirement can be enabled via the requireCPULimit flag parameter
func containerResources(requireCPULimit bool, requireMemoryLimit bool) func(corev1.PodTemplateSpec, metav1.TypeMeta) scorecard.TestScore {
	return func(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {
		pod := podTemplate.Spec

		allContainers := pod.InitContainers
		allContainers = append(allContainers, pod.Containers...)

		hasMissingLimit := false
		hasMissingRequest := false

		for _, container := range allContainers {
			if container.Resources.Limits.Cpu().IsZero() && requireCPULimit {
				score.AddComment(container.Name, "CPU limit is not set", "Resource limits are recommended to avoid resource DDOS. Set resources.limits.cpu")
				hasMissingLimit = true
			}
			if container.Resources.Limits.Memory().IsZero() && requireMemoryLimit {
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

		switch {
		case len(allContainers) == 0:
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "No containers defined", "")
		case hasMissingLimit:
			score.Grade = scorecard.GradeCritical
		case hasMissingRequest:
			score.Grade = scorecard.GradeWarning
		default:
			score.Grade = scorecard.GradeAllOK
		}

		return
	}
}

// containerResourceRequestsEqualLimits checks that all containers have equal requests and limits for CPU and memory resources
func containerResourceRequestsEqualLimits(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {
	cpuScore := containerCPURequestsEqualLimits(podTemplate, typeMeta)
	memoryScore := containerMemoryRequestsEqualLimits(podTemplate, typeMeta)

	score.Grade = scorecard.GradeAllOK
	if cpuScore.Grade == scorecard.GradeCritical {
		score.Grade = scorecard.GradeCritical
		score.Comments = append(score.Comments, cpuScore.Comments...)
	}
	if memoryScore.Grade == scorecard.GradeCritical {
		score.Grade = scorecard.GradeCritical
		score.Comments = append(score.Comments, memoryScore.Comments...)
	}

	return score
}

// containerCPURequestsEqualLimits checks that all containers have equal requests and limits for CPU resources
func containerCPURequestsEqualLimits(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {
	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	resourcesDoNotMatch := false

	for _, container := range allContainers {
		requests := &container.Resources.Requests
		limits := &container.Resources.Limits
		if !requests.Cpu().Equal(*limits.Cpu()) {
			score.AddComment(container.Name, "CPU requests does not match limits", "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.cpu == resources.limits.cpu")
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

// containerMemoryRequestsEqualLimits checks that all containers have equal requests and limits for memory resources
func containerMemoryRequestsEqualLimits(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {
	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	resourcesDoNotMatch := false

	for _, container := range allContainers {
		requests := &container.Resources.Requests
		limits := &container.Resources.Limits
		if !requests.Memory().Equal(*limits.Memory()) {
			score.AddComment(container.Name, "Memory requests does not match limits", "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.memory == resources.limits.memory")
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

// containerImageTag checks that no container is using the ":latest" tag
func containerImageTag(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {
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
		score.Grade = scorecard.GradeCritical
	} else {
		score.Grade = scorecard.GradeAllOK
	}

	return
}

// containerImagePullPolicy checks if the containers ImagePullPolicy is set to PullAlways
func containerImagePullPolicy(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {
	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	// Default to AllOK
	score.Grade = scorecard.GradeAllOK

	for _, container := range allContainers {
		tag := containerTag(container.Image)

		// If the pull policy is not set, and the tag is either empty or latest
		// kubernetes will default to always pull the image
		if container.ImagePullPolicy == corev1.PullPolicy("") && (tag == "" || tag == "latest") {
			continue
		}

		// No defined pull policy
		if container.ImagePullPolicy != corev1.PullAlways || container.ImagePullPolicy == corev1.PullPolicy("") {
			score.AddComment(container.Name, "ImagePullPolicy is not set to Always", "It's recommended to always set the ImagePullPolicy to Always, to make sure that the imagePullSecrets are always correct, and to always get the image you want.")
			score.Grade = scorecard.GradeCritical
		}
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
		}
	}

	return
}

func containerStorageEphemeralRequestEqualsLimit(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {

	allContainers := podTemplate.Spec.InitContainers
	allContainers = append(allContainers, podTemplate.Spec.Containers...)

	score.Grade = scorecard.GradeAllOK

	for _, container := range allContainers {
		if !container.Resources.Limits.StorageEphemeral().IsZero() && !container.Resources.Requests.StorageEphemeral().IsZero() {
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

// List of ports to expose from the container. This is primarily informational. Not specifying a port here
// does not prevent it from being exposed. Specifying it does not expose the port outside the cluster; that
// requires a Service object. However misspecifying elements of this optional Container
func containerPortsCheck(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {

	allContainers := podTemplate.Spec.InitContainers
	allContainers = append(allContainers, podTemplate.Spec.Containers...)

	score.Grade = scorecard.GradeAllOK

	for _, container := range allContainers {
		names := make(map[string]bool)
		for _, port := range container.Ports {
			if len(port.Name) > 0 {
				if _, ok := names[port.Name]; !ok {
					names[port.Name] = true
				} else {
					score.AddComment(container.Name, "Container Port Check", "Container ports.containerPort named ports must be unique")
					score.Grade = scorecard.GradeCritical
				}
			}
			if port.ContainerPort == 0 {
				score.AddComment(container.Name, "Container Port Check", "Container ports.containerPort cannot be empty")
				score.Grade = scorecard.GradeCritical
			}
		}
	}

	return
}

// environmentVariableKeyDuplication checks that no duplicated environment variable keys.
func environmentVariableKeyDuplication(podTemplate corev1.PodTemplateSpec, _ metav1.TypeMeta) (score scorecard.TestScore) {
	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	score.Grade = scorecard.GradeAllOK

	for _, container := range allContainers {
		envs := make(map[string]struct{})
		for _, env := range container.Env {
			if _, duplicated := envs[env.Name]; duplicated {
				msg := fmt.Sprintf("Container environment variable key '%s' is duplicated", env.Name)
				score.AddComment(container.Name, "Environment Variable Key Duplication", msg)
				score.Grade = scorecard.GradeCritical
				continue
			}
			envs[env.Name] = struct{}{}
		}
	}
	return
}
