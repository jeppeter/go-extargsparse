package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

type ArgStruct struct {
	Verbose int
	Removed bool
	Floatv  float64
	Intv    int
	Arrl    []string
	Strv    string
	Args    []string
}

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
	var p *ArgStruct
	var err error

	p = &ArgStruct{}

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

	if len(os.Args[1:]) == 0 {
		_, err = parser.ParseCommandLineEx([]string{"-vvvv", "cc", "-f", "33.2", "--arrl", "wwwe", "-s", "3993"}, nil, p, nil)
	} else {
		_, err = parser.ParseCommandLineEx(nil, nil, p, nil)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parser command err[%s]\n", err.Error())
		os.Exit(5)
		return
	}
	fmt.Fprintf(os.Stdout, "verbose=%d\n", p.Verbose)
	fmt.Fprintf(os.Stdout, "removed=%v\n", p.Removed)
	fmt.Fprintf(os.Stdout, "falotv=%f\n", p.Floatv)
	fmt.Fprintf(os.Stdout, "intv=%d\n", p.Intv)
	fmt.Fprintf(os.Stdout, "arrl=%v\n", p.Arrl)
	fmt.Fprintf(os.Stdout, "strv=%s\n", p.Strv)
	fmt.Fprintf(os.Stdout, "args=%v\n", p.Args)
	return
}
