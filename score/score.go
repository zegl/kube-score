package score

import (
	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/apps"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/score/container"
	"github.com/zegl/kube-score/score/cronjob"
	"github.com/zegl/kube-score/score/disruptionbudget"
	"github.com/zegl/kube-score/score/ingress"
	"github.com/zegl/kube-score/score/networkpolicy"
	"github.com/zegl/kube-score/score/probes"
	"github.com/zegl/kube-score/score/security"
	"github.com/zegl/kube-score/score/service"
	"github.com/zegl/kube-score/score/stable"
	"github.com/zegl/kube-score/scorecard"

	corev1 "k8s.io/api/core/v1"
)

func RegisterAllChecks(allObjects ks.AllTypes, cnf config.Configuration) *checks.Checks {
	allChecks := checks.New()

	ingress.Register(allChecks, allObjects)
	cronjob.Register(allChecks)
	container.Register(allChecks, cnf)
	disruptionbudget.Register(allChecks, allObjects)
	networkpolicy.Register(allChecks, allObjects, allObjects, allObjects)
	probes.Register(allChecks, allObjects)
	security.Register(allChecks)
	service.Register(allChecks, allObjects, allObjects)
	stable.Register(allChecks)
	apps.Register(allChecks)

	return allChecks
}

// Score runs a pre-configured list of tests against the files defined in the configuration, and returns a scorecard.
// Additional configuration and tuning parameters can be provided via the config.
func Score(allObjects ks.AllTypes, cnf config.Configuration) (scorecard.Scorecard, error) {
	allChecks := RegisterAllChecks(allObjects, cnf)
	scoreCard := scorecard.New()

	for _, ingress := range allObjects.Ingresses() {
		o := scoreCard.NewObject(ingress.TypeMeta, ingress.ObjectMeta)
		for _, test := range allChecks.Ingresses() {
			o.Add(test.Fn(ingress), test.Check)
		}
	}

	for _, meta := range allObjects.Metas() {
		o := scoreCard.NewObject(meta.TypeMeta, meta.ObjectMeta)
		for _, test := range allChecks.Metas() {
			o.Add(test.Fn(meta.TypeMeta), test.Check)
		}
	}

	for _, pod := range allObjects.Pods() {
		o := scoreCard.NewObject(pod.TypeMeta, pod.ObjectMeta)
		for _, test := range allChecks.Pods() {
			score := test.Fn(corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			}, pod.TypeMeta)
			o.Add(score, test.Check)
		}
	}

	for _, podspecer := range allObjects.PodSpeccers() {
		o := scoreCard.NewObject(podspecer.GetTypeMeta(), podspecer.GetObjectMeta())
		for _, test := range allChecks.Pods() {
			score := test.Fn(podspecer.GetPodTemplateSpec(), podspecer.GetTypeMeta())
			o.Add(score, test.Check)
		}
	}

	for _, service := range allObjects.Services() {
		o := scoreCard.NewObject(service.TypeMeta, service.ObjectMeta)
		for _, test := range allChecks.Services() {
			o.Add(test.Fn(service), test.Check)
		}
	}

	for _, statefulset := range allObjects.StatefulSets() {
		o := scoreCard.NewObject(statefulset.TypeMeta, statefulset.ObjectMeta)
		for _, test := range allChecks.StatefulSets() {
			o.Add(test.Fn(statefulset), test.Check)
		}
	}

	for _, deployment := range allObjects.Deployments() {
		o := scoreCard.NewObject(deployment.TypeMeta, deployment.ObjectMeta)
		for _, test := range allChecks.Deployments() {
			o.Add(test.Fn(deployment), test.Check)
		}
	}

	for _, netpol := range allObjects.NetworkPolicies() {
		o := scoreCard.NewObject(netpol.TypeMeta, netpol.ObjectMeta)
		for _, test := range allChecks.NetworkPolicies() {
			o.Add(test.Fn(netpol), test.Check)
		}
	}

	for _, cjob := range allObjects.CronJobs() {
		o := scoreCard.NewObject(cjob.TypeMeta, cjob.ObjectMeta)
		for _, test := range allChecks.CronJobs() {
			o.Add(test.Fn(cjob), test.Check)
		}
	}

	return scoreCard, nil
}
