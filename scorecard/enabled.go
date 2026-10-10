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

	matchesAnyAnnotation := func(ann map[string]string, keys []string, checkID string) bool {
		if ann == nil {
			return false
		}
		for _, k := range keys {
			if v, ok := ann[k]; ok && isIn(v, checkID) {
				return true
			}
		}
		return false
	}

	if childAnnotations != nil && so.useIgnoreChecksAnnotation && matchesAnyAnnotation(childAnnotations, ignoredChecksAnnotations, check.ID) {
		return false
	}
	if childAnnotations != nil && so.useOptionalChecksAnnotation && matchesAnyAnnotation(childAnnotations, optionalChecksAnnotations, check.ID) {
		return true
	}
	if so.useIgnoreChecksAnnotation && matchesAnyAnnotation(annotations, ignoredChecksAnnotations, check.ID) {
		return false
	}
	if so.useOptionalChecksAnnotation && matchesAnyAnnotation(annotations, optionalChecksAnnotations, check.ID) {
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
