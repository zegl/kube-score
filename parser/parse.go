package parser

import (
	"bytes"
	ks "github.com/zegl/kube-score"
	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/parser/internal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"log"
)

var scheme = runtime.NewScheme()
var codecs = serializer.NewCodecFactory(scheme)

func init() {
	addToScheme(scheme)
}

func addToScheme(scheme *runtime.Scheme) {
	corev1.AddToScheme(scheme)
	appsv1.AddToScheme(scheme)
	networkingv1.AddToScheme(scheme)
	extensionsv1beta1.AddToScheme(scheme)
	appsv1beta1.AddToScheme(scheme)
	appsv1beta2.AddToScheme(scheme)
	batchv1.AddToScheme(scheme)
	batchv1beta1.AddToScheme(scheme)
	policyv1beta1.AddToScheme(scheme)
}

type detectKind struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

type parsedObjects struct {
	bothMetas            []ks.BothMeta
	pods                 []corev1.Pod
	podspecers           []ks.PodSpecer
	networkPolicies      []networkingv1.NetworkPolicy
	services             []corev1.Service
	podDisruptionBudgets []policyv1beta1.PodDisruptionBudget
	deployments          []appsv1.Deployment
	statefulsets         []appsv1.StatefulSet
	ingresses            []extensionsv1beta1.Ingress
	cronjobs             []batchv1beta1.CronJob
}

func (p *parsedObjects) Services() []corev1.Service {
	return p.services
}

func (p *parsedObjects) Pods() []corev1.Pod {
	return p.pods
}

func (p *parsedObjects) PodSpeccers() []ks.PodSpecer {
	return p.podspecers
}

func (p *parsedObjects) Ingresses() []extensionsv1beta1.Ingress {
	return p.ingresses
}

func (p *parsedObjects) PodDisruptionBudgets() []policyv1beta1.PodDisruptionBudget {
	return p.podDisruptionBudgets
}

func (p *parsedObjects) CronJobs() []batchv1beta1.CronJob {
	return p.cronjobs
}

func (p *parsedObjects) Deployments() []appsv1.Deployment {
	return p.deployments
}

func (p *parsedObjects) StatefulSets() []appsv1.StatefulSet {
	return p.statefulsets
}

func (p *parsedObjects) Metas() []ks.BothMeta {
	return p.bothMetas
}

func (p *parsedObjects) NetworkPolicies() []networkingv1.NetworkPolicy {
	return p.networkPolicies
}

func Empty() ks.AllTypes {
	return &parsedObjects{}
}

func ParseFiles(cnf config.Configuration) (ks.AllTypes, error) {
	s := &parsedObjects{}

	for _, file := range cnf.AllFiles {
		fullFile, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		// Convert to unix style newlines
		fullFile = bytes.Replace(fullFile, []byte("\r\n"), []byte("\n"), -1)

		for _, fileContents := range bytes.Split(fullFile, []byte("\n---\n")) {
			err := detectAndDecode(cnf, s, fileContents)
			if err != nil {
				return nil, err
			}
		}
	}

	return s, nil
}

func detectAndDecode(cnf config.Configuration, s *parsedObjects, raw []byte) error {
	var detect detectKind
	err := yaml.Unmarshal(raw, &detect)
	if err != nil {
		return err
	}

	detectedVersion := schema.FromAPIVersionAndKind(detect.ApiVersion, detect.Kind)

	// Parse lists and their items recursively
	if detectedVersion == corev1.SchemeGroupVersion.WithKind("List") {
		var list corev1.List
		decode(raw, &list)
		for _, listItem := range list.Items {
			err := detectAndDecode(cnf, s, listItem.Raw)
			if err != nil {
				return err
			}
		}
		return nil
	}

	decodeItem(cnf, s, detectedVersion, raw)

	return nil
}

func decode(data []byte, object runtime.Object) {
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(data, nil, object); err != nil {
		panic(err)
	}
}

