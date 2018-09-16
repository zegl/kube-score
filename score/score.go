package score

import (
	"io"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	//"errors"
	//"fmt"
	// "github.com/labstack/echo"
	// "k8s.io/api/admission/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	// batchv1 "k8s.io/api/batch/v1"
	// batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	// "reflect"
)

var scheme = runtime.NewScheme()
var codecs = serializer.NewCodecFactory(scheme)

func init() {
	addToScheme(scheme)
}

func addToScheme(scheme *runtime.Scheme) {
	corev1.AddToScheme(scheme)
	appsv1.AddToScheme(scheme)
	// batchv1.AddToScheme(scheme)
	// batchv1beta1.AddToScheme(scheme)
	// v1beta1.AddToScheme(scheme)
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
	allData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	type detectKind struct {
		Kind string `yaml:"kind"`
	}

	var detect detectKind
	err = yaml.Unmarshal(allData, &detect)
	if err != nil {
		return nil, err
	}

	var pods []corev1.Pod
	var deployments []appsv1.Deployment

	decode := func(data []byte, object runtime.Object) {
		deserializer := codecs.UniversalDeserializer()
		if _, _, err := deserializer.Decode(data, nil, object); err != nil {
			panic(err)
		}
	}

	switch detect.Kind {
	case "Pod":
		var pod corev1.Pod
		decode(allData, &pod)
		pods = append(pods, pod)

	case "Deployment":
		var deployment appsv1.Deployment
		decode(allData, &deployment)
		deployments = append(deployments, deployment)
	}

	podTests := []func(corev1.PodSpec) TestScore{
		scoreContainerLimits,
	}

	scoreCard := Scorecard{}

	for _, pod := range pods {
		for _, podTest := range podTests {
			scoreCard.Scores = append(scoreCard.Scores, podTest(pod.Spec))
		}
	}

	for _, deployment := range deployments {
		for _, podTest := range podTests {
			scoreCard.Scores = append(scoreCard.Scores, podTest(deployment.Spec.Template.Spec))
		}
	}

	return &scoreCard, nil
}

func scoreContainerLimits(pod corev1.PodSpec) (score TestScore) {
	score.Name = "Container Resources"

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
