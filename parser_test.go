package extargsparse

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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
				strings.HasPrefix(k, "EXTARGSPARSE_JSON") ||
				strings.HasPrefix(k, "EXTARGSPARSE_JSONFILE") {
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

func assertGetCmd(sarr []string, cmdname string) bool {
	var matchexpr *regexp.Regexp
	var restr string = fmt.Sprintf(`^\s+\[%s\]\s+.*`, cmdname)
	var ok bool = false
	var s string
	matchexpr = regexp.MustCompile(restr)
	for _, s = range sarr {
		if matchexpr.MatchString(s) {
			ok = true
			break
		}
	}

	return ok
}

func assertOkCmds(t *testing.T, sarr []string, parser *ExtArgsParse, cmdname string) error {
	var cmds []string
	var ok bool
	var c string
	var err error
	cmds, err = parser.GetSubCommands(cmdname)
	check_equal(t, err, nil)
	for _, c = range cmds {
		ok = assertGetCmd(sarr, c)
		if !ok {
			return fmt.Errorf("%s", format_error("cannot found [%s] in {%v}", c, sarr))
		}
	}
	return nil
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
func joinFile(a ...string) string {
	var retf string
	var err error
	retf = filepath.Join(a...)
	retf, err = filepath.Abs(retf)
	if err != nil {
		keyDebug("get [%s] err[%s]", retf, err.Error())
		return ""
	}
	return retf
}

func debug_set_2_args(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	var sarr []string
	if ns == nil {
		return 0, nil
	}
	if (validx + 2) > len(params) {
		return 0, fmt.Errorf("%s", format_error("[%d+2] > len(%d) %v", validx, len(params), params))
	}
	sarr = ns.GetArray(keycls.Optdest())
	sarr = append(sarr, params[validx])
	sarr = append(sarr, params[(validx+1)])
	ns.SetValue(keycls.Optdest(), sarr)
	return 2, nil
}

func Debug_set_2_args(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	var sarr []string
	if ns == nil {
		return 0, nil
	}
	if (validx + 2) > len(params) {
		return 0, fmt.Errorf("%s", format_error("[%d+2] > len(%d) %v", validx, len(params), params))
	}
	sarr = ns.GetArray(keycls.Optdest())
	sarr = append(sarr, strings.ToUpper(params[validx]))
	sarr = append(sarr, strings.ToUpper(params[(validx+1)]))
	ns.SetValue(keycls.Optdest(), sarr)
	return 2, nil
}

func debug_opthelp_set(keycls *ExtKeyParse) string {
	if keycls == nil {
		return ""
	}
	return fmt.Sprintf("opthelp function set [%s] default value (%s)", keycls.Optdest(), keycls.Value())
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

func jsonfunc_2(value interface{}) (vs []string, err error) {
	var vinter []interface{}
	var v interface{}

	vs = make([]string, 0)
	switch value.(type) {
	case []interface{}:
		vinter = value.([]interface{})
		for _, v = range vinter {
			vs = append(vs, fmt.Sprintf("%s", v))
		}
	case []string:
		vs = value.([]string)
	default:
		err = fmt.Errorf("%s", format_error("can not use type [%s]", reflect.ValueOf(value).Type().String()))
		return
	}
	err = nil
	return

}

func Debug_2_jsonfunc(ns *NameSpaceEx, keycls *ExtKeyParse, value interface{}) error {
	var vs []string
	var err error
	var setvs []string
	var idx int
	if ns == nil {
		return nil
	}

	if keycls == nil {
		return fmt.Errorf("%s", format_error("nil keycls"))
	}

	if !keycls.IsFlag() || keycls.TypeName() != "list" {
		return fmt.Errorf("%s", format_error("keycls [%s] not valid", keycls.Format()))
	}

	if value == nil {
		ns.SetValue(keycls.Optdest(), []string{})
	}

	vs, err = jsonfunc_2(value)
	if err != nil {
		return err
	}

	if (len(vs) % 2) != 0 {
		return fmt.Errorf("%s", format_error("%s not even size", vs))
	}

	setvs = make([]string, 0)
	for idx = 0; idx < len(vs); idx += 2 {
		setvs = append(setvs, vs[idx])
	}
	ns.SetValue(keycls.Optdest(), setvs)
	return nil
}

func Debug_upper_jsonfunc(ns *NameSpaceEx, keycls *ExtKeyParse, value interface{}) error {
	var s string
	if ns == nil {
		return nil
	}

	if keycls == nil {
		return fmt.Errorf("%s", format_error("keycls nil"))
	}

	if !keycls.IsFlag() || keycls.TypeName() != "string" {
		return fmt.Errorf("%s", format_error("[%s] not type string", keycls.Format()))
	}
	if value == nil {
		ns.SetValue(keycls.Optdest(), nil)
		return nil
	}

	switch value.(type) {
	case string:
		s = value.(string)
	default:
		return fmt.Errorf("%s", format_error("[%v] not valid string type", value))
	}
	ns.SetValue(keycls.Optdest(), strings.ToUpper(s))
	return nil
}

func init() {
	Debug_2_jsonfunc(nil, nil, nil)
	Debug_upper_jsonfunc(nil, nil, nil)
	debug_opthelp_set(nil)
	Debug_set_2_args(nil, 0, nil, []string{})
	debug_set_2_args(nil, 0, nil, []string{})
	Debug_args_function(nil, nil, nil)
}

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
        }`
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
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"%s" : false}`, OPT_VAR_UPPER_CASE))
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
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"%s" : false}`, OPT_VAR_UPPER_CASE))
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
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"%s" : false}`, OPT_VAR_UPPER_CASE))
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
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
	option, err = NewExtArgsOptions(fmt.Sprintf(`{"%s" : "raise"}`, OPT_ERROR_HANDLER))
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
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
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
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
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
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
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
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
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
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
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
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
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
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
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
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
	rdepjsonfile = makeWriteTempFile(`{"ip": {"list":["rdepjson1","rdepjson3"],"verbose": 5}}`)
	defer func() { safeRemoveFile(rdepjsonfile, "rdepjsonfile", ok) }()

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
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"%s" : "cmd1"}`, OPT_PROG))
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
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"%s" : "cmd1"}`, OPT_PROG))
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
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"%s" : "raise"}`, OPT_ERROR_HANDLER))
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

func Test_parser_A029(t *testing.T) {
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
	var sarr []string
	beforeParser(t)
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"%s" : "nohelp"}`, OPT_HELP_HANDLER))
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "")
	check_equal(t, sarr, []string{"no help information"})
	return
}

