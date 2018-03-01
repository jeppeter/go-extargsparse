package extargsparse

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

func beforeParser(t *testing.T) {
	var sarr []string
	var k string
	var envs []string
	var delone bool
	var err error

	delone = true
	for delone {
		delone = false
		envs = os.Environ()
		for _, k = range envs {
			k = strings.ToUpper(k)
			sarr = strings.Split(k, "=")
			if strings.HasPrefix(k, "EXTARGS_") ||
				strings.HasPrefix(k, "DEP_") ||
				strings.HasPrefix(k, "RDEP_") ||
				strings.HasPrefix(k, "HTTP_") ||
				strings.HasPrefix(k, "SSL_") ||
				strings.HasPrefix(k, "EXTARGSPARSE_JSON") {
				err = os.Unsetenv(sarr[0])
				if err == nil {
					delone = true
					break
				}
			}
		}
	}
	return
}

func makeWriteTempFileInner(s string) (fname string, err error) {
	var f *os.File
	f, err = ioutil.TempFile("", "tmpfile")
	if err != nil {

		return "", err
	}
	defer f.Close()
	_, err = f.WriteString(s)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func makeWriteTempFile(s string) string {
	var err error
	var fname string
	for {
		fname, err = makeWriteTempFileInner(s)
		if err == nil {
			return fname
		}
	}
	return ""
}

func assertGetOpt(opts []*ExtKeyParse, optdest string) *ExtKeyParse {
	for _, curopt := range opts {
		if !curopt.IsFlag() {
			continue
		}
		if curopt.FlagName() == "$" && optdest == "$" {
			return curopt
		}
		if curopt.FlagName() == "$" {
			continue
		}
		if curopt.Optdest() == optdest {
			return curopt
		}
	}
	return nil
}

func assertGetSubCommand(names []string, cmdname string) string {
	for _, c := range names {
		if c == cmdname {
			return cmdname
		}
	}
	return ""
}

func safeRemoveFile(fname string, notice string, ok bool) {
	if len(fname) > 0 {
		if ok {
			os.Remove(fname)
		} else {
			keyDebug("%s %s", notice, fname)
		}
	}
}

type stringIO struct {
	IoWriter
	obuf *bytes.Buffer
}

func newStringIO() *stringIO {
	p := &stringIO{}
	p.obuf = bytes.NewBufferString("")
	return p
}

func (self *stringIO) Write(data []byte) (int, error) {
	return self.obuf.Write(data)
}

func (self *stringIO) WriteString(s string) (int, error) {
	return self.obuf.WriteString(s)
}

func (self *stringIO) String() string {
	return self.obuf.String()
}

func getCmdHelp(parser *ExtArgsParse, cmdname string) []string {
	obuf := newStringIO()
	parser.PrintHelp(obuf, cmdname)
	return strings.Split(obuf.String(), "\n")
}

func getOptOk(t *testing.T, sarr []string, opt *ExtKeyParse) error {
	var exprstr string
	var ex *regexp.Regexp
	var err error
	if opt.FlagName() == "$" {
		return nil
	}
	exprstr = fmt.Sprintf("^\\s+%s", opt.Longopt())
	if len(opt.Shortopt()) > 0 {
		exprstr += fmt.Sprintf("\\|%s", opt.Shortopt())
	}
	if opt.Nargs().(int) != 0 {
		exprstr += fmt.Sprintf("\\s+%s\\s+.*$", opt.Optdest())
	} else {
		exprstr += fmt.Sprintf("\\s+.*$")
	}

	ex, err = regexp.Compile(exprstr)
	if err != nil {
		return fmt.Errorf("%s", format_error("compile [%s] for [%s] err[%s]", exprstr, opt.Format(), err.Error()))
	}
	for _, s := range sarr {
		if ex.MatchString(s) {
			return nil
		}
	}

	return fmt.Errorf("%s", format_error("can not find [%s] for \n%s", exprstr, sarr))
}

func checkAllOptsHelp(t *testing.T, sarr []string, opts []*ExtKeyParse) error {
	var err error
	for _, opt := range opts {
		err = getOptOk(t, sarr, opt)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
type parserTest1 struct {
	Verbose int
	Flag    bool
	Number  int
	List    []string
	String  string
	Args    []string
}

func Test_parser_A001(t *testing.T) {
	var loads = `        {
            "verbose|v##increment verbose mode##" : "+",
            "flag|f## flag set##" : false,
            "number|n" : 0,
            "list|l" : [],
            "string|s" : "string_var",
            "$" : {
                "value" : [],
                "nargs" : "*",
                "type" : "string"
            }
        }
	`
	var params = []string{"-vvvv", "-f", "-n", "30", "-l", "bar1", "-l", "bar2", "var1", "var2"}
	var parser *ExtArgsParse
	var args *NameSpaceEx
	var err error
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetBool("flag"), true)
	check_equal(t, args.GetInt("number"), 30)
	check_equal(t, args.GetArray("list"), []string{"bar1", "bar2"})
	check_equal(t, args.GetString("string"), "string_var")
	check_equal(t, args.GetArray("args"), []string{"var1", "var2"})
	return
}

func Test_parser_A001_2(t *testing.T) {
	var loads = `        {
            "verbose|v##increment verbose mode##" : "+",
            "flag|f## flag set##" : false,
            "number|n" : 0,
            "list|l" : [],
            "string|s" : "string_var",
            "$" : {
                "value" : [],
                "nargs" : "*",
                "type" : "string"
            }
        }
	`
	var params = []string{"-vvvv", "-f", "-n", "30", "-l", "bar1", "-l", "bar2", "var1", "var2"}
	var parser *ExtArgsParse
	var err error
	var p *parserTest1
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	p = &parserTest1{}
	_, err = parser.ParseCommandLineEx(params, nil, p, nil)
	check_equal(t, err, nil)
	check_equal(t, p.Verbose, 4)
	check_equal(t, p.Flag, true)
	check_equal(t, p.Number, 30)
	check_equal(t, p.List, []string{"bar1", "bar2"})
	check_equal(t, p.String, "string_var")
	check_equal(t, p.Args, []string{"var1", "var2"})
	return

}

type parserTest2 struct {
	Verbose int
	Port    int
	Dep     struct {
		List     []string
		String   string
		Subnargs []string
	}
}

func Test_parser_A002(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "port|p" : 3000,
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var parser *ExtArgsParse
	var err error
	var args *NameSpaceEx
	var params []string = []string{"-vvvv", "-p", "5000", "dep", "-l", "arg1", "--dep-list", "arg2", "cc", "dd"}
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetInt("port"), 5000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"arg1", "arg2"})
	check_equal(t, args.GetString("dep_string"), "s_var")
	check_equal(t, args.GetArray("subnargs"), []string{"cc", "dd"})
	return
}

