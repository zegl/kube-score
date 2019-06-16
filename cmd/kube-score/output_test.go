package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"

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
					Grade: scorecard.Grade(9),
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
					Grade: scorecard.Grade(9),
				},
			},
		},
	}

	ignoredTests := make(map[string]struct{})

	// Defaults
	r := outputCi(card, 10, 5, ignoredTests)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `[WARNING] foo/foofoo v1/Testing: (a) summary
[WARNING] foo/foofoo v1/Testing: summary
[WARNING] bar-no-namespace v1/Testing
`, string(all))

	// OK at 9 or higher
	r = outputCi(card, 9, 5, ignoredTests)
	all, err = ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `[OK] foo/foofoo v1/Testing: (a) summary
[OK] foo/foofoo v1/Testing: summary
[OK] bar-no-namespace v1/Testing
`, string(all))

	// OK at 8 or higher
	r = outputCi(card, 8, 5, ignoredTests)
	all, err = ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `[OK] foo/foofoo v1/Testing: (a) summary
[OK] foo/foofoo v1/Testing: summary
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
					Grade: scorecard.Grade(9),
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
					Grade: scorecard.Grade(9),
				},
			},
		},
	}

	ignoredTests := make(map[string]struct{})

	// Defaults
	r := outputHuman(card, 10, 5, ignoredTests)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo
    [WARNING] TestingA
        * a -> summary
             description
        * summary
             description
v1/Testing bar-no-namespace
    [WARNING] TestingA
`, string(all))

	// OK at 9 or higher
	r = outputHuman(card, 9, 5, ignoredTests)
	all, err = ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo
    [OK] TestingA
        * a -> summary
             description
        * summary
             description
v1/Testing bar-no-namespace
    [OK] TestingA
`, string(all))

	// OK at 8 or higher
	r = outputHuman(card, 8, 5, ignoredTests)
	all, err = ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo
    [OK] TestingA
        * a -> summary
             description
        * summary
             description
v1/Testing bar-no-namespace
    [OK] TestingA
`, string(all))
}
