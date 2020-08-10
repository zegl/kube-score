package score

import (
	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/apps"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/score/container"
	"github.com/zegl/kube-score/score/cronjob"
	"github.com/zegl/kube-score/score/disruptionbudget"
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
	apps.Register(allChecks, allObjects.HorizontalPodAutoscalers())
	meta.Register(allChecks)
	hpa.Register(allChecks, allObjects.Metas())

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
		o := newObject(ingress.TypeMeta, ingress.ObjectMeta)
		for _, test := range allChecks.Ingresses() {
			o.Add(test.Fn(ingress), test.Check)
		}
	}

	for _, meta := range allObjects.Metas() {
		o := newObject(meta.TypeMeta, meta.ObjectMeta)
		for _, test := range allChecks.Metas() {
			o.Add(test.Fn(meta), test.Check)
		}
	}

	for _, pod := range allObjects.Pods() {
		o := newObject(pod.TypeMeta, pod.ObjectMeta)
		for _, test := range allChecks.Pods() {
			score := test.Fn(corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			}, pod.TypeMeta)
			o.Add(score, test.Check)
		}
	}

	for _, podspecer := range allObjects.PodSpeccers() {
		o := newObject(podspecer.GetTypeMeta(), podspecer.GetObjectMeta())
		for _, test := range allChecks.Pods() {
			score := test.Fn(podspecer.GetPodTemplateSpec(), podspecer.GetTypeMeta())
			o.Add(score, test.Check)
		}
	}

	for _, service := range allObjects.Services() {
		o := newObject(service.TypeMeta, service.ObjectMeta)
		for _, test := range allChecks.Services() {
			o.Add(test.Fn(service), test.Check)
		}
	}

	for _, statefulset := range allObjects.StatefulSets() {
		o := newObject(statefulset.TypeMeta, statefulset.ObjectMeta)
		for _, test := range allChecks.StatefulSets() {
			res, err := test.Fn(statefulset)
			if err != nil {
				return nil, err
			}
			o.Add(res, test.Check)
		}
	}

	for _, deployment := range allObjects.Deployments() {
		o := newObject(deployment.TypeMeta, deployment.ObjectMeta)
		for _, test := range allChecks.Deployments() {
			res, err := test.Fn(deployment)
			if err != nil {
				return nil, err
			}
			o.Add(res, test.Check)
		}
	}

	for _, netpol := range allObjects.NetworkPolicies() {
		o := newObject(netpol.TypeMeta, netpol.ObjectMeta)
		for _, test := range allChecks.NetworkPolicies() {
			o.Add(test.Fn(netpol), test.Check)
		}
	}

	for _, cjob := range allObjects.CronJobs() {
		o := newObject(cjob.TypeMeta, cjob.ObjectMeta)
		for _, test := range allChecks.CronJobs() {
			o.Add(test.Fn(cjob), test.Check)
		}
	}

	for _, hpa := range allObjects.HorizontalPodAutoscalers() {
		o := newObject(hpa.GetTypeMeta(), hpa.GetObjectMeta())
		for _, test := range allChecks.HorizontalPodAutoscalers() {
			o.Add(test.Fn(hpa), test.Check)
		}
	}

	return &scoreCard, nil
}
