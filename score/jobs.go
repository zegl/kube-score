package score

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type batchv1Job struct {
	batchv1.Job
}

func (d batchv1Job) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d batchv1Job) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d batchv1Job) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.Template
}

type batchv1beta1CronJob struct {
	batchv1beta1.CronJob
}

func (d batchv1beta1CronJob) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d batchv1beta1CronJob) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d batchv1beta1CronJob) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.JobTemplate.Spec.Template
}
