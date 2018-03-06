package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var options *extargsparse.ExtArgsOptions
	var err error
	var confstr = fmt.Sprintf(`{
		"%s" : 90.0		
		}`, extargsparse.OPT_SCREEN_WIDTH)
	options, err = extargsparse.NewExtArgsOptions(confstr)
	if err == nil {
		fmt.Fprintf(os.Stdout, "screenwidth=%d\n", options.GetInt(extargsparse.OPT_SCREEN_WIDTH)) // screenwidth=90
	}
	return
}
