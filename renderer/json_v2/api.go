package json_v2

import (
	"bytes"
	"encoding/json"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ks "github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
)

type Check struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	TargetType string `json:"target_type"`
	Comment    string `json:"comment"`
	Optional   bool   `json:"optional"`
}

type ScoredObject struct {
	ObjectName string            `json:"object_name"`
	TypeMeta   metav1.TypeMeta   `json:"type_meta"`
	ObjectMeta metav1.ObjectMeta `json:"object_meta"`
	Checks     []TestScore       `json:"checks"`
	FileName   string            `json:"file_name"`
	FileRow    int               `json:"file_row"`
}

type TestScore struct {
	Check    Check              `json:"check"`
	Grade    scorecard.Grade    `json:"grade"`
	Skipped  bool               `json:"skipped"`
	Comments []TestScoreComment `json:"comments"`
}

type TestScoreComment struct {
	Path        string `json:"path"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

func Output(input *scorecard.Scorecard) io.Reader {
	var objs []ScoredObject

	for k, v := range *input {
		objs = append(objs, ScoredObject{
			ObjectName: k,
			TypeMeta:   v.TypeMeta,
			ObjectMeta: v.ObjectMeta,
			Checks:     convertTestScore(v.Checks),
			FileName:   v.FileLocation.Name,
			FileRow:    v.FileLocation.Line,
		})
	}

	j, err := json.MarshalIndent(objs, "", "    ")
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(j)
}

func convertTestScore(in []scorecard.TestScore) (res []TestScore) {
	for _, v := range in {
		res = append(res, TestScore{
			Check:    convertCheck(v.Check),
			Grade:    v.Grade,
			Skipped:  v.Skipped,
			Comments: convertComments(v.Comments),
		})
	}
	return
}

func convertComments(in []scorecard.TestScoreComment) (res []TestScoreComment) {
	for _, v := range in {
		res = append(res, TestScoreComment{
			Path:        v.Path,
			Summary:     v.Summary,
			Description: v.Description,
		})
	}
	return
}

func convertCheck(v ks.Check) Check {
	return Check{
		Name:       v.Name,
		ID:         v.ID,
		TargetType: v.TargetType,
		Comment:    v.Comment,
		Optional:   v.Optional,
	}
}
