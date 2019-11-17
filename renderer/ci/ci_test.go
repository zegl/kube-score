package ci

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getTestCard() *scorecard.Scorecard {
	checks := []scorecard.TestScore{
		{
			Check: domain.Check{
				Name: "test-warning-two-comments",
			},
			Grade: scorecard.GradeWarning,
			Comments: []scorecard.TestScoreComment{
				{
					Path:        "a",
					Summary:     "summary",
					Description: "description",
				},
				{
					// No path
					Summary:     "summary",
					Description: "description",
				},
			},
		},
		{
			Check: domain.Check{
				Name: "test-ok-comment",
			},
			Grade: scorecard.GradeAllOK,
			Comments: []scorecard.TestScoreComment{
				{
					Path:        "a",
					Summary:     "summary",
					Description: "description",
				},
			},
		},
		{
			Check: domain.Check{
				Name: "test-skipped-comment",
			},
			Skipped: true,
			Comments: []scorecard.TestScoreComment{
				{
					Path:        "a",
					Summary:     "skipped sum",
					Description: "skipped description",
				},
			},
		},
		{
			Check: domain.Check{
				Name: "test-skipped-no-comment",
			},
			Skipped: true,
		},
	}

	return &scorecard.Scorecard{
		"a": &scorecard.ScoredObject{
			TypeMeta: v1.TypeMeta{
				Kind:       "Testing",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      "foo",
				Namespace: "foofoo",
			},
			Checks: checks,
		},

		// No namespace
		"b": &scorecard.ScoredObject{
			TypeMeta: v1.TypeMeta{
				Kind:       "Testing",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name: "bar-no-namespace",
			},
			Checks: checks,
		},
	}
}

func TestCiOutput(t *testing.T) {
	t.Parallel()
	// Defaults
	r := CI(getTestCard())
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `[WARNING] foo/foofoo v1/Testing: (a) summary
[WARNING] foo/foofoo v1/Testing: summary
[OK] foo/foofoo v1/Testing: (a) summary
[SKIPPED] foo/foofoo v1/Testing: (a) skipped sum
[SKIPPED] foo/foofoo v1/Testing
[WARNING] bar-no-namespace v1/Testing: (a) summary
[WARNING] bar-no-namespace v1/Testing: summary
[OK] bar-no-namespace v1/Testing: (a) summary
[SKIPPED] bar-no-namespace v1/Testing: (a) skipped sum
[SKIPPED] bar-no-namespace v1/Testing
`, string(all))
}
