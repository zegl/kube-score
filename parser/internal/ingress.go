package internal

import (
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
)

var _ ks.Ingress = (*IngressV1)(nil)
var _ ks.Ingress = (*IngressV1beta1)(nil)
var _ ks.Ingress = (*ExtensionsIngressV1beta1)(nil)

type IngressV1 struct {
	networkingv1.Ingress
}

func (i IngressV1) GetObjectMeta() v1.ObjectMeta {
	return i.ObjectMeta
}

func (i IngressV1) GetTypeMeta() v1.TypeMeta {
	return i.TypeMeta
}

func (i IngressV1) Rules() []networkingv1.IngressRule {
	return i.Spec.Rules
}

type IngressV1beta1 struct {
	networkingv1beta1.Ingress
}

func (i IngressV1beta1) GetObjectMeta() v1.ObjectMeta {
	return i.ObjectMeta
}

func (i IngressV1beta1) GetTypeMeta() v1.TypeMeta {
	return i.TypeMeta
}

func (i IngressV1beta1) Rules() []networkingv1.IngressRule {
	var res []networkingv1.IngressRule

	paths := func(in []networkingv1beta1.HTTPIngressPath) (out []networkingv1.HTTPIngressPath) {
		for _, path := range in {
			out = append(out, networkingv1.HTTPIngressPath{
				Path: path.Path,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: path.Backend.ServiceName,
						Port: networkingv1.ServiceBackendPort{
							Name:   path.Backend.ServicePort.StrVal,
							Number: path.Backend.ServicePort.IntVal,
						},
					},
				},
			})
		}
		return
	}

	for _, rule := range i.Spec.Rules {
		res = append(res, networkingv1.IngressRule{
			Host: rule.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: paths(rule.HTTP.Paths),
				},
			},
		})
	}

	return res
}

type ExtensionsIngressV1beta1 struct {
	extensionsv1beta1.Ingress
}

func (i ExtensionsIngressV1beta1) GetObjectMeta() v1.ObjectMeta {
	return i.ObjectMeta
}

func (i ExtensionsIngressV1beta1) GetTypeMeta() v1.TypeMeta {
	return i.TypeMeta
}

func (i ExtensionsIngressV1beta1) Rules() []networkingv1.IngressRule {
	var res []networkingv1.IngressRule

	paths := func(in []extensionsv1beta1.HTTPIngressPath) (out []networkingv1.HTTPIngressPath) {
		for _, path := range in {
			out = append(out, networkingv1.HTTPIngressPath{
				Path: path.Path,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: path.Backend.ServiceName,
						Port: networkingv1.ServiceBackendPort{
							Name:   path.Backend.ServicePort.StrVal,
							Number: path.Backend.ServicePort.IntVal,
						},
					},
				},
			})
		}
		return
	}

	for _, rule := range i.Spec.Rules {
		res = append(res, networkingv1.IngressRule{
			Host: rule.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: paths(rule.HTTP.Paths),
				},
			},
		})
	}

	return res
}
