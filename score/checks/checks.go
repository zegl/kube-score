package checks

import (
	"github.com/zegl/kube-score/scorecard"

	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func New() *Checks {
	return &Checks{
		metas:           make(map[string]func(metav1.TypeMeta) scorecard.TestScore),
		pods:            make(map[string]func(corev1.PodTemplateSpec) scorecard.TestScore),
		services:        make(map[string]func(corev1.Service) scorecard.TestScore),
		statefulsets:    make(map[string]func(appsv1.StatefulSet) scorecard.TestScore),
		deployments:     make(map[string]func(appsv1.Deployment) scorecard.TestScore),
		networkpolicies: make(map[string]func(networkingv1.NetworkPolicy) scorecard.TestScore),
		ingresses:       make(map[string]func(extensionsv1beta1.Ingress) scorecard.TestScore),
		cronjobs:        make(map[string]func(batchv1beta1.CronJob) scorecard.TestScore),
	}
}

type Checks struct {
	metas           map[string]func(metav1.TypeMeta) scorecard.TestScore
	pods            map[string]func(corev1.PodTemplateSpec) scorecard.TestScore
	services        map[string]func(corev1.Service) scorecard.TestScore
	statefulsets    map[string]func(appsv1.StatefulSet) scorecard.TestScore
	deployments     map[string]func(appsv1.Deployment) scorecard.TestScore
	networkpolicies map[string]func(networkingv1.NetworkPolicy) scorecard.TestScore
	ingresses       map[string]func(extensionsv1beta1.Ingress) scorecard.TestScore
	cronjobs        map[string]func(batchv1beta1.CronJob) scorecard.TestScore
}

func (c *Checks) RegisterPodCheck(name string, fn func(corev1.PodTemplateSpec) scorecard.TestScore) {
	c.pods[name] = fn
}

func (c *Checks) RegisterCronJobCheck(name string, fn func(batchv1beta1.CronJob) scorecard.TestScore) {
	c.cronjobs[name] = fn
}

func (c *Checks) RegisterStatefulSetCheck(name string, fn func(appsv1.StatefulSet) scorecard.TestScore) {
	c.statefulsets[name] = fn
}

func (c *Checks) RegisterDeploymentCheck(name string, fn func(appsv1.Deployment) scorecard.TestScore) {
	c.deployments[name] = fn
}

func (c *Checks) RegisterIngressCheck(name string, fn func(extensionsv1beta1.Ingress) scorecard.TestScore) {
	c.ingresses[name] = fn
}

func (c *Checks) RegisterNetworkPolicyCheck(name string, fn func(networkingv1.NetworkPolicy) scorecard.TestScore) {
	c.networkpolicies[name] = fn
}

func (c *Checks) RegisterServiceCheck(name string, fn func(corev1.Service) scorecard.TestScore) {
	c.services[name] = fn
}

func (c *Checks) RegisterMetaCheck(name string, fn func(metav1.TypeMeta) scorecard.TestScore) {
	c.metas[name] = fn
}

func (c *Checks) Metas() map[string]func(metav1.TypeMeta) scorecard.TestScore {
	return c.metas
}
func (c *Checks) Pods() map[string]func(corev1.PodTemplateSpec) scorecard.TestScore {
	return c.pods
}
func (c *Checks) Services() map[string]func(corev1.Service) scorecard.TestScore {
	return c.services
}
func (c *Checks) StatefulSets() map[string]func(appsv1.StatefulSet) scorecard.TestScore {
	return c.statefulsets
}
func (c *Checks) Deployments() map[string]func(appsv1.Deployment) scorecard.TestScore {
	return c.deployments
}
func (c *Checks) NetworkPolicies() map[string]func(networkingv1.NetworkPolicy) scorecard.TestScore {
	return c.networkpolicies
}
func (c *Checks) Ingresses() map[string]func(extensionsv1beta1.Ingress) scorecard.TestScore {
	return c.ingresses
}
func (c *Checks) CronJobs() map[string]func(batchv1beta1.CronJob) scorecard.TestScore {
	return c.cronjobs
}
