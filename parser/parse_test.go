package parser

import (
	"fmt"
	"os"
	"testing"

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
			fmt.Errorf("one or more files failed to parse"),
		}, {
			"testdata/valid-yaml.yaml",
			nil,
		},
	}

	parser, err := New(nil)
	assert.NoError(t, err)

	for _, tc := range cases {
		fp, err := os.Open(tc.fname)
		assert.Nil(t, err)
		_, err = parser.ParseFiles(
			[]ks.NamedReader{fp},
		)
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
