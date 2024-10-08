package scorecard

import (
	"fmt"
	"strings"

	ks "github.com/zegl/kube-score/domain"
	"gopkg.in/yaml.v3"
)

func (so *ScoredObject) isSkipped(allAnnotations []map[string]string) (bool, error) {
	skip := false
	for _, annotations := range allAnnotations {
		if skipAnnotation, ok := annotations["kube-score/skip"]; ok {
			if err := yaml.Unmarshal([]byte(skipAnnotation), &skip); err != nil {
				return false, fmt.Errorf("invalid skip annotation %q, must be boolean", skipAnnotation)
			}
		}
		// if ignoreAnnotation, ok := annotations["kube-score/ignore"]; ok {
		// 	if strings.TrimSpace(ignoreAnnotation) == "*" {
		// 		skip = true
		// 	}
		// }
	}
	return skip, nil
}

func (so *ScoredObject) isEnabled(check ks.Check, annotations, childAnnotations map[string]string) bool {
	isIn := func(csv string, key string) bool {
		for _, v := range strings.Split(csv, ",") {
			v = strings.TrimSpace(v)
			if v == key {
				return true
			}
			if v == "*" {
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
