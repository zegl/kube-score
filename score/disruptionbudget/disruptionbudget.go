package disruptionbudget

import (
	"github.com/zegl/kube-score/scorecard"

	appsv1 "k8s.io/api/apps/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ScoreStatefulSetHas(budgets []policyv1beta1.PodDisruptionBudget) func(appsv1.StatefulSet) scorecard.TestScore {
	return func(statefulset appsv1.StatefulSet) (score scorecard.TestScore) {
		score.Name = "StatefulSet has PodDisruptionBudget"

		hasMatching := false

		for _, budget := range budgets {
			selector, err := metav1.LabelSelectorAsSelector(budget.Spec.Selector)
			if err != nil {
				panic(err)
			}

			if selector.Matches(mapLables(statefulset.Spec.Template.Labels)) {
				hasMatching = true
				break
			}
		}

		if hasMatching {
			score.Grade = scorecard.GradeAllOK
		} else {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "No matching PodDisruptionBudget was found", "It's recommended to define a PodDisruptionBudget to avoid unexpected downtime during Kubernetes maintenance operations, such as when draining a node.")
		}

		return
	}
}

type mapLables map[string]string

func (m mapLables) Has(key string) bool {
	_, ok := m[key]
	return ok
}

func (m mapLables) Get(key string) string {
	return m[key]
}
