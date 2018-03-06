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
		options.SetValue(extargsparse.OPT_SCREEN_WIDTH, float64(100.0))
		fmt.Fprintf(os.Stdout, "screenwidth=%d\n", options.GetInt(extargsparse.OPT_SCREEN_WIDTH)) //screenwidth=100
	}
	return
}
