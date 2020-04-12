package parser

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/zegl/kube-score/config"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	cases := []struct {
		fname    string
		expected error
	}{
		{
			"testdata/invalid-yaml.yaml",
			fmt.Errorf("Failed to parse /v1, Kind=Service: err=v1.Service.Spec: v1.ServiceSpec.Ports: []v1.ServicePort: v1.ServicePort.NodePort: readUint32: unexpected character: \xff, error found in #10 byte of ...|odePort\":\"${PORT}\",\"|..., bigger context ...|\"namespace\":\"test\"},\"spec\":{\"ports\":[{\"nodePort\":\"${PORT}\",\"port\":\"${PORT_EXPOSE}\",\"targetPort\":\"${P|..."),
		}, {
			"testdata/valid-yaml.yaml",
			nil,
		},
	}

	for _, tc := range cases {
		fp, err := os.Open(tc.fname)
		assert.Nil(t, err)
		_, err = ParseFiles(config.Configuration{
			AllFiles: []io.Reader{fp},
		})
		if tc.expected == nil {
			assert.Nil(t, err)
		} else {
			assert.Equal(t, tc.expected.Error(), err.Error())
		}
	}
}
