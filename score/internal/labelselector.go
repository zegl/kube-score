package internal

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// MapLabels implements the Kubernetes Labels interface
type MapLabels map[string]string

func (l MapLabels) Has(key string) bool {
	_, ok := l[key]
	return ok
}

func (l MapLabels) Get(key string) string {
	return l[key]
}

func (l MapLabels) Lookup(label string) (value string, exists bool) {
	value, exists = l[label]
	return value, exists
}

func LabelSelectorMatchesLabels(selectorLabels map[string]string, labels map[string]string) bool {
	labelSelector := &metav1.LabelSelector{
		MatchLabels: selectorLabels,
	}

	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return false
	}

	return selector.Matches(MapLabels(labels))
}
