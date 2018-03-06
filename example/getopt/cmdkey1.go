package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var err error
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
	var confstr = fmt.Sprintf(`        {
            "%s" : true,
            "%s" : true
        }`, extargsparse.OPT_NO_JSON_OPTION, extargsparse.OPT_NO_HELP_OPTION)
	var options *extargsparse.ExtArgsOptions
	var parser *extargsparse.ExtArgsParse
	var keycls *extargsparse.ExtKeyParse
	options, err = extargsparse.NewExtArgsOptions(confstr)
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(loads)
			if err == nil {
				keycls, err = parser.GetCmdKey("")
				if err == nil {
					fmt.Fprintf(os.Stdout, "cmdname=%s\n", keycls.CmdName()) // cmdname=main
					keycls, err = parser.GetCmdKey("dep")
					if err == nil {
						fmt.Fprintf(os.Stdout, "cmdname=%s\n", keycls.CmdName()) // cmdname=dep
						keycls, err = parser.GetCmdKey("rdep.ip")
						if err == nil {
							fmt.Fprintf(os.Stdout, "cmdname=%s\n", keycls.CmdName()) // cmdname=ip  it is the subcommand of subcommand rdep
						}
					}
				}
			}
		}
	}
	return
}
