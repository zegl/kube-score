package checks

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func New(cnf config.Configuration) *Checks {
	return &Checks{
		cnf: cnf,

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

func NewCheck(name, targetType, comment string, optional bool) ks.Check {
	return ks.Check{
		Name:       name,
		ID:         machineFriendlyName(name),
		TargetType: targetType,
		Comment:    comment,
		Optional:   optional,
	}
}

func machineFriendlyName(in string) string {
	in = strings.ToLower(in)
	in = strings.Replace(in, " ", "-", -1)
	return in
}

type MetaCheckFn = func(ks.BothMeta) scorecard.TestScore
type MetaCheck struct {
	ks.Check
	Fn MetaCheckFn
}

type PodCheckFn = func(corev1.PodTemplateSpec, metav1.TypeMeta) scorecard.TestScore
type PodCheck struct {
	ks.Check
	Fn PodCheckFn
}

type ServiceCheckFn = func(corev1.Service) scorecard.TestScore
type ServiceCheck struct {
	ks.Check
	Fn ServiceCheckFn
}

type StatefulSetCheckFn = func(appsv1.StatefulSet) (scorecard.TestScore, error)
type StatefulSetCheck struct {
	ks.Check
	Fn StatefulSetCheckFn
}

type DeploymentCheckFn = func(appsv1.Deployment) (scorecard.TestScore, error)
type DeploymentCheck struct {
	ks.Check
	Fn DeploymentCheckFn
}

type NetworkPolicyCheckFn = func(networkingv1.NetworkPolicy) scorecard.TestScore
type NetworkPolicyCheck struct {
	ks.Check
	Fn NetworkPolicyCheckFn
}

type IngressCheckFn = func(extensionsv1beta1.Ingress) scorecard.TestScore
type IngressCheck struct {
	ks.Check
	Fn IngressCheckFn
}

type CronJobCheckFn = func(batchv1beta1.CronJob) scorecard.TestScore
type CronJobCheck struct {
	ks.Check
	Fn CronJobCheckFn
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

	cnf config.Configuration
}

func (c Checks) isIgnored(id string) bool {
	_, ok := c.cnf.IgnoredTests[id]
	return ok
}

func (c Checks) isEnabled(check ks.Check) bool {
	if c.isIgnored(check.ID) {
		return false
	}

	if !check.Optional {
		return true
	}

	_, ok := c.cnf.EnabledOptionalTests[check.ID]
	return ok
}

func (c *Checks) RegisterMetaCheck(name, comment string, fn MetaCheckFn) {
	ch := NewCheck(name, "all", comment, false)
	c.registerMetaCheck(MetaCheck{ch, fn})
}

func (c *Checks) RegisterOptionalMetaCheck(name, comment string, fn MetaCheckFn) {
	ch := NewCheck(name, "all", comment, true)
	c.registerMetaCheck(MetaCheck{ch, fn})
}

func (c *Checks) registerMetaCheck(ch MetaCheck) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.metas[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) Metas() map[string]MetaCheck {
	return c.metas
}

func (c *Checks) RegisterPodCheck(name, comment string, fn PodCheckFn) {
	ch := NewCheck(name, "Pod", comment, false)
	c.registerPodCheck(PodCheck{ch, fn})
}

func (c *Checks) RegisterOptionalPodCheck(name, comment string, fn PodCheckFn) {
	ch := NewCheck(name, "Pod", comment, true)
	c.registerPodCheck(PodCheck{ch, fn})
}

func (c *Checks) registerPodCheck(ch PodCheck) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.pods[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) Pods() map[string]PodCheck {
	return c.pods
}

func (c *Checks) RegisterCronJobCheck(name, comment string, fn CronJobCheckFn) {
	ch := NewCheck(name, "CronJob", comment, false)
	c.registerCronJobCheck(CronJobCheck{ch, fn})
}

func (c *Checks) RegisterOptionalCronJobCheck(name, comment string, fn CronJobCheckFn) {
	ch := NewCheck(name, "CronJob", comment, true)
	c.registerCronJobCheck(CronJobCheck{ch, fn})
}

func (c *Checks) registerCronJobCheck(ch CronJobCheck) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.cronjobs[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) CronJobs() map[string]CronJobCheck {
	return c.cronjobs
}

func (c *Checks) RegisterStatefulSetCheck(name, comment string, fn StatefulSetCheckFn) {
	ch := NewCheck(name, "StatefulSet", comment, false)
	c.registerStatefulSetCheck(StatefulSetCheck{ch, fn})
}

func (c *Checks) RegisterOptionalStatefulSetCheck(name, comment string, fn StatefulSetCheckFn) {
	ch := NewCheck(name, "StatefulSet", comment, true)
	c.registerStatefulSetCheck(StatefulSetCheck{ch, fn})
}

func (c *Checks) registerStatefulSetCheck(ch StatefulSetCheck) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.statefulsets[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) StatefulSets() map[string]StatefulSetCheck {
	return c.statefulsets
}

func (c *Checks) RegisterDeploymentCheck(name, comment string, fn DeploymentCheckFn) {
	ch := NewCheck(name, "Deployment", comment, false)
	c.registerDeploymentCheck(DeploymentCheck{ch, fn})
}

func (c *Checks) RegisterOptionalDeploymentCheck(name, comment string, fn DeploymentCheckFn) {
	ch := NewCheck(name, "Deployment", comment, true)
	c.registerDeploymentCheck(DeploymentCheck{ch, fn})
}

func (c *Checks) registerDeploymentCheck(ch DeploymentCheck) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.deployments[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) Deployments() map[string]DeploymentCheck {
	return c.deployments
}

func (c *Checks) RegisterIngressCheck(name, comment string, fn IngressCheckFn) {
	ch := NewCheck(name, "Ingress", comment, false)
	c.registerIngressCheck(IngressCheck{ch, fn})
}

func (c *Checks) RegisterOptionalIngressCheck(name, comment string, fn IngressCheckFn) {
	ch := NewCheck(name, "Ingress", comment, true)
	c.registerIngressCheck(IngressCheck{ch, fn})
}

func (c *Checks) registerIngressCheck(ch IngressCheck) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.ingresses[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) Ingresses() map[string]IngressCheck {
	return c.ingresses
}

func (c *Checks) RegisterNetworkPolicyCheck(name, comment string, fn NetworkPolicyCheckFn) {
	ch := NewCheck(name, "NetworkPolicy", comment, false)
	c.registerNetworkPolicyCheck(NetworkPolicyCheck{ch, fn})
}

func (c *Checks) RegisterOptionalNetworkPolicyCheck(name, comment string, fn NetworkPolicyCheckFn) {
	ch := NewCheck(name, "NetworkPolicy", comment, true)
	c.registerNetworkPolicyCheck(NetworkPolicyCheck{ch, fn})
}

func (c *Checks) registerNetworkPolicyCheck(ch NetworkPolicyCheck) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.networkpolicies[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) NetworkPolicies() map[string]NetworkPolicyCheck {
	return c.networkpolicies
}

func (c *Checks) RegisterServiceCheck(name, comment string, fn ServiceCheckFn) {
	ch := NewCheck(name, "Service", comment, false)
	c.registerServiceCheck(ServiceCheck{ch, fn})
}

func (c *Checks) RegisterOptionalServiceCheck(name, comment string, fn ServiceCheckFn) {
	ch := NewCheck(name, "Service", comment, true)
	c.registerServiceCheck(ServiceCheck{ch, fn})
}

func (c *Checks) registerServiceCheck(ch ServiceCheck) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.services[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) Services() map[string]ServiceCheck {
	return c.services
}

func (c *Checks) All() []ks.Check {
	return c.all
}
