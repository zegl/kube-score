package hpa

import (
	v1 "k8s.io/api/autoscaling/v1"

	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks, allTargetableObjs []domain.BothMeta) {
	allChecks.RegisterHorizontalPodAutoscalerCheck("HorizontalPodAutoscaler has target", `Makes sure that the HPA targets a valid object`, hpaHasTarget(allTargetableObjs))
}

func hpaHasTarget(allTargetableObjs []domain.BothMeta) func(v1.HorizontalPodAutoscaler) scorecard.TestScore {
	return func(hpa v1.HorizontalPodAutoscaler) (score scorecard.TestScore) {
		targetRef := hpa.Spec.ScaleTargetRef
		var hasTarget bool
		for _, t := range allTargetableObjs {
			if t.TypeMeta.APIVersion == targetRef.APIVersion &&
				t.TypeMeta.Kind == targetRef.Kind &&
				t.ObjectMeta.Name == targetRef.Name &&
				t.ObjectMeta.Namespace == hpa.Namespace {
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
