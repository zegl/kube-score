package score

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func TestDeploymentHasPodAntiAffinityPreffered(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-host-antiaffinity-preffered.yaml", "Deployment has host PodAntiAffinity", scorecard.GradeAllOK)
}

func TestDeploymentHasPodAntiAffinityPrefferedNoSelectorMatch(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-host-antiaffinity-preffered-selector-no-match.yaml", "Deployment has host PodAntiAffinity", scorecard.GradeWarning)
}

func TestDeploymentHasPodAntiAffinityPrefferedSelectorExpression(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-host-antiaffinity-preffered-selector-expression.yaml", "Deployment has host PodAntiAffinity", scorecard.GradeAllOK)
}

func TestDeploymentHasPodAntiAffinityRequired(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-host-antiaffinity-required.yaml", "Deployment has host PodAntiAffinity", scorecard.GradeAllOK)
}

func TestDeploymentHasPodAntiAffinityNotSet(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-host-antiaffinity-not-set.yaml", "Deployment has host PodAntiAffinity", scorecard.GradeWarning)
}

func TestDeploymentHasPodAntiAffinityOneReplica(t *testing.T) {
	t.Parallel()
	// skipped
	testExpectedScore(t, "deployment-host-antiaffinity-1-replica.yaml", "Deployment has host PodAntiAffinity", 0)
}

func TestStatefulSetHasPodAntiAffinityPreffered(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-host-antiaffinity-preffered.yaml", "StatefulSet has host PodAntiAffinity", scorecard.GradeAllOK)
}

func TestStatefulSetHasPodAntiAffinityRequired(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-host-antiaffinity-required.yaml", "StatefulSet has host PodAntiAffinity", scorecard.GradeAllOK)
}

func TestStatefulSetHasPodAntiAffinityNotSet(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-host-antiaffinity-not-set.yaml", "StatefulSet has host PodAntiAffinity", scorecard.GradeWarning)
}

func TestStatefulSetHasPodAntiAffinityOneReplica(t *testing.T) {
	t.Parallel()
	// skipped
	testExpectedScore(t, "statefulset-host-antiaffinity-1-replica.yaml", "StatefulSet has host PodAntiAffinity", 0)
}

func TestStatefulSetHasPodAntiAffinityUndefinedReplicas(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-host-antiaffinity-undefined-replicas.yaml", "StatefulSet has host PodAntiAffinity", scorecard.GradeWarning)
}

func TestDeploymentWithHPAHasReplicas(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-with-hpa-has-replicas.yaml", "Deployment targeted by HPA does not have replicas configured", scorecard.GradeCritical)
}

func TestDeploymentWithHPANotHasReplicas(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "deployment-with-hpa-not-has-replicas.yaml", "Deployment targeted by HPA does not have replicas configured", scorecard.GradeAllOK)
}

func TestStatefulsetHasServiceName(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-service-name.yaml", "StatefulSet has ServiceName", scorecard.GradeAllOK)
}

func TestStatefulsetHasServiceNameDifferentName(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-service-name-different-name.yaml", "StatefulSet has ServiceName", scorecard.GradeCritical)
}

func TestStatefulsetHasServiceNameDifferentNamespace(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-service-name-not-headless.yaml", "StatefulSet has ServiceName", scorecard.GradeCritical)
}

func TestStatefulsetHasServiceNameNotHeadless(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-service-name-not-headless.yaml", "StatefulSet has ServiceName", scorecard.GradeCritical)
}

func TestStatefulsetHasServiceNameDifferentLabel(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-service-name-different-label.yaml", "StatefulSet has ServiceName", scorecard.GradeCritical)
}

func TestStatefulsetSelectorLabels(t *testing.T) {
	t.Parallel()
	testExpectedScore(t, "statefulset-different-labels.yaml", "StatefulSet Pod Selector labels match template metadata labels", scorecard.GradeCritical)
}

func TestStatefulsetTemplateIgnores(t *testing.T) {
	t.Parallel()
	skipped := wasSkipped(t, []ks.NamedReader{testFile("statefulset-nested-ignores.yaml")}, nil, &config.RunConfiguration{
		UseIgnoreChecksAnnotation:   true,
		UseOptionalChecksAnnotation: true,
	}, "Container Image Tag")
	assert.True(t, skipped)
}

func TestStatefulsetTemplateIgnoresNotIgnoredWhenFlagDisabled(t *testing.T) {
	t.Parallel()
	skipped := wasSkipped(t, []ks.NamedReader{testFile("statefulset-nested-ignores.yaml")}, nil, &config.RunConfiguration{
		UseIgnoreChecksAnnotation:   false,
		UseOptionalChecksAnnotation: true,
	}, "Container Image Tag")
	assert.False(t, skipped)
}

func TestStatefulsetTemplateNestedSkip(t *testing.T) {
	t.Parallel()
	sc, err := testScore(
		[]ks.NamedReader{testFile("statefulset-nested-skip.yaml")},
		nil,
		&config.RunConfiguration{
			UseIgnoreChecksAnnotation:   true,
			UseOptionalChecksAnnotation: true,
		},
	)
	assert.NoError(t, err)

	for _, objectScore := range sc {
		for _, s := range objectScore.Checks {
			t.Logf("id=%s type=%v skipped=%v\n", s.Check.ID, s.Check.TargetType, s.Skipped)
			switch s.Check.TargetType {
			case "StatefulSet", "all":
				assert.False(t, s.Skipped)
			default:
				assert.True(t, s.Skipped)
			}
		}
	}
}

func TestStatefulsetTemplateSkip(t *testing.T) {
	skipped := fileWasSkipped(
		t,
		[]ks.NamedReader{testFile("statefulset-skip.yaml")},
		nil,
		&config.RunConfiguration{
			UseIgnoreChecksAnnotation:   true,
			UseOptionalChecksAnnotation: true,
		},
		"testdata/statefulset-skip.yaml",
	)
	assert.True(t, skipped)
}
