package checks

import (
	ks "github.com/zegl/kube-score"
	"github.com/zegl/kube-score/scorecard"
	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func New() *Checks {
	return &Checks{
		metas:           make(map[string]MetaCheck),
		pods:            make(map[string]PodCheck),
		services:        make(map[string]ServiceCheck),
		statefulsets:    make(map[string]StatefulSetCheck),
		deployments:     make(map[string]DeploymentCheck),
		networkpolicies: make(map[string]NetworkPolicyCheck),
		ingresses:       make(map[string]IngressCheck),
		cronjobs:        make(map[string]CronJobCheck),
	}
}

func NewCheck(name string) ks.Check {
	return ks.Check{
		Name: name,
		ID:   machineFriendlyName(name),
	}
}

type MetaCheck struct {
	ks.Check
	Fn func(metav1.TypeMeta) scorecard.TestScore
}

type PodCheck struct {
	ks.Check
	Fn func(corev1.PodTemplateSpec) scorecard.TestScore
}

type ServiceCheck struct {
	ks.Check
	Fn func(corev1.Service) scorecard.TestScore
}

type StatefulSetCheck struct {
	ks.Check
	Fn func(appsv1.StatefulSet) scorecard.TestScore
}

type DeploymentCheck struct {
	ks.Check
	Fn func(appsv1.Deployment) scorecard.TestScore
}

type NetworkPolicyCheck struct {
	ks.Check
	Fn func(networkingv1.NetworkPolicy) scorecard.TestScore
}

type IngressCheck struct {
	ks.Check
	Fn func(extensionsv1beta1.Ingress) scorecard.TestScore
}

type CronJobCheck struct {
	ks.Check
	Fn func(batchv1beta1.CronJob) scorecard.TestScore
}

type Checks struct {
	metas           map[string]MetaCheck
	pods            map[string]PodCheck
	services        map[string]ServiceCheck
	statefulsets    map[string]StatefulSetCheck
	deployments     map[string]DeploymentCheck
	networkpolicies map[string]NetworkPolicyCheck
	ingresses       map[string]IngressCheck
	cronjobs        map[string]CronJobCheck
}

func machineFriendlyName(in string) string {
	in = strings.ToLower(in)
	in = strings.Replace(in, " ", "-", -1)
	return in
}

func (c *Checks) RegisterMetaCheck(name string, fn func(metav1.TypeMeta) scorecard.TestScore) {
	c.metas[machineFriendlyName(name)] = MetaCheck{NewCheck(name), fn}
}

func (c *Checks) Metas() map[string]MetaCheck {
	return c.metas
}

func (c *Checks) RegisterPodCheck(name string, fn func(corev1.PodTemplateSpec) scorecard.TestScore) {
	c.pods[machineFriendlyName(name)] = PodCheck{NewCheck(name), fn}
}

func (c *Checks) Pods() map[string]PodCheck {
	return c.pods
}

func (c *Checks) RegisterCronJobCheck(name string, fn func(batchv1beta1.CronJob) scorecard.TestScore) {
	c.cronjobs[machineFriendlyName(name)] = CronJobCheck{NewCheck(name), fn}
}

func (c *Checks) CronJobs() map[string]CronJobCheck {
	return c.cronjobs
}

func (c *Checks) RegisterStatefulSetCheck(name string, fn func(appsv1.StatefulSet) scorecard.TestScore) {
	c.statefulsets[machineFriendlyName(name)] = StatefulSetCheck{NewCheck(name), fn}
}

func (c *Checks) StatefulSets() map[string]StatefulSetCheck {
	return c.statefulsets
}

func (c *Checks) RegisterDeploymentCheck(name string, fn func(appsv1.Deployment) scorecard.TestScore) {
	c.deployments[machineFriendlyName(name)] = DeploymentCheck{NewCheck(name), fn}
}

func (c *Checks) Deployments() map[string]DeploymentCheck {
	return c.deployments
}

func (c *Checks) RegisterIngressCheck(name string, fn func(extensionsv1beta1.Ingress) scorecard.TestScore) {
	c.ingresses[machineFriendlyName(name)] = IngressCheck{NewCheck(name), fn}
}

func (c *Checks) Ingresses() map[string]IngressCheck {
	return c.ingresses
}

func (c *Checks) RegisterNetworkPolicyCheck(name string, fn func(networkingv1.NetworkPolicy) scorecard.TestScore) {
	c.networkpolicies[machineFriendlyName(name)] = NetworkPolicyCheck{NewCheck(name), fn}
}

func (c *Checks) NetworkPolicies() map[string]NetworkPolicyCheck {
	return c.networkpolicies
}

func (c *Checks) RegisterServiceCheck(name string, fn func(corev1.Service) scorecard.TestScore) {
	c.services[machineFriendlyName(name)] = ServiceCheck{NewCheck(name), fn}
}

func (c *Checks) Services() map[string]ServiceCheck {
	return c.services
}
