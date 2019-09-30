package scorecard

import (
	"encoding/json"
	"fmt"
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

func (s Scorecard) NewObject(typeMeta metav1.TypeMeta, objectMeta metav1.ObjectMeta) *ScoredObject {
	o := &ScoredObject{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
		Checks:     make([]TestScore, 0),
	}

	// If this object already exists, return the previous version
	if object, ok := s[o.resourceRefKey()]; ok {
		return object
	}

	o.setIgnoredTests()

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

func (s Scorecard) MarshalJSON() ([]byte, error) {
	type EncodedScorecard struct {
		ObjectName string `json:"object_name"`
		*ScoredObject
	}

	var result []EncodedScorecard
	for k, v := range s {
		result = append(result, EncodedScorecard{
			ScoredObject: v,
			ObjectName:   k,
		})
	}

	return json.Marshal(result)
}

type ScoredObject struct {
	TypeMeta   metav1.TypeMeta   `json:"type_meta"`
	ObjectMeta metav1.ObjectMeta `json:"object_meta"`
	Checks     []TestScore       `json:"checks"`

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

func (so *ScoredObject) Add(ts TestScore, check ks.Check) {
	ts.Check = check

	// This test is ignored (via annotations), don't save the score
	if _, ok := so.ignoredChecks[check.ID]; ok {
		ts.Skipped = true
		ts.Comments = []TestScoreComment{{Summary: fmt.Sprintf("Skipped because %s is ignored", check.ID)}}
	}

	so.Checks = append(so.Checks, ts)
}

type TestScore struct {
	Check    ks.Check           `json:"check"`
	Grade    Grade              `json:"grade"`
	Skipped  bool               `json:"skipped"`
	Comments []TestScoreComment `json:"comments"`
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
	Path        string `json:"path"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

func (ts *TestScore) AddComment(path, summary, description string) {
	ts.Comments = append(ts.Comments, TestScoreComment{
		Path:        path,
		Summary:     summary,
		Description: description,
	})
}
