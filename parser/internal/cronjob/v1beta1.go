package cronjob

import (
	ks "github.com/zegl/kube-score/domain"
	"k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CronJobV1beta1 struct {
	Obj      v1beta1.CronJob
	Location ks.FileLocation
}

func (c CronJobV1beta1) StartingDeadlineSeconds() *int64 {
	return c.Obj.Spec.StartingDeadlineSeconds
}

func (c CronJobV1beta1) FileLocation() ks.FileLocation {
	return c.Location
}

func (c CronJobV1beta1) GetTypeMeta() metav1.TypeMeta {
	return c.Obj.TypeMeta
}

func (c CronJobV1beta1) GetObjectMeta() metav1.ObjectMeta {
	return c.Obj.ObjectMeta
}

func (c CronJobV1beta1) GetPodTemplateSpec() corev1.PodTemplateSpec {
	t := c.Obj.Spec.JobTemplate.Spec.Template
	t.ObjectMeta.Namespace = c.Obj.ObjectMeta.Namespace
	return t
}
