package extargsparse

import (
	"encoding/json"
	"fmt"
	"sort"
)

type ExtArgsOptions struct {
	logObject
	values map[string]interface{}
}

var opt_default_VALUE = map[string]interface{}{
	"prog":           "",
	"usage":          "",
	"description":    "",
	"epilog":         "",
	"version":        "0.0.1",
	"errorhandler":   "exit",
	"helphandler":    nil,
	"longprefix":     "--",
	"shortprefix":    "-",
	"nohelpoption":   false,
	"nojsonoption":   false,
	"helplong":       "help",
	"helpshort":      "h",
	"jsonlong":       "json",
	"cmdprefixadded": true,
	"parseall":       true,
	"screenwidth":    80,
	"flagnochange":   false,
}

func (p *ExtArgsOptions) SetValue(k string, v interface{}) error {
	p.values[k] = v
	return nil
}

func (p *ExtArgsOptions) GetValue(k string) interface{} {
	var v interface{}
	var ok bool
	v, ok = p.values[k]
	if !ok {
		return nil
	}
	return v
}

func (p *ExtArgsOptions) Format() string {
	var keys []string
	var s string = ""
	var k string
	var cnt int = 0

	keys = make([]string, 0)

	for k, _ = range p.values {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	s += fmt.Sprintf("{")
	cnt = 0
	for _, k = range keys {
		if cnt > 0 {
			s += fmt.Sprintf(";")
		}
		s += fmt.Sprintf("[%s]=[%v]", k, p.values[k])
		cnt++
	}
	s += fmt.Sprintf("}")
	return s
}

func NewExtArgsOptions(s string) (p *ExtArgsOptions, err error) {
	var vmap map[string]interface{}
	var k string
	var v interface{}

	p = nil
	err = json.Unmarshal([]byte(s), &vmap)
	if err != nil {
		err = fmt.Errorf("%s", format_error(1, "parse [%s] error[%s]", err.Error()))
		return
	}

	p = &ExtArgsOptions{logObject: *newLogObject("extargsparse"), values: make(map[string]interface{})}
	for k, v = range opt_default_VALUE {
		p.SetValue(k, v)
	}

	for k, v = range vmap {
		p.SetValue(k, v)
	}
	return
}
