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
	var appsv1Deployment []appsv1.Deployment
	var appsv1beta1Deployment []appsv1beta1.Deployment
	var appsv1beta2Deployment []appsv1beta2.Deployment
	var extensionsv1beta1Deployment []extensionsv1beta1.Deployment
	var statefulsets []appsv1.StatefulSet
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
				switch detect.ApiVersion {
				case "apps/v1":
					var deployment appsv1.Deployment
					decode(fileContents, &deployment)
					appsv1Deployment = append(appsv1Deployment, deployment)
					typeMetas = append(typeMetas,  bothMeta{deployment.TypeMeta, deployment.ObjectMeta})
				case "apps/v1beta1":
					var deployment appsv1beta1.Deployment
					decode(fileContents, &deployment)
					appsv1beta1Deployment = append(appsv1beta1Deployment, deployment)
					typeMetas = append(typeMetas,  bothMeta{deployment.TypeMeta, deployment.ObjectMeta})
				case "apps/v1beta2":
					var deployment appsv1beta2.Deployment
					decode(fileContents, &deployment)
					appsv1beta2Deployment = append(appsv1beta2Deployment, deployment)
					typeMetas = append(typeMetas,  bothMeta{deployment.TypeMeta, deployment.ObjectMeta})
				case "extensions/v1beta1":
					var deployment extensionsv1beta1.Deployment
					decode(fileContents, &deployment)
					extensionsv1beta1Deployment = append(extensionsv1beta1Deployment, deployment)
					typeMetas = append(typeMetas,  bothMeta{deployment.TypeMeta, deployment.ObjectMeta})
				default:
					log.Printf("Unknown type version of Deployment: %s", detect.ApiVersion)
				}


			case "StatefulSet":
				var statefulSet appsv1.StatefulSet
				decode(fileContents, &statefulSet)
				statefulsets = append(statefulsets, statefulSet)
				typeMetas = append(typeMetas,  bothMeta{statefulSet.TypeMeta, statefulSet.ObjectMeta})

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

	for _, deployment := range appsv1Deployment {
		for _, podTest := range podTests {
			score := podTest(deployment.Spec.Template)
			score.AddMeta(deployment.TypeMeta, deployment.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, deployment := range appsv1beta1Deployment {
		for _, podTest := range podTests {
			score := podTest(deployment.Spec.Template)
			score.AddMeta(deployment.TypeMeta, deployment.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, deployment := range appsv1beta2Deployment {
		for _, podTest := range podTests {
			score := podTest(deployment.Spec.Template)
			score.AddMeta(deployment.TypeMeta, deployment.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, deployment := range extensionsv1beta1Deployment {
		for _, podTest := range podTests {
			score := podTest(deployment.Spec.Template)
			score.AddMeta(deployment.TypeMeta, deployment.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, statefulset := range statefulsets {
		for _, podTest := range podTests {
			score := podTest(statefulset.Spec.Template)
			score.AddMeta(statefulset.TypeMeta, statefulset.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	return scoreCard, nil
}