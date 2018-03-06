package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

type parseArgs struct {
	m_verbose int
	m_removed bool
	m_floatv  float64
	intv      int
	arrl      []string
	strv      []string
	args      []string
}

type parseArgs2 struct {
	M_verbose int
	M_removed bool
	M_floatv  float64
	Intv      int
	Arrl      []string
	Strv      []string
	Args      []string
}

func main() {
	var commandline = `{
		"verbose|v<m_verbose>" : "+",
		"removed|R<m_removed>" : false,
		"floatv|f<m_floatv>" : 3.3,
		"intv|i" : 5,
		"arrl|a" : [],
		"strv|s" : null,
		"$" : "+"
		}`
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var confstr string
	var p1 *parseArgs
	var p2 *parseArgs2
	var opts []*extargsparse.ExtKeyParse
	var flag *extargsparse.ExtKeyParse

	confstr = fmt.Sprintf(`{"%s": false}`, extargsparse.OPT_FUNC_UPPER_CASE)
	options, _ = extargsparse.NewExtArgsOptions(confstr)
	parser, _ = extargsparse.NewExtArgsParse(options, nil)
	parser.LoadCommandLineString(commandline)
	p1 = &parseArgs{}
	parser.ParseCommandLineEx([]string{"-vvv", "-R", "-a", "cc", "-s", "csw", "ww", "ee"}, nil, p1, nil)
	fmt.Fprintf(os.Stdout, "verbose=%d\n", p1.m_verbose)
	fmt.Fprintf(os.Stdout, "removed=%v\n", p1.m_removed)
	fmt.Fprintf(os.Stdout, "floatv=%f\n", p1.m_floatv)
	fmt.Fprintf(os.Stdout, "intv=%d\n", p1.intv)
	fmt.Fprintf(os.Stdout, "arrl=%v\n", p1.arrl)
	fmt.Fprintf(os.Stdout, "strv=%s\n", p1.strv)
	fmt.Fprintf(os.Stdout, "args=%v\n", p1.args)
	opts, _ = parser.GetCmdOpts("")
	for _, flag = range opts {
		if flag.TypeName() == "args" {
			fmt.Fprintf(os.Stdout, "args.varname=%s\n", flag.VarName())
		} else {
			fmt.Fprintf(os.Stdout, "%s.varname=%s\n", flag.FlagName(), flag.VarName())
		}
	}

	confstr = fmt.Sprintf(`{}`)
	options, _ = extargsparse.NewExtArgsOptions(confstr)
	parser, _ = extargsparse.NewExtArgsParse(options, nil)
	parser.LoadCommandLineString(commandline)
	p2 = &parseArgs2{}
	parser.ParseCommandLineEx([]string{"-vvv", "-R", "-a", "cc", "-s", "csw", "ww", "ee"}, nil, p2, nil)
	fmt.Fprintf(os.Stdout, "verbose=%d\n", p2.M_verbose)
	fmt.Fprintf(os.Stdout, "removed=%v\n", p2.M_removed)
	fmt.Fprintf(os.Stdout, "floatv=%f\n", p2.M_floatv)
	fmt.Fprintf(os.Stdout, "intv=%d\n", p2.Intv)
	fmt.Fprintf(os.Stdout, "arrl=%v\n", p2.Arrl)
	fmt.Fprintf(os.Stdout, "strv=%s\n", p2.Strv)
	fmt.Fprintf(os.Stdout, "args=%v\n", p2.Args)
	opts, _ = parser.GetCmdOpts("")
	for _, flag = range opts {
		if flag.TypeName() == "args" {
			fmt.Fprintf(os.Stdout, "args.varname=%s\n", flag.VarName())
		} else {
			fmt.Fprintf(os.Stdout, "%s.varname=%s\n", flag.FlagName(), flag.VarName())
		}
	}
	/*
		Output:
			verbose=0
			removed=false
			floatv=0.000000
			intv=0
			arrl=[]
			strv=[]
			args=[]
			args.varname=args
			arrl.varname=arrl
			floatv.varname=m_floatv
			help.varname=help
			intv.varname=intv
			json.varname=json
			removed.varname=m_removed
			strv.varname=strv
			verbose.varname=m_verbose
			verbose=3
			removed=true
			floatv=3.300000
			intv=5
			arrl=[cc]
			strv=[]
			args=[ww ee]
			args.varname=args
			arrl.varname=arrl
			floatv.varname=m_floatv
			help.varname=help
			intv.varname=intv
			json.varname=json
			removed.varname=m_removed
			strv.varname=strv
			verbose.varname=m_verbose
	*/
	return
}
