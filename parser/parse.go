package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	policyv1 "k8s.io/api/policy/v1"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	autoscalingv2beta1 "k8s.io/api/autoscaling/v2beta1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/parser/internal"
	internalcronjob "github.com/zegl/kube-score/parser/internal/cronjob"
	internalnetpol "github.com/zegl/kube-score/parser/internal/networkpolicy"
	internalpdb "github.com/zegl/kube-score/parser/internal/pdb"
	internalpod "github.com/zegl/kube-score/parser/internal/pod"
	internalservice "github.com/zegl/kube-score/parser/internal/service"
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
	pods                 []ks.Pod
	podspecers           []ks.PodSpecer
	networkPolicies      []ks.NetworkPolicy
	services             []ks.Service
	podDisruptionBudgets []ks.PodDisruptionBudget
	deployments          []ks.Deployment
	statefulsets         []ks.StatefulSet
	ingresses            []ks.Ingress // supports multiple versions of ingress
	cronjobs             []ks.CronJob
	hpaTargeters         []ks.HpaTargeter // all versions of HPAs
}

func (p *parsedObjects) Services() []ks.Service {
	return p.services
}

func (p *parsedObjects) Pods() []ks.Pod {
	return p.pods
}

func (p *parsedObjects) PodSpeccers() []ks.PodSpecer {
	return p.podspecers
}

func (p *parsedObjects) Ingresses() []ks.Ingress {
	return p.ingresses
}

func (p *parsedObjects) PodDisruptionBudgets() []ks.PodDisruptionBudget {
	return p.podDisruptionBudgets
}

func (p *parsedObjects) CronJobs() []ks.CronJob {
	return p.cronjobs
}

func (p *parsedObjects) Deployments() []ks.Deployment {
	return p.deployments
}

func (p *parsedObjects) StatefulSets() []ks.StatefulSet {
	return p.statefulsets
}

func (p *parsedObjects) Metas() []ks.BothMeta {
	return p.bothMetas
}

func (p *parsedObjects) NetworkPolicies() []ks.NetworkPolicy {
	return p.networkPolicies
}

func (p *parsedObjects) HorizontalPodAutoscalers() []ks.HpaTargeter {
	return p.hpaTargeters
}

func Empty() ks.AllTypes {
	return &parsedObjects{}
}

func ParseFiles(cnf config.Configuration) (ks.AllTypes, error) {
	s := &parsedObjects{}

	for _, namedReader := range cnf.AllFiles {
		fullFile, err := ioutil.ReadAll(namedReader)
		if err != nil {
			return nil, err
		}

		// Convert to unix style newlines
		fullFile = bytes.Replace(fullFile, []byte("\r\n"), []byte("\n"), -1)

		offset := 1 // Line numbers are 1 indexed

		// Remove initial "---\n" if present
		if bytes.HasPrefix(fullFile, []byte("---\n")) {
			fullFile = fullFile[4:]
			offset = 2
		}

		for _, fileContents := range bytes.Split(fullFile, []byte("\n---\n")) {

			if len(bytes.TrimSpace(fileContents)) > 0 {
				err := detectAndDecode(cnf, s, namedReader.Name(), offset, fileContents)
				if err != nil {
					return nil, err
				}
			}

			offset += 2 + bytes.Count(fileContents, []byte("\n"))
		}
	}

	return s, nil
}

