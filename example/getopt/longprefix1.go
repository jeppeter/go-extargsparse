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
		if flag.TypeName() == "count" && flag.FlagName() == "verbose" {
			fmt.Fprintf(os.Stdout, "longprefix=%s\n", flag.LongPrefix())
			fmt.Fprintf(os.Stdout, "longopt=%s\n", flag.Longopt())
			fmt.Fprintf(os.Stdout, "shortopt=%s\n", flag.Shortopt())
		}
	}

	confstr = fmt.Sprintf(`{ "%s" : "++", "shortprefix" : "+"}`, extargsparse.OPT_LONG_PREFIX)
	options, _ = extargsparse.NewExtArgsOptions(confstr)
	parser, _ = extargsparse.NewExtArgsParse(options, nil)
	parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	opts, _ = parser.GetCmdOpts("")
	for _, flag = range opts {
		if flag.TypeName() == "count" && flag.FlagName() == "verbose" {
			fmt.Fprintf(os.Stdout, "longprefix=%s\n", flag.LongPrefix())
			fmt.Fprintf(os.Stdout, "longopt=%s\n", flag.Longopt())
			fmt.Fprintf(os.Stdout, "shortopt=%s\n", flag.Shortopt())
		}
	}

	/*
		Output:
			longprefix=--
			longopt=--verbose
			shortopt=-v
			longprefix=++
			longopt=++verbose
			shortopt=+v
	*/
}
