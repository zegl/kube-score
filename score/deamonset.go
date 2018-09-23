package score

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type appsv1DaemonSet struct {
	appsv1.DaemonSet
}

func (d appsv1DaemonSet) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d appsv1DaemonSet) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d appsv1DaemonSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.Template
}

type appsv1beta2DaemonSet struct {
	appsv1beta2.DaemonSet
}

func (d appsv1beta2DaemonSet) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d appsv1beta2DaemonSet) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d appsv1beta2DaemonSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.Template
}

type extensionsv1beta1DaemonSet struct {
	extensionsv1beta1.DaemonSet
}

func (d extensionsv1beta1DaemonSet) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d extensionsv1beta1DaemonSet) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d extensionsv1beta1DaemonSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return d.Spec.Template
}

