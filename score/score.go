package score

import (
	"github.com/zegl/kube-score"
	"github.com/zegl/kube-score/config"
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

// Score runs a pre-configured list of tests against the files defined in the configuration, and returns a scorecard.
// Additional configuration and tuning parameters can be provided via the config.
func Score(allObjects kube_score.AllTypes, cnf config.Configuration) (*scorecard.Scorecard, error) {
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

	scoreCard := scorecard.New()

	for _, ingress := range allObjects.Ingresses() {
		for _, ingressTest := range allChecks.Ingresses() {
			score := ingressTest(ingress)
			score.AddMeta(ingress.TypeMeta, ingress.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, meta := range allObjects.Metas() {
		for _, metaTest := range allChecks.Metas() {
			score := metaTest(meta.TypeMeta)
			score.AddMeta(meta.TypeMeta, meta.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, pod := range allObjects.Pods() {
		for _, podTest := range allChecks.Pods() {
			score := podTest(corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			})
			score.AddMeta(pod.TypeMeta, pod.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, podspecer := range allObjects.PodSpeccers() {
		for _, podTest := range allChecks.Pods() {
			score := podTest(podspecer.GetPodTemplateSpec())
			score.AddMeta(podspecer.GetTypeMeta(), podspecer.GetObjectMeta())
			scoreCard.Add(score)
		}
	}

	for _, service := range allObjects.Services() {
		for _, serviceTest := range allChecks.Services() {
			score := serviceTest(service)
			score.AddMeta(service.TypeMeta, service.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, statefulset := range allObjects.StatefulSets() {
		for _, test := range allChecks.StatefulSets() {
			score := test(statefulset)
			score.AddMeta(statefulset.TypeMeta, statefulset.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, deployment := range allObjects.Deployments() {
		for _, test := range allChecks.Deployments() {
			score := test(deployment)
			score.AddMeta(deployment.TypeMeta, deployment.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, netpol := range allObjects.NetworkPolicies() {
		for _, netpolTest := range allChecks.NetworkPolicies() {
			score := netpolTest(netpol)
			score.AddMeta(netpol.TypeMeta, netpol.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, cjob := range allObjects.CronJobs() {
		for _, cronjobTest := range allChecks.CronJobs() {
			score := cronjobTest(cjob)
			score.AddMeta(cjob.TypeMeta, cjob.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	return scoreCard, nil
}
