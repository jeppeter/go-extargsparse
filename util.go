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

func format_error(level int, fmtstr string, a ...interface{}) string {
	s := format_out_stack(level + 1)
	s += fmt.Sprintf(fmtstr, a...)
	return s
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
