package score

import (
	"bytes"
	"github.com/zegl/kube-score/scorecard"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsbetav1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	extensionsbetav1.AddToScheme(scheme)
}

func Score(file io.Reader) (*scorecard.Scorecard, error) {
	allFiles, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	type detectKind struct {
		Kind string `yaml:"kind"`
	}

	var pods []corev1.Pod
	var deployments []appsv1.Deployment
	var statefulsets []appsv1.StatefulSet
	var networkPolies []networkingv1.NetworkPolicy

	for _, fileContents := range bytes.Split(allFiles, []byte("---\n")) {
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

		case "Deployment":
			var deployment appsv1.Deployment
			decode(fileContents, &deployment)
			deployments = append(deployments, deployment)

		case "StatefulSet":
			var statefulSet appsv1.StatefulSet
			decode(fileContents, &statefulSet)
			statefulsets = append(statefulsets, statefulSet)

		case "NetworkPolicy":
			var netpol networkingv1.NetworkPolicy
			decode(fileContents, &netpol)
			networkPolies = append(networkPolies, netpol)

		default:
			log.Printf("Unknown datatype: %s", detect.Kind)
		}
	}

	podTests := []func(corev1.PodTemplateSpec) scorecard.TestScore{
		scoreContainerLimits,
		scoreContainerImageTag,
		scoreContainerImagePullPolicy,
		scorePodHasNetworkPolicy(networkPolies),
		scoreContainerProbes,
	}

	scoreCard := scorecard.New()

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

	for _, deployment := range deployments {
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