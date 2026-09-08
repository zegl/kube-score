package score

import (
	"testing"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"

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

func TestHorizontalPodAutoscalerV2beta1TargetsDeployment(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "hpa-autoscalingv2beta1-targets-deployment.yaml", "HorizontalPodAutoscaler has target", scorecard.GradeAllOK)
}

func TestHorizontalPodAutoscalerV2beta2TargetsDeployment(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "hpa-autoscalingv2beta2-targets-deployment.yaml", "HorizontalPodAutoscaler has target", scorecard.GradeAllOK)
}

func TestHorizontalPodAutoscalerV2TargetsUnknownCRD(t *testing.T) {
	t.Parallel()
	// Argo Rollouts' "Rollout" (argoproj.io/v1alpha1) is a drop-in replacement for
	// apps/v1 Deployment that is commonly used as an HPA scaleTargetRef. kube-score
	// does not parse this (or any other) CRD kind, so it is never registered in
	// allTargetableObjs, and the HPA-target check reports a false-positive CRITICAL
	// "The HPA target does not match anything" even though the target object is
	// present in the same input. See issue #332.
	testExpectedScore(t, "hpa-targets-argo-rollout.yaml", "HorizontalPodAutoscaler has target", scorecard.GradeAllOK)
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
	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("hpa-min-replicas-nok.yaml")}, nil, &config.RunConfiguration{
		MinReplicasHPA: 2,
	}, "HorizontalPodAutoscaler Replicas", scorecard.GradeWarning)
}

func TestHorizontalPodAutoscalerMinReplicasMin5ReplicasNotOk(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("hpa-min-replicas-nok.yaml")}, nil, &config.RunConfiguration{
		MinReplicasHPA: 5,
	}, "HorizontalPodAutoscaler Replicas", scorecard.GradeWarning)
}

func TestHorizontalPodAutoscalerMinReplicasMin1ReplicaOk(t *testing.T) {
	t.Parallel()
	testExpectedScoreWithConfig(t, []ks.NamedReader{testFile("hpa-min-replicas-nok.yaml")}, nil, &config.RunConfiguration{
		MinReplicasHPA: 1,
	}, "HorizontalPodAutoscaler Replicas", scorecard.GradeAllOK)
}