func decodeItem(cnf config.Configuration, s *parsedObjects, detectedVersion schema.GroupVersionKind, fileContents []byte) {
	addPodSpeccer := func(ps ks.PodSpecer) {
		s.podspecers = append(s.podspecers, ps)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{ps.GetTypeMeta(), ps.GetObjectMeta()})
	}

	switch detectedVersion {
	case corev1.SchemeGroupVersion.WithKind("Pod"):
		var pod corev1.Pod
		decode(fileContents, &pod)
		s.pods = append(s.pods, pod)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{pod.TypeMeta, pod.ObjectMeta})

	case batchv1.SchemeGroupVersion.WithKind("Job"):
		var job batchv1.Job
		decode(fileContents, &job)
		addPodSpeccer(internal.Batchv1Job{job})

	case batchv1beta1.SchemeGroupVersion.WithKind("CronJob"):
		var cronjob batchv1beta1.CronJob
		decode(fileContents, &cronjob)
		addPodSpeccer(internal.Batchv1beta1CronJob{cronjob})
		s.cronjobs = append(s.cronjobs, cronjob)

	case appsv1.SchemeGroupVersion.WithKind("Deployment"):
		var deployment appsv1.Deployment
		decode(fileContents, &deployment)
		addPodSpeccer(internal.Appsv1Deployment{deployment})

		// TODO: Support older versions of Deployment as well?
		s.deployments = append(s.deployments, deployment)
	case appsv1beta1.SchemeGroupVersion.WithKind("Deployment"):
		var deployment appsv1beta1.Deployment
		decode(fileContents, &deployment)
		addPodSpeccer(internal.Appsv1beta1Deployment{deployment})
	case appsv1beta2.SchemeGroupVersion.WithKind("Deployment"):
		var deployment appsv1beta2.Deployment
		decode(fileContents, &deployment)
		addPodSpeccer(internal.Appsv1beta2Deployment{deployment})
	case extensionsv1beta1.SchemeGroupVersion.WithKind("Deployment"):
		var deployment extensionsv1beta1.Deployment
		decode(fileContents, &deployment)
		addPodSpeccer(internal.Extensionsv1beta1Deployment{deployment})

	case appsv1.SchemeGroupVersion.WithKind("StatefulSet"):
		var statefulSet appsv1.StatefulSet
		decode(fileContents, &statefulSet)
		addPodSpeccer(internal.Appsv1StatefulSet{statefulSet})

		// TODO: Support older versions of StatefulSet as well?
		s.statefulsets = append(s.statefulsets, statefulSet)
	case appsv1beta1.SchemeGroupVersion.WithKind("StatefulSet"):
		var statefulSet appsv1beta1.StatefulSet
		decode(fileContents, &statefulSet)
		addPodSpeccer(internal.Appsv1beta1StatefulSet{statefulSet})
	case appsv1beta2.SchemeGroupVersion.WithKind("StatefulSet"):
		var statefulSet appsv1beta2.StatefulSet
		decode(fileContents, &statefulSet)
		addPodSpeccer(internal.Appsv1beta2StatefulSet{statefulSet})

	case appsv1.SchemeGroupVersion.WithKind("DaemonSet"):
		var daemonset appsv1.DaemonSet
		decode(fileContents, &daemonset)
		addPodSpeccer(internal.Appsv1DaemonSet{daemonset})
	case appsv1beta2.SchemeGroupVersion.WithKind("DaemonSet"):
		var daemonset appsv1beta2.DaemonSet
		decode(fileContents, &daemonset)
		addPodSpeccer(internal.Appsv1beta2DaemonSet{daemonset})
	case extensionsv1beta1.SchemeGroupVersion.WithKind("DaemonSet"):
		var daemonset extensionsv1beta1.DaemonSet
		decode(fileContents, &daemonset)
		addPodSpeccer(internal.Extensionsv1beta1DaemonSet{daemonset})

	case networkingv1.SchemeGroupVersion.WithKind("NetworkPolicy"):
		var netpol networkingv1.NetworkPolicy
		decode(fileContents, &netpol)
		s.networkPolicies = append(s.networkPolicies, netpol)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{netpol.TypeMeta, netpol.ObjectMeta})

	case corev1.SchemeGroupVersion.WithKind("Service"):
		var service corev1.Service
		decode(fileContents, &service)
		s.services = append(s.services, service)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{service.TypeMeta, service.ObjectMeta})

	case policyv1beta1.SchemeGroupVersion.WithKind("PodDisruptionBudget"):
		var disruptBudget policyv1beta1.PodDisruptionBudget
		decode(fileContents, &disruptBudget)
		s.podDisruptionBudgets = append(s.podDisruptionBudgets, disruptBudget)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{disruptBudget.TypeMeta, disruptBudget.ObjectMeta})

	case extensionsv1beta1.SchemeGroupVersion.WithKind("Ingress"):
		var ingress extensionsv1beta1.Ingress
		decode(fileContents, &ingress)
		s.ingresses = append(s.ingresses, ingress)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{ingress.TypeMeta, ingress.ObjectMeta})

	default:
		if cnf.VerboseOutput {
			log.Printf("Unknown datatype: %s", detectedVersion.String())
		}
	}
}
