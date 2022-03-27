package score

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
)

func TestFileLocationHelm(t *testing.T) {
	sc, err := testScore(config.Configuration{
		AllFiles:          []ks.NamedReader{testFile("linenumbers-helm.yaml")},
		KubernetesVersion: config.Semver{Major: 1, Minor: 18},
	})
	assert.Nil(t, err)
	for _, c := range sc {
		assert.Equal(t, "app1/templates/deployment.yaml", c.FileLocation.Name)
	}
	assert.Equal(t, 1, sc["Deployment/apps/v1//foo"].FileLocation.Line)
	assert.Equal(t, 1, sc["Deployment/apps/v1//foo2"].FileLocation.Line)
}

func TestFileLocation(t *testing.T) {
	sc, err := testScore(config.Configuration{
		AllFiles:          []ks.NamedReader{testFile("linenumbers.yaml")},
		KubernetesVersion: config.Semver{Major: 1, Minor: 18},
	})
	assert.Nil(t, err)
	for _, c := range sc {
		assert.Equal(t, "testdata/linenumbers.yaml", c.FileLocation.Name)
	}
	assert.Equal(t, 2, sc["Deployment/apps/v1//foo"].FileLocation.Line)
	assert.Equal(t, 12, sc["Deployment/apps/v1//foo2"].FileLocation.Line)
}
