package disruptionbudget

import (
	"github.com/zegl/kube-score/scorecard"

	appsv1 "k8s.io/api/apps/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func hasMatching(budgets []policyv1beta1.PodDisruptionBudget, lables map[string]string) bool {
	for _, budget := range budgets {
		selector, err := metav1.LabelSelectorAsSelector(budget.Spec.Selector)
		if err != nil {
			panic(err)
		}

		if selector.Matches(mapLables(lables)) {
			return true
		}
	}

	return false
}

func ScoreStatefulSetHas(budgets []policyv1beta1.PodDisruptionBudget) func(appsv1.StatefulSet) scorecard.TestScore {
	return func(statefulset appsv1.StatefulSet) (score scorecard.TestScore) {
		score.Name = "StatefulSet has PodDisruptionBudget"

		if hasMatching(budgets, statefulset.Spec.Template.Labels) {
			score.Grade = scorecard.GradeAllOK
		} else {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "No matching PodDisruptionBudget was found", "It's recommended to define a PodDisruptionBudget to avoid unexpected downtime during Kubernetes maintenance operations, such as when draining a node.")
		}

		return
	}
}

func ScoreDeploymentHas(budgets []policyv1beta1.PodDisruptionBudget) func(appsv1.Deployment) scorecard.TestScore {
	return func(deployment appsv1.Deployment) (score scorecard.TestScore) {
		score.Name = "Deployment has PodDisruptionBudget"

		if hasMatching(budgets, deployment.Spec.Template.Labels) {
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
