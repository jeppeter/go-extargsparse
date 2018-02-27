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
				strings.HasPrefix(k, "EXTARGSPARSE_") ||
				strings.HasPrefix(k, "HTTP_") ||
				strings.HasPrefix(k, "SSL_") {
				err = os.Unsetenv(sarr[0])
				if err == nil {
					delone = true
					break
				} else {
					keyDebug("delete [%s] error", sarr[0])
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
	parser := NewExtArgsParse(nil, nil)
	return
}
