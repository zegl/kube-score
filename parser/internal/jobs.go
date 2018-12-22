package internal

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Batchv1Job struct {
	batchv1.Job
}

func (d Batchv1Job) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d Batchv1Job) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d Batchv1Job) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.Template
}

type Batchv1beta1CronJob struct {
	batchv1beta1.CronJob
}

func (d Batchv1beta1CronJob) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d Batchv1beta1CronJob) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d Batchv1beta1CronJob) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.JobTemplate.Spec.Template
}
