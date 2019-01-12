package scorecard

import (
	ks "github.com/zegl/kube-score"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Scorecard struct {
	Objects map[string]*ScoredObject
}

// New creates and initializes a new Scorecard
func New() *Scorecard {
	return &Scorecard{
		Objects: make(map[string]*ScoredObject),
	}
}

func (s *Scorecard) NewObject(typeMeta metav1.TypeMeta, objectMeta metav1.ObjectMeta) *ScoredObject {
	o := &ScoredObject{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
		Checks:     make([]TestScore, 0),
	}

	// If this object already exists, return the previous version
	if object, ok := s.Objects[o.resourceRefKey()]; ok {
		return object
	}

	s.Objects[o.resourceRefKey()] = o
	return o
}

type ScoredObject struct {
	TypeMeta   metav1.TypeMeta
	ObjectMeta metav1.ObjectMeta
	Checks     []TestScore
}

func (so ScoredObject) resourceRefKey() string {
	return so.TypeMeta.Kind + "/" + so.TypeMeta.APIVersion + "/" + so.ObjectMeta.Namespace + "/" + so.ObjectMeta.Name
}

func (so ScoredObject) HumanFriendlyRef() string {
	s := so.ObjectMeta.Name
	if so.ObjectMeta.Namespace != "" {
		s += "/" + so.ObjectMeta.Namespace
	}
	s += " " + so.TypeMeta.APIVersion + "/" + so.TypeMeta.Kind
	return s
}

func (so *ScoredObject) Add(ts TestScore, check ks.Check) {
	ts.Check = check
	so.Checks = append(so.Checks, ts)
}

type TestScore struct {
	ks.Check
	Grade        Grade
	Comments     []TestScoreComment
	IgnoredTests []string
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
