package hpa

import (
	"fmt"
	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
	"k8s.io/utils/ptr"
)

func Register(allChecks *checks.Checks, allTargetableObjs []domain.BothMeta, minReplicas int) {
	allChecks.RegisterHorizontalPodAutoscalerCheck("HorizontalPodAutoscaler has target", `Makes sure that the HPA targets a valid object`, hpaHasTarget(allTargetableObjs))
	allChecks.RegisterHorizontalPodAutoscalerCheck("HorizontalPodAutoscaler Replicas", `Makes sure that the HPA has multiple replicas`, hpaHasMultipleReplicas(minReplicas))
}

func hpaHasTarget(allTargetableObjs []domain.BothMeta) func(hpa domain.HpaTargeter) (scorecard.TestScore, error) {
	return func(hpa domain.HpaTargeter) (score scorecard.TestScore, err error) {
		targetRef := hpa.HpaTarget()
		var hasTarget bool
		for _, t := range allTargetableObjs {
			if t.TypeMeta.APIVersion == targetRef.APIVersion &&
				t.TypeMeta.Kind == targetRef.Kind &&
				t.ObjectMeta.Name == targetRef.Name &&
				t.ObjectMeta.Namespace == hpa.GetObjectMeta().Namespace {
				hasTarget = true
				break
			}
		}

		if hasTarget {
			score.Grade = scorecard.GradeAllOK
		} else {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "The HPA target does not match anything", "")
		}
		return
	}
}

func hpaHasMultipleReplicas(minReplicas int) func(hpa domain.HpaTargeter) (scorecard.TestScore, error) {
	return func(hpa domain.HpaTargeter) (score scorecard.TestScore, err error) {
		if ptr.Deref(hpa.MinReplicas(), 1) >= int32(minReplicas) {
			score.Grade = scorecard.GradeAllOK
		} else {
			score.Grade = scorecard.GradeWarning
			score.AddComment("", "HPA few replicas", fmt.Sprintf("HorizontalPodAutoscalers are recommended to have at least %d replicas to prevent unwanted downtime.", minReplicas))
		}
		return
	}
}
