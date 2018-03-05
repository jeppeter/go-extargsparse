package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

type SubcmdStruct struct {
	Verbose int
	Pair    []string
	Dep     struct {
		List     []string
		Str      string
		Subnargs []string
	}
	Rdep struct {
		List     []string
		Str      string
		Subnargs []string
	}
	Args []string
}

func dep_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
	var p *SubcmdStruct
	if ns == nil {
		return nil
	}
	p = ostruct.(*SubcmdStruct)
	fmt.Fprintf(os.Stdout, "subcommand=%s\n", ns.GetString("subcommand"))
	fmt.Fprintf(os.Stdout, "verbose=%d\n", p.Verbose)
	fmt.Fprintf(os.Stdout, "pair=%v\n", p.Pair)
	fmt.Fprintf(os.Stdout, "dep_list=%v\n", p.Dep.List)
	fmt.Fprintf(os.Stdout, "dep_str=%s\n", p.Dep.Str)
	fmt.Fprintf(os.Stdout, "subnargs=%v\n", p.Dep.Subnargs)
	fmt.Fprintf(os.Stdout, "rdep_list=%v\n", p.Rdep.List)
	fmt.Fprintf(os.Stdout, "rdep_str=%s\n", p.Rdep.Str)
	os.Exit(0)
	return nil
}

func rdep_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
	var p *SubcmdStruct
	if ns == nil {
		return nil
	}
	p = ostruct.(*SubcmdStruct)
	fmt.Fprintf(os.Stdout, "subcommand=%s\n", ns.GetString("subcommand"))
	fmt.Fprintf(os.Stdout, "verbose=%d\n", p.Verbose)
	fmt.Fprintf(os.Stdout, "pair=%v\n", p.Pair)
	fmt.Fprintf(os.Stdout, "dep_list=%v\n", p.Dep.List)
	fmt.Fprintf(os.Stdout, "dep_str=%s\n", p.Dep.Str)
	fmt.Fprintf(os.Stdout, "subnargs=%v\n", p.Dep.Subnargs)
	fmt.Fprintf(os.Stdout, "rdep_list=%v\n", p.Rdep.List)
	fmt.Fprintf(os.Stdout, "rdep_str=%s\n", p.Rdep.Str)
	os.Exit(0)
	return nil
}

func pair_key_handle(ns *extargsparse.NameSpaceEx, validx int, keycls *extargsparse.ExtKeyParse, params []string) (step int, err error) {
	var sarr []string
	if ns == nil {
		return 0, nil
	}
	if len(params) < (validx + 2) {
		return 0, fmt.Errorf("need 2 args")
	}

	sarr = ns.GetArray(keycls.Optdest())
	sarr = append(sarr, params[validx])
	sarr = append(sarr, params[validx+1])

	ns.SetValue(keycls.Optdest(), sarr)

	return 2, nil
}

func init() {
	dep_handler(nil, nil, nil)
	rdep_handler(nil, nil, nil)
	pair_key_handle(nil, 0, nil, []string{})
}

func main() {
	var commandline = `{
		"verbose" : "+",
		"pair|p!optparse=pair_key_handle!##to set pair parameters##" : [],
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
	var p *SubcmdStruct
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

	p = &SubcmdStruct{}
	_, err = parser.ParseCommandLineEx(nil, nil, p, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not parse err[%s]\n", err.Error())
		os.Exit(5)
		return
	}
	return
}
