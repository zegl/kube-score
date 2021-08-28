package internal

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// MapLables implements the Kubernetes Labels interface
type MapLables map[string]string

func (l MapLables) Has(key string) bool {
	_, ok := l[key]
	return ok
}

func (l MapLables) Get(key string) string {
	return l[key]
}

func LabelSelectorMatchesLabels(selectorLabels map[string]string, labels map[string]string) bool {
	labelSelector := &metav1.LabelSelector{
		MatchLabels: selectorLabels,
	}

	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return false
	}

	return selector.Matches(MapLables(labels))
}
