package pod

import (
	policyv1beta1 "k8s.io/api/policy/v1beta1"

	ks "github.com/zegl/kube-score/domain"
)

type PodDisruptionBudget struct {
	Obj      policyv1beta1.PodDisruptionBudget
	Location ks.FileLocation
}

func (p PodDisruptionBudget) PodDisruptionBudget() policyv1beta1.PodDisruptionBudget {
	return p.Obj
}

func (p PodDisruptionBudget) FileLocation() ks.FileLocation {
	return p.Location
}