func detectAndDecode(cnf config.Configuration, s *parsedObjects, fileName string, fileOffset int, raw []byte) error {
	var detect detectKind
	err := yaml.Unmarshal(raw, &detect)
	if err != nil {
		return err
	}

	detectedVersion := schema.FromAPIVersionAndKind(detect.ApiVersion, detect.Kind)

	// Parse lists and their items recursively
	if detectedVersion == corev1.SchemeGroupVersion.WithKind("List") {
		var list corev1.List
		err := decode(raw, &list)
		if err != nil {
			return err
		}
		for _, listItem := range list.Items {
			err := detectAndDecode(cnf, s, fileName, fileOffset, listItem.Raw)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = decodeItem(cnf, s, detectedVersion, fileName, fileOffset, raw)
	if err != nil {
		return err
	}

	return nil
}

func decode(data []byte, object runtime.Object) error {
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(data, nil, object); err != nil {
		gvk := object.GetObjectKind().GroupVersionKind()
		return fmt.Errorf("Failed to parse %s: err=%w", gvk, err)
	}
	return nil
}

func detectFileLocation(fileName string, fileOffset int, fileContents []byte) ks.FileLocation {
	// If the object YAML begins with a Helm style "# Source: " comment
	// Use the information in there as the file name
	firstRow := string(bytes.Split(fileContents, []byte("\n"))[0])
	helmTemplatePrefix := "# Source: "
	if strings.HasPrefix(firstRow, helmTemplatePrefix) {
		return ks.FileLocation{
			Name: firstRow[len(helmTemplatePrefix):],
			Line: 1, // Set line to 1 as the line definition gets lost in Helm
		}
	}

	return ks.FileLocation{
		Name: fileName,
		Line: fileOffset,
	}
}

func decodeItem(cnf config.Configuration, s *parsedObjects, detectedVersion schema.GroupVersionKind, fileName string, fileOffset int, fileContents []byte) error {
	addPodSpeccer := func(ps ks.PodSpecer) {
		s.podspecers = append(s.podspecers, ps)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{ps.GetTypeMeta(), ps.GetObjectMeta(), ps})
	}

	fileLocation := detectFileLocation(fileName, fileOffset, fileContents)

	var errs parseError

	switch detectedVersion {
	case corev1.SchemeGroupVersion.WithKind("Pod"):
		var pod corev1.Pod
		errs.AddIfErr(decode(fileContents, &pod))
		p := internalpod.Pod{pod, fileLocation}
		s.pods = append(s.pods, p)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{pod.TypeMeta, pod.ObjectMeta, p})

	case batchv1.SchemeGroupVersion.WithKind("Job"):
		var job batchv1.Job
		errs.AddIfErr(decode(fileContents, &job))
		addPodSpeccer(internal.Batchv1Job{job, fileLocation})

	case batchv1beta1.SchemeGroupVersion.WithKind("CronJob"):
		var cronjob batchv1beta1.CronJob
		errs.AddIfErr(decode(fileContents, &cronjob))
		cjob := internalcronjob.CronJobV1beta1{cronjob, fileLocation}
		addPodSpeccer(cjob)
		s.cronjobs = append(s.cronjobs, cjob)

	case batchv1.SchemeGroupVersion.WithKind("CronJob"):
		var cronjob batchv1.CronJob
		errs.AddIfErr(decode(fileContents, &cronjob))
		cjob := internalcronjob.CronJobV1{cronjob, fileLocation}
		addPodSpeccer(cjob)
		s.cronjobs = append(s.cronjobs, cjob)

	case appsv1.SchemeGroupVersion.WithKind("Deployment"):
		var deployment appsv1.Deployment
		errs.AddIfErr(decode(fileContents, &deployment))
		deploy := internal.Appsv1Deployment{deployment, fileLocation}
		addPodSpeccer(deploy)

		// TODO: Support older versions of Deployment as well?
		s.deployments = append(s.deployments, deploy)
	case appsv1beta1.SchemeGroupVersion.WithKind("Deployment"):
		var deployment appsv1beta1.Deployment
		errs.AddIfErr(decode(fileContents, &deployment))
		addPodSpeccer(internal.Appsv1beta1Deployment{deployment, fileLocation})
	case appsv1beta2.SchemeGroupVersion.WithKind("Deployment"):
		var deployment appsv1beta2.Deployment
		errs.AddIfErr(decode(fileContents, &deployment))
		addPodSpeccer(internal.Appsv1beta2Deployment{deployment, fileLocation})
	case extensionsv1beta1.SchemeGroupVersion.WithKind("Deployment"):
		var deployment extensionsv1beta1.Deployment
		errs.AddIfErr(decode(fileContents, &deployment))
		addPodSpeccer(internal.Extensionsv1beta1Deployment{deployment, fileLocation})

	case appsv1.SchemeGroupVersion.WithKind("StatefulSet"):
		var statefulSet appsv1.StatefulSet
		errs.AddIfErr(decode(fileContents, &statefulSet))
		sset := internal.Appsv1StatefulSet{statefulSet, fileLocation}
		addPodSpeccer(sset)

		// TODO: Support older versions of StatefulSet as well?
		s.statefulsets = append(s.statefulsets, sset)
	case appsv1beta1.SchemeGroupVersion.WithKind("StatefulSet"):
		var statefulSet appsv1beta1.StatefulSet
		errs.AddIfErr(decode(fileContents, &statefulSet))
		addPodSpeccer(internal.Appsv1beta1StatefulSet{statefulSet, fileLocation})
	case appsv1beta2.SchemeGroupVersion.WithKind("StatefulSet"):
		var statefulSet appsv1beta2.StatefulSet
		errs.AddIfErr(decode(fileContents, &statefulSet))
		addPodSpeccer(internal.Appsv1beta2StatefulSet{statefulSet, fileLocation})

	case appsv1.SchemeGroupVersion.WithKind("DaemonSet"):
		var daemonset appsv1.DaemonSet
		errs.AddIfErr(decode(fileContents, &daemonset))
		addPodSpeccer(internal.Appsv1DaemonSet{daemonset, fileLocation})
	case appsv1beta2.SchemeGroupVersion.WithKind("DaemonSet"):
		var daemonset appsv1beta2.DaemonSet
		errs.AddIfErr(decode(fileContents, &daemonset))
		addPodSpeccer(internal.Appsv1beta2DaemonSet{daemonset, fileLocation})
	case extensionsv1beta1.SchemeGroupVersion.WithKind("DaemonSet"):
		var daemonset extensionsv1beta1.DaemonSet
		errs.AddIfErr(decode(fileContents, &daemonset))
		addPodSpeccer(internal.Extensionsv1beta1DaemonSet{daemonset, fileLocation})

	case networkingv1.SchemeGroupVersion.WithKind("NetworkPolicy"):
		var netpol networkingv1.NetworkPolicy
		errs.AddIfErr(decode(fileContents, &netpol))
		np := internalnetpol.NetworkPolicy{netpol, fileLocation}
		s.networkPolicies = append(s.networkPolicies, np)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{netpol.TypeMeta, netpol.ObjectMeta, np})

	case corev1.SchemeGroupVersion.WithKind("Service"):
		var service corev1.Service
		errs.AddIfErr(decode(fileContents, &service))
		serv := internalservice.Service{service, fileLocation}
		s.services = append(s.services, serv)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{service.TypeMeta, service.ObjectMeta, serv})

	case policyv1beta1.SchemeGroupVersion.WithKind("PodDisruptionBudget"):
		var disruptBudget policyv1beta1.PodDisruptionBudget
		errs.AddIfErr(decode(fileContents, &disruptBudget))
		dbug := internalpdb.PodDisruptionBudgetV1beta1{disruptBudget, fileLocation}
		s.podDisruptionBudgets = append(s.podDisruptionBudgets, dbug)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{disruptBudget.TypeMeta, disruptBudget.ObjectMeta, dbug})
	case policyv1.SchemeGroupVersion.WithKind("PodDisruptionBudget"):
		var disruptBudget policyv1.PodDisruptionBudget
		errs.AddIfErr(decode(fileContents, &disruptBudget))
		dbug := internalpdb.PodDisruptionBudgetV1{disruptBudget, fileLocation}
		s.podDisruptionBudgets = append(s.podDisruptionBudgets, dbug)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{disruptBudget.TypeMeta, disruptBudget.ObjectMeta, dbug})

	case extensionsv1beta1.SchemeGroupVersion.WithKind("Ingress"):
		var ingress extensionsv1beta1.Ingress
		errs.AddIfErr(decode(fileContents, &ingress))
		ing := internal.ExtensionsIngressV1beta1{ingress, fileLocation}
		s.ingresses = append(s.ingresses, ing)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{ingress.TypeMeta, ingress.ObjectMeta, ing})

	case networkingv1beta1.SchemeGroupVersion.WithKind("Ingress"):
		var ingress networkingv1beta1.Ingress
		errs.AddIfErr(decode(fileContents, &ingress))
		ing := internal.IngressV1beta1{ingress, fileLocation}
		s.ingresses = append(s.ingresses, ing)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{ingress.TypeMeta, ingress.ObjectMeta, ing})

	case networkingv1.SchemeGroupVersion.WithKind("Ingress"):
		var ingress networkingv1.Ingress
		errs.AddIfErr(decode(fileContents, &ingress))
		ing := internal.IngressV1{ingress, fileLocation}
		s.ingresses = append(s.ingresses, ing)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{ingress.TypeMeta, ingress.ObjectMeta, ing})

	case autoscalingv1.SchemeGroupVersion.WithKind("HorizontalPodAutoscaler"):
		var hpa autoscalingv1.HorizontalPodAutoscaler
		errs.AddIfErr(decode(fileContents, &hpa))
		h := internal.HPAv1{hpa, fileLocation}
		s.hpaTargeters = append(s.hpaTargeters, h)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{hpa.TypeMeta, hpa.ObjectMeta, h})

	case autoscalingv2beta1.SchemeGroupVersion.WithKind("HorizontalPodAutoscaler"):
		var hpa autoscalingv2beta1.HorizontalPodAutoscaler
		errs.AddIfErr(decode(fileContents, &hpa))
		h := internal.HPAv2beta1{hpa, fileLocation}
		s.hpaTargeters = append(s.hpaTargeters, h)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{hpa.TypeMeta, hpa.ObjectMeta, h})

	case autoscalingv2beta2.SchemeGroupVersion.WithKind("HorizontalPodAutoscaler"):
		var hpa autoscalingv2beta2.HorizontalPodAutoscaler
		errs.AddIfErr(decode(fileContents, &hpa))
		h := internal.HPAv2beta2{hpa, fileLocation}
		s.hpaTargeters = append(s.hpaTargeters, h)
		s.bothMetas = append(s.bothMetas, ks.BothMeta{hpa.TypeMeta, hpa.ObjectMeta, h})

	default:
		if cnf.VerboseOutput > 1 {
			log.Printf("Unknown datatype: %s", detectedVersion.String())
		}
	}

	if errs.Any() {
		return errs
	}
	return nil
}
