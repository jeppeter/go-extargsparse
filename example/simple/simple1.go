package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var commandline = `{
		"verbose|v" : "+",
		"removed|R" : false,
		"floatv|f" : 3.3,
		"intv|i" : 5,
		"arrl|a" : [],
		"strv|s" : null,
		"$" : "+"
		}`
	var parser *extargsparse.ExtArgsParse
	var ns *extargsparse.NameSpaceEx
	var err error

	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not init parser err [%s]\n", err.Error())
		os.Exit(5)
		return
	}

	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse [%s] error[%s]\n", commandline, err.Error())
		os.Exit(5)
		return
	}

	ns, err = parser.ParseCommandLine([]string{"-vvvv", "cc", "-f", "33.2", "--arrl", "wwwe", "-s", "3993"}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parser command err[%s]\n", err.Error())
		os.Exit(5)
		return
	}
	fmt.Fprintf(os.Stdout, "verbose=%d\n", ns.GetInt("verbose"))
	fmt.Fprintf(os.Stdout, "removed=%v\n", ns.GetBool("removed"))
	fmt.Fprintf(os.Stdout, "falotv=%f\n", ns.GetFloat("floatv"))
	fmt.Fprintf(os.Stdout, "intv=%d\n", ns.GetInt("intv"))
	fmt.Fprintf(os.Stdout, "arrl=%v\n", ns.GetArray("arrl"))
	fmt.Fprintf(os.Stdout, "strv=%s\n", ns.GetString("strv"))
	fmt.Fprintf(os.Stdout, "args=%v\n", ns.GetArray("args"))
	return
}
