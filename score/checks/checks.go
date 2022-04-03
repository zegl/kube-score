package checks

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func New(cnf config.Configuration) *Checks {
	return &Checks{
		cnf: cnf,

		all:                      make([]ks.Check, 0),
		metas:                    make(map[string]GenCheck[ks.BothMeta]),
		pods:                     make(map[string]PodCheck),
		services:                 make(map[string]GenCheck[corev1.Service]),
		statefulsets:             make(map[string]GenCheck[appsv1.StatefulSet]),
		deployments:              make(map[string]GenCheck[appsv1.Deployment]),
		networkpolicies:          make(map[string]GenCheck[networkingv1.NetworkPolicy]),
		ingresses:                make(map[string]GenCheck[ks.Ingress]),
		cronjobs:                 make(map[string]GenCheck[ks.CronJob]),
		horizontalPodAutoscalers: make(map[string]GenCheck[ks.HpaTargeter]),
		poddisruptionbudgets:     make(map[string]GenCheck[ks.PodDisruptionBudget]),
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
	in = strings.ReplaceAll(in, " ", "-")
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

type CheckFunc[T any] func(T) (scorecard.TestScore, error)

type GenCheck[T any] struct {
	ks.Check
	Fn CheckFunc[T]
}

type Checks struct {
	all                      []ks.Check
	metas                    map[string]GenCheck[ks.BothMeta]
	pods                     map[string]PodCheck
	services                 map[string]GenCheck[corev1.Service]
	statefulsets             map[string]GenCheck[appsv1.StatefulSet]
	deployments              map[string]GenCheck[appsv1.Deployment]
	networkpolicies          map[string]GenCheck[networkingv1.NetworkPolicy]
	ingresses                map[string]GenCheck[ks.Ingress]
	cronjobs                 map[string]GenCheck[ks.CronJob]
	horizontalPodAutoscalers map[string]GenCheck[ks.HpaTargeter]
	poddisruptionbudgets     map[string]GenCheck[ks.PodDisruptionBudget]

	cnf config.Configuration
}

func (c Checks) isEnabled(check ks.Check) bool {
	_, ok := c.cnf.IgnoredTests[check.ID]
	return !ok
}

func (c *Checks) RegisterMetaCheck(name, comment string, fn CheckFunc[ks.BothMeta]) {
	ch := NewCheck(name, "all", comment, false)
	c.registerMetaCheck(GenCheck[ks.BothMeta]{ch, fn})
}

func (c *Checks) RegisterOptionalMetaCheck(name, comment string, fn CheckFunc[ks.BothMeta]) {
	ch := NewCheck(name, "all", comment, true)
	c.registerMetaCheck(GenCheck[ks.BothMeta]{ch, fn})
}

func (c *Checks) registerMetaCheck(ch GenCheck[ks.BothMeta]) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.metas[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) Metas() map[string]GenCheck[ks.BothMeta] {
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

func (c *Checks) RegisterHorizontalPodAutoscalerCheck(name, comment string, fn CheckFunc[ks.HpaTargeter]) {
	ch := NewCheck(name, "HorizontalPodAutoscaler", comment, false)
	c.registerHorizontalPodAutoscalerCheck(GenCheck[ks.HpaTargeter]{ch, fn})
}

func (c *Checks) RegisterOptionalHorizontalPodAutoscalerCheck(name, comment string, fn CheckFunc[ks.HpaTargeter]) {
	ch := NewCheck(name, "HorizontalPodAutoscaler", comment, true)
	c.registerHorizontalPodAutoscalerCheck(GenCheck[ks.HpaTargeter]{ch, fn})
}

func (c *Checks) registerHorizontalPodAutoscalerCheck(ch GenCheck[ks.HpaTargeter]) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.horizontalPodAutoscalers[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) HorizontalPodAutoscalers() map[string]GenCheck[ks.HpaTargeter] {
	return c.horizontalPodAutoscalers
}

func (c *Checks) RegisterCronJobCheck(name, comment string, fn CheckFunc[ks.CronJob]) {
	ch := NewCheck(name, "CronJob", comment, false)
	c.registerCronJobCheck(GenCheck[ks.CronJob]{ch, fn})
}

func (c *Checks) RegisterOptionalCronJobCheck(name, comment string, fn CheckFunc[ks.CronJob]) {
	ch := NewCheck(name, "CronJob", comment, true)
	c.registerCronJobCheck(GenCheck[ks.CronJob]{ch, fn})
}

func (c *Checks) registerCronJobCheck(ch GenCheck[ks.CronJob]) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.cronjobs[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) CronJobs() map[string]GenCheck[ks.CronJob] {
	return c.cronjobs
}

func (c *Checks) RegisterStatefulSetCheck(name, comment string, fn CheckFunc[appsv1.StatefulSet]) {
	ch := NewCheck(name, "StatefulSet", comment, false)
	c.registerStatefulSetCheck(GenCheck[appsv1.StatefulSet]{ch, fn})
}

func (c *Checks) RegisterOptionalStatefulSetCheck(name, comment string, fn CheckFunc[appsv1.StatefulSet]) {
	ch := NewCheck(name, "StatefulSet", comment, true)
	c.registerStatefulSetCheck(GenCheck[appsv1.StatefulSet]{ch, fn})
}

func (c *Checks) registerStatefulSetCheck(ch GenCheck[appsv1.StatefulSet]) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.statefulsets[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) StatefulSets() map[string]GenCheck[appsv1.StatefulSet] {
	return c.statefulsets
}

func (c *Checks) RegisterDeploymentCheck(name, comment string, fn CheckFunc[appsv1.Deployment]) {
	ch := NewCheck(name, "Deployment", comment, false)
	c.registerDeploymentCheck(GenCheck[appsv1.Deployment]{ch, fn})
}

func (c *Checks) RegisterOptionalDeploymentCheck(name, comment string, fn CheckFunc[appsv1.Deployment]) {
	ch := NewCheck(name, "Deployment", comment, true)
	c.registerDeploymentCheck(GenCheck[appsv1.Deployment]{ch, fn})
}

func (c *Checks) registerDeploymentCheck(ch GenCheck[appsv1.Deployment]) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.deployments[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) Deployments() map[string]GenCheck[appsv1.Deployment] {
	return c.deployments
}

func (c *Checks) RegisterIngressCheck(name, comment string, fn CheckFunc[ks.Ingress]) {
	ch := NewCheck(name, "Ingress", comment, false)
	c.registerIngressCheck(GenCheck[ks.Ingress]{ch, fn})
}

func (c *Checks) RegisterOptionalIngressCheck(name, comment string, fn CheckFunc[ks.Ingress]) {
	ch := NewCheck(name, "Ingress", comment, true)
	c.registerIngressCheck(GenCheck[ks.Ingress]{ch, fn})
}

func (c *Checks) registerIngressCheck(ch GenCheck[ks.Ingress]) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.ingresses[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) Ingresses() map[string]GenCheck[ks.Ingress] {
	return c.ingresses
}

func (c *Checks) RegisterNetworkPolicyCheck(name, comment string, fn CheckFunc[networkingv1.NetworkPolicy]) {
	ch := NewCheck(name, "NetworkPolicy", comment, false)
	c.registerNetworkPolicyCheck(GenCheck[networkingv1.NetworkPolicy]{ch, fn})
}

func (c *Checks) RegisterOptionalNetworkPolicyCheck(name, comment string, fn CheckFunc[networkingv1.NetworkPolicy]) {
	ch := NewCheck(name, "NetworkPolicy", comment, true)
	c.registerNetworkPolicyCheck(GenCheck[networkingv1.NetworkPolicy]{ch, fn})
}

func (c *Checks) registerNetworkPolicyCheck(ch GenCheck[networkingv1.NetworkPolicy]) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.networkpolicies[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) NetworkPolicies() map[string]GenCheck[networkingv1.NetworkPolicy] {
	return c.networkpolicies
}

func (c *Checks) RegisterPodDisruptionBudgetCheck(name, comment string, fn CheckFunc[ks.PodDisruptionBudget]) {
	ch := NewCheck(name, "PodDisruptionBudget", comment, false)
	c.registerPodDisruptionBudgetCheck(GenCheck[ks.PodDisruptionBudget]{ch, fn})
}

func (c *Checks) registerPodDisruptionBudgetCheck(ch GenCheck[ks.PodDisruptionBudget]) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}

	c.poddisruptionbudgets[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) PodDisruptionBudgets() map[string]GenCheck[ks.PodDisruptionBudget] {
	return c.poddisruptionbudgets
}

func (c *Checks) RegisterServiceCheck(name, comment string, fn CheckFunc[corev1.Service]) {
	ch := NewCheck(name, "Service", comment, false)
	c.registerServiceCheck(GenCheck[corev1.Service]{ch, fn})
}

func (c *Checks) RegisterOptionalServiceCheck(name, comment string, fn CheckFunc[corev1.Service]) {
	ch := NewCheck(name, "Service", comment, true)
	c.registerServiceCheck(GenCheck[corev1.Service]{ch, fn})
}

func (c *Checks) registerServiceCheck(ch GenCheck[corev1.Service]) {
	c.all = append(c.all, ch.Check)

	if !c.isEnabled(ch.Check) {
		return
	}
	c.services[machineFriendlyName(ch.Name)] = ch
}

func (c *Checks) Services() map[string]GenCheck[corev1.Service] {
	return c.services
}

func (c *Checks) All() []ks.Check {
	return c.all
}
