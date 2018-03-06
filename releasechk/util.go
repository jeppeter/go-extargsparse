package main

import (
	"fmt"
	"runtime"
)

func format_out_stack(level int) string {
	_, f, l, _ := runtime.Caller(level)
	return fmt.Sprintf("[%s:%d]", f, l)
}

func format_error_ex(level int, fmtstr string, a ...interface{}) string {
	s := format_out_stack(level + 1)
	s += fmt.Sprintf(fmtstr, a...)
	return s
}

func format_error(fmtstr string, a ...interface{}) error {
	return fmt.Errorf("%s", format_error_ex(2, fmtstr, a...))
}

func format_length(l int, fmtstr string, a ...interface{}) string {
	var s = fmt.Sprintf(fmtstr, a...)
	for len(s) < l {
		s += " "
	}
	return s
}
