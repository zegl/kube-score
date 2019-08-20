package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
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
	// Defaults
	r := outputCi(getTestCard())
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

func TestHumanOutputDefault(t *testing.T) {
	r := outputHuman(getTestCard(), 0, 100)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo                                                      🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            description
        · summary
            description
v1/Testing bar-no-namespace                                                   🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            description
        · summary
            description
`, string(all))
}

func TestHumanOutputVerbose1(t *testing.T) {
	r := outputHuman(getTestCard(), 1, 100)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo                                                      🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            description
        · summary
            description
    [OK] test-ok-comment
        · a -> summary
            description
v1/Testing bar-no-namespace                                                   🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            description
        · summary
            description
    [OK] test-ok-comment
        · a -> summary
            description
`, string(all))
}

func TestHumanOutputVerbose2(t *testing.T) {
	r := outputHuman(getTestCard(), 2, 100)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo                                                      🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            description
        · summary
            description
    [OK] test-ok-comment
        · a -> summary
            description
    [SKIPPED] test-skipped-comment
        · a -> skipped sum
            skipped description
    [SKIPPED] test-skipped-no-comment
v1/Testing bar-no-namespace                                                   🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            description
        · summary
            description
    [OK] test-ok-comment
        · a -> summary
            description
    [SKIPPED] test-skipped-comment
        · a -> skipped sum
            skipped description
    [SKIPPED] test-skipped-no-comment
`, string(all))
}

func getTestCardAllOK() *scorecard.Scorecard {
	checks := []scorecard.TestScore{
		{
			Check: domain.Check{
				Name: "test-warning-two-comments",
			},
			Grade: scorecard.GradeAllOK,
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

func TestHumanOutputAllOKDefault(t *testing.T) {
	r := outputHuman(getTestCardAllOK(), 0, 100)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo                                                      ✅
v1/Testing bar-no-namespace                                                   ✅
`, string(all))
}

func getTestCardLongDescription() *scorecard.Scorecard {
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
					Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras elementum sagittis lacus, a dictum tortor lobortis vel. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Nulla eu neque erat. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Maecenas et nisl venenatis, elementum augue a, porttitor libero.",
				},
			},
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
	}
}

func TestHumanOutputLogDescription120Width(t *testing.T) {
	r := outputHuman(getTestCardLongDescription(), 0, 120)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo                                                      🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras elementum sagittis lacus, a dictum tortor
            lobortis vel. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas.
            Nulla eu neque erat. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae;
            Maecenas et nisl venenatis, elementum augue a, porttitor libero.
`, string(all))
}

func TestHumanOutputLogDescription100Width(t *testing.T) {
	r := outputHuman(getTestCardLongDescription(), 0, 100)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo                                                      🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras elementum sagittis lacus,
            a dictum tortor lobortis vel. Pellentesque habitant morbi tristique senectus et netus et
            malesuada fames ac turpis egestas. Nulla eu neque erat. Vestibulum ante ipsum primis in
            faucibus orci luctus et ultrices posuere cubilia Curae; Maecenas et nisl venenatis,
            elementum augue a, porttitor libero.
`, string(all))
}

func TestHumanOutputLogDescription80Width(t *testing.T) {
	r := outputHuman(getTestCardLongDescription(), 0, 80)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo                                                      🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras
            elementum sagittis lacus, a dictum tortor lobortis vel. Pellentesque
            habitant morbi tristique senectus et netus et malesuada fames ac
            turpis egestas. Nulla eu neque erat. Vestibulum ante ipsum primis in
            faucibus orci luctus et ultrices posuere cubilia Curae; Maecenas et
            nisl venenatis, elementum augue a, porttitor libero.
`, string(all))
}

func TestHumanOutputLogDescription0Width(t *testing.T) {
	r := outputHuman(getTestCardLongDescription(), 0, 0)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing foo in foofoo🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            Lorem ipsum dolor sit amet, consectetur
            adipiscing elit. Cras elementum sagittis
            lacus, a dictum tortor lobortis vel.
            Pellentesque habitant morbi tristique
            senectus et netus et malesuada fames ac
            turpis egestas. Nulla eu neque erat.
            Vestibulum ante ipsum primis in faucibus
            orci luctus et ultrices posuere cubilia
            Curae; Maecenas et nisl venenatis,
            elementum augue a, porttitor libero.
`, string(all))
}

func getTestCardLongTitle() *scorecard.Scorecard {
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
					Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras elementum sagittis lacus, a dictum tortor lobortis vel. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Nulla eu neque erat. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Maecenas et nisl venenatis, elementum augue a, porttitor libero.",
				},
			},
		},
	}

	return &scorecard.Scorecard{
		"a": &scorecard.ScoredObject{
			TypeMeta: v1.TypeMeta{
				Kind:       "Testing",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      "this-is-a-very-long-title-this-is-a-very-long-title-this-is-a-very-long-title-this-is-a-very-long-title-this-is-a-very-long-title",
				Namespace: "foofoo",
			},
			Checks: checks,
		},
	}
}

func TestHumanOutputWithLongObjectNames(t *testing.T) {
	r := outputHuman(getTestCardLongTitle(), 0, 80)
	all, err := ioutil.ReadAll(r)
	assert.Nil(t, err)
	assert.Equal(t, `v1/Testing this-is-a-very-long-title-this-is-a-very-long-title-this-is-a-very-long-title-this-is-a-very-long-title-this-is-a-very-long-title in foofoo🤔
    [WARNING] test-warning-two-comments
        · a -> summary
            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras
            elementum sagittis lacus, a dictum tortor lobortis vel. Pellentesque
            habitant morbi tristique senectus et netus et malesuada fames ac
            turpis egestas. Nulla eu neque erat. Vestibulum ante ipsum primis in
            faucibus orci luctus et ultrices posuere cubilia Curae; Maecenas et
            nisl venenatis, elementum augue a, porttitor libero.
`, string(all))
}
