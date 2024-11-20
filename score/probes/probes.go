package probes

import (
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/score/internal"
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
)

// Register registers the pod checks, including the new one for identical probes.
func Register(allChecks *checks.Checks, services ks.Services) {
	allChecks.RegisterPodCheck("Pod Probes", `Makes sure that all Pods have safe probe configurations`, containerProbes(services.Services()))
	allChecks.RegisterPodCheck("Pod Probes Identical", `Container has the same readiness and liveness probe`, containerProbesIdentical(services.Services()))
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
		isTargetedByService := isTargetedByService(allServices, podTemplate)

		// Check probes for each container
		for _, container := range allContainers {
			hasReadinessProbe, hasLivenessProbe = checkBasicProbes(container, hasReadinessProbe, hasLivenessProbe)
		}

		// If pod isn't targeted by a service, skip probe checks
		if !isTargetedByService {
			score.Grade = scorecard.GradeAllOK
			score.Skipped = true
			score.AddComment("", "The pod is not targeted by a service, skipping probe checks.", "")
			return score, nil
		}

		// Evaluate probe checks
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

// containerProbesIdentical checks if the container's readiness and liveness probes are identical.
func containerProbesIdentical(allServices []ks.Service) func(ks.PodSpecer) (scorecard.TestScore, error) {
	return func(ps ks.PodSpecer) (score scorecard.TestScore, err error) {
		typeMeta := ps.GetTypeMeta()
		if typeMeta.Kind == "CronJob" && typeMeta.GroupVersionKind().Group == "batch" || typeMeta.Kind == "Job" && typeMeta.GroupVersionKind().Group == "batch" {
			score.Grade = scorecard.GradeAllOK
			return score, nil
		}

		podTemplate := ps.GetPodTemplateSpec()
		allContainers := podTemplate.Spec.InitContainers
		allContainers = append(allContainers, podTemplate.Spec.Containers...)

		probesAreIdentical := false
		for _, container := range allContainers {
			if container.ReadinessProbe != nil && container.LivenessProbe != nil {
				if areProbesIdentical(container.ReadinessProbe, container.LivenessProbe) {
					probesAreIdentical = true
					break
				}
			}
		}

		// If probes are identical, mark it as a critical issue
		if probesAreIdentical {
			score.Grade = scorecard.GradeCritical
			score.AddCommentWithURL(
				"", "Container has the same readiness and liveness probe",
				"Using the same probe for liveness and readiness is very likely dangerous. It's generally better to avoid re-using the same probe.",
				"https://github.com/zegl/kube-score/blob/master/README_PROBES.md",
			)
			return score, nil
		}

		// No identical probes found, return OK grade
		score.Grade = scorecard.GradeAllOK
		return score, nil
	}
}

// areProbesIdentical checks if readiness and liveness probes are identical.
func areProbesIdentical(r, l *corev1.Probe) bool {
	if r.HTTPGet != nil && l.HTTPGet != nil {
		return r.HTTPGet.Path == l.HTTPGet.Path && r.HTTPGet.Port.IntValue() == l.HTTPGet.Port.IntValue()
	}
	if r.TCPSocket != nil && l.TCPSocket != nil {
		return r.TCPSocket.Port == l.TCPSocket.Port
	}
	if r.Exec != nil && l.Exec != nil {
		if len(r.Exec.Command) == len(l.Exec.Command) {
			for i, v := range r.Exec.Command {
				if l.Exec.Command[i] != v {
					return false
				}
			}
			return true
		}
	}
	return false
}

// checkBasicProbes checks for the presence of readiness and liveness probes.
func checkBasicProbes(container corev1.Container, hasReadinessProbe, hasLivenessProbe bool) (bool, bool) {
	if container.ReadinessProbe != nil {
		hasReadinessProbe = true
	}

	if container.LivenessProbe != nil {
		hasLivenessProbe = true
	}

	return hasReadinessProbe, hasLivenessProbe
}

// isTargetedByService checks if the pod is targeted by any of the services.
func isTargetedByService(allServices []ks.Service, podTemplate corev1.PodTemplateSpec) bool {
	for _, s := range allServices {
		if podIsTargetedByService(podTemplate, s.Service()) {
			return true
		}
	}
	return false
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
