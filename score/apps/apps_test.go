package apps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

type testcase struct {
	replicas        *int32
	affinity        *corev1.Affinity
	expectedGrade   scorecard.Grade
	expectedSkipped bool
}

func i(i int32) *int32 {
	return &i
}

func antiAffinityTestCases() []testcase {
	return []testcase{
		{
			// No affinity configured
			expectedGrade:   scorecard.GradeWarning,
			replicas:        i(5),
			expectedSkipped: false,
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
			expectedSkipped: false,
		},
		{
			// OK! (required) ( topology.kubernetes.io/zone )
			expectedGrade: scorecard.GradeAllOK,
			replicas:      i(5),
			affinity: &corev1.Affinity{
				PodAntiAffinity: &corev1.PodAntiAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
						{
							TopologyKey: "topology.kubernetes.io/zone",
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app": "foo",
								},
							},
						},
					},
				},
			},
			expectedSkipped: false,
		},
		{
			// OK! (required) ( topology.kubernetes.io/region )
			expectedGrade: scorecard.GradeAllOK,
			replicas:      i(5),
			affinity: &corev1.Affinity{
				PodAntiAffinity: &corev1.PodAntiAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
						{
							TopologyKey: "topology.kubernetes.io/region",
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app": "foo",
								},
							},
						},
					},
				},
			},
			expectedSkipped: false,
		},
		{
			// Not OK! (required) ( some other topology key )
			expectedGrade: scorecard.GradeWarning,
			replicas:      i(5),
			affinity: &corev1.Affinity{
				PodAntiAffinity: &corev1.PodAntiAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
						{
							TopologyKey: "topology.kubernetes.io/what-is-this-key",
							LabelSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app": "foo",
								},
							},
						},
					},
				},
			},
			expectedSkipped: false,
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
			expectedSkipped: false,
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
			expectedSkipped: false,
		},
		{
			// Less than two replicas
			expectedGrade:   0,
			replicas:        i(1),
			expectedSkipped: true,
		},
	}
}

func TestStatefulsetHasAntiAffinity(t *testing.T) {
	t.Parallel()
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

		score, err := statefulsetHasAntiAffinity(s)
		assert.Nil(t, err)
		assert.Equal(t, tc.expectedGrade, score.Grade, "caseID=%d", caseID)
	}
}

func TestDeploymentHasAntiAffinity(t *testing.T) {
	t.Parallel()
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

		score, err := deploymentHasAntiAffinity(s)
		assert.Nil(t, err)
		assert.Equal(t, tc.expectedGrade, score.Grade, "unexpected grade caseID=%d", caseID)
		assert.Equal(t, tc.expectedSkipped, score.Skipped, "unexpected skipped, caseID=%d", caseID)
	}
}

func TestDeploymentTargetedByHpaHasNoReplicasAllOK(t *testing.T) {
	t.Parallel()

	deployment := appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: appsv1.DeploymentSpec{
			Replicas: nil,
		},
	}

	hpas := []ks.HpaTargeter{
		hpav1{
			autoscalingv1.HorizontalPodAutoscaler{
				Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
						Kind:       "Deployment",
						Name:       "foo",
						APIVersion: "apps/v1",
					},
				},
			},
		},
	}

	f := hpaDeploymentNoReplicas(hpas)
	score, err := f(deployment)
	assert.Nil(t, err)
	assert.Equal(t, scorecard.GradeAllOK, score.Grade)
	assert.False(t, score.Skipped)
}

func TestDeploymentTargetedByHpaHasSetReplicasCritical(t *testing.T) {
	t.Parallel()

	deployment := appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: appsv1.DeploymentSpec{
			Replicas: i(30),
		},
	}

	hpas := []ks.HpaTargeter{
		hpav1{
			autoscalingv1.HorizontalPodAutoscaler{
				Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
						Kind:       "Deployment",
						Name:       "foo",
						APIVersion: "apps/v1",
					},
				},
			},
		},
	}

	f := hpaDeploymentNoReplicas(hpas)
	score, err := f(deployment)
	assert.Nil(t, err)
	assert.Equal(t, scorecard.GradeCritical, score.Grade)
	assert.False(t, score.Skipped)
}

