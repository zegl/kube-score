package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestIngressTargetsService(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "ingress-targets-service.yaml", "Ingress targets Service", scorecard.GradeAllOK)
}

func TestIngressTargetsServiceNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "ingress-targets-service-no-match.yaml", "Ingress targets Service", scorecard.GradeCritical)
}

func TestNetworkingIngressTargetsService(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "ingress-networkingv1beta1-targets-service.yaml", "Ingress targets Service", scorecard.GradeAllOK)
}

func TestNetworkingIngressTargetsServiceNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "ingress-networkingv1beta1-targets-service-no-match.yaml", "Ingress targets Service", scorecard.GradeCritical)
}
