package score

import (
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

func scorePodHasNetworkPolicy(allNetpols []networkingv1.NetworkPolicy) func(spec corev1.PodTemplateSpec) TestScore {
	return func(podSpec corev1.PodTemplateSpec) (score TestScore) {
		score.Name = "Pod Has NetworkPolicy"

		hasMatchingEgressNetpol := false
		hasMatchingIngressNetpol := false

		for _, netPol := range allNetpols {
			matchLabels := netPol.Spec.PodSelector.MatchLabels

			for labelKey, labelVal := range matchLabels {
				if podLabelVal, ok := podSpec.Labels[labelKey]; ok && podLabelVal == labelVal {

					for _, policyType := range netPol.Spec.PolicyTypes {
						if policyType == networkingv1.PolicyTypeIngress {
							hasMatchingIngressNetpol = true
						}
						if policyType == networkingv1.PolicyTypeEgress {
							hasMatchingEgressNetpol = true
						}
					}

				}
			}
		}

		if hasMatchingEgressNetpol && hasMatchingIngressNetpol {
			score.Grade = 10
		} else if hasMatchingEgressNetpol && !hasMatchingIngressNetpol {
			score.Grade = 5
			score.Comments = append(score.Comments, "The pod does not have a matching ingress network policy")
		} else if hasMatchingIngressNetpol && !hasMatchingEgressNetpol {
			score.Grade = 5
			score.Comments = append(score.Comments, "The pod does not have a matching egress network policy")
		} else {
			score.Grade = 0
			score.Comments = append(score.Comments, "The pod does not have a matching network policy")
		}

		return
	}
}