func TestDeploymentNotTargetedByHpaIsSkippedAllOKK(t *testing.T) {
	t.Parallel()

	deployment := appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: appsv1.DeploymentSpec{
			Replicas: i(30),
		},
	}

	hpas := []ks.HpaTargeter{
		hpav1{
			autoscalingv1.HorizontalPodAutoscaler{
				Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
						Kind:       "Deployment",
						Name:       "some-other-foo",
						APIVersion: "apps/v1",
					},
				},
			},
		},
	}

	f := hpaDeploymentNoReplicas(hpas)
	score, err := f(deployment)
	assert.Nil(t, err)
	assert.Equal(t, scorecard.GradeAllOK, score.Grade)
	assert.True(t, score.Skipped)
}

func TestDeploymentTargetedByHpaHasNoReplicasAllOKCaseInsensitiveKind(t *testing.T) {
	t.Parallel()

	deployment := appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{Name: "foo"},
		Spec: appsv1.DeploymentSpec{
			Replicas: nil,
		},
	}

	hpas := []ks.HpaTargeter{
		hpav1{
			autoscalingv1.HorizontalPodAutoscaler{
				Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
					ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
						Kind:       "deployment",
						Name:       "foo",
						APIVersion: "apps/v1",
					},
				},
			},
		},
	}

	f := hpaDeploymentNoReplicas(hpas)
	score, err := f(deployment)
	assert.Nil(t, err)
	assert.Equal(t, scorecard.GradeAllOK, score.Grade)
	assert.False(t, score.Skipped)
}

type hpav1 struct {
	autoscalingv1.HorizontalPodAutoscaler
}

func (d hpav1) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d hpav1) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d hpav1) MinReplicas() *int32 {
	return d.Spec.MinReplicas
}

func (d hpav1) HpaTarget() autoscalingv1.CrossVersionObjectReference {
	return d.Spec.ScaleTargetRef
}

func (hpav1) FileLocation() ks.FileLocation {
	return ks.FileLocation{}
}

func TestStatefulSetHasServiceName(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		statefulset     appsv1.StatefulSet
		services        []ks.Service
		expectedErr     error
		expectedGrade   scorecard.Grade
		expectedSkipped bool
	}{
		// Match (no namespace)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "foo-svc",
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			services: []ks.Service{
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "None",
							Selector: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			expectedErr:     nil,
			expectedGrade:   scorecard.GradeAllOK,
			expectedSkipped: false,
		},

		// No match (different service name)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "bar-svc",
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			services: []ks.Service{
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "None",
							Selector: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			expectedErr:     nil,
			expectedGrade:   scorecard.GradeCritical,
			expectedSkipped: false,
		},

		// No match (missing service name)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			services: []ks.Service{
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "None",
							Selector: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			expectedErr:     nil,
			expectedGrade:   scorecard.GradeCritical,
			expectedSkipped: false,
		},

		// Match (same namespace)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foo-ns"},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "foo-svc",
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			services: []ks.Service{
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc", Namespace: "foo-ns"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "None",
							Selector: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			expectedErr:     nil,
			expectedGrade:   scorecard.GradeAllOK,
			expectedSkipped: false,
		},

		// No match (different namespace)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foo-ns"},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "foo-svc",
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			services: []ks.Service{
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc", Namespace: "bar-ns"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "None",
							Selector: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			expectedErr:     nil,
			expectedGrade:   scorecard.GradeCritical,
			expectedSkipped: false,
		},

		// Match (multiple namespaces)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foo-ns"},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "foo-svc",
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			services: []ks.Service{
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc", Namespace: "bar-ns"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "None",
							Selector: map[string]string{
								"app": "foo",
							},
						},
					},
				},
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc", Namespace: "foo-ns"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "None",
							Selector: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			expectedErr:     nil,
			expectedGrade:   scorecard.GradeAllOK,
			expectedSkipped: false,
		},

		// Match (multiple namespaces, reversed)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foo-ns"},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "foo-svc",
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			services: []ks.Service{
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc", Namespace: "foo-ns"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "None",
							Selector: map[string]string{
								"app": "foo",
							},
						},
					},
				},
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc", Namespace: "bar-ns"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "None",
							Selector: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			expectedErr:     nil,
			expectedGrade:   scorecard.GradeAllOK,
			expectedSkipped: false,
		},

		// No match (not headless service)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "foo-svc",
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			services: []ks.Service{
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "",
							Selector: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			expectedErr:     nil,
			expectedGrade:   scorecard.GradeCritical,
			expectedSkipped: false,
		},

		// No match (selector)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.StatefulSetSpec{
					ServiceName: "foo-svc",
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			services: []ks.Service{
				service{
					corev1.Service{
						ObjectMeta: metav1.ObjectMeta{Name: "foo-svc"},
						Spec: corev1.ServiceSpec{
							ClusterIP: "None",
							Selector: map[string]string{
								"app": "bar",
							},
						},
					},
				},
			},
			expectedErr:     nil,
			expectedGrade:   scorecard.GradeCritical,
			expectedSkipped: false,
		},
	}

	for _, tc := range testcases {
		fn := statefulsetHasServiceName(tc.services)
		score, err := fn(tc.statefulset)
		assert.Equal(t, tc.expectedErr, err)
		assert.Equal(t, tc.expectedGrade, score.Grade)
		assert.Equal(t, tc.expectedSkipped, score.Skipped)
	}
}

