package pod

import (
	v1 "k8s.io/api/core/v1"

	ks "github.com/zegl/kube-score/domain"
)

type Service struct {
	Obj      v1.Service
	Location ks.FileLocation
}

func (p Service) Service() v1.Service {
	return p.Obj
}

func (p Service) FileLocation() ks.FileLocation {
	return p.Location
}
