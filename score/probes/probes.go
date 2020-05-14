package probes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks, services ks.Services) {
	allChecks.RegisterPodCheck("Pod Probes", `Makes sure that all Pods have safe probe configurations`, containerProbes(services.Services()))
}

// containerProbes returns a function that checks if all probes are defined correctly in the Pod.
// Only one probe of each type is required on the entire pod.
// ReadinessProbes are not required if the pod is not targeted by a Service.
//
// containerProbes takes a slice of all defined Services as input.
func containerProbes(allServices []corev1.Service) func(corev1.PodTemplateSpec, metav1.TypeMeta) scorecard.TestScore {
	return func(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {
		if typeMeta.Kind == "CronJob" && typeMeta.GroupVersionKind().Group == "batch" || typeMeta.Kind == "Job" && typeMeta.GroupVersionKind().Group == "batch" {
			score.Grade = scorecard.GradeAllOK
			return score
		}

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

		if hasLivenessProbe && hasReadinessProbe && probesAreIdentical {
			score.Grade = scorecard.GradeCritical
			score.AddCommentWithURL(
				"", "Container has the same readiness and liveness probe",
				"Using the same probe for liveness and readiness is very likely dangerous. Generally it's better to avoid the livenessProbe than re-using the readinessProbe.",
				"https://github.com/zegl/kube-score/blob/master/README_PROBES.md",
			)
			return score
		}

		if !isTargetedByService {
			score.Grade = scorecard.GradeAllOK
			score.AddComment("", "The pod is not targeted by a service, skipping probe checks.", "")
			return score
		}

		if !hasReadinessProbe {
			score.Grade = scorecard.GradeCritical
			score.AddCommentWithURL("", "Container is missing a readinessProbe",
				"A readinessProbe should be used to indicate when the service is ready to receive traffic. "+
					"Without it, the Pod is risking to receive traffic before it has booted. "+
					"It's also used during rollouts, and can prevent downtime if a new version of the application is failing.",
				"https://github.com/zegl/kube-score/blob/master/README_PROBES.md",
			)
			return score
		}

		if !hasLivenessProbe {
			score.Grade = scorecard.GradeAlmostOK
			score.AddCommentWithURL("", "Container is missing a livenessProbe",
				"A livenessProbe can be used to restart the container if it's deadlocked or has crashed without exiting. "+
					"It's only recommended to setup a livenessProbe if you really need one.",
				"https://github.com/zegl/kube-score/blob/master/README_PROBES.md",
			)
			return score
		}

		score.Grade = scorecard.GradeAllOK

		return score
	}
}
