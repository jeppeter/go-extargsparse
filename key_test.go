package extargsparse

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func format_out_stack(level int) string {
	_, f, l, _ := runtime.Caller(level)
	return fmt.Sprintf("[%s:%d]", f, l)
}

func check_equal(t *testing.T, orig, check interface{}) {
	if !reflect.DeepEqual(orig, check) {
		t.Fatalf("%s[%s] orig [%v] check[%v]", format_out_stack(2), t.Name(), orig, check)
	}
}

func Test_A001(t *testing.T) {
	flags, err := NewExtKeyParse_short("", "$flag|f+type", "string", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "flag")
	check_equal(t, flags.Longopt(), "--type-flag")
	check_equal(t, flags.Shortopt(), "-f")
	check_equal(t, flags.Optdest(), "type_flag")
	check_equal(t, flags.Value(), "string")
	check_equal(t, flags.ShortFlag(), "f")
	check_equal(t, flags.Prefix(), "type")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.VarName(), "type_flag")
	return
}
