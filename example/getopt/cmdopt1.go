package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

func main() {
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
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var opts []*extargsparse.ExtKeyParse
	var flag *extargsparse.ExtKeyParse
	var i int
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{"%s" : "cmd1"}`, extargsparse.OPT_PROG))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				opts, err = parser.GetCmdOpts("")
				if err == nil {
					fmt.Fprintf(os.Stdout, "main cmd opts:\n")
					for i, flag = range opts {
						if flag.TypeName() == "args" {
							fmt.Fprintf(os.Stdout, "[%d].type=args\n", i)
						} else {
							fmt.Fprintf(os.Stdout, "[%d].Longopt=%s;.Shortopt=%s;Optdest=%s;attr=%s\n", i, flag.Longopt(), flag.Shortopt(), flag.Optdest(), flag.Attr(""))
						}
					}
					opts, err = parser.GetCmdOpts("dep")
					if err == nil {
						fmt.Fprintf(os.Stdout, "dep cmd opts:\n")
						for i, flag = range opts {
							if flag.TypeName() == "args" {
								fmt.Fprintf(os.Stdout, "[%d].type=args\n", i)
							} else {
								fmt.Fprintf(os.Stdout, "[%d].Longopt=%s;.Shortopt=%s;Optdest=%s;attr=%s\n", i, flag.Longopt(), flag.Shortopt(), flag.Optdest(), flag.Attr(""))
							}
						}

						opts, err = parser.GetCmdOpts("rdep.ip")
						if err == nil {
							fmt.Fprintf(os.Stdout, "rdep.ip cmd opts:\n")
							for i, flag = range opts {
								if flag.TypeName() == "args" {
									fmt.Fprintf(os.Stdout, "[%d].type=args\n", i)
								} else {
									fmt.Fprintf(os.Stdout, "[%d].Longopt=%s;.Shortopt=%s;Optdest=%s;attr=%s\n", i, flag.Longopt(), flag.Shortopt(), flag.Optdest(), flag.Attr(""))
								}
								/*
									Output:
									main cmd opts:
									[0].type=args
									[1].Longopt=--help;.Shortopt=-h;Optdest=help;attr=
									[2].Longopt=--json;.Shortopt=;Optdest=json;attr=
									[3].Longopt=--port;.Shortopt=-p;Optdest=port;attr=
									[4].Longopt=--http-url;.Shortopt=-u;Optdest=http_url;attr=
									[5].Longopt=--verbose;.Shortopt=-v;Optdest=verbose;attr=
									[6].Longopt=--http-visual-mode;.Shortopt=-V;Optdest=http_visual_mode;attr=
									dep cmd opts:
									[0].type=args
									[1].Longopt=--help;.Shortopt=-h;Optdest=help;attr=
									[2].Longopt=--dep-json;.Shortopt=;Optdest=dep_json;attr=
									[3].Longopt=--dep-list;.Shortopt=-l;Optdest=dep_list;attr=[attr]=[cc]
									[optfunc]=[list_opt_func]

									[4].Longopt=--dep-string;.Shortopt=-s;Optdest=dep_string;attr=
									rdep.ip cmd opts:
									[0].type=args
									[1].Longopt=--rdep-ip-cc;.Shortopt=;Optdest=rdep_ip_cc;attr=
									[2].Longopt=--help;.Shortopt=-h;Optdest=help;attr=
									[3].Longopt=--rdep-ip-json;.Shortopt=;Optdest=rdep_ip_json;attr=
									[4].Longopt=--rdep-ip-list;.Shortopt=;Optdest=rdep_ip_list;attr=
									[5].Longopt=--rdep-ip-verbose;.Shortopt=;Optdest=rdep_ip_verbose;attr=
								*/
							}
						}
					}
				}
			}
		}
	}

	return

}
