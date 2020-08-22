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

func TestNetworkingIngressV1beta1TargetsService(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "ingress-networkingv1beta1-targets-service.yaml", "Ingress targets Service", scorecard.GradeAllOK)
}

func TestNetworkingIngressV1beta1TargetsServiceNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "ingress-networkingv1beta1-targets-service-no-match.yaml", "Ingress targets Service", scorecard.GradeCritical)
}

func TestNetworkingIngressV1TargetsService(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "ingress-networkingv1-targets-service.yaml", "Ingress targets Service", scorecard.GradeAllOK)
}

func TestNetworkingIngressV1TargetsServiceNoMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "ingress-networkingv1-targets-service-no-match.yaml", "Ingress targets Service", scorecard.GradeCritical)
}
