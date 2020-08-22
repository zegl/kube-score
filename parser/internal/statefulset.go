package internal

import (
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
)

type Appsv1StatefulSet struct {
	Obj      appsv1.StatefulSet
	Location ks.FileLocation
}

func (d Appsv1StatefulSet) FileLocation() ks.FileLocation {
	return d.Location
}

func (s Appsv1StatefulSet) GetTypeMeta() metav1.TypeMeta {
	return s.Obj.TypeMeta
}

func (s Appsv1StatefulSet) GetObjectMeta() metav1.ObjectMeta {
	return s.Obj.ObjectMeta
}

func (s Appsv1StatefulSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	s.Obj.Spec.Template.ObjectMeta.Namespace = s.Obj.ObjectMeta.Namespace
	return s.Obj.Spec.Template
}

func (s Appsv1StatefulSet) StatefulSet() appsv1.StatefulSet {
	return s.Obj
}

type Appsv1beta1StatefulSet struct {
	appsv1beta1.StatefulSet
	Location ks.FileLocation
}

func (d Appsv1beta1StatefulSet) FileLocation() ks.FileLocation {
	return d.Location
}

func (s Appsv1beta1StatefulSet) GetTypeMeta() metav1.TypeMeta {
	return s.TypeMeta
}

func (s Appsv1beta1StatefulSet) GetObjectMeta() metav1.ObjectMeta {
	return s.ObjectMeta
}

func (s Appsv1beta1StatefulSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	s.Spec.Template.ObjectMeta.Namespace = s.ObjectMeta.Namespace
	return s.Spec.Template
}

type Appsv1beta2StatefulSet struct {
	appsv1beta2.StatefulSet
	Location ks.FileLocation
}

func (d Appsv1beta2StatefulSet) FileLocation() ks.FileLocation {
	return d.Location
}

func (s Appsv1beta2StatefulSet) GetTypeMeta() metav1.TypeMeta {
	return s.TypeMeta
}

func (s Appsv1beta2StatefulSet) GetObjectMeta() metav1.ObjectMeta {
	return s.ObjectMeta
}

func (s Appsv1beta2StatefulSet) GetPodTemplateSpec() corev1.PodTemplateSpec {
	s.Spec.Template.ObjectMeta.Namespace = s.ObjectMeta.Namespace
	return s.Spec.Template
}
