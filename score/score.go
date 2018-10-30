package score

import (
	"bytes"
	ks "github.com/zegl/kube-score"
	"github.com/zegl/kube-score/score/container"
	"github.com/zegl/kube-score/score/disruptionbudget"
	"github.com/zegl/kube-score/score/ingress"
	"github.com/zegl/kube-score/score/internal"
	"github.com/zegl/kube-score/score/networkpolicy"
	"github.com/zegl/kube-score/score/probes"
	"github.com/zegl/kube-score/score/security"
	"github.com/zegl/kube-score/score/service"
	"github.com/zegl/kube-score/score/stable"
	"github.com/zegl/kube-score/scorecard"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
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
	policyv1beta1.AddToScheme(scheme)
}

type Configuration struct {
	AllFiles                           []io.Reader
	VerboseOutput                      bool
	IgnoreContainerCpuLimitRequirement bool
}

var metaTests []func(metav1.TypeMeta) scorecard.TestScore
var podSpecTests []func(corev1.PodTemplateSpec) scorecard.TestScore
var serviceTests []func(corev1.Service) scorecard.TestScore

type score struct {
	config *Configuration

	typeMetas            []bothMeta
	pods                 []corev1.Pod
	podspecers           []ks.PodSpecer
	networkPolicies      []networkingv1.NetworkPolicy
	services             []corev1.Service
	podDisruptionBudgets []policyv1beta1.PodDisruptionBudget
	deployments          []appsv1.Deployment
	statefulsets         []appsv1.StatefulSet
	ingresses            []extensionsv1beta1.Ingress
}

type detectKind struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

type bothMeta struct {
	typeMeta   metav1.TypeMeta
	objectMeta metav1.ObjectMeta
}

// Score runs a pre-configured list of tests against the files defined in the configuration, and returns a scorecard.
// Additional configuration and tuning parameters can be provided via the config.
func Score(config Configuration) (*scorecard.Scorecard, error) {
	s := &score{
		config: &config,
	}

	addPodSpeccer := func(ps ks.PodSpecer) {
		s.podspecers = append(s.podspecers, ps)
		s.typeMetas = append(s.typeMetas, bothMeta{
			typeMeta:   ps.GetTypeMeta(),
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

		for _, fileContents := range bytes.Split(fullFile, []byte("\n---\n")) {
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
				s.pods = append(s.pods, pod)
				s.typeMetas = append(s.typeMetas, bothMeta{pod.TypeMeta, pod.ObjectMeta})

			case batchv1.SchemeGroupVersion.WithKind("Job"):
				var job batchv1.Job
				decode(fileContents, &job)
				addPodSpeccer(internal.Batchv1Job{job})

			case batchv1beta1.SchemeGroupVersion.WithKind("CronJob"):
				var cronjob batchv1beta1.CronJob
				decode(fileContents, &cronjob)
				addPodSpeccer(internal.Batchv1beta1CronJob{cronjob})

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
				s.typeMetas = append(s.typeMetas, bothMeta{netpol.TypeMeta, netpol.ObjectMeta})

			case corev1.SchemeGroupVersion.WithKind("Service"):
				var service corev1.Service
				decode(fileContents, &service)
				s.services = append(s.services, service)
				s.typeMetas = append(s.typeMetas, bothMeta{service.TypeMeta, service.ObjectMeta})

			case policyv1beta1.SchemeGroupVersion.WithKind("PodDisruptionBudget"):
				var disruptBudget policyv1beta1.PodDisruptionBudget
				decode(fileContents, &disruptBudget)
				s.podDisruptionBudgets = append(s.podDisruptionBudgets, disruptBudget)
				s.typeMetas = append(s.typeMetas, bothMeta{disruptBudget.TypeMeta, disruptBudget.ObjectMeta})

			case extensionsv1beta1.SchemeGroupVersion.WithKind("Ingress"):
				var ingress extensionsv1beta1.Ingress
				decode(fileContents, &ingress)
				s.ingresses = append(s.ingresses, ingress)
				s.typeMetas = append(s.typeMetas, bothMeta{ingress.TypeMeta, ingress.ObjectMeta})

			default:
				if config.VerboseOutput {
					log.Printf("Unknown datatype: %s", detect.Kind)
				}
			}
		}
	}

	return s.runTests()
}

func (s *score) runTests() (*scorecard.Scorecard, error) {
	metaTests := []func(metav1.TypeMeta) scorecard.TestScore{
		stable.ScoreMetaStableAvailable,
	}

	podTests := []func(corev1.PodTemplateSpec) scorecard.TestScore{
		container.ScoreContainerLimits(!s.config.IgnoreContainerCpuLimitRequirement),
		container.ScoreContainerImageTag,
		container.ScoreContainerImagePullPolicy,
		networkpolicy.ScorePodHasNetworkPolicy(s.networkPolicies),
		probes.ScoreContainerProbes(s.services),
		security.ScoreContainerSecurityContext,
	}

	serviceTests := []func(corev1.Service) scorecard.TestScore{
		service.ScoreServiceTargetsPod(s.pods, s.podspecers),
	}

	statefulSetTests := []func(appsv1.StatefulSet) scorecard.TestScore{
		disruptionbudget.ScoreStatefulSetHas(s.podDisruptionBudgets),
	}

	deploymentTests := []func(appsv1.Deployment) scorecard.TestScore{
		disruptionbudget.ScoreDeploymentHas(s.podDisruptionBudgets),
	}

	netpolTests := []func(networkingv1.NetworkPolicy) scorecard.TestScore{
		networkpolicy.ScoreNetworkPolicyTargetsPod(s.pods, s.podspecers),
	}

	ingressTests := []func(extensionsv1beta1.Ingress) scorecard.TestScore{
		ingress.ScoreIngressTargetsService(s.services),
	}

	scoreCard := scorecard.New()

	for _, meta := range s.typeMetas {
		for _, metaTest := range metaTests {
			score := metaTest(meta.typeMeta)
			score.AddMeta(meta.typeMeta, meta.objectMeta)
			scoreCard.Add(score)
		}
	}

	for _, pod := range s.pods {
		for _, podTest := range podTests {
			score := podTest(corev1.PodTemplateSpec{
				ObjectMeta: pod.ObjectMeta,
				Spec:       pod.Spec,
			})
			score.AddMeta(pod.TypeMeta, pod.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, podspecer := range s.podspecers {
		for _, podTest := range podTests {
			score := podTest(podspecer.GetPodTemplateSpec())
			score.AddMeta(podspecer.GetTypeMeta(), podspecer.GetObjectMeta())
			scoreCard.Add(score)
		}
	}

	for _, service := range s.services {
		for _, serviceTest := range serviceTests {
			score := serviceTest(service)
			score.AddMeta(service.TypeMeta, service.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, statefulset := range s.statefulsets {
		for _, test := range statefulSetTests {
			score := test(statefulset)
			score.AddMeta(statefulset.TypeMeta, statefulset.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, deployment := range s.deployments {
		for _, test := range deploymentTests {
			score := test(deployment)
			score.AddMeta(deployment.TypeMeta, deployment.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, netpol := range s.networkPolicies {
		for _, netpolTest := range netpolTests {
			score := netpolTest(netpol)
			score.AddMeta(netpol.TypeMeta, netpol.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	for _, ingress := range s.ingresses {
		for _, ingressTest := range ingressTests {
			score := ingressTest(ingress)
			score.AddMeta(ingress.TypeMeta, ingress.ObjectMeta)
			scoreCard.Add(score)
		}
	}

	return scoreCard, nil
}
