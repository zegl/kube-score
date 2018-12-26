package kube_score

import (
	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Check struct {
	Name       string
	ID         string
	TargetType string
	Comment    string
}

type BothMeta struct {
	TypeMeta   metav1.TypeMeta
	ObjectMeta metav1.ObjectMeta
}

type PodSpecer interface {
	GetTypeMeta() metav1.TypeMeta
	GetObjectMeta() metav1.ObjectMeta
	GetPodTemplateSpec() corev1.PodTemplateSpec
}

type Metas interface {
	Metas() []BothMeta
}

type Pods interface {
	Pods() []corev1.Pod
}

type PodSpeccers interface {
	PodSpeccers() []PodSpecer
}

type Services interface {
	Services() []corev1.Service
}

type StatefulSets interface {
	StatefulSets() []appsv1.StatefulSet
}

type Deployments interface {
	Deployments() []appsv1.Deployment
}

type NetworkPolicies interface {
	NetworkPolicies() []networkingv1.NetworkPolicy
}

type Ingresses interface {
	Ingresses() []extensionsv1beta1.Ingress
}

type CronJobs interface {
	CronJobs() []batchv1beta1.CronJob
}

type PodDisruptionBudgets interface {
	PodDisruptionBudgets() []policyv1beta1.PodDisruptionBudget
}

type AllTypes interface {
	Metas
	Pods
	PodSpeccers
	Services
	StatefulSets
	Deployments
	NetworkPolicies
	Ingresses
	CronJobs
	PodDisruptionBudgets
}
