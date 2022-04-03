package hpa

import (
	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks, allTargetableObjs []domain.BothMeta) {
	allChecks.RegisterHorizontalPodAutoscalerCheck("HorizontalPodAutoscaler has target", `Makes sure that the HPA targets a valid object`, hpaHasTarget(allTargetableObjs))
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
