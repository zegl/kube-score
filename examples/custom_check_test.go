package examples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zegl/kube-score/scorecard"
)

func TestExampleCheckObjectAllOK(t *testing.T) {
	card, err := ExampleCheckObject([]byte(`
apiVersion: apps/v1
kind: Deployment
metadata:
    name: example
spec:
    replicas: 10
    template:
    metadata:
        labels:
            app: foo
    spec:
        containers:
        - name: foobar
          image: foo:bar`))

	assert.NoError(t, err)

	assert.Len(t, *card, 1)

	for _, v := range *card {
		assert.Len(t, v.Checks, 1)
		assert.Equal(t, "custom-deployment-check", v.Checks[0].Check.ID)
		assert.Equal(t, scorecard.GradeAllOK, v.Checks[0].Grade)
	}
}

func TestExampleCheckObjectErrorNameContainsFoo(t *testing.T) {
	card, err := ExampleCheckObject([]byte(`
apiVersion: apps/v1
kind: Deployment
metadata:
    name: example-foo
spec:
    replicas: 10
    template:
    metadata:
        labels:
            app: foo
    spec:
        containers:
        - name: foobar
          image: foo:bar`))

	assert.NoError(t, err)

	assert.Len(t, *card, 1)

	for _, v := range *card {
		assert.Len(t, v.Checks, 1)
		assert.Equal(t, "custom-deployment-check", v.Checks[0].Check.ID)
		assert.Equal(t, scorecard.GradeCritical, v.Checks[0].Grade)
	}
}
