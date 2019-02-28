package internal

import (
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Appsv1Deployment struct {
	appsv1.Deployment
}

func (d Appsv1Deployment) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d Appsv1Deployment) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d Appsv1Deployment) GetPodTemplateSpec() corev1.PodTemplateSpec {
	d.Spec.Template.ObjectMeta.Namespace = d.ObjectMeta.Namespace
	return d.Spec.Template
}

type Appsv1beta1Deployment struct {
	appsv1beta1.Deployment
}

func (d Appsv1beta1Deployment) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d Appsv1beta1Deployment) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d Appsv1beta1Deployment) GetPodTemplateSpec() corev1.PodTemplateSpec {
	d.Spec.Template.ObjectMeta.Namespace = d.ObjectMeta.Namespace
	return d.Spec.Template
}

type Appsv1beta2Deployment struct {
	appsv1beta2.Deployment
}

func (d Appsv1beta2Deployment) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d Appsv1beta2Deployment) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d Appsv1beta2Deployment) GetPodTemplateSpec() corev1.PodTemplateSpec {
	d.Spec.Template.ObjectMeta.Namespace = d.ObjectMeta.Namespace
	return d.Spec.Template
}

type Extensionsv1beta1Deployment struct {
	extensionsv1beta1.Deployment
}

func (d Extensionsv1beta1Deployment) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d Extensionsv1beta1Deployment) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d Extensionsv1beta1Deployment) GetPodTemplateSpec() corev1.PodTemplateSpec {
	d.Spec.Template.ObjectMeta.Namespace = d.ObjectMeta.Namespace
	return d.Spec.Template
}
