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
		for _, test := range allChecks.Ingresses() {
			score := test.Fn(ingress)
			score.AddMeta(ingress.TypeMeta, ingress.ObjectMeta)
			scoreCard.Add(score, test.Check)
		}
	}

	for _, meta := range allObjects.Metas() {
		for _, test := range allChecks.Metas() {
			score := test.Fn(meta.TypeMeta)
			score.AddMeta(meta.TypeMeta, meta.ObjectMeta)
			scoreCard.Add(score, test.Check)
		}
	}

	for _, pod := range allObjects.Pods() {
		for _, test := range allChecks.Pods() {
			score := test.Fn(corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			})
			score.AddMeta(pod.TypeMeta, pod.ObjectMeta)
			scoreCard.Add(score, test.Check)
		}
	}

	for _, podspecer := range allObjects.PodSpeccers() {
		for _, test := range allChecks.Pods() {
			score := test.Fn(podspecer.GetPodTemplateSpec())
			score.AddMeta(podspecer.GetTypeMeta(), podspecer.GetObjectMeta())
			scoreCard.Add(score, test.Check)
		}
	}

	for _, service := range allObjects.Services() {
		for _, test := range allChecks.Services() {
			score := test.Fn(service)
			score.AddMeta(service.TypeMeta, service.ObjectMeta)
			scoreCard.Add(score, test.Check)
		}
	}

	for _, statefulset := range allObjects.StatefulSets() {
		for _, test := range allChecks.StatefulSets() {
			score := test.Fn(statefulset)
			score.AddMeta(statefulset.TypeMeta, statefulset.ObjectMeta)
			scoreCard.Add(score, test.Check)
		}
	}

	for _, deployment := range allObjects.Deployments() {
		for _, test := range allChecks.Deployments() {
			score := test.Fn(deployment)
			score.AddMeta(deployment.TypeMeta, deployment.ObjectMeta)
			scoreCard.Add(score, test.Check)
		}
	}

	for _, netpol := range allObjects.NetworkPolicies() {
		for _, test := range allChecks.NetworkPolicies() {
			score := test.Fn(netpol)
			score.AddMeta(netpol.TypeMeta, netpol.ObjectMeta)
			scoreCard.Add(score, test.Check)
		}
	}

	for _, cjob := range allObjects.CronJobs() {
		for _, test := range allChecks.CronJobs() {
			score := test.Fn(cjob)
			score.AddMeta(cjob.TypeMeta, cjob.ObjectMeta)
			scoreCard.Add(score, test.Check)
		}
	}

	return scoreCard, nil
}