func Test_parser_A002_2(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "port|p" : 3000,
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var parser *ExtArgsParse
	var err error
	var args *NameSpaceEx
	var params []string = []string{"-vvvv", "-p", "5000", "dep", "-l", "arg1", "--dep-list", "arg2", "cc", "dd"}
	var p *parserTest2
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	p = &parserTest2{}
	args, err = parser.ParseCommandLineEx(params, nil, p, nil)
	check_equal(t, err, nil)
	check_equal(t, p.Verbose, 4)
	check_equal(t, p.Port, 5000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, p.Dep.List, []string{"arg1", "arg2"})
	check_equal(t, p.Dep.String, "s_var")
	check_equal(t, p.Dep.Subnargs, []string{"cc", "dd"})
	return
}

func Test_parser_A003(t *testing.T) {
	var loads = `{
            "verbose|v" : "+",
            "port|p" : 3000,
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            },
            "rdep" : {
                "list|L" : [],
                "string|S" : "s_rdep",
                "$" : 2
            }
        }`
	var parser *ExtArgsParse
	var args *NameSpaceEx
	var err error
	var params = []string{"-vvvv", "-p", "5000", "rdep", "-L", "arg1", "--rdep-list", "arg2", "cc", "dd"}
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetInt("port"), 5000)
	check_equal(t, args.GetString("subcommand"), "rdep")
	check_equal(t, args.GetArray("rdep_list"), []string{"arg1", "arg2"})
	check_equal(t, args.GetString("rdep_string"), "s_rdep")
	check_equal(t, args.GetArray("subnargs"), []string{"cc", "dd"})
}

type parserTest3 struct {
	Verbose       int
	Port          int
	Dep_list      []string
	Dep_string    string
	Dep_Subnargs  []string
	Rdep_list     []string
	Rdep_string   string
	Rdep_Subnargs []string
	Args          []string
}