func TestStatefulSetSelectorLabels(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		statefulset   appsv1.StatefulSet
		expectedErr   error
		expectedGrade scorecard.Grade
	}{
		// Match
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "foo",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			expectedErr:   nil,
			expectedGrade: scorecard.GradeAllOK,
		},

		// No match (labels differ)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "foo",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "bar",
							},
						},
					},
				},
			},
			expectedErr:   nil,
			expectedGrade: scorecard.GradeCritical,
		},

		// Match (expression)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "app",
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{"aaa", "bbb", "bar"},
							},
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "bar",
							},
						},
					},
				},
			},
			expectedErr:   nil,
			expectedGrade: scorecard.GradeAllOK,
		},

		// No match (expression)
		{
			statefulset: appsv1.StatefulSet{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "app",
								Operator: metav1.LabelSelectorOpNotIn,
								Values:   []string{"aaa", "bbb", "bar"},
							},
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "bar",
							},
						},
					},
				},
			},
			expectedErr:   nil,
			expectedGrade: scorecard.GradeCritical,
		},
	}

	for _, tc := range testcases {
		score, err := statefulSetSelectorLabelsMatching(tc.statefulset)
		assert.Equal(t, tc.expectedErr, err)
		assert.Equal(t, tc.expectedGrade, score.Grade)
	}
}

func TestDeploymentSelectorLabels(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		statefulset   appsv1.Deployment
		expectedErr   error
		expectedGrade scorecard.Grade
	}{
		// Match
		{
			statefulset: appsv1.Deployment{
				TypeMeta:   metav1.TypeMeta{Kind: "Deployment"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "foo",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "foo",
							},
						},
					},
				},
			},
			expectedErr:   nil,
			expectedGrade: scorecard.GradeAllOK,
		},

		// No match (labels differ)
		{
			statefulset: appsv1.Deployment{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "foo",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "bar",
							},
						},
					},
				},
			},
			expectedErr:   nil,
			expectedGrade: scorecard.GradeCritical,
		},

		// Match (expression)
		{
			statefulset: appsv1.Deployment{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "app",
								Operator: metav1.LabelSelectorOpIn,
								Values:   []string{"aaa", "bbb", "bar"},
							},
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "bar",
							},
						},
					},
				},
			},
			expectedErr:   nil,
			expectedGrade: scorecard.GradeAllOK,
		},

		// No match (expression)
		{
			statefulset: appsv1.Deployment{
				TypeMeta:   metav1.TypeMeta{Kind: "StatefulSet"},
				ObjectMeta: metav1.ObjectMeta{Name: "foo"},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "app",
								Operator: metav1.LabelSelectorOpNotIn,
								Values:   []string{"aaa", "bbb", "bar"},
							},
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "bar",
							},
						},
					},
				},
			},
			expectedErr:   nil,
			expectedGrade: scorecard.GradeCritical,
		},
	}

	for _, tc := range testcases {
		score, err := deploymentSelectorLabelsMatching(tc.statefulset)
		assert.Equal(t, tc.expectedErr, err)
		assert.Equal(t, tc.expectedGrade, score.Grade)
	}
}

type service struct {
	svc corev1.Service
}

func (d service) Service() corev1.Service {
	return d.svc
}

func (d service) FileLocation() ks.FileLocation {
	return ks.FileLocation{}
}
