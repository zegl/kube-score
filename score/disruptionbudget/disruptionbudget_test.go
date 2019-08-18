package disruptionbudget

import (
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	"testing"

	"github.com/zegl/kube-score/scorecard"
)

func TestStatefulSetReplicas(t *testing.T) {
	cases := map[*int32]struct {
		grade   scorecard.Grade
		skipped bool
	}{
		nil:        {scorecard.GradeCritical, false}, // failed
		intptr(1):  {0, true},                        // skipped
		intptr(10): {scorecard.GradeCritical, false}, // failed
	}

	fn := statefulSetHas(nil)

	for replicas, expected := range cases {
		res, err := fn(appsv1.StatefulSet{Spec: appsv1.StatefulSetSpec{Replicas: replicas}})
		assert.Nil(t, err)

		assert.Equal(t, expected.skipped, res.Skipped)

		if replicas == nil {
			assert.Equal(t, expected.grade, res.Grade, "replicas=nil")
		} else {
			assert.Equal(t, expected.grade, res.Grade, "replicas=%+v", *replicas)
		}
	}
}

func TestDeploymentReplicas(t *testing.T) {
	cases := map[*int32]struct {
		grade   scorecard.Grade
		skipped bool
	}{
		nil:        {scorecard.GradeCritical, false}, // failed
		intptr(1):  {0, true},                        // skipped
		intptr(10): {scorecard.GradeCritical, false}, // failed
	}

	fn := deploymentHas(nil)

	for replicas, expected := range cases {
		res, err := fn(appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: replicas}})
		assert.Nil(t, err)

		assert.Equal(t, expected.skipped, res.Skipped)

		if replicas == nil {
			assert.Equal(t, expected.grade, res.Grade, "replicas=nil")
		} else {
			assert.Equal(t, expected.grade, res.Grade, "replicas=%+v", *replicas)
		}
	}
}

func intptr(a int32) *int32 {
	return &a
}
