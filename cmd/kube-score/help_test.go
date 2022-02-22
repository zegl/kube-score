package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecName(t *testing.T) {
	assert.Equal(t, "kube-score", execName("kube-score"))
	assert.Equal(t, "ks", execName("ks"))
	assert.Equal(t, "kubectl score", execName("kubectl-score"))
}