func Test_parser_A030(t *testing.T) {
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
            "dep<dep_handler>!opt=cc!" : {
                "list|l!attr=cc;optfunc=list_opt_func!" : [],
                "string|s" : "s_var",
                "$" : "+",
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            },
            "rdep<rdep_handler>" : {
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            }
        }`
	var err error
	var parser *ExtArgsParse
	//var opts []*ExtKeyParse
	var flag *ExtKeyParse
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	flag, err = parser.GetCmdKey("")
	check_equal(t, err, nil)
	check_equal(t, flag.CmdName(), "main")
	check_equal(t, flag.IsCmd(), true)
	check_equal(t, flag.Function(), "")
	flag, err = parser.GetCmdKey("dep")
	check_equal(t, err, nil)
	check_equal(t, flag.CmdName(), "dep")
	check_equal(t, flag.Function(), "dep_handler")
	check_equal(t, flag.Attr("opt"), "cc")
	flag, err = parser.GetCmdKey("rdep")
	check_equal(t, err, nil)
	check_equal(t, flag.CmdName(), "rdep")
	check_equal(t, flag.Function(), "rdep_handler")
	check_equal(t, flag.Attr(""), "")
	flag, err = parser.GetCmdKey("nosuch")
	check_equal(t, err, nil)
	check_equal(t, flag, (*ExtKeyParse)(nil))
	return
}

func Test_parser_A031(t *testing.T) {
	var loads = `        {
            "verbose|v" : "+",
            "catch|C## to not catch the exception ##" : true,
            "input|i## to specify input default(stdin)##" : null,
            "$caption## set caption ##" : "runcommand",
            "test|t##to test mode##" : false,
            "release|R##to release test mode##" : false,
            "$" : "*"
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
	params = []string{"--test"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetBool("test"), true)
	check_equal(t, args.GetArray("args"), []string{})
	return
}

