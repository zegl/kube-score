package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSemver(t *testing.T) {
	tc := []struct {
		input       string
		expected    Semver
		expectedErr error
	}{
		{"v1.0", Semver{1, 0}, nil},
		{"v1.999", Semver{1, 999}, nil},
		{"1.0", Semver{1, 0}, nil},
		{"1.999", Semver{1, 999}, nil},

		{"foo", Semver{}, errInvalidSemver},
		{"v1.2.3", Semver{}, errInvalidSemver},
		{"v1.foo", Semver{}, errInvalidSemver},
		{"x1.0", Semver{}, errInvalidSemver},
		{"v0x00.123", Semver{}, errInvalidSemver},
		{"v1b.5nn3", Semver{}, errInvalidSemver},
	}

	for d, tc := range tc {
		s, e := ParseSemver(tc.input)
		assert.Equal(t, tc.expected, s, "Case: %d", d)
		assert.Equal(t, tc.expectedErr, e, "Case: %d", d)
	}
}
