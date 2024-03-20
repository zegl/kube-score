package checks

import (
	"strings"

	routev1 "github.com/openshift/api/route/v1"
	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

func New(cnf config.Configuration) *Checks {
	return &Checks{
		cnf: cnf,

		all:                      make([]ks.Check, 0),
		metas:                    make(map[string]GenCheck[ks.BothMeta]),
		pods:                     make(map[string]GenCheck[ks.PodSpecer]),
		services:                 make(map[string]GenCheck[corev1.Service]),
		statefulsets:             make(map[string]GenCheck[appsv1.StatefulSet]),
		deployments:              make(map[string]GenCheck[appsv1.Deployment]),
		networkpolicies:          make(map[string]GenCheck[networkingv1.NetworkPolicy]),
		ingresses:                make(map[string]GenCheck[ks.Ingress]),
		cronjobs:                 make(map[string]GenCheck[ks.CronJob]),
		horizontalPodAutoscalers: make(map[string]GenCheck[ks.HpaTargeter]),
		poddisruptionbudgets:     make(map[string]GenCheck[ks.PodDisruptionBudget]),
		routes:                   make(map[string]GenCheck[routev1.Route]),
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

type CheckFunc[T any] func(T) (scorecard.TestScore, error)

type GenCheck[T any] struct {
	ks.Check
	Fn CheckFunc[T]
}

type Checks struct {
	all                      []ks.Check
	metas                    map[string]GenCheck[ks.BothMeta]
	pods                     map[string]GenCheck[ks.PodSpecer]
	services                 map[string]GenCheck[corev1.Service]
	statefulsets             map[string]GenCheck[appsv1.StatefulSet]
	deployments              map[string]GenCheck[appsv1.Deployment]
	networkpolicies          map[string]GenCheck[networkingv1.NetworkPolicy]
	ingresses                map[string]GenCheck[ks.Ingress]
	cronjobs                 map[string]GenCheck[ks.CronJob]
	horizontalPodAutoscalers map[string]GenCheck[ks.HpaTargeter]
	poddisruptionbudgets     map[string]GenCheck[ks.PodDisruptionBudget]
	routes                   map[string]GenCheck[routev1.Route]

	cnf config.Configuration
}

func (c Checks) isEnabled(check ks.Check) bool {
	_, ok := c.cnf.IgnoredTests[check.ID]
	return !ok
}

func (c *Checks) RegisterMetaCheck(name, comment string, fn CheckFunc[ks.BothMeta]) {
	reg(c, "all", name, comment, false, fn, c.metas)
}

func (c *Checks) RegisterOptionalMetaCheck(name, comment string, fn CheckFunc[ks.BothMeta]) {
	reg(c, "all", name, comment, true, fn, c.metas)
}

func (c *Checks) Metas() map[string]GenCheck[ks.BothMeta] {
	return c.metas
}

func reg[T any](c *Checks, targetType, name, comment string, optional bool, fn CheckFunc[T], mp map[string]GenCheck[T]) {
	ch := NewCheck(name, targetType, comment, optional)
	check := GenCheck[T]{Check: ch, Fn: fn}
	c.all = append(c.all, check.Check)
	if !c.isEnabled(check.Check) {
		return
	}
	mp[machineFriendlyName(ch.Name)] = check
}

func (c *Checks) RegisterPodCheck(name, comment string, fn CheckFunc[ks.PodSpecer]) {
	reg(c, "Pod", name, comment, false, fn, c.pods)
}

func (c *Checks) RegisterOptionalPodCheck(name, comment string, fn CheckFunc[ks.PodSpecer]) {
	reg(c, "Pod", name, comment, true, fn, c.pods)
}

func (c *Checks) Pods() map[string]GenCheck[ks.PodSpecer] {
	return c.pods
}

func (c *Checks) RegisterHorizontalPodAutoscalerCheck(name, comment string, fn CheckFunc[ks.HpaTargeter]) {
	reg(c, "HorizontalPodAutoscaler", name, comment, false, fn, c.horizontalPodAutoscalers)
}

func (c *Checks) RegisterOptionalHorizontalPodAutoscalerCheck(name, comment string, fn CheckFunc[ks.HpaTargeter]) {
	reg(c, "HorizontalPodAutoscaler", name, comment, true, fn, c.horizontalPodAutoscalers)
}

func (c *Checks) HorizontalPodAutoscalers() map[string]GenCheck[ks.HpaTargeter] {
	return c.horizontalPodAutoscalers
}

func (c *Checks) RegisterCronJobCheck(name, comment string, fn CheckFunc[ks.CronJob]) {
	reg(c, "CronJob", name, comment, false, fn, c.cronjobs)
}

func (c *Checks) RegisterOptionalCronJobCheck(name, comment string, fn CheckFunc[ks.CronJob]) {
	reg(c, "CronJob", name, comment, true, fn, c.cronjobs)
}

func (c *Checks) CronJobs() map[string]GenCheck[ks.CronJob] {
	return c.cronjobs
}

func (c *Checks) RegisterStatefulSetCheck(name, comment string, fn CheckFunc[appsv1.StatefulSet]) {
	reg(c, "StatefulSet", name, comment, false, fn, c.statefulsets)
}

func (c *Checks) RegisterOptionalStatefulSetCheck(name, comment string, fn CheckFunc[appsv1.StatefulSet]) {
	reg(c, "StatefulSet", name, comment, true, fn, c.statefulsets)
}

func (c *Checks) StatefulSets() map[string]GenCheck[appsv1.StatefulSet] {
	return c.statefulsets
}

func (c *Checks) RegisterDeploymentCheck(name, comment string, fn CheckFunc[appsv1.Deployment]) {
	reg(c, "Deployment", name, comment, false, fn, c.deployments)
}

func (c *Checks) RegisterOptionalDeploymentCheck(name, comment string, fn CheckFunc[appsv1.Deployment]) {
	reg(c, "Deployment", name, comment, true, fn, c.deployments)
}

func (c *Checks) Deployments() map[string]GenCheck[appsv1.Deployment] {
	return c.deployments
}

func (c *Checks) RegisterIngressCheck(name, comment string, fn CheckFunc[ks.Ingress]) {
	reg(c, "Ingress", name, comment, false, fn, c.ingresses)
}

func (c *Checks) RegisterOptionalIngressCheck(name, comment string, fn CheckFunc[ks.Ingress]) {
	reg(c, "Ingress", name, comment, true, fn, c.ingresses)
}

func (c *Checks) Ingresses() map[string]GenCheck[ks.Ingress] {
	return c.ingresses
}

func (c *Checks) RegisterNetworkPolicyCheck(name, comment string, fn CheckFunc[networkingv1.NetworkPolicy]) {
	reg(c, "NetworkPolicy", name, comment, false, fn, c.networkpolicies)
}

func (c *Checks) RegisterOptionalNetworkPolicyCheck(name, comment string, fn CheckFunc[networkingv1.NetworkPolicy]) {
	reg(c, "NetworkPolicy", name, comment, true, fn, c.networkpolicies)
}

func (c *Checks) NetworkPolicies() map[string]GenCheck[networkingv1.NetworkPolicy] {
	return c.networkpolicies
}

func (c *Checks) RegisterPodDisruptionBudgetCheck(name, comment string, fn CheckFunc[ks.PodDisruptionBudget]) {
	reg(c, "PodDisruptionBudget", name, comment, false, fn, c.poddisruptionbudgets)
}

func (c *Checks) PodDisruptionBudgets() map[string]GenCheck[ks.PodDisruptionBudget] {
	return c.poddisruptionbudgets
}

func (c *Checks) RegisterServiceCheck(name, comment string, fn CheckFunc[corev1.Service]) {
	reg(c, "Service", name, comment, false, fn, c.services)
}

func (c *Checks) RegisterOptionalServiceCheck(name, comment string, fn CheckFunc[corev1.Service]) {
	reg(c, "Service", name, comment, true, fn, c.services)
}

func (c *Checks) Services() map[string]GenCheck[corev1.Service] {
	return c.services
}

func (c *Checks) RegisterRouteCheck(name, comment string, fn CheckFunc[routev1.Route]) {
	reg(c, "Route", name, comment, false, fn, c.routes)
}

func (c *Checks) RegisterOptionalRouteCheck(name, comment string, fn CheckFunc[routev1.Route]) {
	reg(c, "Route", name, comment, true, fn, c.routes)
}

func (c *Checks) Routes() map[string]GenCheck[routev1.Route] {
	return c.routes
}

func (c *Checks) All() []ks.Check {
	return c.all
}
