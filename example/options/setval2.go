package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var options *extargsparse.ExtArgsOptions
	var err error
	options, err = extargsparse.NewExtArgsOptions(`{}`)
	if err == nil {
		options.SetValue("screenwidth", float64(100.0))
		fmt.Fprintf(os.Stdout, "screenwidth=%d\n", options.GetInt("screenwidth")) //screenwidth=100
	}
	return
}
