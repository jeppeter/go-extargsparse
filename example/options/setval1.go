package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var options *extargsparse.ExtArgsOptions
	var err error
	var confstr = `{
		"screenwidth" : 90.0		
		}`
	options, err = extargsparse.NewExtArgsOptions(confstr)
	if err == nil {
		fmt.Fprintf(os.Stdout, "screenwidth=%d\n", options.GetInt("screenwidth")) // screenwidth=90
	}
	return

}
