package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var loads = `{
		"dep" : {

		},
		"rdep": {

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
				fmt.Fprintf(os.Stdout, "cmdname=%s\n", flag.CmdName())
				flag, _ = parser.GetCmdKey("dep")
				fmt.Fprintf(os.Stdout, "cmdname=%s\n", flag.CmdName())
				flag, _ = parser.GetCmdKey("rdep")
				fmt.Fprintf(os.Stdout, "cmdname=%s\n", flag.CmdName())
				/*
					Output:
					cmdname=main
					cmdname=dep
					cmdname=ip
				*/

			}
		}
	}
}
