package internal

// Implements the Kubernetes Labels interface
type MapLables map[string]string

func (l MapLables) Has(key string) bool {
	_, ok := l[key]
	return ok
}

func (l MapLables) Get(key string) string {
	return l[key]
}
