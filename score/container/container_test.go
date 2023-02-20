package container

import (
	"testing"

	ks "github.com/zegl/kube-score/domain"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zegl/kube-score/scorecard"
)

type podSpeccer struct {
	typeMeta   metav1.TypeMeta
	objectMeta metav1.ObjectMeta
	spec       corev1.PodTemplateSpec
}

func (p *podSpeccer) GetTypeMeta() metav1.TypeMeta {
	return p.typeMeta
}

func (p *podSpeccer) GetObjectMeta() metav1.ObjectMeta {
	return p.objectMeta
}

func (p *podSpeccer) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return p.spec
}

func (p *podSpeccer) FileLocation() ks.FileLocation {
	return ks.FileLocation{}
}

func TestOkAllTheSameContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerResourceRequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestOkMultipleContainersContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerResourceRequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
						{
							Name: "foo2",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestOkSameQuantityContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerResourceRequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1000m"),
									"memory": resource.MustParse("0.25Gi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestFailBothContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerResourceRequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("2"),
									"memory": resource.MustParse("512Mi"),
								},
							},
						},
					},
				},
			},
		})

	assert.Equal(t, scorecard.GradeCritical, s.Grade)
	assert.Len(t, s.Comments, 2)
	assert.Equal(t, "foo", s.Comments[0].Path)
	assert.Equal(t, "CPU requests does not match limits", s.Comments[0].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.cpu == resources.limits.cpu", s.Comments[0].Description)
	assert.Equal(t, "foo", s.Comments[1].Path)
	assert.Equal(t, "Memory requests does not match limits", s.Comments[1].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.memory == resources.limits.memory", s.Comments[1].Description)
}

func TestFailCpuInitContainerResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerResourceRequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name: "init",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("2"),
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeCritical, s.Grade)
	assert.Len(t, s.Comments, 1)
	assert.Equal(t, "init", s.Comments[0].Path)
	assert.Equal(t, "CPU requests does not match limits", s.Comments[0].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.cpu == resources.limits.cpu", s.Comments[0].Description)
}

func TestOkAllCPURequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerCPURequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestOkMultipleContainersContainerCPURequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerCPURequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1"),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1"),
								},
							},
						},
						{
							Name: "foo2",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestOkSameQuantityContainerCPURequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerCPURequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1000m"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestFailContainerCPURequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerCPURequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("2"),
									"memory": resource.MustParse("512Mi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeCritical, s.Grade)
	assert.Len(t, s.Comments, 1)
	assert.Equal(t, "foo", s.Comments[0].Path)
	assert.Equal(t, "CPU requests does not match limits", s.Comments[0].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.cpu == resources.limits.cpu", s.Comments[0].Description)

}

func TestFailInitContainerCPURequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerCPURequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name: "init",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("1"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu": resource.MustParse("2"),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeCritical, s.Grade)
	assert.Len(t, s.Comments, 1)
	assert.Equal(t, "init", s.Comments[0].Path)
	assert.Equal(t, "CPU requests does not match limits", s.Comments[0].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.cpu == resources.limits.cpu", s.Comments[0].Description)
}

func TestOkContainerMemoryResourceRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerMemoryRequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestOkMultipleContainersContainerMemoryRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerMemoryRequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
						{
							Name: "foo2",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestOkSameQuantityContainerMemoryRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerMemoryRequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("0.25Gi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
	assert.Len(t, s.Comments, 0)
}

func TestFailContainerMemoryRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerMemoryRequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("2"),
									"memory": resource.MustParse("512Mi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeCritical, s.Grade)
	assert.Len(t, s.Comments, 1)
	assert.Equal(t, "foo", s.Comments[0].Path)
	assert.Equal(t, "Memory requests does not match limits", s.Comments[0].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.memory == resources.limits.memory", s.Comments[0].Description)
}

func TestFailInitContainerMemoryRequestsEqualLimits(t *testing.T) {
	t.Parallel()
	s, _ := containerMemoryRequestsEqualLimits(
		&podSpeccer{
			spec: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name: "init",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"memory": resource.MustParse("512Mi"),
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name: "foo",
							Resources: corev1.ResourceRequirements{
								Requests: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
								Limits: map[corev1.ResourceName]resource.Quantity{
									"cpu":    resource.MustParse("1"),
									"memory": resource.MustParse("256Mi"),
								},
							},
						},
					},
				},
			},
		},
	)

	assert.Equal(t, scorecard.GradeCritical, s.Grade)
	assert.Len(t, s.Comments, 1)
	assert.Equal(t, "init", s.Comments[0].Path)
	assert.Equal(t, "Memory requests does not match limits", s.Comments[0].Summary)
	assert.Equal(t, "Having equal requests and limits is recommended to avoid resource DDOS of the node during spikes. Set resources.requests.memory == resources.limits.memory", s.Comments[0].Description)
}
