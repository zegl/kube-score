package score

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
)

type Deployment interface {
	GetTypeMeta() metav1.TypeMeta
	GetObjectMeta() metav1.ObjectMeta
	GetPodTemplateSpec() corev1.PodTemplateSpec
}

type appsv1Deployment struct {
	appsv1.Deployment
}

func (d appsv1Deployment) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d appsv1Deployment) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d appsv1Deployment) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.Template
}

type appsv1beta1Deployment struct {
	appsv1beta1.Deployment
}

func (d appsv1beta1Deployment) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d appsv1beta1Deployment) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d appsv1beta1Deployment) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.Template
}

type appsv1beta2Deployment struct {
	appsv1beta2.Deployment
}

func (d appsv1beta2Deployment) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d appsv1beta2Deployment) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d appsv1beta2Deployment) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.Template
}

type extensionsv1beta1Deployment struct {
	extensionsv1beta1.Deployment
}

func (d extensionsv1beta1Deployment) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d extensionsv1beta1Deployment) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d extensionsv1beta1Deployment) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.Template
}