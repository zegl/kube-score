package disruptionbudget

import (
	"fmt"

	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/score/internal"
	"github.com/zegl/kube-score/scorecard"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Register(allChecks *checks.Checks, budgets ks.PodDisruptionBudgets) {
	allChecks.RegisterStatefulSetCheck("StatefulSet has PodDisruptionBudget", `Makes sure that all StatefulSets are targeted by a PDB`, statefulSetHas(budgets.PodDisruptionBudgets()))
	allChecks.RegisterDeploymentCheck("Deployment has PodDisruptionBudget", `Makes sure that all Deployments are targeted by a PDB`, deploymentHas(budgets.PodDisruptionBudgets()))
	allChecks.RegisterPodDisruptionBudgetCheck("PodDisruptionBudget has policy", `Makes sure that PodDisruptionBudgets specify minAvailable or maxUnavailable`, hasPolicy)
}

func hasMatching(budgets []ks.PodDisruptionBudget, namespace string, labels map[string]string) (bool, error) {
	for _, budget := range budgets {
		if budget.Namespace() != namespace {
			score.AddComment("", "PodDisruptionBudget namespace didn't match Deployment namespace", "Cannot associate the PodDisruptionBudget to the Deployment because they are not explicitly in the same namespace")
			continue
		}

		selector, err := metav1.LabelSelectorAsSelector(budget.PodDisruptionBudgetSelector())
		if err != nil {
			return false, fmt.Errorf("failed to create selector: %w", err)
		}

		if selector.Matches(internal.MapLabels(labels)) {
			return true, nil
		}
	}

	return false, nil
}

func statefulSetHas(budgets []ks.PodDisruptionBudget) func(appsv1.StatefulSet) (scorecard.TestScore, error) {
	return func(statefulset appsv1.StatefulSet) (score scorecard.TestScore, err error) {
		if statefulset.Spec.Replicas != nil && *statefulset.Spec.Replicas < 2 {
			score.Skipped = true
			score.AddComment("", "Skipped because the statefulset has less than 2 replicas", "")
			return
		}

		match, matchErr := hasMatching(budgets, statefulset.Namespace, statefulset.Spec.Template.Labels)
		if matchErr != nil {
			err = matchErr
			return
		}

		if match {
			score.Grade = scorecard.GradeAllOK
		} else {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "No matching PodDisruptionBudget was found", "It's recommended to define a PodDisruptionBudget to avoid unexpected downtime during Kubernetes maintenance operations, such as when draining a node.")
		}

		return
	}
}

func deploymentHas(budgets []ks.PodDisruptionBudget) func(appsv1.Deployment) (scorecard.TestScore, error) {
	return func(deployment appsv1.Deployment) (score scorecard.TestScore, err error) {
		if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas < 2 {
			score.Skipped = true
			score.AddComment("", "Skipped because the deployment has less than 2 replicas", "")
			return
		}

		match, matchErr := hasMatching(budgets, deployment.Namespace, deployment.Spec.Template.Labels)
		if matchErr != nil {
			err = matchErr
			return
		}

		if match {
			score.Grade = scorecard.GradeAllOK
		} else {
			score.Grade = scorecard.GradeCritical
			score.AddComment("", "No matching PodDisruptionBudget was found", "It's recommended to define a PodDisruptionBudget to avoid unexpected downtime during Kubernetes maintenance operations, such as when draining a node.")
		}

		return
	}
}

func hasPolicy(pdb ks.PodDisruptionBudget) (score scorecard.TestScore, err error) {
	spec := pdb.Spec()
	if spec.MinAvailable == nil && spec.MaxUnavailable == nil {
		score.AddComment("", "PodDisruptionBudget missing policy", "PodDisruptionBudget should specify minAvailable or maxUnavailable.")
		score.Grade = scorecard.GradeCritical
	} else {
		score.Grade = scorecard.GradeAllOK
	}

	return
}
