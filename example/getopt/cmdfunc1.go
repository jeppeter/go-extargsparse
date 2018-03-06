package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var loads = `{
		"verbose|v" : "+",
		"dep<dep_handler>" : {
			"cc|c" : ""
		},
		"rdep<rdep_handler>": {
			"dd|C" : ""
		}
	}`
	var err error
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var flag *extargsparse.ExtKeyParse
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{}`))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				flag, _ = parser.GetCmdKey("")
				fmt.Fprintf(os.Stdout, "main funcion:%s\n", flag.Function())

				flag, _ = parser.GetCmdKey("dep")
				fmt.Fprintf(os.Stdout, "dep funcion:%s\n", flag.Function())

				flag, _ = parser.GetCmdKey("rdep")
				fmt.Fprintf(os.Stdout, "rdep funcion:%s\n", flag.Function())
			}
		}
	}
	/*
		Notice:
			if options not set fmt.Sprintf(`{"%s" : false}`,extargsparse.OPT_FUNC_UPPER_CASE)
			the real function is Dep_handler and Rdep_handler
		Output:
			main funcion:
			dep funcion:dep_handler
			rdep funcion:rdep_handler
	*/
}
