package extargsparse

import (
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"unsafe"
)

func assert_test(ischeck bool, fmtstr string, a ...interface{}) {
	if !ischeck {
		s := fmt.Sprintf(fmtstr, a...)
		panic(s)
	}
}

func format_out_stack(level int) string {
	_, f, l, _ := runtime.Caller(level)
	return fmt.Sprintf("[%s:%d]", f, l)
}

func format_error_ex(level int, fmtstr string, a ...interface{}) string {
	s := format_out_stack(level + 1)
	s += fmt.Sprintf(fmtstr, a...)
	return s
}

func format_error(fmtstr string, a ...interface{}) string {
	return format_error_ex(2, fmtstr, a...)
}

func formatMap(kattr map[string]string) string {
	var s string = ""
	var ks []string
	var k, v string
	ks = make([]string, 0)
	for k, _ = range kattr {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k = range ks {
		v = kattr[k]
		s += fmt.Sprintf("[%s]=[%s]\n", k, v)
	}
	return s
}

func check_in_array(sarr []string, s string) bool {
	for _, k := range sarr {
		if k == s {
			return true
		}
	}
	return false
}

const (
	COMMAND_SET              = 10
	SUB_COMMAND_JSON_SET     = 20
	COMMAND_JSON_SET         = 30
	ENVIRONMENT_SET          = 40
	ENV_SUB_COMMAND_JSON_SET = 50
	ENV_COMMAND_JSON_SET     = 60
	DEFAULT_SET              = 70
)

func keyDebug(fmtstr string, a ...interface{}) {
	s := format_out_stack(2)
	s += fmt.Sprintf(fmtstr, a...)
	fmt.Printf("%s\n", s)
	return
}

func findInterfaceField(a *reflect.Value, key string) int {
	var maxfld int
	var pt reflect.Type
	var i int
	maxfld = a.NumField()
	pt = a.Type()
	for i = 0; i < maxfld; i++ {
		if pt.Field(i).Name == key {
			return i
		}
	}
	return -1
}

func setMemberValueInner(a *reflect.Value, name string, value interface{}) error {
	var sarr []string
	var idx int
	var rf, vrf reflect.Value
	var s, rs string
	sarr = strings.Split(name, ".")
	if len(sarr) == 0 {
		return fmt.Errorf("%s", format_error("not set name []"))
	}

	if len(sarr) == 1 {
		idx = findInterfaceField(a, name)
		if idx < 0 {
			return fmt.Errorf("%s", format_error("can not find [%s]", name))
		}
		rf = a.Field(idx)
		rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
		s = rf.Type().String()
		vrf = reflect.ValueOf(value)
		if vrf.Kind() == reflect.Interface || vrf.Kind() == reflect.Ptr {
			rs = vrf.Elem().Type().String()
		} else {
			rs = vrf.Type().String()
		}

		if s != rs {
			return fmt.Errorf("%s", format_error("rf type[%s] != value[%v] type[%s]", s, value, rs))
		}
		rf.Set(reflect.ValueOf(value))
		return nil
	}

	idx = findInterfaceField(a, sarr[0])
	if idx < 0 {
		return fmt.Errorf("%s", format_error("can not find [%s] for [%s]", sarr[0], name))
	}
	rf = a.Field(idx)
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	return setMemberValueInner(&rf, strings.Join(sarr[1:], "."), value)
}

func setMemberValue(a interface{}, name string, value interface{}) error {
	var rf reflect.Value
	rf = reflect.ValueOf(a).Elem()
	return setMemberValueInner(&rf, name, value)
}
