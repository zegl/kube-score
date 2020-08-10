package stable

import (
	"fmt"

	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(kubernetesVersion config.Semver, allChecks *checks.Checks) {
	allChecks.RegisterMetaCheck("Stable version", `Checks if the object is using a deprecated apiVersion`, metaStableAvailable(kubernetesVersion))
}

// ScoreMetaStableAvailable checks if the supplied TypeMeta is an unstable object type, that has a stable(r) replacement
func metaStableAvailable(kubernetsVersion config.Semver) func(meta domain.BothMeta) (score scorecard.TestScore) {
	return func(meta domain.BothMeta) (score scorecard.TestScore) {
		type recommendedApi struct {
			newAPI         string
			availableSince config.Semver
		}

		withStable := map[string]map[string]recommendedApi{
			"extensions/v1beta1": {
				"Deployment": recommendedApi{"apps/v1", config.Semver{1, 9}},
				"DaemonSet":  recommendedApi{"apps/v1", config.Semver{1, 9}},
				"Ingress":    recommendedApi{"networking.k8s.io/v1beta1", config.Semver{1, 14}},
			},
			"apps/v1beta1": {
				"Deployment":  recommendedApi{"apps/v1", config.Semver{1, 9}},
				"StatefulSet": recommendedApi{"apps/v1", config.Semver{1, 9}},
			},
			"apps/v1beta2": {
				"Deployment":  recommendedApi{"apps/v1", config.Semver{1, 9}},
				"StatefulSet": recommendedApi{"apps/v1", config.Semver{1, 9}},
				"DaemonSet":   recommendedApi{"apps/v1", config.Semver{1, 9}},
			},
		}

		score.Grade = scorecard.GradeAllOK

		if inVersion, ok := withStable[meta.TypeMeta.APIVersion]; ok {
			if recAPI, ok := inVersion[meta.TypeMeta.Kind]; ok {

				// The recommended replacement is not available in the version of Kubernetes
				// that the user is using
				if kubernetsVersion.LessThan(recAPI.availableSince) {
					return
				}

				score.Grade = scorecard.GradeWarning
				score.AddComment("",
					fmt.Sprintf("The apiVersion and kind %s/%s is deprecated", meta.TypeMeta.APIVersion, meta.TypeMeta.Kind),
					fmt.Sprintf("It's recommended to use %s instead which has been available since Kubernetes %s", recAPI.newAPI, recAPI.availableSince.String()),
				)
				return
			}
		}

		return
	}
}