func Test_parser_A032(t *testing.T) {
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
            "dep<dep_handler>!opt=cc!" : {
                "list|l!attr=cc;optfunc=list_opt_func!" : [],
                "string|s" : "s_var",
                "$" : "+",
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            },
            "rdep<rdep_handler>" : {
                "ip" : {
                    "verbose" : "+",
                    "list" : [],
                    "cc" : []
                }
            }
        }`
	var err error
	var ok bool = false
	var cl *compileExec
	var setvars map[string]string
	var parser *ExtArgsParse
	var opts []*ExtKeyParse
	beforeParser(t)
	cl = newComileExec()
	check_not_equal(t, cl, (*compileExec)(nil))
	defer func() {
		if cl != nil {
			cl.Release(ok)
		}
		cl = nil
	}()
	err = cl.WriteScript("{}", loads, []string{}, "", nil, false, "ns", "pp")
	check_equal(t, err, nil)
	setvars = make(map[string]string)
	err = cl.Compile()
	check_equal(t, err, nil)
	err = cl.RunCmd(setvars, []string{}, "-h")
	check_equal(t, err, nil)

	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)

	opts, err = parser.GetCmdOpts("")
	check_equal(t, err, nil)

	err = checkAllOptsHelp(t, cl.GetOut(), opts)
	check_equal(t, err, nil)

	err = cl.RunCmd(setvars, []string{}, "dep", "-h")
	check_equal(t, err, nil)

	opts, err = parser.GetCmdOpts("dep")
	check_equal(t, err, nil)

	err = checkAllOptsHelp(t, cl.GetOut(), opts)
	check_equal(t, err, nil)

	err = cl.RunCmd(setvars, []string{}, "rdep", "-h")
	check_equal(t, err, nil)

	opts, err = parser.GetCmdOpts("rdep")
	check_equal(t, err, nil)

	err = checkAllOptsHelp(t, cl.GetOut(), opts)
	check_equal(t, err, nil)

	ok = true
	return
}

func Test_parser_A033(t *testing.T) {
	var cmd1_fmt = `        {
            "%s" : true
        }`
	var cmd2_fmt = `        {
            "+%s" : {
                "reserve": true
            }
        }`
	var cmd3_fmt = `        {
            "%s" : {
                "function" : 30
            }
        }`
	var test_reserve_args = []string{"subcommand", "subnargs", "nargs", "extargs", "args"}
	var cmd_fmts = []string{cmd1_fmt, cmd2_fmt, cmd3_fmt}
	var fmtstr string
	var k string
	var err error
	var parser *ExtArgsParse
	var loads string
	beforeParser(t)
	for _, fmtstr = range cmd_fmts {
		for _, k = range test_reserve_args {
			loads = fmt.Sprintf(fmtstr, k)
			parser, err = NewExtArgsParse(nil, nil)
			check_equal(t, err, nil)
			err = parser.LoadCommandLineString(loads)
			check_not_equal(t, err, nil)
		}
	}
	return
}

func Test_parser_A034(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads = `        {
            "dep" : {
                "string|S" : "stringval"
            }
        }`
	var depjson string = ""
	var ok = false
	var params []string
	var args *NameSpaceEx
	beforeParser(t)
	depjson = makeWriteTempFile(`{"dep_string":null}`)
	defer func() { safeRemoveFile(depjson, "depjson", ok) }()
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"--json", depjson, "dep"}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetString("dep_string"), "")
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("subnargs"), []string{})
	ok = true
	return
}

func Test_parser_A035(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads = `        {
            "float1|f" : 3.633 ,
            "float2" : 6422.22,
            "float3" : 44463.23,
            "verbose|v" : "+",
            "dep" : {
                "float3" : 3332.233
            },
            "rdep" : {
                "ip" : {
                    "float4" : 3377.33,
                    "float6" : 33.22,
                    "float7" : 0.333
                }
            }

        }`
	var depjson string = ""
	var rdepjson string = ""
	var rdepipjson string = ""
	var jsonfile string = ""
	var ok = false
	var params []string
	var args *NameSpaceEx
	beforeParser(t)
	depjson = makeWriteTempFile(`{"float3":33.221}`)
	defer func() { safeRemoveFile(depjson, "depjson", ok) }()
	rdepjson = makeWriteTempFile(`{"ip" : { "float4" : 40.3}}`)
	defer func() { safeRemoveFile(rdepjson, "rdepjson", ok) }()
	jsonfile = makeWriteTempFile(`{"verbose": 30,"float3": 77.1}`)
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	rdepipjson = makeWriteTempFile(`{"float7" : 11.22,"float4" : 779.2}`)
	defer func() { safeRemoveFile(rdepipjson, "rdepipjson", ok) }()
	os.Setenv("EXTARGSPARSE_JSON", jsonfile)
	os.Setenv("DEP_JSON", depjson)
	os.Setenv("RDEP_JSON", rdepjson)
	os.Setenv("DEP_FLOAT3", fmt.Sprintf("33.52"))
	os.Setenv("RDEP_IP_FLOAT7", fmt.Sprintf("99.3"))
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	params = []string{"-vvfvv", "33.21", "rdep", "ip", "--json", jsonfile, "--rdep-ip-json", rdepipjson}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetArray("subnargs"), []string{})
	check_equal(t, args.GetString("subcommand"), "rdep.ip")
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetFloat("float1"), 33.21)
	check_equal(t, args.GetFloat("dep_float3"), 33.52)
	check_equal(t, args.GetFloat("float2"), 6422.22)
	check_equal(t, args.GetFloat("float3"), 77.1)
	check_equal(t, args.GetFloat("rdep_ip_float4"), 779.2)
	check_equal(t, args.GetFloat("rdep_ip_float6"), 33.22)
	check_equal(t, args.GetFloat("rdep_ip_float7"), 11.22)
	ok = true
	return
}

func Test_parser_A037(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads = `        {
            "jsoninput|j##input json default stdin##" : null,
            "input|i##input file to get default nothing - for stdin##" : null,
            "output|o##output c file##" : null,
            "verbose|v##verbose mode default(0)##" : "+",
            "cmdpattern|c" : "%EXTARGS_CMDSTRUCT%",
            "optpattern|O" : "%EXTARGS_STRUCT%",
            "structname|s" : "args_options_t",
            "funcname|F" : "debug_extargs_output",
            "releasename|R" : "release_extargs_output",
            "funcpattern" : "%EXTARGS_DEBUGFUNC%",
            "prefix|p" : "",
            "test" : {
                "$" : 0
            },
            "optstruct" : {
                "$" : 0
            },
            "cmdstruct" : {
                "$" : 0
            },
            "debugfunc" : {
                "$" : 0
            },
            "all" : {
                "$" : 0
            }
        }`
	var subcmds []string
	var opts []*ExtKeyParse
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	subcmds, err = parser.GetSubCommands("")
	check_equal(t, err, nil)
	check_equal(t, len(subcmds), 5)
	check_equal(t, subcmds[0], "all")
	check_equal(t, subcmds[1], "cmdstruct")
	check_equal(t, subcmds[2], "debugfunc")
	check_equal(t, subcmds[3], "optstruct")
	check_equal(t, subcmds[4], "test")
	opts, err = parser.GetCmdOpts("")
	check_equal(t, err, nil)
	check_equal(t, len(opts), 14)
	check_equal(t, opts[0].FlagName(), "$")
	check_equal(t, opts[1].Longopt(), "--cmdpattern")
	check_equal(t, opts[2].Optdest(), "funcname")
	check_equal(t, opts[3].VarName(), "funcpattern")
	check_equal(t, opts[4].TypeName(), "help")
	return
}

func Test_parser_A038(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads = `        {
            "verbose|v" : "+",
            "kernel|K" : "/boot/",
            "initrd|I" : "/boot/",
            "encryptfile|e" : null,
            "encryptkey|E" : null,
            "setupsectsoffset" : 0x1f1,
            "ipxe<ipxe_handler>" : {
                "$" : "+"
            }
        }`
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_not_equal(t, err, nil)
	return
}

func Test_parser_A039(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads = `        {
            "verbose|v" : "+",
            "kernel|K" : "/boot/",
            "initrd|I" : "/boot/",
            "encryptfile|e" : null,
            "encryptkey|E" : null,
            "setupsectsoffset" : 451
        }`
	var args *NameSpaceEx
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	os.Setenv("EXTARGS_VERBOSE", "4")
	os.Setenv("EXTARGS_SETUPSECTSOFFSET", "0x612")
	args, err = parser.ParseCommandLine([]string{}, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetInt("setupsectsoffset"), 0x612)
	return
}

type parserTest40 struct {
	Tce_mirror       string
	Tce_root         string
	Tce_tceversion   string
	Tce_wget         string
	Tce_cat          string
	Tce_rm           string
	Tce_sudoprefix   string
	Tce_optional_dir string
	Tce_trymode      bool
	Tce_platform     string
	Tce_mount        string
	Tce_umount       string
	Tce_chroot       string
	Tce_chown        string
	Tce_mkdir        string
	Tce_rollback     bool
	Tce_cp           string
	Tce_jsonfile     string
	Tce_perspace     int
	Tce_depmapfile   string
	Tce_timeout      int
	Tce_listsfile    string
	Tce_maxtries     int
	Args             []string
}

func Test_parser_A040(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads = `        {
            "+tce": {
                "mirror": "http://repo.tinycorelinux.net",
                "root": "/",
                "tceversion": "7.x",
                "wget": "wget",
                "cat": "cat",
                "rm": "rm",
                "sudoprefix": "sudo",
                "optional_dir": "/cde",
                "trymode": false,
                "platform": "x86_64",
                "mount": "mount",
                "umount": "umount",
                "chroot": "chroot",
                "chown": "chown",
                "chmod": "chmod",
                "mkdir": "mkdir",
                "rollback": true,
                "cp": "cp",
                "jsonfile": null,
                "perspace": 3,
                "depmapfile": null,
                "timeout": 10,
                "listsfile": null,
                "maxtries": 5
            }
        }`
	var p *parserTest40
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	p = &parserTest40{}
	_, err = parser.ParseCommandLineEx([]string{"--tce-root", "/home/"}, nil, p, nil)
	check_equal(t, err, nil)
	check_equal(t, p.Tce_mirror, "http://repo.tinycorelinux.net")
	check_equal(t, p.Tce_root, "/home/")
	check_equal(t, p.Tce_listsfile, "")
	check_equal(t, p.Tce_maxtries, 5)
	check_equal(t, p.Tce_timeout, 10)
	return
}

func Test_parser_A041(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads string
	var loads_fmt = `        {
            "countryname|N" : "CN",
            "statename|S" : "ZJ",
            "localityname" : "HZ",
            "organizationname|O" : ["BT"],
            "organizationunitname" : "BT R&D",
            "commonname|C" : "bingte.com",
            "+ssl" : {
                "chain" : true,
                "dir" : "%s",
                "bits" : 4096,
                "md" : "sha256",
                "utf8" : true,
                "name" : "ipxe",
                "days" : 3650,
                "crl-days": 365,
                "emailaddress" : "bt@bingte.com",
                "aia_url" : "http://bingte.com/sec/aia",
                "crl_url" : "http://bingte.com/sec/crl",
                "ocsp_url" : "http://bingte.com/sec/ocsp",
                "dns_url" : ["bingte.com"],
                "excluded_ip" : ["0.0.0.0/0.0.0.0","0:0:0:0:0:0:0:0/0:0:0:0:0:0:0:0"],
                "password|P" : null,
                "copy_extensions" : "none",
                "subca" : false,
                "comment": ""
            }
        }`
	var args *NameSpaceEx
	var fdir string
	var ok bool = false
	var jsonfile string
	beforeParser(t)
	fdir = getCallerFilename(1)
	check_not_equal(t, fdir, "")
	fdir = joinFile(filepath.Dir(fdir), "certs")
	check_not_equal(t, fdir, "")
	fdir = strings.Replace(fdir, "\\", "\\\\", -1)
	loads = fmt.Sprintf(loads_fmt, fdir)

	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	jsonfile = makeWriteTempFile(`{"emailaddress" : "unit@bingte.com","organizationname" : "BT RD","ssl" :{ "dir" : "./certs/bingte","name" : "bingte","subca" : true,"copy_extensions" : "copy","days" : 375,"crl_days" : 30,"bits" : 4096}}`)
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	args, err = parser.ParseCommandLine([]string{"--json", jsonfile}, parser)
	check_not_equal(t, err, nil)
	args = args
	ok = true
	return
}

func Test_parser_A042(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads = `        {
            "verbose|v" : "+",
            "kernel|K" : "/boot/",
            "initrd|I" : "/boot/",
            "encryptfile|e" : null,
            "encryptkey|E" : null,
            "setupsectsoffset" : 663,
            "ipxe" : {
                "$" : "+"
            }
        }`
	var args *NameSpaceEx
	beforeParser(t)
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine([]string{"-vvvK", "kernel", "--initrd", "initrd", "cc", "dd", "-E", "encryptkey", "-e", "encryptfile", "ipxe"}, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetString("subcommand"), "ipxe")
	check_equal(t, args.GetArray("subnargs"), []string{"cc", "dd"})
	return
}

func Test_parser_A043(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads = `        {
            "verbose|v" : "+",
            "kernel|K" : "/boot/",
            "initrd|I" : "/boot/",
            "encryptfile|e" : null,
            "encryptkey|E" : null,
            "setupsectsoffset" : 663,
            "ipxe" : {
                "$" : "+"
            }
        }`
	var args *NameSpaceEx
	var options *ExtArgsOptions
	beforeParser(t)
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"parseall": true,"%s" : "-", "%s" : "-"}`, OPT_LONG_PREFIX, OPT_SHORT_PREFIX))
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine([]string{"-K", "kernel", "-initrd", "initrd", "cc", "dd", "-E", "encryptkey", "-e", "encryptfile", "ipxe"}, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetString("subcommand"), "ipxe")
	check_equal(t, args.GetArray("subnargs"), []string{"cc", "dd"})
	return
}

