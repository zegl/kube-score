package parser

import (
	"strings"
)

type parseErrors []error

func (p parseErrors) Error() string {
	var s []string
	for _, e := range p {
		s = append(s, e.Error())
	}
	return strings.Join(s, "\n")
}

func (p *parseErrors) AddIfErr(err error) {
	if err != nil {
		*p = append(*p, err)
	}
}

func (p parseErrors) Any() bool {
	return len(p) > 0
}
