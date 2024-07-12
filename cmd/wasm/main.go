//go:build wasm
// +build wasm

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall/js"

	t2html "github.com/buildkite/terminal-to-html"

	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/renderer/ci"
	"github.com/zegl/kube-score/renderer/human"
	"github.com/zegl/kube-score/renderer/json_v2"
	"github.com/zegl/kube-score/renderer/junit"
	"github.com/zegl/kube-score/renderer/sarif"
	"github.com/zegl/kube-score/score"
	"github.com/zegl/kube-score/score/checks"
	"golang.org/x/term"
)

func main() {
	js.Global().Set("handleScore", js.FuncOf(handleScore))
	select {}
}

type inputReader struct {
	*strings.Reader
}

func (inputReader) Name() string {
	return "input"
}

func handleScore(this js.Value, inputs []js.Value) interface{} {
	if len(inputs) != 2 {
		fmt.Println("Unexpected number of arguments")
		return "Unexpected number of arguments"
	}

	inputYaml := inputs[0].String()
	format := inputs[1].String()

	reader := &inputReader{
		Reader: strings.NewReader(inputYaml),
	}

	files := []domain.NamedReader{reader}

	p, err := parser.New(&parser.Config{})
	if err != nil {
		fmt.Println(err)
		return string(err.Error())
	}

	allObjs, err := p.ParseFiles(files)
	if err != nil {
		fmt.Println(err)
		return string(err.Error())
	}

	checksConfig := &checks.Config{}
	runConfig := &config.RunConfiguration{}

	allChecks := score.RegisterAllChecks(allObjs, checksConfig, runConfig)

	scoreCard, err := score.Score(allObjs, allChecks, runConfig)
	if err != nil {
		fmt.Println(err)
		return string(err.Error())
	}

	var r io.Reader
	switch format {
	case "json":
		r = json_v2.Output(scoreCard)
	case "human":
		termWidth, _, err := term.GetSize(int(os.Stdin.Fd()))
		// Assume a width of 80 if it can't be detected
		if err != nil {
			termWidth = 80
		}
		body, err := human.Human(scoreCard, 0, termWidth, true)
		if err != nil {
			fmt.Println(err)
			return string(err.Error())
		}

		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			fmt.Println(err)
			return string(err.Error())
		}

		htmlBody := t2html.Render(bodyBytes)
		r = bytes.NewReader(htmlBody)
	case "ci":
		r = ci.CI(scoreCard)
	case "sarif":
		r = sarif.Output(scoreCard)
	case "junit":
		r = junit.JUnit(scoreCard)
	default:
		return fmt.Errorf("error: unknown format")
	}

	body, err := io.ReadAll(r)
	fmt.Println("body", body)
	if err != nil {
		fmt.Println(err)
		return string(err.Error())
	}

	return string(body)
}
