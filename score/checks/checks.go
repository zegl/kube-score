package checks

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func New() *Checks {
	return &Checks{
		all:             make([]ks.Check, 0),
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

func NewCheck(name, targetType, comment string) ks.Check {
	return ks.Check{
		Name:       name,
		ID:         machineFriendlyName(name),
		TargetType: targetType,
		Comment:    comment,
	}
}

func machineFriendlyName(in string) string {
	in = strings.ToLower(in)
	in = strings.Replace(in, " ", "-", -1)
	return in
}

type MetaCheck struct {
	ks.Check
	Fn func(ks.BothMeta) scorecard.TestScore
}

type PodCheck struct {
	ks.Check
	Fn func(corev1.PodTemplateSpec, metav1.TypeMeta) scorecard.TestScore
}

type ServiceCheck struct {
	ks.Check
	Fn func(corev1.Service) scorecard.TestScore
}

type StatefulSetCheck struct {
	ks.Check
	Fn func(appsv1.StatefulSet) (scorecard.TestScore, error)
}

type DeploymentCheck struct {
	ks.Check
	Fn func(appsv1.Deployment) (scorecard.TestScore, error)
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
	all             []ks.Check
	metas           map[string]MetaCheck
	pods            map[string]PodCheck
	services        map[string]ServiceCheck
	statefulsets    map[string]StatefulSetCheck
	deployments     map[string]DeploymentCheck
	networkpolicies map[string]NetworkPolicyCheck
	ingresses       map[string]IngressCheck
	cronjobs        map[string]CronJobCheck
}

func (c *Checks) RegisterMetaCheck(name, comment string, fn func(meta ks.BothMeta) scorecard.TestScore) {
	ch := NewCheck(name, "all", comment)
	c.all = append(c.all, ch)
	c.metas[machineFriendlyName(name)] = MetaCheck{ch, fn}
}

func (c *Checks) Metas() map[string]MetaCheck {
	return c.metas
}

func (c *Checks) RegisterPodCheck(name, comment string, fn func(corev1.PodTemplateSpec, metav1.TypeMeta) scorecard.TestScore) {
	ch := NewCheck(name, "Pod", comment)
	c.all = append(c.all, ch)
	c.pods[machineFriendlyName(name)] = PodCheck{ch, fn}
}

func (c *Checks) Pods() map[string]PodCheck {
	return c.pods
}

func (c *Checks) RegisterCronJobCheck(name, comment string, fn func(batchv1beta1.CronJob) scorecard.TestScore) {
	ch := NewCheck(name, "CronJob", comment)
	c.all = append(c.all, ch)
	c.cronjobs[machineFriendlyName(name)] = CronJobCheck{ch, fn}
}

func (c *Checks) CronJobs() map[string]CronJobCheck {
	return c.cronjobs
}

func (c *Checks) RegisterStatefulSetCheck(name, comment string, fn func(appsv1.StatefulSet) (scorecard.TestScore, error)) {
	ch := NewCheck(name, "StatefulSet", comment)
	c.all = append(c.all, ch)
	c.statefulsets[machineFriendlyName(name)] = StatefulSetCheck{ch, fn}
}

func (c *Checks) StatefulSets() map[string]StatefulSetCheck {
	return c.statefulsets
}

func (c *Checks) RegisterDeploymentCheck(name, comment string, fn func(appsv1.Deployment) (scorecard.TestScore, error)) {
	ch := NewCheck(name, "Deployment", comment)
	c.all = append(c.all, ch)
	c.deployments[machineFriendlyName(name)] = DeploymentCheck{ch, fn}
}

func (c *Checks) Deployments() map[string]DeploymentCheck {
	return c.deployments
}

func (c *Checks) RegisterIngressCheck(name, comment string, fn func(extensionsv1beta1.Ingress) scorecard.TestScore) {
	ch := NewCheck(name, "Ingress", comment)
	c.all = append(c.all, ch)
	c.ingresses[machineFriendlyName(name)] = IngressCheck{ch, fn}
}

func (c *Checks) Ingresses() map[string]IngressCheck {
	return c.ingresses
}

func (c *Checks) RegisterNetworkPolicyCheck(name, comment string, fn func(networkingv1.NetworkPolicy) scorecard.TestScore) {
	ch := NewCheck(name, "NetworkPolicy", comment)
	c.all = append(c.all, ch)
	c.networkpolicies[machineFriendlyName(name)] = NetworkPolicyCheck{ch, fn}
}

func (c *Checks) NetworkPolicies() map[string]NetworkPolicyCheck {
	return c.networkpolicies
}

func (c *Checks) RegisterServiceCheck(name, comment string, fn func(corev1.Service) scorecard.TestScore) {
	ch := NewCheck(name, "Service", comment)
	c.all = append(c.all, ch)
	c.services[machineFriendlyName(name)] = ServiceCheck{ch, fn}
}

func (c *Checks) Services() map[string]ServiceCheck {
	return c.services
}

func (c *Checks) All() []ks.Check {
	return c.all
}
