package internal

import (
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Appsv1StatefulSet struct {
	appsv1.StatefulSet
}

func (s Appsv1StatefulSet) GetTypeMeta() metav1.TypeMeta {
	return s.TypeMeta
}

func (s Appsv1StatefulSet) GetObjectMeta() metav1.ObjectMeta {
	return s.ObjectMeta
}

func (s Appsv1StatefulSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return s.Spec.Template
}

type Appsv1beta1StatefulSet struct {
	appsv1beta1.StatefulSet
}

func (s Appsv1beta1StatefulSet) GetTypeMeta() metav1.TypeMeta {
	return s.TypeMeta
}

func (s Appsv1beta1StatefulSet) GetObjectMeta() metav1.ObjectMeta {
	return s.ObjectMeta
}

func (s Appsv1beta1StatefulSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return s.Spec.Template
}

type Appsv1beta2StatefulSet struct {
	appsv1beta2.StatefulSet
}

func (s Appsv1beta2StatefulSet) GetTypeMeta() metav1.TypeMeta {
	return s.TypeMeta
}

func (s Appsv1beta2StatefulSet) GetObjectMeta() metav1.ObjectMeta {
	return s.ObjectMeta
}

func (s Appsv1beta2StatefulSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	return s.Spec.Template
}
