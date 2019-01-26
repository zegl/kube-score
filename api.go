package kube_score

import (
	"fmt"
	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/score"
	"github.com/zegl/kube-score/scorecard"
	"io"
)

// Score runs all kube-score tests on all files in input
// The result is a single Scorecard that contains the result of all checked objects
// The purpose of this API is to make it easier for other tools to integrate with kube-score
func Score(input []io.Reader) (*scorecard.Scorecard, error) {
	cnf := config.Configuration{
		AllFiles: input,
	}

	parsedFiles, err := parser.ParseFiles(cnf)
	if err != nil {
		return nil, fmt.Errorf("failed to parse files: %s", err)
	}

	res, err := score.Score(parsedFiles, cnf)
	if err != nil {
		return nil, fmt.Errorf("failed to score files: %s", err)
	}

	return res, nil
}