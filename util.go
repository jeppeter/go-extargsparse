package extargsparse

import (
	"fmt"
	"runtime"
	"sort"
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
