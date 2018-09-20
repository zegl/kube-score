package score

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"

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
}

type Scorecard struct {
	Scores []TestScore
}

type TestScore struct {
	Name        string
	Description string
	Grade       int
	Comments    []string
}

func Score(file io.Reader) (*Scorecard, error) {
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
			log.Panicf("Unknown datatype: %s", detect.Kind)
		}
	}

	podTests := []func(corev1.PodTemplateSpec) TestScore{
		scoreContainerLimits,
		scoreContainerImageTag,
		scoreContainerImagePullPolicy,
		scorePodHasNetworkPolicy(networkPolies),
	}

	scoreCard := Scorecard{}

	for _, pod := range pods {
		for _, podTest := range podTests {
			scoreCard.Scores = append(scoreCard.Scores, podTest(corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec: pod.Spec,
			}))
		}
	}

	for _, deployment := range deployments {
		for _, podTest := range podTests {
			scoreCard.Scores = append(scoreCard.Scores, podTest(deployment.Spec.Template))
		}
	}

	for _, statefulset := range statefulsets {
		for _, podTest := range podTests {
			scoreCard.Scores = append(scoreCard.Scores, podTest(statefulset.Spec.Template))
		}
	}

	return &scoreCard, nil
}

func scoreContainerLimits(podTemplate corev1.PodTemplateSpec) (score TestScore) {
	score.Name = "Container Resources"

	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	hasMissingLimit := false
	hasMissingRequest := false

	for _, container := range allContainers {
		if container.Resources.Limits.Cpu().IsZero() {
			score.Comments = append(score.Comments, "CPU limit is not set")
			hasMissingLimit = true
		}
		if container.Resources.Limits.Memory().IsZero() {
			score.Comments = append(score.Comments, "Memory limit is not set")
			hasMissingLimit = true
		}
		if container.Resources.Requests.Cpu().IsZero() {
			score.Comments = append(score.Comments, "CPU request is not set")
			hasMissingRequest = true
		}
		if container.Resources.Requests.Memory().IsZero() {
			score.Comments = append(score.Comments, "Memory request is not set")
			hasMissingRequest = true
		}
	}

	if len(allContainers) == 0 {
		score.Grade = 0
		score.Comments = append(score.Comments, "No containers defined")
	} else if hasMissingLimit {
		score.Grade = 0
	} else if hasMissingRequest {
		score.Grade = 5
	} else {
		score.Grade = 10
	}

	return
}

func scoreContainerImageTag(podTemplate corev1.PodTemplateSpec) (score TestScore) {
	score.Name = "Container Image Tag"

	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	hasTagLatest := false

	for _, container := range allContainers{
		imageParts := strings.Split(container.Image, ":")
		imageVersion := imageParts[len(imageParts)-1]

		if imageVersion == "latest" {
			score.Comments = append(score.Comments, "Image with latest tag")
			hasTagLatest = true
		}
	}

	if hasTagLatest {
		score.Grade = 0
	} else {
		score.Grade = 10
	}

	return
}

func scoreContainerImagePullPolicy(podTemplate corev1.PodTemplateSpec) (score TestScore) {
	score.Name = "Container Image Pull Policy"

	pod := podTemplate.Spec

	allContainers := pod.InitContainers
	allContainers = append(allContainers, pod.Containers...)

	hasNonAlways := false

	for _, container := range allContainers{
		if container.ImagePullPolicy != corev1.PullAlways {
			score.Comments = append(score.Comments, "ImagePullPolicy is not set to PullAlways")
			hasNonAlways = true
		}
	}

	if hasNonAlways {
		score.Grade = 0
	} else {
		score.Grade = 10
	}

	return
}