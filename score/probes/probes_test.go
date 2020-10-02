package probes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPodIsTargetedByService(t *testing.T) {
	t.Run("single label match", func(t *testing.T) {
		res := podIsTargetedByService(v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{"foo": "bar"},
			},
		},
			v1.Service{
				Spec: v1.ServiceSpec{
					Selector: map[string]string{"foo": "bar"},
				},
			},
		)

		assert.True(t, res)
	})

	t.Run("single label mismatch", func(t *testing.T) {
		res := podIsTargetedByService(v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{"foo": "bar"},
			},
		},
			v1.Service{
				Spec: v1.ServiceSpec{
					Selector: map[string]string{"foo": "baz"},
				},
			},
		)

		assert.False(t, res)
	})

	t.Run("multi label match", func(t *testing.T) {
		res := podIsTargetedByService(v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"foo1": "bar1",
					"foo2": "bar2",
				},
			},
		},
			v1.Service{
				Spec: v1.ServiceSpec{
					Selector: map[string]string{
						"foo1": "bar1",
						"foo2": "bar2",
					},
				},
			},
		)

		assert.True(t, res)
	})

	t.Run("multi non full match", func(t *testing.T) {
		res := podIsTargetedByService(v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"foo1": "bar1",
					"foo2": "bar2",
				},
			},
		},
			v1.Service{
				Spec: v1.ServiceSpec{
					Selector: map[string]string{
						"foo1": "bar1",
						"foo2": "bar-whatever",
					},
				},
			},
		)

		assert.False(t, res)
	})

	t.Run("multi label match same namespace", func(t *testing.T) {
		res := podIsTargetedByService(v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "foospace",
				Labels: map[string]string{
					"foo1": "bar1",
					"foo2": "bar2",
				},
			},
		},
			v1.Service{
				ObjectMeta: metav1.ObjectMeta{Namespace: "foospace"},
				Spec: v1.ServiceSpec{
					Selector: map[string]string{
						"foo1": "bar1",
						"foo2": "bar2",
					},
				},
			},
		)

		assert.True(t, res)
	})

	t.Run("multi label match different namespace", func(t *testing.T) {
		res := podIsTargetedByService(v1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "foospace",
				Labels: map[string]string{
					"foo1": "bar1",
					"foo2": "bar2",
				},
			},
		},
			v1.Service{
				ObjectMeta: metav1.ObjectMeta{Namespace: "someOtherNamespace"},
				Spec: v1.ServiceSpec{
					Selector: map[string]string{
						"foo1": "bar1",
						"foo2": "bar2",
					},
				},
			},
		)

		assert.False(t, res)
	})
}
