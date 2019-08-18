package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestHorizontalPodAutoscalerTargetsDeployment(t *testing.T) {
	testExpectedScore(t, "hpa-targets-deployment.yaml", "HorizontalPodAutoscaler has target", scorecard.GradeAllOK)
}

func TestHorizontalPodAutoscalerHasNoTarget(t *testing.T) {
	testExpectedScore(t, "hpa-has-no-target.yaml", "HorizontalPodAutoscaler has target", scorecard.GradeCritical)
}
