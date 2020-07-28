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