func Test_parser_A044(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads = `        {
            "verbose|v" : "+",
            "kernel|K" : "/boot/",
            "initrd|I" : "/boot/",
            "encryptfile|e" : null,
            "encryptkey|E" : null,
            "setupsectsoffset" : 663,
            "ipxe" : {
                "$" : "+"
            }
        }`
	var args *NameSpaceEx
	var options *ExtArgsOptions
	beforeParser(t)
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"parseall": true,"%s" : "++", "%s" : "+"}`, OPT_LONG_PREFIX, OPT_SHORT_PREFIX))
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine([]string{"+K", "kernel", "++initrd", "initrd", "cc", "dd", "+E", "encryptkey", "+e", "encryptfile", "ipxe"}, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetString("subcommand"), "ipxe")
	check_equal(t, args.GetArray("subnargs"), []string{"cc", "dd"})
	return
}

func Test_parser_A045(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads_fmt = `        {
            "verbose|v" : "+",
            "kernel|K" : "/boot/",
            "initrd|I" : "/boot/",
            "pair|P!optparse=%s.debug_set_2_args!" : [],
            "encryptfile|e" : null,
            "encryptkey|E" : null,
            "setupsectsoffset" : 663,
            "ipxe" : {
                "$" : "+"
            }
        }`
	var loads string
	var args *NameSpaceEx
	var options *ExtArgsOptions
	var pkgname string
	beforeParser(t)
	pkgname = getCallerPackage(0)
	check_not_equal(t, pkgname, "")
	loads = fmt.Sprintf(loads_fmt, pkgname)
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"parseall": true,"%s" : "++", "%s" : "+","%s": false}`, OPT_LONG_PREFIX, OPT_SHORT_PREFIX, OPT_FUNC_UPPER_CASE))
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine([]string{"+K", "kernel", "++pair", "initrd", "cc", "dd", "+E", "encryptkey", "+e", "encryptfile", "ipxe"}, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetString("subcommand"), "ipxe")
	check_equal(t, args.GetArray("pair"), []string{"initrd", "cc"})
	check_equal(t, args.GetArray("subnargs"), []string{"dd"})
	return
}

