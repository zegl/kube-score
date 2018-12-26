package service

import (
	ks "github.com/zegl/kube-score"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
)

func Register(allChecks *checks.Checks, pods ks.Pods, podspeccers ks.PodSpeccers) {
	allChecks.RegisterServiceCheck("Service Targets Pod", `Makes sure that all Services targets a Pod`, serviceTargetsPod(pods.Pods(), podspeccers.PodSpeccers()))
	allChecks.RegisterServiceCheck("Service Type", `Makes sure that the Service type is not NodePort`, serviceType)
}

// serviceTargetsPod checks if a Service targets a pod and issues a critical warning if no matching pod
// could be found
func serviceTargetsPod(pods []corev1.Pod, podspecers []ks.PodSpecer) func(corev1.Service) scorecard.TestScore {
	podsInNamespace := make(map[string][]map[string]string)
	for _, pod := range pods {
		if _, ok := podsInNamespace[pod.Namespace]; !ok {
			podsInNamespace[pod.Namespace] = []map[string]string{}
		}
		podsInNamespace[pod.Namespace] = append(podsInNamespace[pod.Namespace], pod.Labels)
	}
	for _, podSpec := range podspecers {
		if _, ok := podsInNamespace[podSpec.GetObjectMeta().Namespace]; !ok {
			podsInNamespace[podSpec.GetObjectMeta().Namespace] = []map[string]string{}
		}
		podsInNamespace[podSpec.GetObjectMeta().Namespace] = append(podsInNamespace[podSpec.GetObjectMeta().Namespace], podSpec.GetPodTemplateSpec().Labels)
	}

	return func(service corev1.Service) (score scorecard.TestScore) {
		// Services of type ExternalName does not have a selector
		if service.Spec.Type == corev1.ServiceTypeExternalName {
			score.Grade = scorecard.GradeAllOK
			return
		}

		hasMatch := false

		for _, podLables := range podsInNamespace[service.Namespace] {
			matchCount := 0
			for selectorKey, selectorVal := range service.Spec.Selector {
				if labelVal, ok := podLables[selectorKey]; ok && labelVal == selectorVal {
					matchCount++
				}
			}
			if matchCount == len(service.Spec.Selector) {
				hasMatch = true
			}
		}

		if hasMatch {
			score.Grade = scorecard.GradeAllOK
		} else {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "The services selector does not match any pods", "")
		}

		return
	}
}

func serviceType(service corev1.Service) (score scorecard.TestScore) {
	if service.Spec.Type == corev1.ServiceTypeNodePort {
		score.Grade = scorecard.GradeWarning
		score.AddComment("", "The service is of type NodePort", "NodePort services should be avoided as they are insecure, and can not be used together with NetworkPolicies. LoadBalancers or use of an Ingress is recommended over NodePorts.")
		return
	}

	score.Grade = scorecard.GradeAllOK
	return
}
