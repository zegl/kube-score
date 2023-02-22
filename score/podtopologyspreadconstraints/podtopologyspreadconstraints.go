package podtopologyconstraints

import (
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks) {
	allChecks.RegisterPodCheck("Pod Topology Spread Constraints", "Pod Topology Spread Constraints", podTopologySpreadConstraints)
}

func podTopologySpreadConstraints(pod ks.PodSpecer) (score scorecard.TestScore, err error) {
	spreads := pod.GetPodTemplateSpec().Spec.TopologySpreadConstraints

	if spreads == nil {
		score.Grade = scorecard.GradeAllOK
		score.AddComment("", "Pod Topology Spread Constraints", "No Pod Topology Spread Constraints set, kube-scheduler defaults assumed")
		return
	}

	for _, spread := range spreads {
		if spread.LabelSelector == nil {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "Pod Topology Spread Constraints", "No labelSelector detected. A label selector is needed determine the number of pods in a topology domain")
			return
		}

		if spread.MaxSkew == 0 {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "Pod Topology Spread Constraints", "MaxSkew is set to zero. This is not allowed.")
			return
		}

		if spread.MinDomains != nil && *spread.MinDomains == 0 {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "Pod Topology Spread Constraints", "MaxDomain set to zero. This is not allowed. Constraint behaves if minDomains is set to 1 if nil")
			return
		}

		if spread.TopologyKey == "" {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "Pod Topology Spread Constraints", "TopologyKey is not set. This is the key of node labels used to bucket nodes into a domain")
			return
		}

		if spread.WhenUnsatisfiable != "DoNotSchedule" && spread.WhenUnsatisfiable != "ScheduleAnyway" {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "Pod Topology Spread Constraints", "Invalid WhenUnsatisfiable setting detected")
			return
		}
	}

	score.Grade = scorecard.GradeAllOK
	score.AddComment("", "Pod Topology Spread Constraints", "Pod Topology Spread Constraints")
	return
}