func Test_parser_A003_2(t *testing.T) {
	var loads = `{
            "verbose|v" : "+",
            "port|p" : 3000,
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            },
            "rdep" : {
                "list|L" : [],
                "string|S" : "s_rdep",
                "$" : 2
            }
        }`
	var parser *ExtArgsParse
	var args *NameSpaceEx
	var err error
	var params = []string{"-vvvv", "-p", "5000", "rdep", "-L", "arg1", "--rdep-list", "arg2", "cc", "dd"}
	var p *parserTest3
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	p = &parserTest3{}
	args, err = parser.ParseCommandLineEx(params, nil, p, nil)
	check_equal(t, err, nil)
	check_equal(t, p.Verbose, 4)
	check_equal(t, p.Port, 5000)
	check_equal(t, args.GetString("subcommand"), "rdep")
	check_equal(t, p.Rdep_list, []string{"arg1", "arg2"})
	check_equal(t, p.Rdep_string, "s_rdep")
	check_equal(t, p.Rdep_Subnargs, []string{"cc", "dd"})
	check_equal(t, p.Dep_Subnargs, []string{})
	check_equal(t, p.Dep_list, []string{})
	check_equal(t, p.Dep_string, "s_var")
	check_equal(t, p.Args, []string{})
}

func Test_parser_A004(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "port|p" : 3000,
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            },
            "rdep" : {
                "list|L" : [],
                "string|S" : "s_rdep",
                "$" : 2
            }
        }`
	var parser *ExtArgsParse
	var err error
	var params = []string{"-vvvv", "-p", "5000", "rdep", "-L", "arg1", "--rdep-list", "arg2", "cc", "dd"}
	var args *NameSpaceEx
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetInt("port"), 5000)
	check_equal(t, args.GetString("subcommand"), "rdep")
	check_equal(t, args.GetArray("rdep_list"), []string{"arg1", "arg2"})
	check_equal(t, args.GetString("rdep_string"), "s_rdep")
	check_equal(t, args.GetArray("subnargs"), []string{"cc", "dd"})
	return
}

type parserTest4 struct {
	Verbose int
	Port    int
	Dep     struct {
		List     []string
		String   string
		Subnargs []string
	}
	Rdep struct {
		List     []string
		String   string
		Subnargs []string
	}
	Args []string
}

func Test_parser_A004_2(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "port|p" : 3000,
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            },
            "rdep" : {
                "list|L" : [],
                "string|S" : "s_rdep",
                "$" : 2
            }
        }`
	var parser *ExtArgsParse
	var err error
	var params = []string{"-vvvv", "-p", "5000", "rdep", "-L", "arg1", "--rdep-list", "arg2", "cc", "dd"}
	var args *NameSpaceEx
	var p *parserTest4
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	p = &parserTest4{}
	args, err = parser.ParseCommandLineEx(params, nil, p, nil)
	check_equal(t, err, nil)
	check_equal(t, p.Verbose, 4)
	check_equal(t, p.Port, 5000)
	check_equal(t, args.GetString("subcommand"), "rdep")
	check_equal(t, p.Rdep.List, []string{"arg1", "arg2"})
	check_equal(t, p.Rdep.String, "s_rdep")
	check_equal(t, p.Rdep.Subnargs, []string{"cc", "dd"})
	check_equal(t, p.Dep.List, []string{})
	check_equal(t, p.Dep.String, "s_var")
	check_equal(t, p.Dep.Subnargs, []string{})
	check_equal(t, p.Args, []string{})
	return
}

type parserTest5Ctx struct {
	has_called_args string
}

func Debug_args_function(ns *NameSpaceEx, ostruct interface{}, Context interface{}) error {
	var p *parserTest5Ctx
	if Context == nil || ns == nil {
		return nil
	}
	p = Context.(*parserTest5Ctx)
	if ns.GetString("subcommand") != "" {
		p.has_called_args = ns.GetString("subcommand")
	} else {
		p.has_called_args = ""
	}
	return nil
}

