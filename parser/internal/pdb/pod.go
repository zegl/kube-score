package pod

import (
	policyv1 "k8s.io/api/policy/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
)

type PodDisruptionBudgetV1beta1 struct {
	Obj      policyv1beta1.PodDisruptionBudget
	Location ks.FileLocation
}

func (p PodDisruptionBudgetV1beta1) PodDisruptionBudgetSelector() *metav1.LabelSelector {
	return p.Obj.Spec.Selector
}

func (p PodDisruptionBudgetV1beta1) Namespace() string {
	return p.Obj.Namespace
}

func (p PodDisruptionBudgetV1beta1) FileLocation() ks.FileLocation {
	return p.Location
}

type PodDisruptionBudgetV1 struct {
	Obj      policyv1.PodDisruptionBudget
	Location ks.FileLocation
}

func (p PodDisruptionBudgetV1) PodDisruptionBudgetSelector() *metav1.LabelSelector {
	return p.Obj.Spec.Selector
}

func (p PodDisruptionBudgetV1) FileLocation() ks.FileLocation {
	return p.Location
}

func (p PodDisruptionBudgetV1) Namespace() string {
	return p.Obj.Namespace
}
