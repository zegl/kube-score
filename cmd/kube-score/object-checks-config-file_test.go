package main

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// functionally the same as having no configuration file
func TestKubeScoreConfigDefaultChecks(t *testing.T) {

	if cfg, err := loadConfigFile("testdata/kube-score-default.yml"); err == nil {
		assert.Equal(t, cfg.DisableAll, false)
		assert.Equal(t, cfg.EnableAll, false)
		assert.Equal(t, len(cfg.DisableChecks), 0)
		assert.Equal(t, len(cfg.EnableChecks), 0)

		include := includeChecks(&cfg)
		exclude := excludeChecks(&cfg)
		dflt, _ := registeredChecks()
		assert.Equal(t, len(include), len(dflt))
		assert.Equal(t, len(exclude), 0)
	}
}

// functionally means include all optional tests
func TestKubeScoreConfigIncludeAllOptionalChecks(t *testing.T) {

	if cfg, err := loadConfigFile("testdata/kube-score-enable-all.yml"); err == nil {
		allChecks := allRegisteredChecks()
		include := includeChecks(&cfg)
		assert.Equal(t, cfg.EnableAll, true)
		assert.Equal(t, len(include), len(allChecks))
	}
}

// enable most tests, but disable a select few
func TestKubeScoreConfigExcludeSelectDefaultChecks(t *testing.T) {

	if cfg, err := loadConfigFile("testdata/kube-score-enable-all-select-disable.yml"); err == nil {
		assert.Equal(t, cfg.EnableAll, true)
		assert.True(t, len(cfg.DisableChecks) > 0)
		excludeChecks := excludeChecks(&cfg)
		idx := rand.Intn(len(excludeChecks))
		assert.Contains(t, cfg.DisableChecks, excludeChecks[idx])
	}
}

func TestKubeScoreConfigNoDefaultChecksIncludeSelectChecks(t *testing.T) {

	if cfg, err := loadConfigFile("testdata/kube-score-disable-all-select-enable.yml"); err == nil {
		allChecks := allRegisteredChecks()
		includeChecks := includeChecks(&cfg)
		excludeChecks := excludeChecks(&cfg)
		assert.Equal(t, cfg.DisableAll, true)
		assert.Equal(t, len(excludeChecks), len(allChecks))
		idx := rand.Intn(len(includeChecks))
		assert.Contains(t, cfg.EnableChecks, includeChecks[idx])
	}
}

func TestKubeScoreConfigOnlyEnableAndDisableSelectChecks(t *testing.T) {

	if cfg, err := loadConfigFile("testdata/kube-score-only-enable-disable.yml"); err == nil {
		includeChecks := includeChecks(&cfg)
		excludeChecks := excludeChecks(&cfg)
		assert.Equal(t, cfg.DisableAll, true)
		assert.Equal(t, len(excludeChecks), len(cfg.DisableChecks))
		assert.Equal(t, len(includeChecks), len(cfg.EnableChecks))
	}
}

func TestKubeScoreConfigBadFileNameReturnsError(t *testing.T) {

	// badfilename.yml does not exist
	if _, err := loadConfigFile("testdata/badfilename.yml"); err != nil {
		assert.ErrorContains(t, err, "Failed to read file testdata/badfilename.yml")
	}
}
