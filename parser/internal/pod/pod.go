package pod

import (
	corev1 "k8s.io/api/core/v1"

	ks "github.com/zegl/kube-score/domain"
)

type Pod struct {
	Obj      corev1.Pod
	Location ks.FileLocation
}

func (p Pod) Pod() corev1.Pod {
	return p.Obj
}

func (p Pod) FileLocation() ks.FileLocation {
	return p.Location
}
