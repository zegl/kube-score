package networkpolicy

import (
	ks "github.com/zegl/kube-score"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/score/internal"
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Register(allChecks *checks.Checks, netpols ks.NetworkPolicies, pods ks.Pods, podspecers ks.PodSpeccers) {
	allChecks.RegisterPodCheck("Pod NetworkPolicy", `Makes sure that all Pods are targeted by a NetworkPolicy`, podHasNetworkPolicy(netpols.NetworkPolicies()))
	allChecks.RegisterNetworkPolicyCheck("NetworkPolicy targets Pod", `Makes sure that all NetworkPolicies targets at least one Pod`, networkPolicyTargetsPod(pods.Pods(), podspecers.PodSpeccers()))
}

// podHasNetworkPolicy returns a function that tests that all pods have matching NetworkPolicies
// podHasNetworkPolicy takes a list of all defined NetworkPolicies as input
func podHasNetworkPolicy(allNetpols []networkingv1.NetworkPolicy) func(spec corev1.PodTemplateSpec) scorecard.TestScore {
	return func(podSpec corev1.PodTemplateSpec) (score scorecard.TestScore) {
		hasMatchingEgressNetpol := false
		hasMatchingIngressNetpol := false

		for _, netPol := range allNetpols {

			// Make sure that the pod and networkpolicy is in the same namespace
			if podSpec.Namespace != netPol.Namespace {
				continue
			}

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

func networkPolicyTargetsPod(pods []corev1.Pod, podspecers []ks.PodSpecer) func(networkingv1.NetworkPolicy) scorecard.TestScore {
	return func(netpol networkingv1.NetworkPolicy) (score scorecard.TestScore) {
		hasMatch := false

		for _, pod := range pods {
			if pod.Namespace != netpol.Namespace {
				continue
			}

			if selector, err := metav1.LabelSelectorAsSelector(&netpol.Spec.PodSelector); err == nil {
				if selector.Matches(internal.MapLables(pod.Labels)) {
					hasMatch = true
					break
				}
			}
		}

		if !hasMatch {
			for _, pod := range podspecers {
				if pod.GetObjectMeta().Namespace != netpol.Namespace {
					continue
				}

				if selector, err := metav1.LabelSelectorAsSelector(&netpol.Spec.PodSelector); err == nil {
					if selector.Matches(internal.MapLables(pod.GetPodTemplateSpec().Labels)) {
						hasMatch = true
						break
					}
				}
			}
		}

		if hasMatch {
			score.Grade = scorecard.GradeAllOK
		} else {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "The NetworkPolicys selector doesn't match any pods", "")
		}

		return
	}
}