func Test_parser_A005(t *testing.T) {
	var loads_fmt = `        {
            "verbose|v" : "+",
            "port|p" : 3000,
            "dep<%s.debug_args_function>" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            },
            "rdep" : {
                "list|L" : [],
                "string|S" : "s_rdep",
                "$" : 2
            }
        }`
	var parser *ExtArgsParse
	var err error
	var loads string
	var pkgname string
	var params = []string{"-p", "7003", "-vvvvv", "dep", "-l", "foo1", "-s", "new_var", "zz"}
	var args *NameSpaceEx
	var pc *parserTest5Ctx
	Debug_args_function(nil, nil, nil) // we call this function here because this function will compiled when call
	pc = &parserTest5Ctx{}
	pkgname = getCallerPackage(1)
	beforeParser(t)
	loads = fmt.Sprintf(loads_fmt, pkgname)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine(params, pc)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("port"), 7003)
	check_equal(t, args.GetInt("verbose"), 5)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"foo1"})
	check_equal(t, args.GetString("dep_string"), "new_var")
	check_equal(t, pc.has_called_args, "dep")
	check_equal(t, args.GetArray("subnargs"), []string{"zz"})
	return

}

func Test_parser_A006(t *testing.T) {
	var load1 = `        {
            "verbose|v" : "+",
            "port|p" : 3000,
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var load2 = `        {
            "rdep" : {
                "list|L" : [],
                "string|S" : "s_rdep",
                "$" : 2
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params = []string{"-p", "7003", "-vvvvv", "rdep", "-L", "foo1", "-S", "new_var", "zz", "64"}
	var args *NameSpaceEx
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", load1))
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(load2)
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("port"), 7003)
	check_equal(t, args.GetInt("verbose"), 5)
	check_equal(t, args.GetString("subcommand"), "rdep")
	check_equal(t, args.GetArray("rdep_list"), []string{"foo1"})
	check_equal(t, args.GetString("rdep_string"), "new_var")
	check_equal(t, args.GetArray("subnargs"), []string{"zz", "64"})
	return
}

func Test_parser_A007(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "port|p+http" : 3000,
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params = []string{"-vvvv", "dep", "-l", "cc", "--dep-string", "ee", "ww"}
	var args *NameSpaceEx
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetInt("http_port"), 3000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"cc"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	return
}

type parserTest7 struct {
	verbose    int
	http_port  int
	dep_list   []string
	dep_string string
	subnargs   []string
}

func Test_parser_A007_2(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "port|p+http" : 3000,
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params = []string{"-vvvv", "dep", "-l", "cc", "--dep-string", "ee", "ww"}
	var args *NameSpaceEx
	var options *ExtArgsOptions
	var p *parserTest7
	beforeParser(t)
	options, err = NewExtArgsOptions(`{"varuppercase" : false}`)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	p = &parserTest7{}
	args, err = parser.ParseCommandLineEx(params, nil, p, nil)
	check_equal(t, err, nil)
	check_equal(t, p.verbose, 4)
	check_equal(t, p.http_port, 3000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, p.dep_list, []string{"cc"})
	check_equal(t, p.dep_string, "ee")
	check_equal(t, p.subnargs, []string{"ww"})
	return
}

func Test_parser_A008(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "+http" : {
                "port|p" : 3000,
                "visual_mode|V" : false
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params = []string{"-vvvv", "--http-port", "9000", "--http-visual-mode", "dep", "-l", "cc", "--dep-string", "ee", "ww"}
	var args *NameSpaceEx
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetInt("http_port"), 9000)
	check_equal(t, args.GetBool("http_visual_mode"), true)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"cc"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	return
}

type parserTest8 struct {
	verbose          int
	http_port        int
	http_visual_mode bool
	dep_list         []string
	dep_string       string
	subnargs         []string
}

func Test_parser_A008_2(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "+http" : {
                "port|p" : 3000,
                "visual_mode|V" : false
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params = []string{"-vvvv", "--http-port", "9000", "--http-visual-mode", "dep", "-l", "cc", "--dep-string", "ee", "ww"}
	var args *NameSpaceEx
	var p *parserTest8
	var options *ExtArgsOptions
	beforeParser(t)
	options, err = NewExtArgsOptions(`{"varuppercase" : false}`)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	beforeParser(t)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	p = &parserTest8{}
	args, err = parser.ParseCommandLineEx(params, nil, p, nil)
	check_equal(t, err, nil)
	check_equal(t, p.verbose, 4)
	check_equal(t, p.http_port, 9000)
	check_equal(t, p.http_visual_mode, true)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, p.dep_list, []string{"cc"})
	check_equal(t, p.dep_string, "ee")
	check_equal(t, p.subnargs, []string{"ww"})
	return
}

func Test_parser_A009(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params = []string{"-vvvv", "-p", "9000", "dep", "-l", "cc", "--dep-string", "ee", "ww"}
	var args *NameSpaceEx
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"cc"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	return
}

type parserTest9 struct {
	verbose int
	port    int
	dep     struct {
		list     []string
		strv     string
		subnargs []string
	}
	args []string
}

func Test_parser_A009_2(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s<dep.strv>" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params = []string{"-vvvv", "-p", "9000", "dep", "-l", "cc", "--dep-string", "ee", "ww"}
	var args *NameSpaceEx
	var p *parserTest9
	var options *ExtArgsOptions
	beforeParser(t)
	options, err = NewExtArgsOptions(`{"varuppercase" : false}`)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	beforeParser(t)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	p = &parserTest9{}
	args, err = parser.ParseCommandLineEx(params, nil, p, nil)
	check_equal(t, err, nil)
	check_equal(t, p.verbose, 4)
	check_equal(t, p.port, 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, p.dep.list, []string{"cc"})
	check_equal(t, p.dep.strv, "ee")
	check_equal(t, p.dep.subnargs, []string{"ww"})
	return
}

type parserTest10 struct {
	Verbose int
	Port    int
	Dep     struct {
		List     []string
		String   string
		Subnargs []string
	}
	Args []string
}

func Test_parser_A010(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	var option *ExtArgsOptions
	var depjsonfile string = ""
	var ok bool = false
	var p *parserTest10
	beforeParser(t)

	depjsonfile = makeWriteTempFile(`{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"}`)
	defer safeRemoveFile(depjsonfile, "depjsonfile", ok)
	option, err = NewExtArgsOptions(`{"errorhandler" : "raise"}`)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(option, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-vvvv", "-p", "9000", "dep", "--dep-json", depjsonfile, "--dep-string", "ee", "ww"}
	p = &parserTest10{}
	args, err = parser.ParseCommandLineEx(params, nil, p, nil)
	check_equal(t, err, nil)
	check_equal(t, p.Verbose, 4)
	check_equal(t, p.Port, 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, p.Dep.List, []string{"jsonval1", "jsonval2"})
	check_equal(t, p.Dep.String, "ee")
	check_equal(t, p.Dep.Subnargs, []string{"ww"})
	ok = true
	return
}

type parserTest11 struct {
	Verbose int
	Port    int
	Dep     struct {
		List     []string
		String   string
		Subnargs []string
	}
	Args []string
}

func Test_parser_A011(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	var depjsonfile string = ""
	var ok bool = false
	beforeParser(t)

	depjsonfile = makeWriteTempFile(`{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"}`)
	defer safeRemoveFile(depjsonfile, "depjsonfile", ok)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-vvvv", "-p", "9000", "dep", "--dep-json", depjsonfile, "--dep-string", "ee", "ww"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"jsonval1", "jsonval2"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A012(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	var jsonfile string = ""
	var ok bool = false
	beforeParser(t)

	jsonfile = makeWriteTempFile(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer safeRemoveFile(jsonfile, "jsonfile", ok)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-p", "9000", "--json", jsonfile, "dep", "--dep-string", "ee", "ww"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 3)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"jsonval1", "jsonval2"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A013(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	var jsonfile string = ""
	var ok bool = false
	beforeParser(t)

	jsonfile = makeWriteTempFile(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer safeRemoveFile(jsonfile, "jsonfile", ok)
	os.Setenv("EXTARGSPARSE_JSON", jsonfile)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-p", "9000", "dep", "--dep-string", "ee", "ww"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 3)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"jsonval1", "jsonval2"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A014(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	var jsonfile string = ""
	var depjsonfile string = ""
	var ok bool = false
	beforeParser(t)

	jsonfile = makeWriteTempFile(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer safeRemoveFile(jsonfile, "jsonfile", ok)
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer safeRemoveFile(depjsonfile, "depjsonfile", ok)
	os.Setenv("EXTARGSPARSE_JSON", jsonfile)
	os.Setenv("DEP_JSON", depjsonfile)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-p", "9000", "dep", "--dep-string", "ee", "ww"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 3)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"depjson1", "depjson2"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A015(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	var jsonfile string = ""
	var depjsonfile string = ""
	var ok bool = false
	beforeParser(t)

	jsonfile = makeWriteTempFile(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer safeRemoveFile(jsonfile, "jsonfile", ok)
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer safeRemoveFile(depjsonfile, "depjsonfile", ok)
	os.Setenv("DEP_JSON", depjsonfile)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-p", "9000", "--json", jsonfile, "dep", "--dep-string", "ee", "ww"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 3)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"jsonval1", "jsonval2"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A016(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	var jsonfile string = ""
	var depjsonfile string = ""
	var ok bool = false
	var depstrval string
	var depliststr string
	beforeParser(t)

	depstrval = "newval"
	depliststr = `["depenv1","depenv2"]`
	jsonfile = makeWriteTempFile(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer safeRemoveFile(jsonfile, "jsonfile", ok)
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer safeRemoveFile(depjsonfile, "depjsonfile", ok)
	os.Setenv("EXTARGSPARSE_JSON", jsonfile)
	os.Setenv("DEP_JSON", depjsonfile)

	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	os.Setenv("DEP_STRING", depstrval)
	os.Setenv("DEP_LIST", depliststr)
	params = []string{"-p", "9000", "dep", "--dep-string", "ee", "ww"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 3)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"depenv1", "depenv2"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A017(t *testing.T) {
	var loads = `        {
            "+dpkg" : {
                "dpkg" : "dpkg"
            },
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	beforeParser(t)

	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 0)
	check_equal(t, args.GetInt("port"), 3000)
	check_equal(t, args.GetString("dpkg_dpkg"), "dpkg")
	check_equal(t, args.GetArray("args"), []string{})
	return
}

func Test_parser_A018(t *testing.T) {
	var loads = `        {
            "+dpkg" : {
                "dpkg" : "dpkg"
            },
            "verbose|v" : "+",
            "rollback|r": true,
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	beforeParser(t)

	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-vvrvv"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetBool("rollback"), false)
	check_equal(t, args.GetInt("port"), 3000)
	check_equal(t, args.GetString("dpkg_dpkg"), "dpkg")
	check_equal(t, args.GetArray("args"), []string{})
	return
}

func Test_parser_A019(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	var jsonfile string = ""
	var depjsonfile string = ""
	var ok bool = false
	var depstrval string
	var depliststr string
	beforeParser(t)

	depstrval = "newval"
	depliststr = `["depenv1","depenv2"]`
	jsonfile = makeWriteTempFile(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer safeRemoveFile(jsonfile, "jsonfile", ok)
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer safeRemoveFile(depjsonfile, "depjsonfile", ok)
	os.Setenv("EXTARGSPARSE_JSON", jsonfile)
	os.Setenv("DEP_JSON", depjsonfile)

	parser, err = NewExtArgsParse(nil, []int{ENV_COMMAND_JSON_SET, ENVIRONMENT_SET, ENV_SUB_COMMAND_JSON_SET})
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	os.Setenv("DEP_STRING", depstrval)
	os.Setenv("DEP_LIST", depliststr)
	params = []string{"-p", "9000", "dep", "--dep-string", "ee", "ww"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 3)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"jsonval1", "jsonval2"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A020(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "rollback|R" : true,
            "$port|P" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	beforeParser(t)

	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-P", "9000", "--no-rollback", "dep", "--dep-string", "ee", "ww"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 0)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetBool("rollback"), false)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("args"), []string{})
	return
}

func Test_parser_A021(t *testing.T) {
	var loads = `        {
            "maxval|m" : 392244922
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	beforeParser(t)

	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-m", "0xffcc"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("maxval"), 0xffcc)
	return
}

func Test_parser_A022(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+"
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var opts []*ExtKeyParse
	var curopt *ExtKeyParse
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params, err = parser.GetSubCommands("")
	check_equal(t, err, nil)
	check_equal(t, params, []string{})
	opts, err = parser.GetCmdOpts("")
	check_equal(t, err, nil)
	check_equal(t, len(opts), 4)
	curopt = assertGetOpt(opts, "verbose")
	check_not_equal(t, curopt, nil)
	check_equal(t, curopt.Optdest(), "verbose")
	check_equal(t, curopt.Longopt(), "--verbose")
	check_equal(t, curopt.Shortopt(), "-v")
	curopt = assertGetOpt(opts, "noflag")
	check_equal(t, curopt, (*ExtKeyParse)(nil))
	curopt = assertGetOpt(opts, "json")
	check_not_equal(t, curopt, nil)
	check_equal(t, curopt.Value(), nil)
	curopt = assertGetOpt(opts, "help")
	check_not_equal(t, curopt, nil)
	check_equal(t, curopt.Longopt(), "--help")
	check_equal(t, curopt.Shortopt(), "-h")
	check_equal(t, curopt.TypeName(), "help")
	return
}

func Test_parser_A023(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "dep" : {
                "new|n" : false,
                "$<NARGS>" : "+"
            },
            "rdep" : {
                "new|n" : true,
                "$<NARGS>" : "?"
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var opts []*ExtKeyParse
	var curopt *ExtKeyParse
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params, err = parser.GetSubCommands("")
	check_equal(t, err, nil)
	check_equal(t, params, []string{"dep", "rdep"})
	opts, err = parser.GetCmdOpts("")
	check_equal(t, err, nil)
	check_equal(t, len(opts), 4)
	curopt = assertGetOpt(opts, "$")
	check_not_equal(t, curopt, (*ExtKeyParse)(nil))
	check_equal(t, curopt.Nargs().(string), "*")
	curopt = assertGetOpt(opts, "verbose")
	check_not_equal(t, curopt, (*ExtKeyParse)(nil))
	check_equal(t, curopt.TypeName(), "count")
	curopt = assertGetOpt(opts, "json")
	check_not_equal(t, curopt, (*ExtKeyParse)(nil))
	check_equal(t, curopt.TypeName(), "jsonfile")
	curopt = assertGetOpt(opts, "help")
	check_not_equal(t, curopt, (*ExtKeyParse)(nil))
	check_equal(t, curopt.TypeName(), "help")
	opts, err = parser.GetCmdOpts("dep")
	check_equal(t, err, nil)
	check_equal(t, len(opts), 4)
	curopt = assertGetOpt(opts, "$")
	check_not_equal(t, curopt, (*ExtKeyParse)(nil))
	check_equal(t, curopt.VarName(), "NARGS")
	curopt = assertGetOpt(opts, "help")
	check_not_equal(t, curopt, (*ExtKeyParse)(nil))
	check_equal(t, curopt.TypeName(), "help")
	curopt = assertGetOpt(opts, "dep_json")
	check_not_equal(t, curopt, (*ExtKeyParse)(nil))
	check_equal(t, curopt.TypeName(), "jsonfile")
	curopt = assertGetOpt(opts, "dep_new")
	check_not_equal(t, curopt, (*ExtKeyParse)(nil))
	check_equal(t, curopt.TypeName(), "bool")
	return
}

func Test_parser_A024(t *testing.T) {
	var loads = `        {
            "rdep" : {
                "ip" : {
                    "modules" : [],
                    "called" : true,
                    "setname" : null,
                    "$" : 2
                }
            },
            "dep" : {
                "port" : 5000,
                "cc|C" : true
            },
            "verbose|v" : "+"
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	beforeParser(t)

	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"rdep", "ip", "--verbose", "--rdep-ip-modules", "cc", "--rdep-ip-setname", "bb", "xx", "bb"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetString("subcommand"), "rdep.ip")
	check_equal(t, args.GetInt("verbose"), 1)
	check_equal(t, args.GetArray("rdep_ip_modules"), []string{"cc"})
	check_equal(t, args.GetString("rdep_ip_setname"), "bb")
	check_equal(t, args.GetArray("subnargs"), []string{"xx", "bb"})
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"dep", "--verbose", "--verbose", "-vvC"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetInt("dep_port"), 5000)
	check_equal(t, args.GetBool("dep_cc"), false)
	check_equal(t, args.GetArray("subnargs"), []string{})
	return
}

func Test_parser_A025(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "+http" : {
                "url|u" : "http://www.google.com",
                "visual_mode|V": false
            },
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+",
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            },
            "rdep" : {
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            }
        }`
	var err error
	var parser *ExtArgsParse
	var params []string
	var args *NameSpaceEx
	var jsonfile string = ""
	var depjsonfile string = ""
	var rdepjsonfile string = ""
	var ok bool = false
	beforeParser(t)

	jsonfile = makeWriteTempFile(`{ "http" : { "url" : "http://www.github.com"} ,"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer safeRemoveFile(jsonfile, "jsonfile", ok)
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer safeRemoveFile(depjsonfile, "depjsonfile", ok)
	rdepjsonfile = makeWriteTempFile(`{"ip": {"list":["rdepjson1","rdepjson3"],"verbose": 5}}`)
	defer safeRemoveFile(rdepjsonfile, "rdepjsonfile", ok)

	os.Setenv("EXTARGSPARSE_JSON", jsonfile)
	os.Setenv("DEP_JSON", depjsonfile)
	os.Setenv("RDEP_JSON", rdepjsonfile)

	parser, err = NewExtArgsParse(nil, []int{ENV_COMMAND_JSON_SET, ENVIRONMENT_SET, ENV_SUB_COMMAND_JSON_SET})
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-p", "9000", "rdep", "ip", "--rdep-ip-verbose", "--rdep-ip-cc", "ee", "ww"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 3)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("dep_string"), "jsonstring")
	check_equal(t, args.GetArray("dep_list"), []string{"jsonval1", "jsonval2"})
	check_equal(t, args.GetBool("http_visual_mode"), false)
	check_equal(t, args.GetString("http_url"), "http://www.github.com")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	check_equal(t, args.GetString("subcommand"), "rdep.ip")
	check_equal(t, args.GetInt("rdep_ip_verbose"), 1)
	check_equal(t, args.GetArray("rdep_ip_cc"), []string{"ee"})
	check_equal(t, args.GetArray("rdep_ip_list"), []string{"rdepjson1", "rdepjson3"})
	ok = true
	return
}

func Test_parser_A026(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "+http" : {
                "url|u" : "http://www.google.com",
                "visual_mode|V": false
            },
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+",
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            },
            "rdep" : {
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            }
        }`
	var err error
	var parser *ExtArgsParse
	var options *ExtArgsOptions
	var sarr []string
	var opts []*ExtKeyParse
	beforeParser(t)
	options, err = NewExtArgsOptions(`{"prog" : "cmd1"}`)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "")
	opts, err = parser.GetCmdOpts("")
	check_equal(t, err, nil)
	err = checkAllOptsHelp(t, sarr, opts)
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "rdep")
	opts, err = parser.GetCmdOpts("rdep")
	check_equal(t, err, nil)
	err = checkAllOptsHelp(t, sarr, opts)
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "rdep.ip")
	opts, err = parser.GetCmdOpts("rdep.ip")
	check_equal(t, err, nil)
	err = checkAllOptsHelp(t, sarr, opts)
	check_equal(t, err, nil)
	return
}

func Test_parser_A027(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "+http" : {
                "url|u" : "http://www.google.com",
                "visual_mode|V": false
            },
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l!attr=cc;optfunc=list_opt_func!" : [],
                "string|s" : "s_var",
                "$" : "+",
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            },
            "rdep" : {
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            }
        }`
	var err error
	var parser *ExtArgsParse
	var options *ExtArgsOptions
	var opts []*ExtKeyParse
	var flag *ExtKeyParse
	beforeParser(t)
	options, err = NewExtArgsOptions(`{"prog" : "cmd1"}`)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	opts, err = parser.GetCmdOpts("dep")
	check_equal(t, err, nil)
	flag = nil
	for _, f := range opts {
		if f.TypeName() == "args" {
			continue
		}
		if f.FlagName() == "list" {
			flag = f
			break
		}
	}
	check_not_equal(t, flag, (*ExtKeyParse)(nil))
	check_equal(t, flag.Attr("attr"), "cc")
	check_equal(t, flag.Attr("optfunc"), "list_opt_func")
	return
}
*/

func Test_parser_A028(t *testing.T) {
	var loads = `        {
            "verbose<VAR1>|v" : "+",
            "+http" : {
                "url|u<VAR1>" : "http://www.google.com",
                "visual_mode|V": false
            },
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 , 
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l!attr=cc;optfunc=list_opt_func!" : [],
                "string|s" : "s_var",
                "$" : "+",
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            },
            "rdep" : {
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            }
        }`
	var err error
	var parser *ExtArgsParse
	var options *ExtArgsOptions
	var params []string
	beforeParser(t)
	options, err = NewExtArgsOptions(`{"errorhandler" : "raise"}`)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"dep", "cc"}
	_, err = parser.ParseCommandLine(params, nil)
	check_not_equal(t, err, nil)
	return
}
