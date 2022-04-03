package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := []struct {
		fname    string
		expected error
	}{
		{
			"testdata/invalid-yaml.yaml",
			fmt.Errorf("Failed to parse /v1, Kind=Service: err=json: cannot unmarshal string into Go struct field ServicePort.spec.ports.nodePort of type int32"),
		}, {
			"testdata/valid-yaml.yaml",
			nil,
		},
	}

	parser, err := New()
	assert.NoError(t, err)

	for _, tc := range cases {
		fp, err := os.Open(tc.fname)
		assert.Nil(t, err)
		_, err = parser.ParseFiles(config.Configuration{
			AllFiles: []ks.NamedReader{fp},
		})
		if tc.expected == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, tc.expected.Error(), err.Error())
		}
	}
}

func TestFileLocationHelm(t *testing.T) {
	doc := `# Source: app1/templates/deployment.yaml
kind: Deployment
apiVersion: apps/v1
metadata:
  name: foo
spec:
  template:
    metadata:
      labels:
        foo: bar`

	fl := detectFileLocation("someName", 1, []byte(doc))
	assert.Equal(t, "app1/templates/deployment.yaml", fl.Name)
	assert.Equal(t, 1, fl.Line)
}

func TestFileLocation(t *testing.T) {
	doc := `kind: Deployment
apiVersion: apps/v1
metadata:
  name: foo
spec:
  template:
    metadata:
      labels:
        foo: bar`

	fl := detectFileLocation("someName", 123, []byte(doc))
	assert.Equal(t, "someName", fl.Name)
	assert.Equal(t, 123, fl.Line)
}
