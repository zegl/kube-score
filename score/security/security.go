package security

import (
	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
	corev1 "k8s.io/api/core/v1"
)

func Register(allChecks *checks.Checks) {
	allChecks.RegisterPodCheck("Container Security Context", containerSecurityContext)
}

// containerSecurityContext checks that the recommended securityPolicy options are set
func containerSecurityContext(podTemplate corev1.PodTemplateSpec) (score scorecard.TestScore) {
	allContainers := podTemplate.Spec.InitContainers
	allContainers = append(allContainers, podTemplate.Spec.Containers...)

	hasPrivileged := false
	hasWritableRootFS := false
	hasLowUserID := false
	hasLowGroupID := false

	for _, container := range allContainers {

		if container.SecurityContext == nil {
			continue
		}

		sec := container.SecurityContext

		if sec.Privileged != nil && *sec.Privileged {
			hasPrivileged = true
			score.AddComment(container.Name, "The container is privileged", "Set securityContext.Privileged to false")
		}

		if sec.ReadOnlyRootFilesystem != nil && *sec.ReadOnlyRootFilesystem == false {
			hasWritableRootFS = true
			score.AddComment(container.Name, "The pod has a container with a writable root filesystem", "Set securityContext.ReadOnlyFileSystem to true")
		}

		if sec.RunAsUser != nil && *sec.RunAsUser < 10000 {
			hasLowUserID = true
			score.AddComment(container.Name, "The container is running with a low user ID", "A userid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.RunAsUser to a value > 10000")
		}

		if sec.RunAsGroup != nil && *sec.RunAsGroup < 10000 {
			hasLowGroupID = true
			score.AddComment(container.Name, "The container running with a low group ID", "A groupid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.RunAsGroup to a value > 10000")
		}
	}

	if hasPrivileged || hasWritableRootFS || hasLowUserID || hasLowGroupID {
		score.Grade = scorecard.GradeCritical
	} else {
		score.Grade = scorecard.GradeAllOK
	}

	return
}
