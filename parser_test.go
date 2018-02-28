package extargsparse

import (
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
*/

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
