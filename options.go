package extargsparse

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

type ExtArgsOptions struct {
	logObject
	strValue  map[string]string
	boolValue map[string]bool
	intValue  map[string]int
}

var opt_default_VALUE = map[string]interface{}{
	"prog":  "",
	"usage": "",
}

func (p *ExtArgsOptions) SetValue(k string, v interface{}) error {
	fmt.Printf("[%s]=[%v]\n", k, v)
	switch v.(type) {
	case string:
		if v != nil {
			p.strValue[k] = fmt.Sprintf("%v", v)
		} else {
			p.strValue[k] = ""
		}
	case int:
		p.intValue[k] = v.(int)
	case bool:
		p.boolValue[k] = v.(bool)
	default:
		p.strValue[k] = fmt.Sprintf("%v", v)
	}
	return nil
}

func (p *ExtArgsOptions) GetValue(k string) interface{} {
	var v interface{}
	var k2 string
	for k2, _ = range p.strValue {
		if k2 == k {
			v = p.strValue[k]
			return v
		}
	}

	for k2, _ = range p.intValue {
		if k2 == k {
			v = p.intValue[k]
			return v
		}
	}

	for k2, _ = range p.boolValue {
		if k2 == k {
			v = p.boolValue[k]
			return v
		}
	}

	return nil
}

func (p *ExtArgsOptions) Format() string {
	var skeys []string
	var ikeys []string
	var bkeys []string
	var s string = ""
	var k string
	var cnt int = 0

	skeys = make([]string, 0)
	ikeys = make([]string, 0)
	bkeys = make([]string, 0)

	for k, _ = range p.strValue {
		skeys = append(skeys, k)
	}

	for k, _ = range p.intValue {
		ikeys = append(ikeys, k)
	}

	for k, _ = range p.boolValue {
		bkeys = append(bkeys, k)
	}

	sort.Strings(skeys)
	sort.Strings(ikeys)
	sort.Strings(bkeys)

	s += fmt.Sprintf("{")
	cnt = 0
	for _, k = range skeys {
		if cnt > 0 {
			s += fmt.Sprintf(";")
		}
		s += fmt.Sprintf("[%s]=[%s]", k, p.strValue[k])
		cnt++
	}

	for _, k = range ikeys {
		if cnt > 0 {
			s += fmt.Sprintf(";")
		}
		s += fmt.Sprintf("[%s]=[%d]", k, p.intValue[k])
		cnt++
	}

	for _, k = range bkeys {
		if cnt > 0 {
			s += fmt.Sprintf(";")
		}
		s += fmt.Sprintf("[%s]=[%v]", k, p.boolValue[k])
		cnt++
	}
	s += fmt.Sprintf("}")
	return s
}

func NewExtArgsOptions(s string) (p *ExtArgsOptions, err error) {
	var v interface{}
	var vmap map[string]interface{}
	var k string
	var emptymap interface{}

	p = nil
	err = json.Unmarshal([]byte(`{}`), &emptymap)
	if err != nil {
		err = fmt.Errorf("%s", format_error(1, "parse [%s] error[%s]", err.Error()))
		return
	}
	err = json.Unmarshal([]byte(s), &v)
	if err != nil {
		err = fmt.Errorf("%s", format_error(1, "parse [%s] error[%s]", err.Error()))
		return
	}

	p = &ExtArgsOptions{logObject: *newLogObject("extargsparse"), strValue: make(map[string]string), intValue: make(map[string]int), boolValue: make(map[string]bool)}
	for k, v = range opt_default_VALUE {
		p.SetValue(k, v)
	}

	fmt.Printf("%s", format_error(1, "\n"))
	switch v.(type) {
	case map[string]interface{}:
		vmap = v.(map[string]interface{})
		fmt.Printf("%s", format_error(1, "\n"))
	default:
		fmt.Printf("%s", format_error(1, "emptymap [%s]\n", reflect.ValueOf(emptymap).Kind().String()))
		if reflect.DeepEqual(emptymap, v) {
			return p, nil
		}
		return nil, fmt.Errorf("%s", format_error(1, "type [%s] not supported for [%s]", reflect.ValueOf(v).Kind().String(), s))
	}

	for k, v = range vmap {
		p.SetValue(k, v)
	}
	return
}
