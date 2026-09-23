package sarif

import (
	"encoding/json"
	"io"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/scorecard"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// emptySnippetRE matches a "snippet" key with an empty-object value,
// tolerant of json.Marshal's compact form ("snippet":{}) and
// json.MarshalIndent's spaced form ("snippet": {}).
var emptySnippetRE = regexp.MustCompile(`"snippet":\s*\{\s*\}`)

// getTestCard mirrors the fixture pattern used by the other renderer test
// suites (see renderer/junit/junit_test.go): one critical check with a
// comment, on an object with a known FileLocation.
func getTestCard() *scorecard.Scorecard {
	return &scorecard.Scorecard{
		"a": &scorecard.ScoredObject{
			TypeMeta: v1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      "foo",
				Namespace: "foofoo",
			},
			FileLocation: domain.FileLocation{
				Name: "dev.yaml",
				Line: 15,
			},
			Checks: []scorecard.TestScore{
				{
					Check: domain.Check{
						ID:   "pod-probes-identical",
						Name: "pod-probes-identical",
					},
					Grade: scorecard.GradeCritical,
					Comments: []scorecard.TestScoreComment{
						{
							Summary:     "summary",
							Description: "description",
						},
					},
				},
			},
		},
	}
}

// TestOutput_NoEmptySnippetObject is a live repro/regression test for
// https://github.com/zegl/kube-score/issues/698: the SARIF 2.1.0 schema
// types "snippet" as an artifactContent object that must carry "text" or
// "binary" when present at all. kube-score never fills in either field
// (see Output() above, which only ever sets ContextRegion.StartLine), but
// Region/ContextRegion.Snippet is a plain (non-pointer) struct, so Go's
// `encoding/json` omitempty has no effect on it -- omitempty only skips
// the recognized "empty" kinds (false, 0, nil pointer/interface, empty
// array/slice/map/string), never a struct value, zero or not. The result
// is a literal `"snippet":{}` in every result's region/contextRegion,
// which schema validators (GitLab's SARIF ingestion, GitHub's
// `upload-sarif`) reject outright, silently dropping the whole report.
func TestOutput_NoEmptySnippetObject(t *testing.T) {
	t.Parallel()

	r := Output(getTestCard())
	raw, err := io.ReadAll(r)
	assert.Nil(t, err)

	// The bug is specifically the presence of the key "snippet" with an
	// empty object value anywhere in the document -- assert on the raw
	// JSON text so the test fails the same way a schema validator would,
	// not just on Go-side struct shape.
	if emptySnippetRE.Match(raw) {
		t.Fatalf("output contains an empty \"snippet\":{} object, which is invalid under the SARIF 2.1.0 schema (artifactContent requires \"text\" or \"binary\"); got:\n%s", raw)
	}

	// Sanity check: the document must still parse and carry the result.
	var parsed map[string]interface{}
	assert.Nil(t, json.Unmarshal(raw, &parsed))
}
