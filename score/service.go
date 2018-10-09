package score

import "github.com/zegl/kube-score/scorecard"
import corev1 "k8s.io/api/core/v1"

// scoreServiceTargetsPod checks if a Service targets a pod and issues a critical warning if no matching pod
// could be found
func scoreServiceTargetsPod(pods []corev1.Pod, podspecers []PodSpecer) func(spec corev1.ServiceSpec) scorecard.TestScore {
	var allPodsWithLabels []map[string]string

	for _, pod := range pods {
		allPodsWithLabels = append(allPodsWithLabels, pod.Labels)
	}

	for _, podSpec := range podspecers {
		allPodsWithLabels = append(allPodsWithLabels, podSpec.GetPodTemplateSpec().Labels)
	}

	return func(spec corev1.ServiceSpec) (score scorecard.TestScore) {
		score.Name = "Service Targets Pod"

		hasMatch := false

		for _, podLables := range allPodsWithLabels {
			matchCount := 0
			for selectorKey, selectorVal := range spec.Selector {
				if labelVal, ok := podLables[selectorKey]; ok && labelVal == selectorVal {
					matchCount++
				}
			}
			if matchCount == len(spec.Selector) {
				hasMatch = true
			}
		}

		if hasMatch {
			score.Grade = 10
		} else {
			score.AddComment("" , "The services selector does not match any pods", "")
		}

		return
	}
}
