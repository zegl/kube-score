//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"io"
	"strings"
	"syscall/js"

	t2html "github.com/buildkite/terminal-to-html"
	"github.com/fatih/color"
	"github.com/zegl/kube-score/config"
	"github.com/zegl/kube-score/domain"
	"github.com/zegl/kube-score/parser"
	"github.com/zegl/kube-score/renderer/human"
	"github.com/zegl/kube-score/score"
	"github.com/zegl/kube-score/score/checks"
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
	if len(inputs) == 0 {
		fmt.Println("Unexpected number of arguments")
		return "Unexpected number of arguments"
	}

	fmt.Println(inputs[0].String())

	reader := &inputReader{
		Reader: strings.NewReader(inputs[0].String()),
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

	card, err := score.Score(allObjs, allChecks, runConfig)
	if err != nil {
		fmt.Println(err)
		return string(err.Error())
	}

	color.NoColor = false
	output, err := human.Human(card, 0, 110, true)
	if err != nil {
		fmt.Println(err)
		return string(err.Error())
	}

	body, err := io.ReadAll(output)
	fmt.Println("body", body)
	if err != nil {
		fmt.Println(err)
		return string(err.Error())
	}

	htmlBody := t2html.Render(body)

	return string(htmlBody)
}
