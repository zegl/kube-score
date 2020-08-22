package internal

import (
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
)

type Appsv1Deployment struct {
	Obj      appsv1.Deployment
	Location ks.FileLocation
}

func (d Appsv1Deployment) FileLocation() ks.FileLocation {
	return d.Location
}

func (d Appsv1Deployment) GetTypeMeta() metav1.TypeMeta {
	return d.Obj.TypeMeta
}

func (d Appsv1Deployment) GetObjectMeta() metav1.ObjectMeta {
	return d.Obj.ObjectMeta
}

func (d Appsv1Deployment) GetPodTemplateSpec() corev1.PodTemplateSpec {
	d.Obj.Spec.Template.ObjectMeta.Namespace = d.Obj.ObjectMeta.Namespace
	return d.Obj.Spec.Template
}

func (d Appsv1Deployment) Deployment() appsv1.Deployment {
	return d.Obj
}

type Appsv1beta1Deployment struct {
	appsv1beta1.Deployment
	Location ks.FileLocation
}

func (d Appsv1beta1Deployment) FileLocation() ks.FileLocation {
	return d.Location
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
	Location ks.FileLocation
}

func (d Appsv1beta2Deployment) FileLocation() ks.FileLocation {
	return d.Location
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
	Location ks.FileLocation
}

func (d Extensionsv1beta1Deployment) FileLocation() ks.FileLocation {
	return d.Location
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
