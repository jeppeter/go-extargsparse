package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var err error
	var loads = `        {
            "verbose|v" : "+",
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var confstr = fmt.Sprintf(fmt.Sprintf(`{"%s": "usage", "%s" : "?" , "%s" : "++", "%s" : "+"}`, extargsparse.OPT_HELP_LONG, extargsparse.OPT_HELP_SHORT, extargsparse.OPT_LONG_PREFIX, extargsparse.OPT_SHORT_PREFIX))
	var options *extargsparse.ExtArgsOptions
	var parser *extargsparse.ExtArgsParse
	var ns *extargsparse.NameSpaceEx

	options, err = extargsparse.NewExtArgsOptions(confstr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "options[%s] err[%s]\n", confstr, err.Error())
		os.Exit(5)
		return
	}
	parser, err = extargsparse.NewExtArgsParse(options, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "new args err[%s]\n", err.Error())
		os.Exit(5)
		return
	}
	err = parser.LoadCommandLineString(loads)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load[%s] err[%s]\n", loads, err.Error())
		os.Exit(5)
		return
	}

	ns, err = parser.ParseCommandLine(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse err[%s]\n", err.Error())
		os.Exit(5)
		return
	}

	fmt.Fprintf(os.Stdout, "subcommand=%s\n", ns.GetString("subcommand"))
	fmt.Fprintf(os.Stdout, "verbose=%d\n", ns.GetInt("verbose"))
	fmt.Fprintf(os.Stdout, "dep_list=%v\n", ns.GetArray("dep_list"))
	fmt.Fprintf(os.Stdout, "dep_string=%s\n", ns.GetString("dep_string"))
	fmt.Fprintf(os.Stdout, "subnargs=%v\n", ns.GetArray("subnargs"))
	fmt.Fprintf(os.Stdout, "args=%v\n", ns.GetArray("args"))

	return
}
