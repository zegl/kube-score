package score

import "github.com/zegl/kube-score/scorecard"
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func scoreMetaStableAvailable(meta metav1.TypeMeta) (score scorecard.TestScore) {
	score.Name = "Unstable version"

	withStable := map[string]map[string]string {
		"extensions/v1beta1": {
			"Deployment": "apps/v1",
		},
		"apps/v1beta1": {
			"Deployment": "apps/v1",
		},
		"apps/v1beta2": {
			"Deployment": "apps/v1",
		},
	}

	if inVersion, ok := withStable[meta.APIVersion]; ok  {
		if recommendedVersion, ok := inVersion[meta.Kind]; ok {
			score.Grade = 5
			score.Comments = []string{"Recommended replacement: " + recommendedVersion}
			return
		}
	}

	score.Grade = 10
	return
}