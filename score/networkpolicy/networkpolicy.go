package networkpolicy

import (
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

// ScorePodHasNetworkPolicy returns a function that tests that all pods have matching NetworkPolicies
// ScorePodHasNetworkPolicy takes a list of all defined NetworkPolicies as input
func ScorePodHasNetworkPolicy(allNetpols []networkingv1.NetworkPolicy) func(spec corev1.PodTemplateSpec) scorecard.TestScore {
	return func(podSpec corev1.PodTemplateSpec) (score scorecard.TestScore) {
		score.Name = "Pod NetworkPolicy"

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
			score.Grade = scorecard.GradeAllOK
		} else if hasMatchingEgressNetpol && !hasMatchingIngressNetpol {
			score.Grade = scorecard.GradeWarning
			score.AddComment("", "The pod does not have a matching ingress network policy", "Add a egress policy to the pods NetworkPolicy")
		} else if hasMatchingIngressNetpol && !hasMatchingEgressNetpol {
			score.Grade = scorecard.GradeWarning
			score.AddComment("", "The pod does not have a matching egress network policy", "Add a ingress policy to the pods NetworkPolicy")
		} else {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "The pod does not have a matching network policy", "Create a NetworkPolicy that targets this pod")
		}

		return
	}
}
