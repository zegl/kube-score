package scorecard

import (
	"fmt"
	"log"

	"github.com/zegl/kube-score/config"
	ks "github.com/zegl/kube-score/domain"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ignoredChecksAnnotation  = "kube-score/ignore"
	optionalChecksAnnotation = "kube-score/enable"
)

// if this, then that
var impliedIgnoreAnnotations = map[string][]string{
	"container-resources": {"container-ephemeral-storage-request-and-limit"},
}

type Scorecard map[string]*ScoredObject

// New creates and initializes a new Scorecard
func New() Scorecard {
	return make(Scorecard)
}

func (s Scorecard) NewObject(typeMeta metav1.TypeMeta, objectMeta metav1.ObjectMeta, cnf *config.RunConfiguration) *ScoredObject {
	if cnf == nil {
		cnf = &config.RunConfiguration{}
	}

	o := &ScoredObject{
		TypeMeta:   typeMeta,
		ObjectMeta: objectMeta,
		Checks:     make([]TestScore, 0),

		useIgnoreChecksAnnotation:   cnf.UseIgnoreChecksAnnotation,
		useOptionalChecksAnnotation: cnf.UseOptionalChecksAnnotation,
		enabledOptionalTests:        cnf.EnabledOptionalTests,
	}

	// If this object already exists, return the previous version
	if object, ok := s[o.resourceRefKey()]; ok {
		return object
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

	useIgnoreChecksAnnotation   bool
	useOptionalChecksAnnotation bool
	enabledOptionalTests        map[string]struct{}
}

func (so *ScoredObject) AnyBelowOrEqualToGrade(threshold Grade) bool {
	for _, o := range so.Checks {
		if !o.Skipped && o.Grade <= threshold {
			return true
		}
	}
	return false
}

func (so *ScoredObject) resourceRefKey() string {
	return so.TypeMeta.Kind + "/" + so.TypeMeta.APIVersion + "/" + so.ObjectMeta.Namespace + "/" + so.ObjectMeta.Name
}

func (so *ScoredObject) HumanFriendlyRef() string {
	s := so.ObjectMeta.Name
	if so.ObjectMeta.Namespace != "" {
		s += "/" + so.ObjectMeta.Namespace
	}
	s += " " + so.TypeMeta.APIVersion + "/" + so.TypeMeta.Kind
	return s
}

func (so *ScoredObject) Add(ts TestScore, check ks.Check, locationer ks.FileLocationer, annotations ...map[string]string) {
	ts.Check = check
	so.FileLocation = locationer.FileLocation()

	skipAll := so.FileLocation.Skip
	skip := skipAll
	if !skip && annotations != nil {
		var err error
		skipAll, err = so.isSkipped(annotations)
		if err != nil {
			log.Printf(
				"failed to parse %s#L%d",
				so.FileLocation.Name,
				so.FileLocation.Line,
			)
		}
		// if skipAll {
		// 	log.Printf("skip all for %s#L%d %v\n",
		// 		so.FileLocation.Name,
		// 		so.FileLocation.Line,
		// 		annotations,
		// 	)
		// }
		// skip = skipAll
		if len(annotations) == 1 && !so.isEnabled(check, annotations[0], nil) {
			skip = true
		}
		if len(annotations) == 2 && !so.isEnabled(check, annotations[0], annotations[1]) {
			skip = true
		}
	}

	// This test is ignored (via annotations), don't save the score
	// ts.Skipped = skip || skipAll
	if skipAll {
		ts.Skipped = true
		ts.Comments = []TestScoreComment{{Summary: fmt.Sprintf(
			"Skipped because %s#L%d is ignored",
			so.FileLocation.Name, so.FileLocation.Line,
		)}}
	} else if skip {
		ts.Skipped = true
		ts.Comments = []TestScoreComment{{Summary: fmt.Sprintf("Skipped because %s is ignored", check.ID)}}
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
