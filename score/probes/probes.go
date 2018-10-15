package probes

import (
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
)

// ScoreContainerProbes returns a function that checks if all probes are defined correctly in the Pod.
// Only one probe of each type is required on the entire pod.
// ReadinessProbes are not required if the pod is not targeted by a Service.
//
// ScoreContainerProbes takes a slice of all defined Services as input.
func ScoreContainerProbes(allServices []corev1.Service) func(corev1.PodTemplateSpec) scorecard.TestScore {
	return func(podTemplate corev1.PodTemplateSpec) (score scorecard.TestScore) {
		score.Name = "Pod Probes"

		allContainers := podTemplate.Spec.InitContainers
		allContainers = append(allContainers, podTemplate.Spec.Containers...)

		hasReadinessProbe := false
		hasLivenessProbe := false
		probesAreIdentical := false
		isTargetedByService := false

		for _, service := range allServices {
			if podTemplate.Namespace == service.Namespace {
				for selectorKey, selectorVal := range service.Spec.Selector {
					if podLabelVal, ok := podTemplate.Labels[selectorKey]; ok && podLabelVal == selectorVal {
						isTargetedByService = true
					}
				}
			}
		}

		for _, container := range allContainers {
			if container.ReadinessProbe != nil {
				hasReadinessProbe = true
			}

			if container.LivenessProbe != nil {
				hasLivenessProbe = true
			}

			if container.ReadinessProbe != nil && container.LivenessProbe != nil {

				r := container.ReadinessProbe
				l := container.LivenessProbe

				if r.HTTPGet != nil && l.HTTPGet != nil {
					if r.HTTPGet.Path == l.HTTPGet.Path &&
						r.HTTPGet.Port.IntValue() == l.HTTPGet.Port.IntValue() {
						probesAreIdentical = true
					}
				}

				if r.TCPSocket != nil && l.TCPSocket != nil {
					if r.TCPSocket.Port == l.TCPSocket.Port {
						probesAreIdentical = true
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
						}
					}
				}

			}
		}

		if hasLivenessProbe && (hasReadinessProbe || !isTargetedByService) {
			if !probesAreIdentical {
				score.Grade = scorecard.GradeAllOK
			} else {
				score.Grade = scorecard.GradeAlmostOK
				score.AddComment("", "Pod has the same readiness and liveness probe", "It's recommended to have different probes for the two different purposes.")
			}
		} else if !hasReadinessProbe && !hasLivenessProbe {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "Container is missing a readinessProbe", "Without a readinessProbe Services will start sending traffic to this pod before it's ready")
			score.AddComment("", "Container is missing a livenessProbe", "Without a livenessProbe kubelet can not restart the Pod if it has crashed")
		} else if isTargetedByService && !hasReadinessProbe {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "Container is missing a readinessProbe", "Without a readinessProbe Services will start sending traffic to this pod before it's ready")
		} else if !hasLivenessProbe {
			score.Grade = scorecard.GradeWarning
			score.AddComment("", "Pod is missing a livenessProbe", "Without a livenessProbe kubelet can not restart the Pod if it has crashed")
		}

		return score
	}
}
