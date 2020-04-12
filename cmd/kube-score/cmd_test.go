package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCli(t *testing.T) {
	cmds := map[string]cmdFunc{
		"a":     func(string, []string) {},
		"b":     func(string, []string) {},
		"score": func(string, []string) {},
	}

	cmd, offset, err := parse([]string{"kube-score", "a"}, cmds)
	assert.Equal(t, "a", cmd)
	assert.Equal(t, 2, offset)
	assert.Nil(t, err)

	cmd, offset, err = parse([]string{"kube-score", "unknown"}, cmds)
	assert.Equal(t, "", cmd)
	assert.Equal(t, 0, offset)
	assert.Nil(t, err)

	cmd, offset, err = parse([]string{"kubectl-score", "a"}, cmds)
	assert.Equal(t, "a", cmd)
	assert.Equal(t, 2, offset)
	assert.Nil(t, err)

	cmd, offset, err = parse([]string{"kubectl-score", "xyz.yaml"}, cmds)
	assert.Equal(t, "score", cmd)
	assert.Equal(t, 1, offset)
	assert.Nil(t, err)

	cmd, offset, err = parse([]string{"kubectl-score", "score", "xyz.yaml"}, cmds)
	assert.Equal(t, "score", cmd)
	assert.Equal(t, 2, offset)
	assert.Nil(t, err)
}
