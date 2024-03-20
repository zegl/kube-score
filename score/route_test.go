package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestRouteTargetsService(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "route-targets-service.yaml", "Route targets Service", scorecard.GradeAllOK)
}

func TestRouteTargetsServiceNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "route-targets-service-no-match.yaml", "Route targets Service", scorecard.GradeCritical)
}

func TestRouteTargetsServiceNumberedPort(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "route-targets-service-numbered-port.yaml", "Route targets Service", scorecard.GradeAllOK)
}

func TestRouteTargetsServiceNoPortOk(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "route-targets-service-no-port-ok.yaml", "Route targets Service", scorecard.GradeAllOK)
}

func TestRouteTargetsServiceNoPortNok(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "route-targets-service-no-port-nok.yaml", "Route targets Service", scorecard.GradeAlmostOK)
}
