package score

import (
	"bytes"
	"io"
	"log"
	"io/ioutil"

	"github.com/zegl/kube-score/scorecard"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
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
}

func Score(files []io.Reader) (*scorecard.Scorecard, error) {
	type detectKind struct {
		ApiVersion string `yaml:"apiVersion"`
		Kind string `yaml:"kind"`
	}

	type bothMeta struct {
		typeMeta metav1.TypeMeta
		objectMeta metav1.ObjectMeta
	}

	var typeMetas []bothMeta
	var pods []corev1.Pod
	var genericDeployments []Deployment
	var genericStatefulsets []StatefulSet
	var networkPolies []networkingv1.NetworkPolicy

	for _, file := range files {
		fullFile, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

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

			case "Deployment":

				var genericDeployment Deployment

				switch detect.ApiVersion {
				case "apps/v1":
					var deployment appsv1.Deployment
					decode(fileContents, &deployment)
					genericDeployment = appsv1Deployment{deployment}
				case "apps/v1beta1":
					var deployment appsv1beta1.Deployment
					decode(fileContents, &deployment)
					genericDeployment = appsv1beta1Deployment{deployment}
				case "apps/v1beta2":
					var deployment appsv1beta2.Deployment
					decode(fileContents, &deployment)
					genericDeployment = appsv1beta2Deployment{deployment}
				case "extensions/v1beta1":
					var deployment extensionsv1beta1.Deployment
					decode(fileContents, &deployment)
					genericDeployment = extensionsv1beta1Deployment{deployment}
				default:
					log.Printf("Unknown type version of Deployment: %s", detect.ApiVersion)
				}

				genericDeployments = append(genericDeployments, genericDeployment)
				typeMetas = append(typeMetas,  bothMeta{
					genericDeployment.GetTypeMeta(),
					genericDeployment.GetObjectMeta(),
				})

			case "StatefulSet":
				var genericStatefulset StatefulSet

				switch detect.ApiVersion {
				case "apps/v1":
					var statefulSet appsv1.StatefulSet
					decode(fileContents, &statefulSet)
					genericStatefulset = appsv1StatefulSet{statefulSet}
				case "apps/v1beta1":
					var statefulSet appsv1beta1.StatefulSet
					decode(fileContents, &statefulSet)
					genericStatefulset = appsv1beta1StatefulSet{statefulSet}
				case "apps/v1beta2":
					var statefulSet appsv1beta2.StatefulSet
					decode(fileContents, &statefulSet)
					genericStatefulset = appsv1beta2StatefulSet{statefulSet}
				}

				genericStatefulsets = append(genericStatefulsets, genericStatefulset)
				typeMetas = append(typeMetas,  bothMeta{
					genericStatefulset.GetTypeMeta(),
					genericStatefulset.GetObjectMeta(),
				})

			case "NetworkPolicy":
				var netpol networkingv1.NetworkPolicy
				decode(fileContents, &netpol)
				networkPolies = append(networkPolies, netpol)
				typeMetas = append(typeMetas,  bothMeta{netpol.TypeMeta, netpol.ObjectMeta})

			default:
				log.Printf("Unknown datatype: %s", detect.Kind)
			}
		}
	}

	metaTests := []func(metav1.TypeMeta) scorecard.TestScore {
		scoreMetaStableAvailable,
	}

	podTests := []func(corev1.PodTemplateSpec) scorecard.TestScore{
		scoreContainerLimits,
		scoreContainerImageTag,
		scoreContainerImagePullPolicy,
		scorePodHasNetworkPolicy(networkPolies),
		scoreContainerProbes,
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
				Spec: pod.Spec,
			})
			score.AddMeta(pod.TypeMeta, pod.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, deployment := range genericDeployments {
		for _, podTest := range podTests {
			score := podTest(deployment.GetPodTemplateSpec())
			score.AddMeta(deployment.GetTypeMeta(), deployment.GetObjectMeta())
			scoreCard.Add(score)
		}
	}

	for _, statefulset := range genericStatefulsets {
		for _, podTest := range podTests {
			score := podTest(statefulset.GetPodTemplateSpec())
			score.AddMeta(statefulset.GetTypeMeta(), statefulset.GetObjectMeta())
			scoreCard.Add(score)
		}
	}

	return scoreCard, nil
}