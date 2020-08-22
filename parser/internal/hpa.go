package internal

import (
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
)

type HPAv1 struct {
	autoscalingv1.HorizontalPodAutoscaler
	Location ks.FileLocation
}

func (d HPAv1) FileLocation() ks.FileLocation {
	return d.Location
}

func (d HPAv1) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d HPAv1) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d HPAv1) HpaTarget() autoscalingv1.CrossVersionObjectReference {
	return d.Spec.ScaleTargetRef
}

type HPAv2beta1 struct {
	autoscalingv2beta1.HorizontalPodAutoscaler
	Location ks.FileLocation
}

func (d HPAv2beta1) FileLocation() ks.FileLocation {
	return d.Location
}

func (d HPAv2beta1) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d HPAv2beta1) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d HPAv2beta1) HpaTarget() autoscalingv1.CrossVersionObjectReference {
	return autoscalingv1.CrossVersionObjectReference(d.Spec.ScaleTargetRef)
}

type HPAv2beta2 struct {
	autoscalingv2beta2.HorizontalPodAutoscaler
	Location ks.FileLocation
}

func (d HPAv2beta2) FileLocation() ks.FileLocation {
	return d.Location
}

func (d HPAv2beta2) GetTypeMeta() metav1.TypeMeta {
	return d.TypeMeta
}

func (d HPAv2beta2) GetObjectMeta() metav1.ObjectMeta {
	return d.ObjectMeta
}

func (d HPAv2beta2) HpaTarget() autoscalingv1.CrossVersionObjectReference {
	return autoscalingv1.CrossVersionObjectReference(d.Spec.ScaleTargetRef)
}
