package extargsparse

import (
	"fmt"
	"os"
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
				strings.HasPrefix(k, "SSL_") {
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

type parserTest1 struct {
	Verbose int
	Flag    bool
	Number  int
	List    []string
	String  string
	Args    []string
}

/*
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
	beforeParser()
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
	beforeParser()
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
*/

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
