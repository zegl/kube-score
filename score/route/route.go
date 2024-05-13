package route

import (
	routev1 "github.com/openshift/api/route/v1"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks, services ks.Services) {
	allChecks.RegisterRouteCheck("Route targets Service", `Makes sure that the Route targets a Service`, routeTargetsService(services.Services()))
}

// routeTargetsService checks if a Service targets a pod and issues a critical warning if no matching pod
func routeTargetsService(svcs []ks.Service) func(routev1.Route) (scorecard.TestScore, error) {
	return func(route routev1.Route) (score scorecard.TestScore, err error) {
		hasMatchService := false
		hasMatchPort := false

		if route.Spec.Port != nil && route.Spec.Port.TargetPort.IntValue() != 0 {
			// We consider this as a match as this now matches the pod port which is very difficult to determine
			hasMatchPort = true
			hasMatchService = true
		} else {
			for _, s := range svcs {
				svc := s.Service()
				if svc.Namespace != route.Namespace {
					break
				}
				if route.Spec.To.Name == svc.Name {
					hasMatchService = true
					if route.Spec.Port == nil && len(svc.Spec.Ports) == 1 {
						hasMatchPort = true
						break
					}
					if route.Spec.Port != nil && route.Spec.Port.TargetPort.String() != "" {
						for _, p := range svc.Spec.Ports {
							if route.Spec.Port.TargetPort.String() == p.Name {
								hasMatchPort = true
								break
							}
						}
					}
				}
			}
		}

		if hasMatchService && hasMatchPort {
			score.Grade = scorecard.GradeAllOK
		} else if hasMatchService && !hasMatchPort {
			score.Grade = scorecard.GradeAlmostOK
			score.AddComment("", "The route does not match any port on the service", "")
		} else {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "The route does not reference a service", "")
		}

		return
	}
}
