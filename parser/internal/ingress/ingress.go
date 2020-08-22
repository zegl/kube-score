package ingress

import (
	"k8s.io/api/extensions/v1beta1"

	ks "github.com/zegl/kube-score/domain"
)

type Ingress struct {
	Obj      v1beta1.Ingress
	Location ks.FileLocation
}

func (i Ingress) Ingress() v1beta1.Ingress {
	return i.Obj
}

func (i Ingress) FileLocation() ks.FileLocation {
	return i.Location
}