func Test_parser_A046(t *testing.T) {
	var err error
	var parser *ExtArgsParse
	var loads_fmt = `        {
            "verbose|v" : "+",
            "kernel|K" : "/boot/",
            "initrd|I" : "/boot/",
            "pair|P!optparse=%s.debug_set_2_args;opthelp=%s.debug_opthelp_set!" : [],
            "encryptfile|e" : null,
            "encryptkey|E" : null,
            "setupsectsoffset" : 663,
            "ipxe" : {
                "$" : "+"
            }
        }`
	var loads string
	var options *ExtArgsOptions
	var pkgname string
	var sarr []string
	var expr *regexp.Regexp
	var ok bool = false
	var c string
	beforeParser(t)
	pkgname = getCallerPackage(0)
	check_not_equal(t, pkgname, "")
	loads = fmt.Sprintf(loads_fmt, pkgname, pkgname)
	options, err = NewExtArgsOptions(fmt.Sprintf(`{"parseall": true,"%s" : "++", "%s" : "+", "%s": false}`, OPT_LONG_PREFIX, OPT_SHORT_PREFIX, OPT_FUNC_UPPER_CASE))
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "")
	expr = regexp.MustCompile(`.*opthelp function set \[pair\].*`)
	for _, c = range sarr {
		if expr.MatchString(c) {
			ok = true
			break
		}
	}
	check_equal(t, ok, true)
	return
}

func Test_parser_A047(t *testing.T) {
	var err error
	var loads = `        {
            "verbose|v" : "+",
            "kernel|K" : "/boot/",
            "initrd|I" : "/boot/",
            "pair|P!optparse=debug_set_2_args;opthelp=debug_opthelp_set!" : [],
            "encryptfile|e" : null,
            "encryptkey|E" : null,
            "setupsectsoffset" : 663,
            "ipxe" : {
                "$" : "+"
            }
        }`
	var sarr []string
	var expr *regexp.Regexp
	var ok bool = false
	var cl *compileExec
	var c string
	var setvars map[string]string
	var confstr = fmt.Sprintf(`{"parseall": true,"%s" : "++", "%s" : "+" , "%s" : "?" , "%s": "usage" , "%S" : "jsonfile", "%s": false}`, OPT_LONG_PREFIX, OPT_SHORT_PREFIX, OPT_HELP_SHORT, OPT_HELP_LONG, OPT_JSON_LONG, OPT_FUNC_UPPER_CASE)
	var codestr = `func format_error(fmtstr string, a ...interface{}) string {
	return fmt.Sprintf(fmtstr, a...)
}

func debug_set_2_args(ns *extargsparse.NameSpaceEx, validx int, keycls *extargsparse.ExtKeyParse, params []string) (step int, err error) {
	var sarr []string
	if ns == nil {
		return 0, nil
	}
	if (validx + 2) > len(params) {
		return 0, fmt.Errorf("%s", format_error("[%d+2] > len(%d) %v", validx, len(params), params))
	}
	sarr = ns.GetArray(keycls.Optdest())
	sarr = append(sarr, params[validx])
	sarr = append(sarr, params[(validx+1)])
	ns.SetValue(keycls.Optdest(), sarr)
	return 2, nil
}

func debug_opthelp_set(keycls *extargsparse.ExtKeyParse) string {
	if keycls == nil {
		return ""
	}
	return fmt.Sprintf("opthelp function set [%s] default value (%s)", keycls.Optdest(), keycls.Value())
}

func init(){
	debug_set_2_args(nil,1,nil,[]string{})
	debug_opthelp_set(nil)
}`
	beforeParser(t)
	cl = newComileExec()
	check_not_equal(t, cl, (*compileExec)(nil))
	defer func() { cl.Release(ok) }()
	err = cl.WriteScript(confstr, loads, []string{}, codestr, nil, false, "args", "ppc")
	check_equal(t, err, nil)
	err = cl.Compile()
	check_equal(t, err, nil)
	setvars = make(map[string]string)
	err = cl.RunCmd(setvars, []string{}, "++usage")
	check_equal(t, err, nil)

	sarr = cl.GetOut()
	expr = regexp.MustCompile(`.*opthelp function set \[pair\].*`)
	for _, c = range sarr {
		if expr.MatchString(c) {
			ok = true
			break
		}
	}
	check_equal(t, ok, true)
	return
}

