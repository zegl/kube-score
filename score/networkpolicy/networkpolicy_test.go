package networkpolicy

import (
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestPodHasNetworkPolicy(t *testing.T) {
	t.Parallel()
	cases := []struct {
		polTypes    []v1.PolicyType
		expected    scorecard.Grade
		ingress     []v1.NetworkPolicyIngressRule
		egress      []v1.NetworkPolicyEgressRule
		selectorVal string
	}{
		{
			polTypes:    []v1.PolicyType{v1.PolicyTypeIngress},
			expected:    scorecard.GradeWarning, // has no egress
			selectorVal: "test-a",
		},
		{
			polTypes:    []v1.PolicyType{v1.PolicyTypeEgress},
			expected:    scorecard.GradeWarning, // has no ingress
			selectorVal: "test-a",
		},
		{
			polTypes:    []v1.PolicyType{v1.PolicyTypeEgress, v1.PolicyTypeIngress},
			expected:    scorecard.GradeAllOK,
			selectorVal: "test-a",
		},
		{
			polTypes:    []v1.PolicyType{},
			expected:    scorecard.GradeWarning, // has no egress
			selectorVal: "test-a",
		},
		{
			polTypes:    nil,
			expected:    scorecard.GradeWarning, // has no egress
			selectorVal: "test-a",
		},
		{
			polTypes:    []v1.PolicyType{},
			egress:      []v1.NetworkPolicyEgressRule{{}, {}},
			expected:    scorecard.GradeAllOK,
			selectorVal: "test-a",
		},
		{
			polTypes:    []v1.PolicyType{},
			ingress:     []v1.NetworkPolicyIngressRule{{}, {}},
			expected:    scorecard.GradeWarning, // has no ingress
			selectorVal: "test-a",
		},
		{
			polTypes:    []v1.PolicyType{},
			ingress:     []v1.NetworkPolicyIngressRule{{}, {}},
			egress:      []v1.NetworkPolicyEgressRule{{}, {}},
			expected:    scorecard.GradeAllOK,
			selectorVal: "test-a",
		},
		{
			polTypes:    nil,
			ingress:     []v1.NetworkPolicyIngressRule{{}, {}},
			egress:      []v1.NetworkPolicyEgressRule{{}, {}},
			expected:    scorecard.GradeAllOK,
			selectorVal: "test-a",
		},
		{
			polTypes:    []v1.PolicyType{v1.PolicyTypeEgress, v1.PolicyTypeIngress},
			expected:    scorecard.GradeCritical, // pod has no ingress matching
			selectorVal: "test-not-matching",
		},
	}

	for caseID, tc := range cases {
		pol := v1.NetworkPolicy{
			Spec: v1.NetworkPolicySpec{
				PodSelector: metav1.LabelSelector{
					MatchLabels: map[string]string{"test": tc.selectorVal},
				},
				Ingress:     tc.ingress,
				Egress:      tc.egress,
				PolicyTypes: tc.polTypes,
			},
		}

		pod := corev1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"test": "test-a",
				},
			},
			Spec: corev1.PodSpec{},
		}

		fn := podHasNetworkPolicy([]v1.NetworkPolicy{pol})
		score := fn(corev1.PodTemplateSpec{ObjectMeta: pod.ObjectMeta, Spec: pod.Spec}, pod.TypeMeta)
		assert.Equal(t, tc.expected, score.Grade, "caseID = %d", caseID)
	}
}
