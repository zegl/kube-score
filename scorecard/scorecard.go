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
	Name string

	ResourceRef struct {
		Kind      string
		Name      string
		Namespace string
		Version   string
	}

	Grade    Grade
	Comments []TestScoreComment
}

type Grade int

const (
	GradeCritical = 1
	GradeWarning  = 5
	GradeAlmostOK = 7
	GradeAllOK    = 10
)

type TestScoreComment struct {
	Path        string
	Summary     string
	Description string
}

func (ts *TestScore) AddComment(path, summary, description string) {
	ts.Comments = append(ts.Comments, TestScoreComment{
		Path:        path,
		Summary:     summary,
		Description: description,
	})
}

func (ts *TestScore) AddMeta(typeMeta metav1.TypeMeta, objectMeta metav1.ObjectMeta) {
	ts.ResourceRef.Name = objectMeta.Name
	ts.ResourceRef.Namespace = objectMeta.Namespace
	ts.ResourceRef.Kind = typeMeta.Kind
	ts.ResourceRef.Version = typeMeta.APIVersion
}

func (ts TestScore) resourceRefKey() string {
	return ts.ResourceRef.Kind + "/" + ts.ResourceRef.Namespace + "/" + ts.ResourceRef.Name
}
