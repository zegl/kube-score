package score

import (
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
)

func scoreContainerProbes(allServices []corev1.Service) func(corev1.PodTemplateSpec) scorecard.TestScore {
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
			} else {
				if isTargetedByService {
					score.AddComment(container.Name, "Container is missing a readinessProbe", "Without a readinessProbe Services will start sending traffic to this pod before it's ready")
				}
			}

			if container.LivenessProbe != nil {
				hasLivenessProbe = true
			} else {
				score.AddComment(container.Name, "Container is missing a livenessProbe", "Without a livenessProbe kubelet can not restart the Pod if it has crashed")
			}

			if container.ReadinessProbe != nil && container.LivenessProbe != nil {

				r := container.ReadinessProbe
				l := container.LivenessProbe

				if r.HTTPGet != nil && l.HTTPGet != nil {
					if r.HTTPGet.Path == l.HTTPGet.Path &&
						r.HTTPGet.Port.IntValue() == l.HTTPGet.Port.IntValue() {
						probesAreIdentical = true
						score.AddComment(container.Name, "Container has the same readiness and liveness probe", "It's recommended to have different probes for the two different purposes.")
					}
				}

				if r.TCPSocket != nil && l.TCPSocket != nil {
					if r.TCPSocket.Port == l.TCPSocket.Port {
						probesAreIdentical = true
						score.AddComment(container.Name, "Container has the same readiness and liveness probe", "It's recommended to have different probes for the two different purposes.")
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
							score.AddComment(container.Name, "Container has the same readiness and liveness probe", "It's recommended to have different probes for the two different purposes.")
						}
					}
				}

			}
		}

		if hasLivenessProbe && (hasReadinessProbe || !isTargetedByService) {
			if !probesAreIdentical {
				score.Grade = 10
			} else {
				score.Grade = 7
			}
		} else if !hasReadinessProbe && !hasLivenessProbe {
			score.Grade = 0
		} else if isTargetedByService && !hasReadinessProbe {
			score.Grade = 0
		} else if !hasLivenessProbe {
			score.Grade = 5
		} else {
			score.Grade = 0
		}

		return score
	}
}
