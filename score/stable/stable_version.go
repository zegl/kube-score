package stable

import (
	"fmt"

	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks) {
	allChecks.RegisterMetaCheck("Stable version", `Checks if the object is using a deprecated apiVersion`, metaStableAvailable)
}

// ScoreMetaStableAvailable checks if the supplied TypeMeta is an unstable object type, that has a stable(r) replacement
func metaStableAvailable(meta domain.BothMeta) (score scorecard.TestScore) {
	withStable := map[string]map[string]string{
		"extensions/v1beta1": {
			"Deployment": "apps/v1",
			"DaemonSet":  "apps/v1",
		},
		"apps/v1beta1": {
			"Deployment":  "apps/v1",
			"StatefulSet": "apps/v1",
		},
		"apps/v1beta2": {
			"Deployment":  "apps/v1",
			"StatefulSet": "apps/v1",
			"DaemonSet":   "apps/v1",
		},
	}

	if inVersion, ok := withStable[meta.TypeMeta.APIVersion]; ok {
		if recommendedVersion, ok := inVersion[meta.TypeMeta.Kind]; ok {
			score.Grade = scorecard.GradeWarning
			score.AddComment("",
				fmt.Sprintf("The apiVersion and kind %s/%s is deprecated", meta.TypeMeta.APIVersion, meta.TypeMeta.Kind),
				fmt.Sprintf("It's recommended to use %s instead", recommendedVersion),
			)
			return
		}
	}

	score.Grade = scorecard.GradeAllOK
	return
}
