package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var loads = `        {
            "dep" : {
                "ip" : {
                	"$" : "*"
                },
                "mip" : {
                	"$" : "*"
                }
            },
            "rdep" : {
                "ip" : {
                },
                "rmip" : {                	
                }
            }
        }`
	var err error
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var subcmds []string
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{"%s" : "cmd1"}`, extargsparse.OPT_PROG))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				subcmds, err = parser.GetSubCommands("")
				if err == nil {
					fmt.Fprintf(os.Stdout, "main cmd subcmds:%v\n", subcmds)
					subcmds, err = parser.GetSubCommands("dep")
					if err == nil {
						fmt.Fprintf(os.Stdout, "dep cmd subcmds:%v\n", subcmds)
						subcmds, err = parser.GetSubCommands("rdep.ip")
						if err == nil {
							fmt.Fprintf(os.Stdout, "rdep.ip cmd subcmds:%v\n", subcmds)
							/*
								Output:
								main cmd subcmds:[dep rdep]
								dep cmd subcmds:[ip mip]
								rdep.ip cmd subcmds:[]
							*/
						}
					}
				}
			}
		}
	}

	return

}
