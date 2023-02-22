package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestPodTopologySpreadContraintsWithOneConstraint(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-topology-spread-constraints-one-constraint.yaml", "Pod Topology Spread Constraints", scorecard.GradeAllOK)
}

func TestPodTopologySpreadContraintsWithTwoConstraints(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-topology-spread-constraints-two-constraints.yaml", "Pod Topology Spread Constraints", scorecard.GradeAllOK)
}

func TestPodTopologySpreadContraintsNoLabelSelector(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-topology-spread-constraints-no-labelselector.yaml", "Pod Topology Spread Constraints", scorecard.GradeCritical)
}

func TestPodTopologySpreadContraintsInvalidMaxSkew(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-topology-spread-constraints-invalid-maxskew.yaml", "Pod Topology Spread Constraints", scorecard.GradeCritical)
}

func TestPodTopologySpreadContraintsInvalidMinDomains(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-topology-spread-constraints-invalid-mindomains.yaml", "Pod Topology Spread Constraints", scorecard.GradeCritical)
}

func TestPodTopologySpreadContraintsNoTopologyKey(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-topology-spread-constraints-no-topologykey.yaml", "Pod Topology Spread Constraints", scorecard.GradeCritical)
}

func TestPodTopologySpreadContraintsInvalidDirective(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "pod-topology-spread-constraints-invalid-whenunsatisfiable.yaml", "Pod Topology Spread Constraints", scorecard.GradeCritical)
}
