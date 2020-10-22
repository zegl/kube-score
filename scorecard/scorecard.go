package scorecard

import (
	"fmt"
	"github.com/zegl/kube-score/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"

	ks "github.com/zegl/kube-score/domain"
)

const (
	ignoredChecksAnnotation = "kube-score/ignore"
)

type Scorecard map[string]*ScoredObject

// New creates and initializes a new Scorecard
func New() Scorecard {
	return make(Scorecard)
}

func (s Scorecard) NewObject(typeMeta metav1.TypeMeta, objectMeta metav1.ObjectMeta, useIgnoreChecksAnnotation bool) *ScoredObject {
	o := &ScoredObject{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
		Checks:     make([]TestScore, 0),
	}

	// If this object already exists, return the previous version
	if object, ok := s[o.resourceRefKey()]; ok {
		return object
	}

	if useIgnoreChecksAnnotation {
		o.setIgnoredTests()
	}

	s[o.resourceRefKey()] = o
	return o
}

func (s Scorecard) AnyBelowOrEqualToGrade(threshold Grade) bool {
	for _, o := range s {
		if o.AnyBelowOrEqualToGrade(threshold) {
			return true
		}
	}
	return false
}

type ScoredObject struct {
	TypeMeta     metav1.TypeMeta
	ObjectMeta   metav1.ObjectMeta
	FileLocation ks.FileLocation
	Checks       []TestScore

	ignoredChecks map[string]struct{}
}

func (s ScoredObject) AnyBelowOrEqualToGrade(threshold Grade) bool {
	for _, o := range s.Checks {
		if o.Skipped == false && o.Grade <= threshold {
			return true
		}
	}
	return false
}

func (so *ScoredObject) setIgnoredTests() {
	ignoredMap := make(map[string]struct{})
	if ignoredCSV, ok := so.ObjectMeta.Annotations[ignoredChecksAnnotation]; ok {
		for _, ignored := range strings.Split(ignoredCSV, ",") {
			ignoredMap[strings.TrimSpace(ignored)] = struct{}{}
		}
	}
	so.ignoredChecks = ignoredMap
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

func (so *ScoredObject) Add(ts TestScore, check ks.Check, locationer ks.FileLocationer, cnf config.Configuration) {
	ts.Check = check
	so.FileLocation = locationer.FileLocation()

	// This test is ignored (via annotations), don't save the score
	if _, ok := so.ignoredChecks[check.ID]; ok {
		ts.Skipped = true
		ts.Comments = []TestScoreComment{{Summary: fmt.Sprintf("Skipped because %s is ignored", check.ID)}}
	}

	if _, ok := cnf.IgnoredNamespaces[so.ObjectMeta.Namespace]; ok {
		ts.Skipped = true
		ts.Comments = []TestScoreComment{{Summary: fmt.Sprintf("Skipped because %s namespace is ignored", so.ObjectMeta.Namespace)}}
	}

	so.Checks = append(so.Checks, ts)
}

type TestScore struct {
	Check    ks.Check
	Grade    Grade
	Skipped  bool
	Comments []TestScoreComment
}

type Grade int

const (
	GradeCritical Grade = 1
	GradeWarning  Grade = 5
	GradeAlmostOK Grade = 7
	GradeAllOK    Grade = 10
)

func (g Grade) String() string {
	switch g {
	case GradeCritical:
		return "CRITICAL"
	case GradeWarning:
		return "WARNING"
	case GradeAlmostOK, GradeAllOK:
		return "OK"
	default:
		panic("Unknown grade")
	}
}

type TestScoreComment struct {
	Path             string
	Summary          string
	Description      string
	DocumentationURL string
}

func (ts *TestScore) AddComment(path, summary, description string) {
	ts.Comments = append(ts.Comments, TestScoreComment{
		Path:        path,
		Summary:     summary,
		Description: description,
	})
}

func (ts *TestScore) AddCommentWithURL(path, summary, description, documentationURL string) {
	ts.Comments = append(ts.Comments, TestScoreComment{
		Path:             path,
		Summary:          summary,
		Description:      description,
		DocumentationURL: documentationURL,
	})
}