func Test_parser_A048(t *testing.T) {
	var err error
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
	var ok bool = false
	var confstr = fmt.Sprintf(`{"%s": "jsonfile"}`, OPT_JSON_LONG)
	var jsonfile string
	var depjsonfile string
	var depstrval string = `newval`
	var depliststr string = `["depenv1","depenv2"]`
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var args *NameSpaceEx

	beforeParser(t)

	jsonfile = makeWriteTempFile(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
	os.Setenv("EXTARGSPARSE_JSONFILE", jsonfile)
	os.Setenv("DEP_JSONFILE", depjsonfile)
	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, []int{ENV_COMMAND_JSON_SET, ENVIRONMENT_SET, ENV_SUB_COMMAND_JSON_SET})
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	os.Setenv("DEP_STRING", depstrval)
	os.Setenv("DEP_LIST", depliststr)
	args, err = parser.ParseCommandLine([]string{"-p", "9000", "dep", "--dep-string", "ee", "ww"}, nil)
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

func Test_parser_A049(t *testing.T) {
	var err error
	var loads = `        {
            "verbose|v##very long very long very long very long very long very long very long very long very long very long very long very long very long very long very long very long very long very long##" : "+",
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
	var confstr = fmt.Sprintf(`{"screenwidth": 60}`)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var sio *stringIO
	var sarr []string
	var overlength bool
	var idx int
	var c string

	beforeParser(t)

	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	sio = newStringIO()
	check_not_equal(t, sio, (*stringIO)(nil))
	sarr = getCmdHelp(parser, "")
	overlength = false
	for idx, c = range sarr {
		if len(c) > 64 && idx > 0 {
			overlength = true
			break
		}
	}
	check_equal(t, overlength, false)
	confstr = fmt.Sprintf(`{"screenwidth" : 80}`)
	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	sio = newStringIO()
	check_not_equal(t, sio, (*stringIO)(nil))
	sarr = getCmdHelp(parser, "")
	overlength = false
	for idx, c = range sarr {
		if len(c) > 64 && idx > 0 {
			overlength = true
			break
		}
	}
	check_equal(t, overlength, true)
	return
}

func Test_parser_A050(t *testing.T) {
	var err error
	var loads = `        {
            "verbose|v" : "+",
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var confstr = fmt.Sprintf(`{"%s": "usage", "%s" : "?" , "%s" : "++", "%s" : "+"}`, OPT_HELP_LONG, OPT_HELP_SHORT, OPT_LONG_PREFIX, OPT_SHORT_PREFIX)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var sarr []string
	var ok bool
	var c string
	var matchexpr *regexp.Regexp

	beforeParser(t)

	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "")
	ok = false
	matchexpr = regexp.MustCompile(`^\s+\+\+usage|\+\?\s+to display.*`)
	for _, c = range sarr {
		if matchexpr.MatchString(c) {
			ok = true
			break
		}
	}
	check_equal(t, ok, true)
	return
}

func Test_parser_A051(t *testing.T) {
	var err error
	var loads = `        {
            "verbose|v" : "+",
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var confstr = fmt.Sprintf(`        {
            "%s" : "usage",
            "%s" : null,
            "%s" : "++",
            "%s" : "+"
        }`, OPT_HELP_LONG, OPT_HELP_SHORT, OPT_LONG_PREFIX, OPT_SHORT_PREFIX)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var sarr []string
	var ok bool
	var c string
	var matchexpr *regexp.Regexp

	beforeParser(t)

	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "")
	ok = false
	matchexpr = regexp.MustCompile(`^\s+\+\+usage|\+\?\s+to display.*`)
	for _, c = range sarr {
		if matchexpr.MatchString(c) {
			ok = true
			break
		}
	}
	check_equal(t, ok, true)
	return
}

func Test_parser_A052(t *testing.T) {
	var err error
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
	var confstr = fmt.Sprintf(`        {
            "%s" : true,
            "%s" : true
        }`, OPT_NO_JSON_OPTION, OPT_NO_HELP_OPTION)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var args *NameSpaceEx
	var sarr []string
	var ok bool
	var jsonfile string
	var depjsonfile string
	var c string
	var depstrval string = `newval`
	var depliststr string = `["depenv1","depenv2"]`
	var helpexpr, jsonexpr *regexp.Regexp
	var helpok, jsonok bool

	beforeParser(t)
	ok = false
	jsonfile = makeWriteTempFile(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
	os.Setenv("EXTARGSPARSE_JSONFILE", jsonfile)
	os.Setenv("DEP_JSONFILE", depjsonfile)

	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, []int{ENV_COMMAND_JSON_SET, ENVIRONMENT_SET, ENV_SUB_COMMAND_JSON_SET})
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	os.Setenv("DEP_STRING", depstrval)
	os.Setenv("DEP_LIST", depliststr)
	sarr = getCmdHelp(parser, "")
	helpok = false
	jsonok = false
	helpexpr = regexp.MustCompile(`^\s+--help.*`)
	jsonexpr = regexp.MustCompile(`^\s+--json.*`)
	for _, c = range sarr {
		if helpexpr.MatchString(c) {
			helpok = true
		}
		if jsonexpr.MatchString(c) {
			jsonok = true
		}
	}
	check_equal(t, helpok, false)
	check_equal(t, jsonok, false)
	args, err = parser.ParseCommandLine([]string{"-p", "9000", "dep", "--dep-string", "ee", "ww"}, nil)
	check_equal(t, args.GetInt("verbose"), 0)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"depenv1", "depenv2"})
	check_equal(t, args.GetString("dep_string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A053(t *testing.T) {
	var err error
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
            },
            "rdep" : {
                "list|l" : [],
                "string|s" : "s_rdep",
                "$" : "+"
            }
        }`
	var confstr = fmt.Sprintf(`        {
            "cmdprefixadded" : false
        }`)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var args *NameSpaceEx
	var sarr []string
	var ok bool
	var jsonfile string
	var depjsonfile string
	var c string
	var depstrval string = `newval`
	var depliststr string = `["depenv1","depenv2"]`
	var helpexpr, jsonexpr, listexpr, stringexpr *regexp.Regexp
	var helpok, jsonok, listok, stringok bool

	beforeParser(t)
	ok = false
	jsonfile = makeWriteTempFile(`{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring","port":6000,"verbose":3}`)
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
	os.Setenv("EXTARGSPARSE_JSON", jsonfile)
	os.Setenv("DEP_JSON", depjsonfile)

	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, []int{ENV_COMMAND_JSON_SET, ENVIRONMENT_SET, ENV_SUB_COMMAND_JSON_SET})
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	os.Setenv("DEP_STRING", depstrval)
	os.Setenv("DEP_LIST", depliststr)
	sarr = getCmdHelp(parser, "dep")
	helpok = false
	jsonok = false
	listok = false
	stringok = false
	helpexpr = regexp.MustCompile(`^\s+--help.*`)
	jsonexpr = regexp.MustCompile(`^\s+--dep-json.*`)
	listexpr = regexp.MustCompile(`^\s+--list.*`)
	stringexpr = regexp.MustCompile(`^\s+--string.*`)
	for _, c = range sarr {
		if helpexpr.MatchString(c) {
			helpok = true
		}
		if jsonexpr.MatchString(c) {
			jsonok = true
		}
		if listexpr.MatchString(c) {
			listok = true
		}

		if stringexpr.MatchString(c) {
			stringok = true
		}
	}
	check_equal(t, helpok, true)
	check_equal(t, jsonok, true)
	check_equal(t, listok, true)
	check_equal(t, stringok, true)
	sarr = getCmdHelp(parser, "rdep")
	helpok = false
	jsonok = false
	listok = false
	stringok = false
	helpexpr = regexp.MustCompile(`^\s+--help.*`)
	jsonexpr = regexp.MustCompile(`^\s+--rdep-json.*`)
	listexpr = regexp.MustCompile(`^\s+--list.*`)
	stringexpr = regexp.MustCompile(`^\s+--string.*`)
	for _, c = range sarr {
		if helpexpr.MatchString(c) {
			helpok = true
		}
		if jsonexpr.MatchString(c) {
			jsonok = true
		}
		if listexpr.MatchString(c) {
			listok = true
		}

		if stringexpr.MatchString(c) {
			stringok = true
		}
	}
	check_equal(t, helpok, true)
	check_equal(t, jsonok, true)
	check_equal(t, listok, true)
	check_equal(t, stringok, true)
	args, err = parser.ParseCommandLine([]string{"-p", "9000", "dep", "--string", "ee", "ww"}, nil)
	check_equal(t, args.GetInt("verbose"), 3)
	check_equal(t, args.GetInt("port"), 9000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("list"), []string{"jsonval1", "jsonval2"})
	check_equal(t, args.GetString("string"), "ee")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A054(t *testing.T) {
	var err error
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
	var confstr = fmt.Sprintf(`        {
            "%s" : "jsonfile"
        }`, OPT_JSON_LONG)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var args *NameSpaceEx
	var ok bool
	var jsonfile string
	var depjsonfile string
	var depstrval string = `newval`
	var depliststr string = `["depenv1","depenv2"]`

	beforeParser(t)
	ok = false
	jsonfile = makeWriteTempFile(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
	os.Setenv("EXTARGSPARSE_JSONFILE", jsonfile)
	os.Setenv("DEP_JSONFILE", depjsonfile)

	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, []int{ENV_COMMAND_JSON_SET, ENVIRONMENT_SET, ENV_SUB_COMMAND_JSON_SET})
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	os.Setenv("DEP_STRING", depstrval)
	os.Setenv("DEP_LIST", depliststr)
	args, err = parser.ParseCommandLine([]string{"--jsonfile", jsonfile, "dep", "ww"}, nil)
	check_equal(t, args.GetInt("verbose"), 3)
	check_equal(t, args.GetInt("port"), 6000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"jsonval1", "jsonval2"})
	check_equal(t, args.GetString("dep_string"), "jsonstring")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A055(t *testing.T) {
	var err error
	var loads_fmt = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l!jsonfunc=%s.debug_2_jsonfunc!" : [],
                "string|s!jsonfunc=%s.debug_upper_jsonfunc!" : "s_var",
                "$" : "+"
            }
        }`
	var loads string
	var confstr = fmt.Sprintf(`        {
            "%s" : "jsonfile"
        }`, OPT_JSON_LONG)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var args *NameSpaceEx
	var ok bool
	var jsonfile string
	var depjsonfile string
	var depstrval string = `newval`
	var depliststr string = `["depenv1","depenv2"]`
	var pkgname string

	beforeParser(t)
	ok = false
	pkgname = getCallerPackage(1)
	check_not_equal(t, pkgname, "")
	loads = fmt.Sprintf(loads_fmt, pkgname, pkgname)
	jsonfile = makeWriteTempFile(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`)
	defer func() { safeRemoveFile(jsonfile, "jsonfile", ok) }()
	depjsonfile = makeWriteTempFile(`{"list":["depjson1","depjson2"]}`)
	defer func() { safeRemoveFile(depjsonfile, "depjsonfile", ok) }()
	os.Setenv("EXTARGSPARSE_JSONFILE", jsonfile)
	os.Setenv("DEP_JSONFILE", depjsonfile)

	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, []int{ENV_COMMAND_JSON_SET, ENVIRONMENT_SET, ENV_SUB_COMMAND_JSON_SET})
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	os.Setenv("DEP_STRING", depstrval)
	os.Setenv("DEP_LIST", depliststr)
	args, err = parser.ParseCommandLine([]string{"--jsonfile", jsonfile, "dep", "ww"}, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 3)
	check_equal(t, args.GetInt("port"), 6000)
	check_equal(t, args.GetString("subcommand"), "dep")
	check_equal(t, args.GetArray("dep_list"), []string{"jsonval1"})
	check_equal(t, args.GetString("dep_string"), "JSONSTRING")
	check_equal(t, args.GetArray("subnargs"), []string{"ww"})
	ok = true
	return
}

