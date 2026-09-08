package internal

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
)

// GenericObject represents a Kubernetes-style object whose kind kube-score
// does not have a dedicated parser for (most commonly a CRD, for example
// Argo Rollouts' "Rollout" or Knative's "Service"). Only TypeMeta and
// ObjectMeta are known; there is no PodSpec or other typed data available.
//
// This lets checks that only need to know "does an object with this
// apiVersion/kind/name/namespace exist" (for example the HPA
// "scaleTargetRef" check, see domain.BothMeta and score/hpa) work correctly
// even when the target is a CRD, instead of silently dropping it and
// producing a false-positive "target does not match anything" result.
type GenericObject struct {
	TypeMetaVal   metav1.TypeMeta
	ObjectMetaVal metav1.ObjectMeta
	Location      ks.FileLocation
}

func (g GenericObject) FileLocation() ks.FileLocation {
	return g.Location
}

func (g GenericObject) GetTypeMeta() metav1.TypeMeta {
	return g.TypeMetaVal
}

func (g GenericObject) GetObjectMeta() metav1.ObjectMeta {
	return g.ObjectMetaVal
}
