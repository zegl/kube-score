package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func TestInvalidLabel(t *testing.T) {
	t.Parallel()
	s := validateLabelValues(domain.BothMeta{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"foo": "engineering/kustomize", // label values can't contain slashes
				"bar": "baribar",
			},
		},
	})
	assert.Equal(t, scorecard.GradeCritical, s.Grade)
	assert.Len(t, s.Comments, 1)
	assert.Equal(t, "foo", s.Comments[0].Path)
	assert.Equal(t, "Invalid label value", s.Comments[0].Summary)
	assert.Equal(t, "The label value is invalid, and will not be accepted by Kubernetes", s.Comments[0].Description)
}
func TestOKLabel(t *testing.T) {
	t.Parallel()
	s := validateLabelValues(domain.BothMeta{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"foo": "foo-bar",
			},
		},
	})
	assert.Equal(t, scorecard.GradeAllOK, s.Grade)
}
