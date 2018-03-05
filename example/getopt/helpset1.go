package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var loads = `{
		"verbose|v##we used verbose##" : "+",
		"$##this is args help##" : "*",
		"dep##dep help set##" : {
			"cc|c##cc sss##" : "",
			"$##this is dep subnargs help##" : "*"
		},
		"rdep##can not set rdep help##": {
			"dd|C##capital C##" : "",
			"$##this is rdep subnargs help##" : "*"
		}
	}`
	var err error
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var flag *extargsparse.ExtKeyParse
	var opts []*extargsparse.ExtKeyParse
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{}`))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				flag, _ = parser.GetCmdKey("")
				fmt.Fprintf(os.Stdout, "main helpinfo:%s\n", flag.HelpInfo())

				flag, _ = parser.GetCmdKey("dep")
				fmt.Fprintf(os.Stdout, "dep helpinfo:%s\n", flag.HelpInfo())

				flag, _ = parser.GetCmdKey("rdep")
				fmt.Fprintf(os.Stdout, "rdep helpinfo:%s\n", flag.HelpInfo())

				opts, _ = parser.GetCmdOpts("")
				for _, flag = range opts {
					if flag.TypeName() == "args" {
						fmt.Fprintf(os.Stdout, "main.args.HelpInfo=%s\n", flag.HelpInfo())
					} else if flag.FlagName() == "verbose" {
						fmt.Fprintf(os.Stdout, "verbose.HelpInfo=%s\n", flag.HelpInfo())
					}
				}
				opts, _ = parser.GetCmdOpts("dep")
				for _, flag = range opts {
					if flag.TypeName() == "args" {
						fmt.Fprintf(os.Stdout, "dep.subnargs.HelpInfo=%s\n", flag.HelpInfo())
					} else if flag.FlagName() == "cc" {
						fmt.Fprintf(os.Stdout, "dep.cc.HelpInfo=%s\n", flag.HelpInfo())
					}
				}
				opts, _ = parser.GetCmdOpts("rdep")
				for _, flag = range opts {
					if flag.TypeName() == "args" {
						fmt.Fprintf(os.Stdout, "rdep.subnargs.HelpInfo=%s\n", flag.HelpInfo())
					} else if flag.FlagName() == "dd" {
						fmt.Fprintf(os.Stdout, "rdep.dd.HelpInfo=%s\n", flag.HelpInfo())
					}
				}

			}
		}
	}
	/*
		Notice:
			HelpInfo is part of between ##(information)##
		Output:
			main helpinfo:
			dep helpinfo:dep help set
			rdep helpinfo:can not set rdep help
			main.args.HelpInfo=this is args help
			verbose.HelpInfo=we used verbose
			dep.subnargs.HelpInfo=this is dep subnargs help
			dep.cc.HelpInfo=cc sss
			rdep.subnargs.HelpInfo=this is rdep subnargs help
			rdep.dd.HelpInfo=capital C
	*/
}
