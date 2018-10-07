package score

import (
	"bytes"
	"io"
	"io/ioutil"
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

func Score(files []io.Reader) (*scorecard.Scorecard, error) {
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

	for _, file := range files {
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

			switch detect.Kind {
			case "Pod":
				var pod corev1.Pod
				decode(fileContents, &pod)
				pods = append(pods, pod)
				typeMetas = append(typeMetas, bothMeta{pod.TypeMeta, pod.ObjectMeta})

			case "Job":
				fallthrough
			case "CronJob":
				fallthrough
			case "Deployment":
				fallthrough
			case "DaemonSet":
				fallthrough
			case "StatefulSet":
				var podspecer PodSpecer

				kindAndVersion := detect.Kind + "-" + detect.ApiVersion

				switch kindAndVersion {
				case "Deployment-apps/v1":
					var deployment appsv1.Deployment
					decode(fileContents, &deployment)
					podspecer = appsv1Deployment{deployment}
				case "Deployment-apps/v1beta1":
					var deployment appsv1beta1.Deployment
					decode(fileContents, &deployment)
					podspecer = appsv1beta1Deployment{deployment}
				case "Deployment-apps/v1beta2":
					var deployment appsv1beta2.Deployment
					decode(fileContents, &deployment)
					podspecer = appsv1beta2Deployment{deployment}
				case "Deployment-extensions/v1beta1":
					var deployment extensionsv1beta1.Deployment
					decode(fileContents, &deployment)
					podspecer = extensionsv1beta1Deployment{deployment}

				case "StatefulSet-apps/v1":
					var statefulSet appsv1.StatefulSet
					decode(fileContents, &statefulSet)
					podspecer = appsv1StatefulSet{statefulSet}
				case "StatefulSet-apps/v1beta1":
					var statefulSet appsv1beta1.StatefulSet
					decode(fileContents, &statefulSet)
					podspecer = appsv1beta1StatefulSet{statefulSet}
				case "StatefulSet-apps/v1beta2":
					var statefulSet appsv1beta2.StatefulSet
					decode(fileContents, &statefulSet)
					podspecer = appsv1beta2StatefulSet{statefulSet}

				case "DaemonSet-apps/v1":
					var daemonset appsv1.DaemonSet
					decode(fileContents, &daemonset)
					podspecer = appsv1DaemonSet{daemonset}
				case "DaemonSet-apps/v1beta2":
					var daemonset appsv1beta2.DaemonSet
					decode(fileContents, &daemonset)
					podspecer = appsv1beta2DaemonSet{daemonset}
				case "DaemonSet-extensions/v1beta1":
					var daemonset extensionsv1beta1.DaemonSet
					decode(fileContents, &daemonset)
					podspecer = extensionsv1beta1DaemonSet{daemonset}

				case "Job-batch/v1":
					var job batchv1.Job
					decode(fileContents, &job)
					podspecer = batchv1Job{job}

				case "CronJob-batch/v1beta1":
					var cronjob batchv1beta1.CronJob
					decode(fileContents, &cronjob)
					podspecer = batchv1beta1CronJob{cronjob}

				default:
					log.Printf("Unknown type %s %s", detect.ApiVersion, detect.Kind)
					continue
				}

				podspecers = append(podspecers, podspecer)
				typeMetas = append(typeMetas, bothMeta{
					podspecer.GetTypeMeta(),
					podspecer.GetObjectMeta(),
				})

			case "NetworkPolicy":
				var netpol networkingv1.NetworkPolicy
				decode(fileContents, &netpol)
				networkPolies = append(networkPolies, netpol)
				typeMetas = append(typeMetas, bothMeta{netpol.TypeMeta, netpol.ObjectMeta})

			case "Service":
				var service corev1.Service
				decode(fileContents, &service)
				services = append(services, service)
				typeMetas = append(typeMetas, bothMeta{service.TypeMeta, service.ObjectMeta})

			default:
				log.Printf("Unknown datatype: %s", detect.Kind)
			}
		}
	}

	metaTests := []func(metav1.TypeMeta) scorecard.TestScore{
		scoreMetaStableAvailable,
	}

	podTests := []func(corev1.PodTemplateSpec) scorecard.TestScore{
		scoreContainerLimits,
		scoreContainerImageTag,
		scoreContainerImagePullPolicy,
		scorePodHasNetworkPolicy(networkPolies),
		scoreContainerProbes(services),
		scoreContainerSecurityContext,
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

	return scoreCard, nil
}
