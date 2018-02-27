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
	parser, err = NewExtArgsParse(nil, nil)
	check_equal(t, err, nil)
	err = parser.LoadCommandLineString(loads)
	check_equal(t, err, nil)
	args, err = parser.ParseCommandLine(params, nil, nil, nil)
	check_equal(t, err, nil)
	check_equal(t, args.GetInt("verbose"), 4)
	check_equal(t, args.GetBool("flag"), true)
	check_equal(t, args.GetInt("number"), 30)
	check_equal(t, args.GetArray("list"), []string{"bar1", "bar2"})
	check_equal(t, args.GetString("string"), "string_var")
	check_equal(t, args.GetArray("args"), []string{"var1", "var2"})
	return
}
