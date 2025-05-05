package scorecard

import (
	"strings"

	ks "github.com/zegl/kube-score/domain"
)

func (so *ScoredObject) isEnabled(check ks.Check, annotations, childAnnotations map[string]string) bool {
	isIn := func(csv string, key string) bool {
		for _, v := range strings.Split(csv, ",") {
			v = strings.TrimSpace(v)
			if v == key {
				return true
			}
			if v == "*" {
				// "*" wildcard matches all checks
				return true
			}
			if vals, ok := impliedIgnoreAnnotations[v]; ok {
				for i := range vals {
					if vals[i] == key {
						return true
					}
				}
			}
		}
		return false
	}

	if childAnnotations != nil && so.useIgnoreChecksAnnotation && isIn(childAnnotations[ignoredChecksAnnotation], check.ID) {
		return false
	}
	if childAnnotations != nil && so.useOptionalChecksAnnotation && isIn(childAnnotations[optionalChecksAnnotation], check.ID) {
		return true
	}
	if so.useIgnoreChecksAnnotation && isIn(annotations[ignoredChecksAnnotation], check.ID) {
		return false
	}
	if so.useOptionalChecksAnnotation && isIn(annotations[optionalChecksAnnotation], check.ID) {
		return true
	}

	// Enabled optional test from command line arguments
	if _, ok := so.enabledOptionalTests[check.ID]; ok {
		return true
	}

	// Optional checks are disabled unless explicitly allowed above
	if check.Optional {
		return false
	}

	// Enabled by default
	return true
}
