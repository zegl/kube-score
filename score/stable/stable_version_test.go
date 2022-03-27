package stable

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

func TestStableVersionOldKubernetesVersion(t *testing.T) {
	oldKubernetes := metaStableAvailable(config.Semver{Major: 1, Minor: 4})
	scoreOld := oldKubernetes(ks.BothMeta{TypeMeta: v1.TypeMeta{Kind: "Deployment", APIVersion: "extensions/v1beta1"}})
	assert.Equal(t, scorecard.GradeAllOK, scoreOld.Grade)
	assert.Equal(t, []scorecard.TestScoreComment(nil), scoreOld.Comments)
}

func TestStableVersionNewKubernetesVersion(t *testing.T) {
	newKubernetes := metaStableAvailable(config.Semver{Major: 1, Minor: 18})
	scoreNew := newKubernetes(ks.BothMeta{TypeMeta: v1.TypeMeta{Kind: "Deployment", APIVersion: "extensions/v1beta1"}})
	assert.Equal(t, scorecard.GradeWarning, scoreNew.Grade)
	assert.Equal(t, []scorecard.TestScoreComment{{Path: "", Summary: "The apiVersion and kind extensions/v1beta1/Deployment is deprecated", Description: "It's recommended to use apps/v1 instead which has been available since Kubernetes v1.9", DocumentationURL: ""}}, scoreNew.Comments)
}

func TestStableVersionIngress(t *testing.T) {
	newKubernetes := metaStableAvailable(config.Semver{Major: 1, Minor: 20})
	scoreNew := newKubernetes(ks.BothMeta{TypeMeta: v1.TypeMeta{Kind: "Ingress", APIVersion: "extensions/v1beta1"}})
	assert.Equal(t, scorecard.GradeWarning, scoreNew.Grade)
	assert.Equal(t, []scorecard.TestScoreComment{{Path: "", Summary: "The apiVersion and kind extensions/v1beta1/Ingress is deprecated", Description: "It's recommended to use networking.k8s.io/v1 instead which has been available since Kubernetes v1.19", DocumentationURL: ""}}, scoreNew.Comments)
}

func TestStableVersionPodDisruptionBudget(t *testing.T) {
	newKubernetes := metaStableAvailable(config.Semver{Major: 1, Minor: 21})
	scoreNew := newKubernetes(ks.BothMeta{TypeMeta: v1.TypeMeta{Kind: "PodDisruptionBudget", APIVersion: "policy/v1beta1"}})
	assert.Equal(t, scorecard.GradeWarning, scoreNew.Grade)
	assert.Equal(t, []scorecard.TestScoreComment{{Path: "", Summary: "The apiVersion and kind policy/v1beta1/PodDisruptionBudget is deprecated", Description: "It's recommended to use policy/v1 instead which has been available since Kubernetes v1.21", DocumentationURL: ""}}, scoreNew.Comments)
}

func TestStableNetworkingIngress(t *testing.T) {
	newKubernetes := metaStableAvailable(config.Semver{Major: 1, Minor: 21})
	scoreNew := newKubernetes(ks.BothMeta{TypeMeta: v1.TypeMeta{Kind: "Ingress", APIVersion: "networking.k8s.io/v1beta1"}})
	assert.Equal(t, scorecard.GradeWarning, scoreNew.Grade)
	assert.Equal(t, []scorecard.TestScoreComment{{Path: "", Summary: "The apiVersion and kind networking.k8s.io/v1beta1/Ingress is deprecated", Description: "It's recommended to use networking.k8s.io/v1 instead which has been available since Kubernetes v1.19", DocumentationURL: ""}}, scoreNew.Comments)
}
