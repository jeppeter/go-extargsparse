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
			"cc|c" : "",
			"$" : "+"
		},
		"rdep<rdep_handler>": {
			"dd|C" : "",
			"$" : "?"
		},
		"$port" : {
			"nargs" : 1,
			"type" : "int",
			"value" : 9000
		}
	}`
	var confstr string
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var flag *extargsparse.ExtKeyParse
	var opts []*extargsparse.ExtKeyParse
	confstr = `{}`
	options, _ = extargsparse.NewExtArgsOptions(confstr)
	parser, _ = extargsparse.NewExtArgsParse(options, nil)
	parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	opts, _ = parser.GetCmdOpts("")
	for _, flag = range opts {
		if flag.TypeName() == "args" {
			fmt.Fprintf(os.Stdout, "args.nargs=%v\n", flag.Nargs())
		} else if flag.FlagName() == "port" {
			fmt.Fprintf(os.Stdout, "port.nargs=%d\n", flag.Nargs().(int))
		}
	}

	opts, _ = parser.GetCmdOpts("dep")
	for _, flag = range opts {
		if flag.TypeName() == "args" {
			fmt.Fprintf(os.Stdout, "dep.args.nargs=%v\n", flag.Nargs())
		}
	}
	opts, _ = parser.GetCmdOpts("rdep")
	for _, flag = range opts {
		if flag.TypeName() == "args" {
			fmt.Fprintf(os.Stdout, "rdep.args.nargs=%v\n", flag.Nargs())
		}
	}
	/*
		Output:
			args.nargs=*
			port.nargs=1
			dep.args.nargs=+
			rdep.args.nargs=?
	*/
}
