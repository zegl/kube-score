package disruptionbudget

import (
	"github.com/stretchr/testify/assert"
	"github.com/zegl/kube-score/scorecard"
	appsv1 "k8s.io/api/apps/v1"
	"testing"
)

func TestStatefulSetReplicas(t *testing.T) {
	cases := map[*int32]scorecard.Grade{
		nil:        scorecard.GradeCritical, // failed
		intptr(1):  scorecard.GradeAllOK,    // skipped
		intptr(10): scorecard.GradeCritical, // failed
	}

	fn := statefulSetHas(nil)

	for replicas, expected := range cases {
		res, err := fn(appsv1.StatefulSet{Spec: appsv1.StatefulSetSpec{Replicas: replicas}})
		assert.Nil(t, err)

		if replicas == nil {
			assert.Equal(t, expected, res.Grade, "replicas=nil")
		} else {
			assert.Equal(t, expected, res.Grade, "replicas=%+v", *replicas)
		}
	}
}

func TestDeploymentReplicas(t *testing.T) {
	cases := map[*int32]scorecard.Grade{
		nil:        scorecard.GradeCritical, // failed
		intptr(1):  scorecard.GradeAllOK,    // skipped
		intptr(10): scorecard.GradeCritical, // failed
	}

	fn := deploymentHas(nil)

	for replicas, expected := range cases {
		res, err := fn(appsv1.Deployment{Spec: appsv1.DeploymentSpec{Replicas: replicas}})
		assert.Nil(t, err)

		if replicas == nil {
			assert.Equal(t, expected, res.Grade, "replicas=nil")
		} else {
			assert.Equal(t, expected, res.Grade, "replicas=%+v", *replicas)
		}
	}
}

func intptr(a int32) *int32 {
	return &a
}
