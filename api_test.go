package kube_score

import (
	"github.com/zegl/kube-score/scorecard"
	"io"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestScoreAPI(t *testing.T) {
	res, err := Score([]io.Reader{strings.NewReader(`
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  startingDeadlineSeconds: 100
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: hello
              image: busybox
              args:
                - /bin/sh
                - -c
                - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
`)})
	assert.Nil(t, err)
	assert.Len(t, res.Objects, 1)

	checkedScores := 0

	for _, o := range res.Objects {
		for _, c := range o.Checks {
			if c.ID == "pod-probes" {
				checkedScores++
				assert.Equal(t, scorecard.GradeCritical, c.Grade)
			}
			if c.ID == "stable-version" {
				checkedScores++
				assert.Equal(t, scorecard.GradeAllOK, c.Grade)
			}
		}
	}

	assert.Equal(t, 2, checkedScores)
}
