package apps

import (
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

type testcase struct {
	replicas      *int32
	affinity      *corev1.Affinity
	expectedGrade scorecard.Grade
}

func antiAffinityTestCases() []testcase {
	i := func(i int32) *int32 {
		return &i
	}

	return []testcase{
		{
			// No affinity configured
			expectedGrade: scorecard.GradeWarning,
			replicas:      i(5),
		},
		{
			// OK! (required)
			expectedGrade: scorecard.GradeAllOK,
			replicas:      i(5),
			affinity: &corev1.Affinity{
				PodAntiAffinity: &corev1.PodAntiAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
						{
							TopologyKey: "kubernetes.io/hostname",
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app": "foo",
								},
							},
						},
					},
				},
			},
		},
		{
			// OK (preferred)
			expectedGrade: scorecard.GradeAllOK,
			replicas:      i(5),
			affinity: &corev1.Affinity{
				PodAntiAffinity: &corev1.PodAntiAffinity{
					PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{
						{
							Weight: 100,
							PodAffinityTerm: corev1.PodAffinityTerm{
								TopologyKey: "kubernetes.io/hostname",
								LabelSelector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"app": "foo",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			// Not matching app label
			expectedGrade: scorecard.GradeWarning,
			replicas:      i(5),
			affinity: &corev1.Affinity{
				PodAntiAffinity: &corev1.PodAntiAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
						{
							TopologyKey: "kubernetes.io/hostname",
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app": "not-foo",
								},
							},
						},
					},
				},
			},
		},
		{
			// Less than two replicas
			expectedGrade: scorecard.GradeAllOK,
			replicas:      i(1),
		},
	}
}

func TestStatefulsetHasAntiAffinity(t *testing.T) {
	for caseID, tc := range antiAffinityTestCases() {
		s := appsv1.StatefulSet{
			Spec: appsv1.StatefulSetSpec{
				Replicas: tc.replicas,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "foo",
						},
					},
					Spec: corev1.PodSpec{
						Affinity: tc.affinity,
					},
				},
			},
		}

		score := statefulsetHasAntiAffinity(s)
		assert.Equal(t, tc.expectedGrade, score.Grade, "caseID=%d", caseID)
	}
}

func TestDeploymentHasAntiAffinity(t *testing.T) {
	for caseID, tc := range antiAffinityTestCases() {
		s := appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Replicas: tc.replicas,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "foo",
						},
					},
					Spec: corev1.PodSpec{
						Affinity: tc.affinity,
					},
				},
			},
		}

		score := deploymentHasAntiAffinity(s)
		assert.Equal(t, tc.expectedGrade, score.Grade, "caseID=%d", caseID)
	}
}
