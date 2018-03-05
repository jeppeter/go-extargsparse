package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func Flag_parse(ns *extargsparse.NameSpaceEx, validx int, keycls *extargsparse.ExtKeyParse, params []string) (step int, err error) {
	if ns == nil {
		return 0, nil
	}
	fmt.Fprintf(os.Stdout, "Attr=%s\n", keycls.Attr(""))
	fmt.Fprintf(os.Stdout, "opthelp=%s\n", keycls.Attr("opthelp"))
	fmt.Fprintf(os.Stdout, "optparse=%s\n", keycls.Attr("optparse"))
	ns.SetValue(keycls.Optdest(), []string{params[validx]})
	return 1, nil
}

func Flag_help(keycls *extargsparse.ExtKeyParse) string {
	if keycls == nil {
		return ""
	}
	return fmt.Sprintf("flag special set []")
}

func init() {
	Flag_help(nil)
	Flag_parse(nil, 0, nil, []string{})
}

func main() {
	var loads = `{
		"flag|f!optparse=flag_parse;opthelp=flag_help!" : []
	}`
	var err error
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{}`))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				parser.ParseCommandLine(nil, nil)
			}
		}
	}
}

/*
	when call ./cmd -h will see the  flag special set []
	when call ./cmd -f cc  see the flag
*/
