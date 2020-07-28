package security

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zegl/kube-score/score/checks"
	"github.com/zegl/kube-score/scorecard"
)

func Register(allChecks *checks.Checks) {
	allChecks.RegisterPodCheck("Container Security Context", `Makes sure that all pods have good securityContexts configured`, containerSecurityContext)
	allChecks.RegisterOptionalPodCheck("Container Seccomp Profile", `Makes sure that all pods have at a seccomp policy configured.`, podSeccompProfile)
}

// containerSecurityContext checks that the recommended securityPolicy options are set
func containerSecurityContext(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {
	allContainers := podTemplate.Spec.InitContainers
	allContainers = append(allContainers, podTemplate.Spec.Containers...)

	noContextSet := false
	hasPrivileged := false
	hasWritableRootFS := false
	hasLowUserID := false
	hasLowGroupID := false

	podSecurityContext := podTemplate.Spec.SecurityContext

	for _, container := range allContainers {

		if container.SecurityContext == nil && podSecurityContext == nil {
			noContextSet = true
			score.AddComment(container.Name, "Container has no configured security context", "Set securityContext to run the container in a more secure context.")
			continue
		}

		sec := container.SecurityContext

		if sec == nil {
			sec = &corev1.SecurityContext{}
		}

		// Forward values from PodSecurityContext to the (container level) SecurityContext if not set
		if podSecurityContext != nil {
			if sec.RunAsGroup == nil {
				sec.RunAsGroup = podSecurityContext.RunAsGroup
			}
			if sec.RunAsUser == nil {
				sec.RunAsUser = podSecurityContext.RunAsUser
			}
		}

		if sec.Privileged != nil && *sec.Privileged {
			hasPrivileged = true
			score.AddComment(container.Name, "The container is privileged", "Set securityContext.privileged to false. Privileged containers can access all devices on the host, and grants almost the same access as non-containerized processes on the host.")
		}

		if sec.ReadOnlyRootFilesystem == nil || *sec.ReadOnlyRootFilesystem == false {
			hasWritableRootFS = true
			score.AddComment(container.Name, "The pod has a container with a writable root filesystem", "Set securityContext.readOnlyRootFilesystem to true")
		}

		if sec.RunAsUser == nil || *sec.RunAsUser < 10000 {
			hasLowUserID = true
			score.AddComment(container.Name, "The container is running with a low user ID", "A userid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.runAsUser to a value > 10000")
		}

		if sec.RunAsGroup == nil || *sec.RunAsGroup < 10000 {
			hasLowGroupID = true
			score.AddComment(container.Name, "The container running with a low group ID", "A groupid above 10 000 is recommended to avoid conflicts with the host. Set securityContext.runAsGroup to a value > 10000")
		}
	}

	if noContextSet || hasPrivileged || hasWritableRootFS || hasLowUserID || hasLowGroupID {
		score.Grade = scorecard.GradeCritical
	} else {
		score.Grade = scorecard.GradeAllOK
	}

	return
}

// podSeccompProfile checks if the any Seccommp profile is configured for the pod
func podSeccompProfile(podTemplate corev1.PodTemplateSpec, typeMeta metav1.TypeMeta) (score scorecard.TestScore) {
	metadata := podTemplate.ObjectMeta

	seccompAnnotated := false
	if metadata.Annotations != nil {
		if _, ok := metadata.Annotations["seccomp.security.alpha.kubernetes.io/defaultProfileName"]; ok {
			seccompAnnotated = true
		}
	}

	if !seccompAnnotated {
		score.Grade = scorecard.GradeWarning
		score.AddComment(metadata.Name, "The pod has not configured Seccomp for its containers", "Running containers with Seccomp is recommended to reduce the kernel attack surface")
	} else {
		score.Grade = scorecard.GradeAllOK
	}

	return
}
