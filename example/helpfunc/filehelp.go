package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
	var fout *extargsparse.FileIoWriter
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
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var f *os.File
	f, err = os.OpenFile("help.out", os.O_RDWR|os.O_CREATE, 0640)
	if err == nil {
		defer f.Close()
		fout = extargsparse.NewFileWriter(f)

		options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{"prog" : "cmd1"}`))
		if err == nil {
			parser, err = extargsparse.NewExtArgsParse(options, nil)
			if err == nil {
				err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
				if err == nil {

					parser.PrintHelp(fout, "")
					/*
						in help.out
							OUTPUT:
								cmd1 0.0.1  [OPTIONS] [SUBCOMMANDS] [args...]'

								[OPTIONS]
								    --json                 json
								                    json input file to get the value set
								    --help|-h
								                    to display this help information
								    --verbose|-v           verbose
								                    verbose set default(0)
								    --http-url|-u          http_url
								                    http_url set default(http://www.google.com)
								    --http-visual-mode|-V
								                    http_visual_mode set true default(False)
								    --port|-p              port
								                    port to connect

								[SUBCOMMANDS]
								    [rdep]  rdep handler
								    [dep]   dep handler
					*/
					parser.PrintHelp(fout, "dep")
					/*
						in help.out
							OUTPUT:
								cmd1 0.0.1  dep [OPTIONS] [SUBCOMMANDS] args...

								[OPTIONS]
								    --dep-json       dep_json    json input file to get the value set
								    --help|-h                    to display this help information
								    --dep-list|-l    dep_list    dep_list set default([])
								    --dep-string|-s  dep_string  dep_string set default(s_var)

								[SUBCOMMANDS]
								    [ip]   ip handler

					*/
					parser.PrintHelp(fout, "rdep")
					/*
						in help.out
							OUTPUT:
								cmd1 0.0.1  rdep [OPTIONS] [SUBCOMMANDS] [args...]'

								[OPTIONS]
								    --rdep-json  rdep_json  json input file to get the value set
								    --help|-h               to display this help information

								[SUBCOMMANDS]
								    [ip]    ip handler
					*/
				}
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "create help.out err[%s]\n", err.Error())
	}

}
