package ingress

import (
	"fmt"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
)

func Register(allChecks *checks.Checks, services ks.Services) {
	allChecks.RegisterIngressCheck("Ingress targets Service", `Makes sure that the Ingress targets a Service`, ingressTargetsService(services.Services()))
}

func ingressTargetsService(allServices []corev1.Service) func(extensionsv1beta1.Ingress) scorecard.TestScore {
	return func(ingress extensionsv1beta1.Ingress) (score scorecard.TestScore) {
		allRulesHaveMatches := true

		for _, rule := range ingress.Spec.Rules {
			for _, path := range rule.IngressRuleValue.HTTP.Paths {

				pathHasMatch := false

				for _, service := range allServices {
					if service.Namespace != ingress.Namespace {
						continue
					}

					if service.Name == path.Backend.ServiceName {
						for _, servicePort := range service.Spec.Ports {
							if servicePort.Port == path.Backend.ServicePort.IntVal {
								pathHasMatch = true
							}
						}
					}
				}

				if !pathHasMatch {
					allRulesHaveMatches = false
					score.AddComment(path.Path, "No service match was found", fmt.Sprintf("No service with name %s and port %d was found", path.Backend.ServiceName, path.Backend.ServicePort.IntVal))
				}
			}

		}

		if allRulesHaveMatches {
			score.Grade = scorecard.GradeAllOK
		} else {
			score.Grade = scorecard.GradeCritical
		}

		return
	}
}
