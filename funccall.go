package extargsparse

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

// getFunc gets the function defined by the given fully-qualified name. The
// outFuncPtr parameter should be a pointer to a function with the appropriate
// type (e.g. the address of a local variable), and is set to a new function
// value that calls the specified function. If the specified function does not
// exist, outFuncPtr is not set and an error is returned.
func getFunc(outFuncPtr interface{}, name string) error {
	codePtr, err := findFuncWithName(name)
	if err != nil {
		return err
	}
	createFuncForCodePtr(outFuncPtr, codePtr)
	return nil
}

// Convenience struct for modifying the underlying code pointer of a function
// value. The actual struct has other values, but always starts with a code
// pointer.
type funcType struct {
	codePtr uintptr
}

// createFuncForCodePtr is given a code pointer and creates a function value
// that uses that pointer. The outFun argument should be a pointer to a function
// of the proper type (e.g. the address of a local variable), and will be set to
// the result function value.
func createFuncForCodePtr(outFuncPtr interface{}, codePtr uintptr) {
	outFuncVal := reflect.ValueOf(outFuncPtr).Elem()
	// Use reflect.MakeFunc to create a well-formed function value that's
	// guaranteed to be of the right type and guaranteed to be on the heap
	// (so that we can modify it). We give a nil delegate function because
	// it will never actually be called.
	newFuncVal := reflect.MakeFunc(outFuncVal.Type(), nil)
	// Use reflection on the reflect.Value (yep!) to grab the underling
	// function value pointer. Trying to call newFuncVal.Pointer() wouldn't
	// work because it gives the code pointer rather than the function value
	// pointer. The function value is a struct that starts with its code
	// pointer, so we can swap out the code pointer with our desired value.
	funcValuePtr := reflect.ValueOf(newFuncVal).FieldByName("ptr").Pointer()
	funcPtr := (*funcType)(unsafe.Pointer(funcValuePtr))
	funcPtr.codePtr = codePtr
	outFuncVal.Set(newFuncVal)
}

// findFuncWithName searches through the moduledata table created by the linker
// and returns the function's code pointer. If the function was not found, it
// returns an error. Since the data structures here are not exported, we copy
// them below (and they need to stay in sync or else things will fail
// catastrophically).
func findFuncWithName(name string) (uintptr, error) {
	for moduleData := &wFirstmoduledata; moduleData != nil; moduleData = moduleData.next {
		for i, ftab := range moduleData.ftab {
			if i < (len(moduleData.ftab) - 1) {
				f := (*runtime.Func)(unsafe.Pointer(&moduleData.pclntable[ftab.funcoff]))
				if f.Name() == name {
					return f.Entry(), nil
				}
			}
		}
	}
	return 0, fmt.Errorf("Invalid function name: %s", name)
}

// Everything below is taken from the runtime package, and must stay in sync
// with it.

//go:linkname wFirstmoduledata runtime.firstmoduledata
var wFirstmoduledata moduledata

type moduledata struct {
	pclntable    []byte
	ftab         []functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr

	// Original type was []*_type
	typelinks []interface{}

	modulename string
	// Original type was []modulehash
	modulehashes []interface{}

	gcdatamask, gcbssmask bitvector

	next *moduledata
}

type functab struct {
	entry   uintptr
	funcoff uintptr
}

type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}
