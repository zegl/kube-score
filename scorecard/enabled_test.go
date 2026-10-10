package scorecard

import (
	"testing"

	"github.com/stretchr/testify/assert"
	ks "github.com/zegl/kube-score/domain"
)

func TestIsEnabled(t *testing.T) {
	t.Parallel()

	defaultCheck := ks.Check{ID: "test-check", Optional: false}
	optionalCheck := ks.Check{ID: "test-optional-check", Optional: true}
	impliedParentCheck := ks.Check{ID: "container-resources", Optional: false}
	impliedChildCheck := ks.Check{ID: "container-ephemeral-storage-request-and-limit", Optional: false}

	tests := []struct {
		name             string
		check            ks.Check
		annotations      map[string]string
		childAnnotations map[string]string
		useIgnore        bool
		useOptional      bool
		cliOptional      map[string]struct{}
		expected         bool
	}{
		{
			name:        "default check is enabled with no annotations",
			check:       defaultCheck,
			useIgnore:   true,
			useOptional: true,
			expected:    true,
		},
		{
			name:        "optional check is disabled by default",
			check:       optionalCheck,
			useIgnore:   true,
			useOptional: true,
			expected:    false,
		},
		{
			name:        "optional check enabled via cli flag",
			check:       optionalCheck,
			useIgnore:   true,
			useOptional: true,
			cliOptional: map[string]struct{}{"test-optional-check": {}},
			expected:    true,
		},
		{
			name:        "legacy annotation ignores check",
			check:       defaultCheck,
			annotations: map[string]string{"kube-score/ignore": "test-check"},
			useIgnore:   true,
			useOptional: true,
			expected:    false,
		},
		{
			name:        "dns-prefixed annotation ignores check",
			check:       defaultCheck,
			annotations: map[string]string{"kube-score.com/ignore": "test-check"},
			useIgnore:   true,
			useOptional: true,
			expected:    false,
		},
		{
			name:        "dns-prefixed annotation with whitespace and multiple items",
			check:       defaultCheck,
			annotations: map[string]string{"kube-score.com/ignore": " other-check , test-check , third-check "},
			useIgnore:   true,
			useOptional: true,
			expected:    false,
		},
		{
			name:        "legacy annotation ignores check via child annotations",
			check:       defaultCheck,
			childAnnotations: map[string]string{"kube-score/ignore": "test-check"},
			useIgnore:   true,
			useOptional: true,
			expected:    false,
		},
		{
			name:        "dns-prefixed annotation ignores check via child annotations",
			check:       defaultCheck,
			childAnnotations: map[string]string{"kube-score.com/ignore": "test-check"},
			useIgnore:   true,
			useOptional: true,
			expected:    false,
		},
		{
			name:        "implied ignore works with dns-prefixed annotation",
			check:       impliedChildCheck,
			annotations: map[string]string{"kube-score.com/ignore": "container-resources"},
			useIgnore:   true,
			useOptional: true,
			expected:    false,
		},
		{
			name:        "implied ignore works with legacy annotation",
			check:       impliedChildCheck,
			annotations: map[string]string{"kube-score/ignore": "container-resources"},
			useIgnore:   true,
			useOptional: true,
			expected:    false,
		},
		{
			name:        "ignore annotation ignored when useIgnoreChecksAnnotation is false",
			check:       defaultCheck,
			annotations: map[string]string{"kube-score.com/ignore": "test-check"},
			useIgnore:   false,
			useOptional: true,
			expected:    true,
		},
		{
			name:        "legacy annotation enables optional check",
			check:       optionalCheck,
			annotations: map[string]string{"kube-score/enable": "test-optional-check"},
			useIgnore:   true,
			useOptional: true,
			expected:    true,
		},
		{
			name:        "dns-prefixed annotation enables optional check",
			check:       optionalCheck,
			annotations: map[string]string{"kube-score.com/enable": "test-optional-check"},
			useIgnore:   true,
			useOptional: true,
			expected:    true,
		},
		{
			name:        "dns-prefixed annotation enables optional check via child annotations",
			check:       optionalCheck,
			childAnnotations: map[string]string{"kube-score.com/enable": "test-optional-check"},
			useIgnore:   true,
			useOptional: true,
			expected:    true,
		},
		{
			name:        "optional check annotation ignored when useOptionalChecksAnnotation is false",
			check:       optionalCheck,
			annotations: map[string]string{"kube-score.com/enable": "test-optional-check"},
			useIgnore:   true,
			useOptional: false,
			expected:    false,
		},
		{
			name: "both legacy and dns annotations present for different checks",
			check: defaultCheck,
			annotations: map[string]string{
				"kube-score/ignore":     "other-check",
				"kube-score.com/ignore": "test-check",
			},
			useIgnore:   true,
			useOptional: true,
			expected:    false,
		},
		{
			name: "both legacy and dns annotations present where legacy matches",
			check: defaultCheck,
			annotations: map[string]string{
				"kube-score/ignore":     "test-check",
				"kube-score.com/ignore": "other-check",
			},
			useIgnore:   true,
			useOptional: true,
			expected:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			so := &ScoredObject{
				useIgnoreChecksAnnotation:   tt.useIgnore,
				useOptionalChecksAnnotation: tt.useOptional,
				enabledOptionalTests:        tt.cliOptional,
			}
			result := so.isEnabled(tt.check, tt.annotations, tt.childAnnotations)
			assert.Equal(t, tt.expected, result)
		})
	}
	_ = impliedParentCheck
}
