package score

import (
	"bytes"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"

	"github.com/zegl/kube-score/scorecard"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
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
}

type PodSpecer interface {
	GetTypeMeta() metav1.TypeMeta
	GetObjectMeta() metav1.ObjectMeta
	GetPodTemplateSpec() corev1.PodTemplateSpec
}

type Configuration struct {
	AllFiles []io.Reader
	VerboseOutput bool

	IgnoreContainerCpuLimitRequirement bool
}

func Score(config Configuration) (*scorecard.Scorecard, error) {
	type detectKind struct {
		ApiVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
	}

	type bothMeta struct {
		typeMeta   metav1.TypeMeta
		objectMeta metav1.ObjectMeta
	}

	var typeMetas []bothMeta
	var pods []corev1.Pod
	var podspecers []PodSpecer
	var networkPolies []networkingv1.NetworkPolicy
	var services []corev1.Service

	addPodSpeccer := func(ps PodSpecer) {
		podspecers = append(podspecers, ps)
		typeMetas = append(typeMetas, bothMeta{
			typeMeta: ps.GetTypeMeta(),
			objectMeta: ps.GetObjectMeta(),
		})
	}

	for _, file := range config.AllFiles {
		fullFile, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		// Convert to unix style newlines
		fullFile = bytes.Replace(fullFile, []byte("\r\n"), []byte("\n"), -1)

		for _, fileContents := range bytes.Split(fullFile, []byte("---\n")) {
			var detect detectKind
			err = yaml.Unmarshal(fileContents, &detect)
			if err != nil {
				return nil, err
			}

			decode := func(data []byte, object runtime.Object) {
				deserializer := codecs.UniversalDeserializer()
				if _, _, err := deserializer.Decode(data, nil, object); err != nil {
					panic(err)
				}
			}

			detectedVersion := schema.FromAPIVersionAndKind(detect.ApiVersion, detect.Kind)

			switch detectedVersion {
			case corev1.SchemeGroupVersion.WithKind("Pod"):
				var pod corev1.Pod
				decode(fileContents, &pod)
				pods = append(pods, pod)
				typeMetas = append(typeMetas, bothMeta{pod.TypeMeta, pod.ObjectMeta})

			case batchv1.SchemeGroupVersion.WithKind("Job"):
				var job batchv1.Job
				decode(fileContents, &job)
				addPodSpeccer(batchv1Job{job})

			case batchv1beta1.SchemeGroupVersion.WithKind("CronJob"):
				var cronjob batchv1beta1.CronJob
				decode(fileContents, &cronjob)
				addPodSpeccer(batchv1beta1CronJob{cronjob})

			case appsv1.SchemeGroupVersion.WithKind("Deployment"):
				var deployment appsv1.Deployment
				decode(fileContents, &deployment)
				addPodSpeccer(appsv1Deployment{deployment})
			case appsv1beta1.SchemeGroupVersion.WithKind("Deployment"):
				var deployment appsv1beta1.Deployment
				decode(fileContents, &deployment)
				addPodSpeccer(appsv1beta1Deployment{deployment})
			case appsv1beta2.SchemeGroupVersion.WithKind("Deployment"):
				var deployment appsv1beta2.Deployment
				decode(fileContents, &deployment)
				addPodSpeccer(appsv1beta2Deployment{deployment})
			case extensionsv1beta1.SchemeGroupVersion.WithKind("Deployment"):
				var deployment extensionsv1beta1.Deployment
				decode(fileContents, &deployment)
				addPodSpeccer(extensionsv1beta1Deployment{deployment})

			case appsv1.SchemeGroupVersion.WithKind("StatefulSet"):
				var statefulSet appsv1.StatefulSet
				decode(fileContents, &statefulSet)
				addPodSpeccer(appsv1StatefulSet{statefulSet})
			case appsv1beta1.SchemeGroupVersion.WithKind("StatefulSet"):
				var statefulSet appsv1beta1.StatefulSet
				decode(fileContents, &statefulSet)
				addPodSpeccer(appsv1beta1StatefulSet{statefulSet})
			case appsv1beta2.SchemeGroupVersion.WithKind("StatefulSet"):
				var statefulSet appsv1beta2.StatefulSet
				decode(fileContents, &statefulSet)
				addPodSpeccer(appsv1beta2StatefulSet{statefulSet})

			case appsv1.SchemeGroupVersion.WithKind("DaemonSet"):
				var daemonset appsv1.DaemonSet
				decode(fileContents, &daemonset)
				addPodSpeccer(appsv1DaemonSet{daemonset})
			case appsv1beta2.SchemeGroupVersion.WithKind("DaemonSet"):
				var daemonset appsv1beta2.DaemonSet
				decode(fileContents, &daemonset)
				addPodSpeccer(appsv1beta2DaemonSet{daemonset})
			case extensionsv1beta1.SchemeGroupVersion.WithKind("DaemonSet"):
				var daemonset extensionsv1beta1.DaemonSet
				decode(fileContents, &daemonset)
				addPodSpeccer(extensionsv1beta1DaemonSet{daemonset})

			case networkingv1.SchemeGroupVersion.WithKind("NetworkPolicy"):
				var netpol networkingv1.NetworkPolicy
				decode(fileContents, &netpol)
				networkPolies = append(networkPolies, netpol)
				typeMetas = append(typeMetas, bothMeta{netpol.TypeMeta, netpol.ObjectMeta})

			case corev1.SchemeGroupVersion.WithKind("Service"):
				var service corev1.Service
				decode(fileContents, &service)
				services = append(services, service)
				typeMetas = append(typeMetas, bothMeta{service.TypeMeta, service.ObjectMeta})

			default:
				if config.VerboseOutput {
					log.Printf("Unknown datatype: %s", detect.Kind)
				}
			}
		}
	}

	metaTests := []func(metav1.TypeMeta) scorecard.TestScore{
		scoreMetaStableAvailable,
	}

	podTests := []func(corev1.PodTemplateSpec) scorecard.TestScore{
		scoreContainerLimits(!config.IgnoreContainerCpuLimitRequirement),
		scoreContainerImageTag,
		scoreContainerImagePullPolicy,
		scorePodHasNetworkPolicy(networkPolies),
		scoreContainerProbes(services),
		scoreContainerSecurityContext,
	}

	serviceTests := []func(corev1.Service) scorecard.TestScore{
		scoreServiceTargetsPod(pods, podspecers),
	}

	scoreCard := scorecard.New()

	for _, meta := range typeMetas {
		for _, metaTest := range metaTests {
			score := metaTest(meta.typeMeta)
			score.AddMeta(meta.typeMeta, meta.objectMeta)
			scoreCard.Add(score)
		}
	}

	for _, pod := range pods {
		for _, podTest := range podTests {
			score := podTest(corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			})
			score.AddMeta(pod.TypeMeta, pod.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, podspecer := range podspecers {
		for _, podTest := range podTests {
			score := podTest(podspecer.GetPodTemplateSpec())
			score.AddMeta(podspecer.GetTypeMeta(), podspecer.GetObjectMeta())
			scoreCard.Add(score)
		}
	}

	for _, service := range services {
		for _, serviceTest := range serviceTests {
			score := serviceTest(service)
			score.AddMeta(service.TypeMeta, service.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	return scoreCard, nil
}
