package extargsparse

import (
	"fmt"
)

type helpSize struct {
	logObject
	intvalue map[string]int
}

var helpSizeKeywords = []string{"optnamesize", "optexprsize", "opthelpsize", "cmdnamesize", "cmdhelpsize"}

func newHelpSize() *helpSize {
	p := &helpSize{logObject: *newLogObject("extargsparse"), intvalue: make(map[string]int)}
	for _, k := range helpSizeKeywords {
		p.intvalue[k] = 0
	}
	return p
}

func (p *helpSize) GetValue(name string) int {
	var i int = 0
	var ok bool
	i, ok = p.intvalue[name]
	if !ok {
		return 0
	}

	return i
}

func (p *helpSize) SetValue(name string, value int) {
	var i int = 0
	var ok bool
	i, ok = p.intvalue[name]
	if !ok {
		return
	}
	if i < value {
		p.intvalue[name] = value
	}
	return

}

func (p *helpSize) Format() string {
	var s string
	var cnt int
	s += "{"
	cnt = 0
	for _, k := range helpSizeKeywords {
		if cnt > 0 {
			s += fmt.Sprintf(";")
		}
		s += fmt.Sprintf("%s=%d", k, p.intvalue[k])
	}
	s += "}"
	return s
}
