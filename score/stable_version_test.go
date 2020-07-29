package score

import (
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestStatefulSetAppsv1beta1(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-appsv1beta1.yaml", "Stable version", scorecard.GradeWarning)
}

func TestStatefulSetAppsv1beta2(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-appsv1beta2.yaml", "Stable version", scorecard.GradeWarning)
}

func TestDeploymentExtensionsv1beta1(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-extensions-v1beta1.yaml", "Stable version", scorecard.GradeWarning)
}

func TestDeploymentAppsv1beta1(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-appsv1beta1.yaml", "Stable version", scorecard.GradeWarning)
}

func TestDeploymentAppsv1beta2(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-appsv1beta2.yaml", "Stable version", scorecard.GradeWarning)
}

func TestDaemonSetAppsv1(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "daemonset-appsv1.yaml", "Stable version", scorecard.GradeAllOK)
}

func TestDaemonSetAppsv1beta2(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "daemonset-appsv1beta2.yaml", "Stable version", scorecard.GradeWarning)
}

func TestDaemonSetExtensionsv1beta1(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "daemonset-extensionsv1beta1.yaml", "Stable version", scorecard.GradeWarning)
}

func TestCronJobBatchv1beta1(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "cronjob-batchv1beta1.yaml", "Stable version", scorecard.GradeAllOK)
}

func TestJobBatchv1(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "job-batchv1.yaml", "Stable version", scorecard.GradeAllOK)
}
