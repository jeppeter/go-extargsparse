package extargsparse

import (
	"encoding/json"
	"fmt"
	"sort"
)

type ExtArgsOptions struct {
	logger *logObject
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
	VAR_UPPER_CASE:   true,
	FUNC_UPPER_CASE:  true,
}

// to set the value of k and v
//    it almost the direct but one case with float in the type of no small part ,it will return int set
func (p *ExtArgsOptions) SetValue(k string, v interface{}) error {
	var retv interface{}
	var iv int
	var bret bool
	retv = v
	switch v.(type) {
	case float64:
		iv, bret = isFloatToInt(v.(float64))
		if bret {
			retv = iv
		}
	case float32:
		iv, bret = isFloatToInt(float64(v.(float32)))
		if bret {
			retv = iv
		}
	}

	p.values[k] = retv
	return nil
}

// get the value of key ,if not set return nil, otherwise return the interface{}
//    it gives the caller to check the type
func (p *ExtArgsOptions) GetValue(k string) interface{} {
	var v interface{}
	var ok bool
	v, ok = p.values[k]
	if !ok {
		return nil
	}
	return v
}

// get the value of key , if it is not set or not string type it will return ""
func (p *ExtArgsOptions) GetString(k string) string {
	var v interface{}
	v = p.GetValue(k)
	if v == nil {
		return ""
	}
	switch v.(type) {
	case string:
		return v.(string)
	}
	return fmt.Sprintf("%v", v)
}

// get the value of key ,if it is not set or not bool type it will return false
func (p *ExtArgsOptions) GetBool(k string) bool {
	var v interface{}
	v = p.GetValue(k)
	if v == nil {
		return false
	}

	switch v.(type) {
	case bool:
		return v.(bool)
	}
	return false
}

// get the value of key ,if it is not set or not int type it will return 0
func (p *ExtArgsOptions) GetInt(k string) int {
	var v interface{}
	v = p.GetValue(k)
	if v == nil {
		return 0
	}

	switch v.(type) {
	case int:
		return v.(int)
	case int32:
		return int(v.(int32))
	case int64:
		return int(v.(int64))
	case uint32:
		return int(v.(uint32))
	case uint64:
		return int(v.(uint64))
	}
	return 0
}

// Format to give the value in the string format split by ; it is by debug used
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

// NewExtArgsOptions for create new options for *ExtArgsParse
//    s is the json file ,
//    key               default value
//    "prog":           ""
//    "usage":          ""
//    "description":    ""
//    "epilog":         ""
//    "version":        "0.0.1"
//    "errorhandler":   "exit"
//    "helphandler":    nil
//    "longprefix":     "--"
//    "shortprefix":    "-"
//    "nohelpoption":   false
//    "nojsonoption":   false
//    "helplong":       "help"
//    "helpshort":      "h"
//    "jsonlong":       "json"
//    "cmdprefixadded": true
//    "parseall":       true
//    "screenwidth":    80
//    "flagnochange":   false
//    VAR_UPPER_CASE:   true
//    FUNC_UPPER_CASE:  true
func NewExtArgsOptions(s string) (p *ExtArgsOptions, err error) {
	var vmap map[string]interface{}
	var k string
	var v interface{}

	p = nil
	err = json.Unmarshal([]byte(s), &vmap)
	if err != nil {
		err = fmt.Errorf("%s", format_error("parse [%s] error[%s]", err.Error()))
		return
	}

	p = &ExtArgsOptions{logger: newLogObject("extargsparse"), values: make(map[string]interface{})}
	for k, v = range opt_default_VALUE {
		p.SetValue(k, v)
	}

	for k, v = range vmap {
		p.SetValue(k, v)
	}
	return
}
