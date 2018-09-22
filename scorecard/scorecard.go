package scorecard

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Scorecard struct {
	Scores map[string][]TestScore
}

// New creates and initializes a new Scorecard
func New() *Scorecard {
	return &Scorecard{
		Scores: make(map[string][]TestScore),
	}
}

// Add adds a TestScore to the Scorecard
func (s *Scorecard) Add(ts TestScore) {
	if existingScores, ok := s.Scores[ts.resourceRefKey()]; ok {
		existingScores = append(existingScores, ts)
		s.Scores[ts.resourceRefKey()] = existingScores
	} else {
		s.Scores[ts.resourceRefKey()] = []TestScore{ts}
	}
}

type TestScore struct {
	Name        string

	ResourceRef struct{
		Kind string
		Name string
		Namespace string
	}

	Grade       int
	Comments    []string
}

func (ts *TestScore) AddMeta(typeMeta metav1.TypeMeta, objectMeta metav1.ObjectMeta) {
	ts.ResourceRef.Name = objectMeta.Name
	ts.ResourceRef.Namespace = objectMeta.Namespace
	ts.ResourceRef.Kind = typeMeta.Kind
}

func (ts TestScore) resourceRefKey() string {
	return ts.ResourceRef.Kind + "/" + ts.ResourceRef.Namespace + "/" + ts.ResourceRef.Name
}

