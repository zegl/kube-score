package probes

import (
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/score/internal"
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
)

func Register(allChecks *checks.Checks, services ks.Services) {
	allChecks.RegisterPodCheck("Pod Probes", `Makes sure that all Pods have safe probe configurations`, containerProbes(services.Services()))
}

// containerProbes returns a function that checks if all probes are defined correctly in the Pod.
// Only one probe of each type is required on the entire pod.
// ReadinessProbes are not required if the pod is not targeted by a Service.
//
// containerProbes takes a slice of all defined Services as input.
func containerProbes(allServices []ks.Service) func(ks.PodSpecer) (scorecard.TestScore, error) {
	return func(ps ks.PodSpecer) (score scorecard.TestScore, err error) {
		typeMeta := ps.GetTypeMeta()
		if typeMeta.Kind == "CronJob" && typeMeta.GroupVersionKind().Group == "batch" || typeMeta.Kind == "Job" && typeMeta.GroupVersionKind().Group == "batch" {
			score.Grade = scorecard.GradeAllOK
			return score, nil
		}

		podTemplate := ps.GetPodTemplateSpec()
		allContainers := podTemplate.Spec.InitContainers
		allContainers = append(allContainers, podTemplate.Spec.Containers...)

		hasReadinessProbe := false
		hasLivenessProbe := false
		probesAreIdentical := false
		isTargetedByService := false

		for _, s := range allServices {
			if podIsTargetedByService(podTemplate, s.Service()) {
				isTargetedByService = true
				break
			}
		}

		if podTemplate.Spec.ServiceAccountName != "" {
			isTargetedByService = true
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
			return score, nil
		}

		if !isTargetedByService {
			score.Grade = scorecard.GradeAllOK
			score.Skipped = true
			score.AddComment("", "The pod is not targeted by a service, skipping probe checks.", "")
			return score, nil
		}

		if !hasReadinessProbe {
			score.Grade = scorecard.GradeCritical
			score.AddCommentWithURL("", "Container is missing a readinessProbe",
				"A readinessProbe should be used to indicate when the service is ready to receive traffic. "+
					"Without it, the Pod is risking to receive traffic before it has booted. "+
					"It's also used during rollouts, and can prevent downtime if a new version of the application is failing.",
				"https://github.com/zegl/kube-score/blob/master/README_PROBES.md",
			)
			return score, nil
		}

		if !hasLivenessProbe {
			score.Grade = scorecard.GradeAlmostOK
			score.AddCommentWithURL("", "Container is missing a livenessProbe",
				"A livenessProbe can be used to restart the container if it's deadlocked or has crashed without exiting. "+
					"It's only recommended to setup a livenessProbe if you really need one.",
				"https://github.com/zegl/kube-score/blob/master/README_PROBES.md",
			)
			return score, nil
		}

		score.Grade = scorecard.GradeAllOK

		return score, nil
	}
}

func podIsTargetedByService(pod corev1.PodTemplateSpec, service corev1.Service) bool {
	if pod.Namespace != service.Namespace {
		return false
	}

	return internal.LabelSelectorMatchesLabels(
		service.Spec.Selector,
		pod.GetObjectMeta().GetLabels(),
	)
}
