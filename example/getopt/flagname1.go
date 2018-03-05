package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var loads = `{
		"verbose|v" : "+",
		"dep" : {
			"cc|c" : ""
		},
		"rdep": {
			"dd|C" : ""
		}
	}`
	var err error
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var flag *extargsparse.ExtKeyParse
	var opts []*extargsparse.ExtKeyParse
	var finded bool
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{}`))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				opts, _ = parser.GetCmdOpts("")
				for _, flag = range opts {
					if flag.FlagName() == "verbose" {
						fmt.Fprintf(os.Stdout, "flagname=%s\n", flag.FlagName())
					}
				}

				opts, _ = parser.GetCmdOpts("dep")
				for _, flag = range opts {
					if flag.FlagName() == "cc" {
						fmt.Fprintf(os.Stdout, "flagname=%s\n", flag.FlagName())
					}
				}
				opts, _ = parser.GetCmdOpts("rdep")
				for _, flag = range opts {
					if flag.FlagName() == "dd" {
						fmt.Fprintf(os.Stdout, "flagname=%s\n", flag.FlagName())
					}
				}

				opts, _ = parser.GetCmdOpts("rdep")
				finded = false
				for _, flag = range opts {
					if flag.FlagName() == "cc" {
						fmt.Fprintf(os.Stdout, "flagname=%s\n", flag.FlagName())
						finded = true
					}
				}

				if !finded {
					fmt.Fprintf(os.Stdout, "can not found cc for rdep cmd\n")
				}

			}
		}
	}
	/*
		Output:
		flagname=verbose
		flagname=cc
		flagname=dd
		can not found cc for rdep cmd
	*/
}
