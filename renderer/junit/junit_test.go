package junit

import (
	"io"
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

func TestJUnitOutput(t *testing.T) {
	t.Parallel()
	r := JUnit(getTestCard())
	all, err := io.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `<testsuites name="kube-score" tests="10" failures="4" skipped="4">
	<testsuite name="foo/foofoo v1/Testing" tests="5" failures="2" errors="0" id="0" skipped="2" time="">
		<testcase name="test-warning-two-comments" classname="foo/foofoo v1/Testing">
			<failure message="(a) summary: description"></failure>
		</testcase>
		<testcase name="test-warning-two-comments" classname="foo/foofoo v1/Testing">
			<failure message="summary: description"></failure>
		</testcase>
		<testcase name="test-ok-comment" classname="foo/foofoo v1/Testing"></testcase>
		<testcase name="test-skipped-comment" classname="foo/foofoo v1/Testing">
			<skipped message="(a) skipped sum: skipped description"></skipped>
		</testcase>
		<testcase name="test-skipped-no-comment" classname="foo/foofoo v1/Testing">
			<skipped message=""></skipped>
		</testcase>
	</testsuite>
	<testsuite name="bar-no-namespace v1/Testing" tests="5" failures="2" errors="0" id="0" skipped="2" time="">
		<testcase name="test-warning-two-comments" classname="bar-no-namespace v1/Testing">
			<failure message="(a) summary: description"></failure>
		</testcase>
		<testcase name="test-warning-two-comments" classname="bar-no-namespace v1/Testing">
			<failure message="summary: description"></failure>
		</testcase>
		<testcase name="test-ok-comment" classname="bar-no-namespace v1/Testing"></testcase>
		<testcase name="test-skipped-comment" classname="bar-no-namespace v1/Testing">
			<skipped message="(a) skipped sum: skipped description"></skipped>
		</testcase>
		<testcase name="test-skipped-no-comment" classname="bar-no-namespace v1/Testing">
			<skipped message=""></skipped>
		</testcase>
	</testsuite>
</testsuites>
`, string(all))
}
