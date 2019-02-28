package internal

import (
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Appsv1DaemonSet struct {
	appsv1.DaemonSet
}

func (d Appsv1DaemonSet) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d Appsv1DaemonSet) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d Appsv1DaemonSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	d.Spec.Template.ObjectMeta.Namespace = d.ObjectMeta.Namespace
	return d.Spec.Template
}

type Appsv1beta2DaemonSet struct {
	appsv1beta2.DaemonSet
}

func (d Appsv1beta2DaemonSet) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d Appsv1beta2DaemonSet) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d Appsv1beta2DaemonSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	d.Spec.Template.ObjectMeta.Namespace = d.ObjectMeta.Namespace
	return d.Spec.Template
}

type Extensionsv1beta1DaemonSet struct {
	extensionsv1beta1.DaemonSet
}

func (d Extensionsv1beta1DaemonSet) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d Extensionsv1beta1DaemonSet) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d Extensionsv1beta1DaemonSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	d.Spec.Template.ObjectMeta.Namespace = d.ObjectMeta.Namespace
	return d.Spec.Template
}
