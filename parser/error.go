package parser

import (
	"fmt"
	"strings"
)

type parseError []error

func (p parseError) Error() string {
	var s []string
	for _, e := range p {
		s = append(s, e.Error())
	}
	return fmt.Sprintf(strings.Join(s, "\n"))
}

func (p *parseError) AddIfErr(err error) {
	if err != nil {
		*p = append(*p, err)
	}
}

func (p parseError) Any() bool {
	return len(p) > 0
}
