package parser

import (
	"fmt"
	"io"
	"os"
	"strings"
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
			fmt.Errorf("Failed to parse /v1, Kind=Service: err=json: cannot unmarshal string into Go struct field ServicePort.spec.ports.nodePort of type int32"),
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

type namedReader struct {
	io.Reader
	name string
}

func (n namedReader) Name() string {
	return n.name
}

func parse(t *testing.T, doc, name string) ks.AllTypes {
	p, err := New(nil)
	assert.NoError(t, err)
	parsedFiles, err := p.ParseFiles([]ks.NamedReader{
		namedReader{Reader: strings.NewReader(doc), name: name},
	})
	assert.NoError(t, err)
	return parsedFiles
}

func TestSkipNo(t *testing.T) {
	t.Parallel()
	doc := `kind: Deployment
apiVersion: apps/v1
metadata:
  name: foo
  annotations:
    kube-score/skip:  "No"
spec:
  template:
    metadata:
      labels:
        foo: bar`

	location := parse(t, doc, "skip-yes.yaml").Deployments()[0].FileLocation()
	assert.Equal(t, "skip-yes.yaml", location.Name)
	assert.Equal(t, false, location.Skip)
}

func TestSkipYes(t *testing.T) {
	t.Parallel()
	doc := `kind: Deployment
apiVersion: apps/v1
metadata:
  name: foo
  annotations:
    kube-score/skip:  " yes  "
spec:
  template:
    metadata:
      labels:
        foo: bar`

	location := parse(t, doc, "skip-yes.yaml").Deployments()[0].FileLocation()
	assert.Equal(t, "skip-yes.yaml", location.Name)
	assert.Equal(t, true, location.Skip)
}

func TestSkipTrueUppercase(t *testing.T) {
	t.Parallel()
	doc := `kind: Deployment
apiVersion: apps/v1
metadata:
  name: foo
  annotations:
    "kube-score/skip": "True"
spec:
  template:
    metadata:
      labels:
        foo: bar`

	location := parse(t, doc, "skip-true-uppercase.yaml").Deployments()[0].FileLocation()
	assert.Equal(t, "skip-true-uppercase.yaml", location.Name)
	assert.Equal(t, true, location.Skip)
}

func TestSkipTrue(t *testing.T) {
	t.Parallel()
	doc := `kind: Deployment
apiVersion: apps/v1
metadata:
  name: foo
  annotations:
    "kube-score/skip": "true"
spec:
  template:
    metadata:
      labels:
        foo: bar`

	location := parse(t, doc, "skip-true.yaml").Deployments()[0].FileLocation()
	assert.Equal(t, "skip-true.yaml", location.Name)
	assert.Equal(t, true, location.Skip)
}

func TestSkipFalse(t *testing.T) {
	t.Parallel()
	doc := `kind: Deployment
apiVersion: apps/v1
metadata:
  name: foo
  annotations:
    "kube-score/skip": "false"
spec:
  template:
    metadata:
      labels:
        foo: bar`

	location := parse(t, doc, "skip-false.yaml").Deployments()[0].FileLocation()
	assert.Equal(t, "skip-false.yaml", location.Name)
	assert.Equal(t, false, location.Skip)
}
