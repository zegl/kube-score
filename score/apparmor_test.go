package score

import "testing"

func TestApparmorAnnotation(t *testing.T) {
	testExpectedScore(t, "deployment-sets-apparmor.yaml", "Deployment sets apparmor annotation", 10)
}
func TestNoApparmorAnnotation(t *testing.T) {
        testExpectedScore(t, "deployment-sets-no-apparmor.yaml", "Deployment sets apparmor annotation", 1)
}
func TestDefaultApparmorAnnotation(t *testing.T) {
        testExpectedScore(t, "deployment-sets-default-apparmor.yaml", "Deployment sets apparmor annotation", 7)
}