func Test_parser_A056(t *testing.T) {
	var err error
	var loads = `        {
            "asn1parse" : {
                "$" : 0,
                "$inform!optparse=inform_optparse;completefunc=inform_complete!" : null,
                "$in" : null,
                "$out" : null,
                "$noout" : false,
                "$offset" : 0,
                "$length" : -1,
                "$dump" : false,
                "$dlimit" : -1,
                "$oid" : null,
                "$strparse" : 0,
                "$genstr" : null,
                "$genconf" : null
            }
        }`
	var confstr = fmt.Sprintf(`        {
            "%s" : "-",
            "%s" : "-",
            "%s" : true,
            "cmdprefixadded" : false
        }`, OPT_LONG_PREFIX, OPT_SHORT_PREFIX, OPT_NO_JSON_OPTION)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var sarr []string
	var optnames []string = []string{"inform", "in", "out", "noout", "offset", "length", "dump", "dlimit", "oid", "strparse", "genstr", "genconf", "help"}
	var c string
	var ok bool
	var opts []*ExtKeyParse
	var opt *ExtKeyParse

	beforeParser(t)

	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	sarr, err = parser.GetSubCommands("")
	check_equal(t, err, nil)
	check_equal(t, sarr, []string{"asn1parse"})
	opts, err = parser.GetCmdOpts("asn1parse")
	check_equal(t, err, nil)
	for _, opt = range opts {
		if !opt.IsFlag() || opt.TypeName() == "args" {
			continue
		}
		ok = false
		if opt.TypeName() == "help" {
			if opt.Longopt() == "-help" && opt.Shortopt() == "-h" {
				ok = true
			}

		} else {
			for _, c = range optnames {
				if c == opt.Optdest() && fmt.Sprintf("-%s", c) == opt.Longopt() {
					ok = true
				}
			}

		}
		check_equal(t, ok, true)
	}

	return
}

