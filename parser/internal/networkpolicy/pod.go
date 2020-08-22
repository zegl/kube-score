package pod

import (
	networkingv1 "k8s.io/api/networking/v1"

	ks "github.com/zegl/kube-score/domain"
)

type NetworkPolicy struct {
	Obj      networkingv1.NetworkPolicy
	Location ks.FileLocation
}

func (p NetworkPolicy) NetworkPolicy() networkingv1.NetworkPolicy {
	return p.Obj
}

func (p NetworkPolicy) FileLocation() ks.FileLocation {
	return p.Location
}
