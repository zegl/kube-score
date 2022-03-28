package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubeScoreConfigExcludeAllDefaultChecks(t *testing.T) {

	cfg := loadConfigFile("testdata/kube-score.yml")
	cfg.AddAllDefaultChecks = false
	excludeThese := excludeChecks(&cfg)

	assert.Equal(t, len(excludeThese), len(cfg.DefaultChecks))
}

func TestKubeScoreConfigIncludeAllOptionalChecks(t *testing.T) {

	cfg := loadConfigFile("testdata/kube-score.yml")
	cfg.AddAllOptionalChecks = true
	includeThese := includeChecks(&cfg)

	assert.Equal(t, len(includeThese), len(cfg.OptionalChecks))
}

func TestKubeScoreConfigExcludeSelectDefaultChecks(t *testing.T) {

	cfg := loadConfigFile("testdata/kube-score.yml")
	cfg.AddAllDefaultChecks = true
	cfg.ExcludeChecks = append(cfg.ExcludeChecks, "pod-probes")
	excludeThese := excludeChecks(&cfg)

	assert.Contains(t, cfg.ExcludeChecks, "pod-probes")
	assert.Equal(t, len(excludeThese), 1)
}

func TestKubeScoreConfigNoDefaultChecksIncludeSelectChecks(t *testing.T) {

	cfg := loadConfigFile("testdata/kube-score.yml")
	cfg.AddAllDefaultChecks = false

	onlyThese := []string{"container-resources", "image-tag", "image-pull-policy"}

	cfg.IncludeChecks = append(cfg.IncludeChecks, onlyThese...)
	includeThese := includeChecks(&cfg)

	for _, v := range onlyThese {
		assert.Contains(t, cfg.IncludeChecks, v)
	}

	assert.NotContains(t, cfg.IncludeChecks, "pod-networkpolicy")

	assert.Equal(t, len(includeThese), len(onlyThese))
}
