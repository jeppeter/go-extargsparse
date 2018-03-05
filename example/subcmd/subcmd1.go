package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func dep_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
	if ns == nil {
		return nil
	}
	fmt.Fprintf(os.Stdout, "subcommand=%s\n", ns.GetString("subcommand"))
	fmt.Fprintf(os.Stdout, "verbose=%d\n", ns.GetInt("verbose"))
	fmt.Fprintf(os.Stdout, "dep_list=%v\n", ns.GetArray("dep_list"))
	fmt.Fprintf(os.Stdout, "dep_str=%s\n", ns.GetString("dep_str"))
	fmt.Fprintf(os.Stdout, "subnargs=%v\n", ns.GetArray("subnargs"))
	fmt.Fprintf(os.Stdout, "rdep_list=%v\n", ns.GetArray("rdep_list"))
	fmt.Fprintf(os.Stdout, "rdep_str=%s\n", ns.GetString("rdep_str"))
	os.Exit(0)
	return nil
}

func rdep_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
	if ns == nil {
		return nil
	}
	fmt.Fprintf(os.Stdout, "subcommand=%s\n", ns.GetString("subcommand"))
	fmt.Fprintf(os.Stdout, "verbose=%d\n", ns.GetInt("verbose"))
	fmt.Fprintf(os.Stdout, "dep_list=%v\n", ns.GetArray("dep_list"))
	fmt.Fprintf(os.Stdout, "dep_str=%s\n", ns.GetString("dep_str"))
	fmt.Fprintf(os.Stdout, "subnargs=%v\n", ns.GetArray("subnargs"))
	fmt.Fprintf(os.Stdout, "rdep_list=%v\n", ns.GetArray("rdep_list"))
	fmt.Fprintf(os.Stdout, "rdep_str=%s\n", ns.GetString("rdep_str"))
	os.Exit(0)
	return nil
}

func init() {
	dep_handler(nil, nil, nil)
	rdep_handler(nil, nil, nil)
}

func main() {
	var commandline = `{
		"verbose|v" : "+",
		"dep<dep_handler>" : {
			"$" : "*",
			"list|L" :  [],
			"str|S" : ""
		},
		"rdep<rdep_handler>" : {
			"$" : "*",
			"list|l" : [],
			"str|s" : ""
		}
		}`
	var parser *extargsparse.ExtArgsParse
	var err error
	var options *extargsparse.ExtArgsOptions
	var confstr = fmt.Sprintf(`{ "%s" : false}`, extargsparse.FUNC_UPPER_CASE)
	options, err = extargsparse.NewExtArgsOptions(confstr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not make string [%s] err[%s]\n", confstr, err.Error())
		os.Exit(5)
		return
	}
	parser, err = extargsparse.NewExtArgsParse(options, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not make parser err[%s]\n", err.Error())
		os.Exit(5)
		return
	}

	err = parser.LoadCommandLineString(commandline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not load string [%s] err[%s]\n", commandline, err.Error())
		os.Exit(5)
		return
	}

	_, err = parser.ParseCommandLineEx(nil, nil, nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse err[%s]\n", err.Error())
		os.Exit(5)
		return
	}
	return
}
