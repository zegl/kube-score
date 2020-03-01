package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecName(t *testing.T) {
	assert.Equal(t, "kube-score", execName("kube-score"))
	assert.Equal(t, "ks", execName("ks"))
	assert.Equal(t, "kubectl score", execName("kubectl-score"))
}
