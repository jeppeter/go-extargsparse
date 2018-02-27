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

func Test_parser_A001(t *testing.T) {
	return
}
