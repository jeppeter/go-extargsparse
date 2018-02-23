package extargsparse

import (
	"reflect"
	"testing"
)

func check_equal(t *testing.T, orig, check interface{}) {
	if !reflect.DeepEqual(orig, check) {
		t.Fatalf("%s[%s] orig [%v] != check[%v]", format_out_stack(2), t.Name(), orig, check)
	}
}

func check_not_equal(t *testing.T, orig, check interface{}) {
	if reflect.DeepEqual(orig, check) {
		t.Fatalf("%s[%s] orig [%v] == check[%v]", format_out_stack(2), t.Name(), orig, check)
	}
}
