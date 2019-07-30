package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func TestCiOutput(t *testing.T) {
	card := &scorecard.Scorecard{
		"a": &scorecard.ScoredObject{
			TypeMeta: v1.TypeMeta{
				Kind:       "Testing",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      "foo",
				Namespace: "foofoo",
			},
			Checks: []scorecard.TestScore{
				{
					Check: domain.Check{
						Name: "TestingA",
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
			},
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
			Checks: []scorecard.TestScore{
				{
					Check: domain.Check{
						Name: "TestingA",
					},
					Grade: scorecard.GradeAllOK,
				},
			},
		},
	}

	// Defaults
	r := outputCi(card)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `[WARNING] foo/foofoo v1/Testing: (a) summary
[WARNING] foo/foofoo v1/Testing: summary
[OK] bar-no-namespace v1/Testing
`, string(all))
}

func TestHumanOutput(t *testing.T) {
	card := &scorecard.Scorecard{
		"a": &scorecard.ScoredObject{
			TypeMeta: v1.TypeMeta{
				Kind:       "Testing",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      "foo",
				Namespace: "foofoo",
			},
			Checks: []scorecard.TestScore{
				{
					Check: domain.Check{
						Name: "TestingA",
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
			},
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
			Checks: []scorecard.TestScore{
				{
					Check: domain.Check{
						Name: "TestingA",
					},
					Grade: scorecard.GradeAllOK,
				},
			},
		},
	}

	// Defaults
	r := outputHuman(card)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo
    [WARNING] TestingA
        * a -> summary
             description
        * summary
             description
v1/Testing bar-no-namespace
    [OK] TestingA
`, string(all))

}
