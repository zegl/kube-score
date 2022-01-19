package score

import (
	"kube-score/score/ephemeralstorage"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/apps"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/score/container"
	"github.com/zegl/kube-score/score/cronjob"
	"github.com/zegl/kube-score/score/disruptionbudget"
	"github.com/zegl/kube-score/score/ephemeralstorage"
	"github.com/zegl/kube-score/score/hpa"
	"github.com/zegl/kube-score/score/ingress"
	"github.com/zegl/kube-score/score/meta"
	"github.com/zegl/kube-score/score/networkpolicy"
	"github.com/zegl/kube-score/score/probes"
	"github.com/zegl/kube-score/score/security"
	"github.com/zegl/kube-score/score/service"
	"github.com/zegl/kube-score/score/stable"
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RegisterAllChecks(allObjects ks.AllTypes, cnf config.Configuration) *checks.Checks {
	allChecks := checks.New(cnf)

	ingress.Register(allChecks, allObjects)
	cronjob.Register(allChecks)
	container.Register(allChecks, cnf)
	disruptionbudget.Register(allChecks, allObjects)
	networkpolicy.Register(allChecks, allObjects, allObjects, allObjects)
	probes.Register(allChecks, allObjects)
	security.Register(allChecks)
	service.Register(allChecks, allObjects, allObjects)
	stable.Register(cnf.KubernetesVersion, allChecks)
	apps.Register(allChecks, allObjects.HorizontalPodAutoscalers(), allObjects.Services())
	meta.Register(allChecks)
	hpa.Register(allChecks, allObjects.Metas())
	ephemeralstorage.Register(allChecks)

	return allChecks
}

// Score runs a pre-configured list of tests against the files defined in the configuration, and returns a scorecard.
// Additional configuration and tuning parameters can be provided via the config.
func Score(allObjects ks.AllTypes, cnf config.Configuration) (*scorecard.Scorecard, error) {
	allChecks := RegisterAllChecks(allObjects, cnf)
	scoreCard := scorecard.New()

	newObject := func(typeMeta metav1.TypeMeta, objectMeta metav1.ObjectMeta) *scorecard.ScoredObject {
		return scoreCard.NewObject(typeMeta, objectMeta, cnf.UseIgnoreChecksAnnotation)
	}

	for _, ingress := range allObjects.Ingresses() {
		o := newObject(ingress.GetTypeMeta(), ingress.GetObjectMeta())
		for _, test := range allChecks.Ingresses() {
			o.Add(test.Fn(ingress), test.Check, ingress)
		}
	}

	for _, meta := range allObjects.Metas() {
		o := newObject(meta.TypeMeta, meta.ObjectMeta)
		for _, test := range allChecks.Metas() {
			o.Add(test.Fn(meta), test.Check, meta)
		}
	}

	for _, pod := range allObjects.Pods() {
		o := newObject(pod.Pod().TypeMeta, pod.Pod().ObjectMeta)
		for _, test := range allChecks.Pods() {
			score := test.Fn(corev1.PodTemplateSpec{
				ObjectMeta: pod.Pod().ObjectMeta,
				Spec:       pod.Pod().Spec,
			}, pod.Pod().TypeMeta)
			o.Add(score, test.Check, pod)
		}
	}

	for _, podspecer := range allObjects.PodSpeccers() {
		o := newObject(podspecer.GetTypeMeta(), podspecer.GetObjectMeta())
		for _, test := range allChecks.Pods() {
			score := test.Fn(podspecer.GetPodTemplateSpec(), podspecer.GetTypeMeta())
			o.Add(score, test.Check, podspecer)
		}
	}

	for _, service := range allObjects.Services() {
		o := newObject(service.Service().TypeMeta, service.Service().ObjectMeta)
		for _, test := range allChecks.Services() {
			o.Add(test.Fn(service.Service()), test.Check, service)
		}
	}

	for _, statefulset := range allObjects.StatefulSets() {
		o := newObject(statefulset.StatefulSet().TypeMeta, statefulset.StatefulSet().ObjectMeta)
		for _, test := range allChecks.StatefulSets() {
			res, err := test.Fn(statefulset.StatefulSet())
			if err != nil {
				return nil, err
			}
			o.Add(res, test.Check, statefulset)
		}
	}

	for _, deployment := range allObjects.Deployments() {
		o := newObject(deployment.Deployment().TypeMeta, deployment.Deployment().ObjectMeta)
		for _, test := range allChecks.Deployments() {
			res, err := test.Fn(deployment.Deployment())
			if err != nil {
				return nil, err
			}
			o.Add(res, test.Check, deployment)
		}
	}

	for _, netpol := range allObjects.NetworkPolicies() {
		o := newObject(netpol.NetworkPolicy().TypeMeta, netpol.NetworkPolicy().ObjectMeta)
		for _, test := range allChecks.NetworkPolicies() {
			o.Add(test.Fn(netpol.NetworkPolicy()), test.Check, netpol)
		}
	}

	for _, cjob := range allObjects.CronJobs() {
		o := newObject(cjob.GetTypeMeta(), cjob.GetObjectMeta())
		for _, test := range allChecks.CronJobs() {
			o.Add(test.Fn(cjob), test.Check, cjob)
		}
	}

	for _, hpa := range allObjects.HorizontalPodAutoscalers() {
		o := newObject(hpa.GetTypeMeta(), hpa.GetObjectMeta())
		for _, test := range allChecks.HorizontalPodAutoscalers() {
			o.Add(test.Fn(hpa), test.Check, hpa)
		}
	}

	for _, pdb := range allObjects.PodDisruptionBudgets() {
		o := newObject(pdb.GetTypeMeta(), pdb.GetObjectMeta())
		for _, test := range allChecks.PodDisruptionBudgets() {
			o.Add(test.Fn(pdb), test.Check, pdb)
		}
	}

	return &scoreCard, nil
}