func Test_parser_A057(t *testing.T) {
	var err error
	var loads = `        {
            "asn1parse" : {
                "$" : 0,
                "$inform!optparse=inform_optparse;completefunc=inform_complete!" : null,
                "$in" : null,
                "$out" : null,
                "$noout" : false,
                "$offset" : 0,
                "$length" : -1,
                "$dump" : false,
                "$dlimit" : -1,
                "$oid" : null,
                "$strparse" : 0,
                "$genstr" : null,
                "$genconf" : null
            },
            "ca" : {
                "$" : 0,
                "$config" : null,
                "$name" : null,
                "$in" : null,
                "$ss_cert" : null,
                "$spkac" : null,
                "$infiles" : null,
                "$out" : null,
                "$outdir" : null,
                "$cert" : null,
                "$keyfile" : null,
                "$keyform!optparse=inform_optparse;completefunc=inform_complete!" : null,
                "$key" : null,
                "$selfsign" : false,
                "$passin" : null,
                "$verbose" : "+",
                "$notext" : false,
                "$startdate" : null,
                "$enddate" : null,
                "$days" : 30,
                "$md" : null,
                "$policy" : null,
                "$preserveDN" : false,
                "$msie_hack" : false,
                "$noemailDN" : false,
                "$batch" : false,
                "$extensions" : null,
                "$extfile" : null,
                "$engine" : null,
                "$subj" : null,
                "$utf8" : false,
                "$multivalue-rdn" : false,
                "$gencrl" : false,
                "$crldays" : 30,
                "$crlhours" : -1,
                "$revoke" : null,
                "$status" : null,
                "$updatedb" : false,
                "$crl_reason" : null,
                "$crl_hold" : null,
                "$crl_compromise" : null,
                "$crl_CA_compromise" : null,
                "$crlexts" : null
            }
        }`
	var confstr = fmt.Sprintf(`        {
            "%s" : "-",
            "%s" : "-",
            "%s" : true,
            "cmdprefixadded" : false,
            "flagnochange" : true
        }`, OPT_LONG_PREFIX, OPT_SHORT_PREFIX, OPT_NO_JSON_OPTION)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var sarr []string
	var optnames []string = []string{"config", "name", "in", "ss_cert", "spkac", "infiles", "out", "outdir", "cert", "keyfile", "keyform", "key", "selfsign", "passin", "verbose", "notext", "startdate", "enddate", "days", "md", "policy", "preserveDN", "msie_hack", "noemailDN", "batch", "extensions", "extfile", "engine", "subj", "utf8", "gencrl", "crldays", "crlhours", "revoke", "status", "updatedb", "crl_reason", "crl_hold", "crl_compromise", "crl_CA_compromise", "crlexts"}
	var c string
	var ok bool
	var opts []*ExtKeyParse
	var opt *ExtKeyParse

	beforeParser(t)

	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	sarr, err = parser.GetSubCommands("")
	check_equal(t, err, nil)
	check_equal(t, sarr, []string{"asn1parse", "ca"})
	opts, err = parser.GetCmdOpts("ca")
	check_equal(t, err, nil)
	for _, opt = range opts {
		if !opt.IsFlag() || opt.TypeName() == "args" {
			continue
		}
		ok = false
		if opt.TypeName() == "help" {
			if opt.Longopt() == "-help" && opt.Shortopt() == "-h" {
				ok = true
			}

		} else if opt.Longopt() == "-multivalue-rdn" {
			if opt.Optdest() == "multivalue_rdn" {
				ok = true
			}
		} else {
			for _, c = range optnames {
				if c == opt.Optdest() && fmt.Sprintf("-%s", c) == opt.Longopt() {
					ok = true
				}
			}
		}
		check_equal(t, ok, true)
	}

	return
}

func Test_parser_A058(t *testing.T) {
	var err error
	var loads = `        {
            "verbose" : "+",
            "dep" : {
                "$" : "*"
            },
            "rdep" : {
                "$" : "*"
            }
        }`
	var confstr = fmt.Sprintf(`{}`)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var sarr []string
	var matchexpr *regexp.Regexp

	beforeParser(t)

	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "")
	matchexpr = regexp.MustCompile(`.*\[OPTIONS\]\s+\[SUBCOMMANDS\]\s+.*`)
	check_equal(t, matchexpr.MatchString(sarr[0]), true)
	return
}

func Test_parser_A059(t *testing.T) {
	var err error
	var loads_fmt = `        {
            "verbose|v" : "+",
            "kernel|K" : "/boot/",
            "initrd|I" : "/boot/",
            "pair|P!optparse=%s.Debug_set_2_args!" : [],
            "encryptfile|e" : null,
            "encryptkey|E" : null,
            "setupsectsoffset" : 663,
            "ipxe" : {
                "$" : "+"
            }
        }`
	var loads string
	var confstr = fmt.Sprintf(`{"parseall" : true,"%s" : "++", "%s" : "+"}`, OPT_LONG_PREFIX, OPT_SHORT_PREFIX)
	var options *ExtArgsOptions
	var parser *ExtArgsParse
	var args *NameSpaceEx
	var pkgname string

	beforeParser(t)
	pkgname = getCallerPackage(1)
	loads = fmt.Sprintf(loads_fmt, pkgname)
	options, err = NewExtArgsOptions(confstr)
	check_equal(t, err, nil)
	parser, err = NewExtArgsParse(options, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine([]string{"+K", "kernel", "++pair", "initrd", "cc", "dd", "+E", "encryptkey", "+e", "encryptfile", "ipxe"}, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetString("subcommand"), "ipxe")
	check_equal(t, args.GetArray("subnargs"), []string{"dd"})
	check_equal(t, args.GetArray("pair"), []string{"INITRD", "CC"})
	return
}

func Test_parser_A060(t *testing.T) {
	var loads = `{
    "dep": {
        "$": "*",
        "ip": {
            "$": "*"
        }
    },
    "rdep" : {
        "$" : "*",
        "ip" : {
            "$" : "*"
        }
    }
}`
	var parser *ExtArgsParse
	var sarr []string
	var err error
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "")
	err = assertOkCmds(t, sarr, parser, "")
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "dep")
	err = assertOkCmds(t, sarr, parser, "dep")
	check_equal(t, err, nil)
	sarr = getCmdHelp(parser, "rdep")
	err = assertOkCmds(t, sarr, parser, "rdep")
	check_equal(t, err, nil)
	return
}
