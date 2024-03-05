package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestHorizontalPodAutoscalerV1TargetsDeployment(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "hpa-autoscalingv1-targets-deployment.yaml", "HorizontalPodAutoscaler has target", scorecard.GradeAllOK)
}

func TestHorizontalPodAutoscalerV2TargetsDeployment(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "hpa-autoscalingv2-targets-deployment.yaml", "HorizontalPodAutoscaler has target", scorecard.GradeAllOK)
}

func TestHorizontalPodAutoscalerHasNoTarget(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "hpa-has-no-target.yaml", "HorizontalPodAutoscaler has target", scorecard.GradeCritical)
}

func TestHorizontalPodAutoscalerMinReplicasOk(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "hpa-min-replicas-ok.yaml", "HorizontalPodAutoscaler Replicas", scorecard.GradeAllOK)
}

func TestHorizontalPodAutoscalerMinReplicasNok(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "hpa-min-replicas-nok.yaml", "HorizontalPodAutoscaler Replicas", scorecard.GradeWarning)
}
