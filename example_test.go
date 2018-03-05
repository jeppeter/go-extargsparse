package extargsparse_test

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func ExampleNewExtArgsOptions() {
	var options *extargsparse.ExtArgsOptions
	var err error
	var confstr = `{
		"screenwidth" : 80.0		
		}`
	options, err = extargsparse.NewExtArgsOptions(confstr)
	if err == nil {
		fmt.Fprintf(os.Stdout, "screenwidth=%d\n", options.GetInt("screenwidth"))
		// Output screenwidth=80
	}
	return
}

func ExampleExtArgsOptions_SetValue() {
	var options *extargsparse.ExtArgsOptions
	var err error
	options, err = extargsparse.NewExtArgsOptions(`{}`)
	if err == nil {
		options.SetValue("screenwidth", float64(90.0))
		fmt.Fprintf(os.Stdout, "screenwidth=%d\n", options.GetInt("screenwidth"))
		// Output screenwidth=90
	}
	return
}
