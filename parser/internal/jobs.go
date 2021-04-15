package internal

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
)

type Batchv1Job struct {
	batchv1.Job
	Location ks.FileLocation
}

func (d Batchv1Job) FileLocation() ks.FileLocation {
	return d.Location
}

func (d Batchv1Job) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d Batchv1Job) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d Batchv1Job) GetPodTemplateSpec() corev1.PodTemplateSpec {
	d.Spec.Template.ObjectMeta.Namespace = d.ObjectMeta.Namespace
	return d.Spec.Template
}
