package extargsparse_test

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"io/ioutil"
	"os"
)

func ExampleNewExtArgsOptions() {
	var options *extargsparse.ExtArgsOptions
	var err error
	var confstr = `{
		"screenwidth" : 90.0		
		}`
	options, err = extargsparse.NewExtArgsOptions(confstr)
	if err == nil {
		fmt.Fprintf(os.Stdout, "screenwidth=%d\n", options.GetInt("screenwidth")) // screenwidth=90
	}
	return
}

func ExampleExtArgsOptions_SetValue() {
	var options *extargsparse.ExtArgsOptions
	var err error
	options, err = extargsparse.NewExtArgsOptions(`{}`)
	if err == nil {
		options.SetValue("screenwidth", float64(100.0))
		fmt.Fprintf(os.Stdout, "screenwidth=%d\n", options.GetInt("screenwidth")) //screenwidth=100
	}
	return
}

func ExampleNewExtArgsParse() {
	var parser *extargsparse.ExtArgsParse
	var err error
	var loads = `{}`
	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err == nil {
		parser.LoadCommandLineString(loads)
		parser.ParseCommandLine([]string{"-h"}, nil)
		/*
			Output:
			cmd 0.0.1  [OPTIONS] [args...]'

			[OPTIONS]
			    --json     json  json input file to get the value set
			    --help|-h        to display this help information
		*/
	}
}

func ExampleNewExtArgsParse_withnohelp() {
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var confstr = `{"nohelpoption" : true}`
	var err error
	var loads = `{}`
	options, err = extargsparse.NewExtArgsOptions(confstr)
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			parser.LoadCommandLineString(loads)
			// simplest the parser without help option
		}
	}
	return
}

func ExampleNewExtArgsParse_priority() {
	var err error
	var loads = `        {
            "verbose|v" : "+",
            "$port|p" : {
                "value" : 3000,
                "type" : "int",
                "nargs" : 1 ,
                "helpinfo" : "port to connect"
            },
            "dep" : {
                "list|l" : [],
                "string|s" : "s_var",
                "$" : "+"
            }
        }`
	var confstr = fmt.Sprintf(`        {
            "nojsonoption" : true,
            "nohelpoption" : true
        }`)
	var options *extargsparse.ExtArgsOptions
	var parser *extargsparse.ExtArgsParse
	var args *extargsparse.NameSpaceEx
	var jsonfile string
	var depjsonfile string
	var depstrval string = `newval`
	var depliststr string = `["depenv1","depenv2"]`
	var f *os.File

	f, _ = ioutil.TempFile("", "jsonfile")
	jsonfile = f.Name()
	f.Close()
	ioutil.WriteFile(jsonfile, []byte(`{"dep":{"list" : ["jsonval1","jsonval2"],"string" : "jsonstring"},"port":6000,"verbose":3}`), 0600)
	defer os.RemoveAll(jsonfile)

	f, _ = ioutil.TempFile("", "jsonfile")
	depjsonfile = f.Name()
	f.Close()
	ioutil.WriteFile(depjsonfile, []byte(`{"list":["depjson1","depjson2"]}`), 0600)
	defer os.RemoveAll(depjsonfile)

	os.Setenv("EXTARGSPARSE_JSONFILE", jsonfile)
	os.Setenv("DEP_JSONFILE", depjsonfile)

	options, err = extargsparse.NewExtArgsOptions(confstr)
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, []int{extargsparse.ENV_COMMAND_JSON_SET, extargsparse.ENVIRONMENT_SET, extargsparse.ENV_SUB_COMMAND_JSON_SET})
		if err == nil {
			err = parser.LoadCommandLineString(loads)
			if err == nil {
				os.Setenv("DEP_STRING", depstrval)
				os.Setenv("DEP_LIST", depliststr)
				args, err = parser.ParseCommandLine([]string{"-p", "9000", "dep", "--dep-string", "ee", "ww"}, nil)
				fmt.Fprintf(os.Stdout, "verbose=%d\n", args.GetInt("verbose"))          // verbose=0
				fmt.Fprintf(os.Stdout, "port=%d\n", args.GetInt("port"))                // port=9000
				fmt.Fprintf(os.Stdout, "subcommand=%s\n", args.GetString("subcommand")) // subcommand=dep
				fmt.Fprintf(os.Stdout, "dep_list=%v\n", args.GetArray("dep_list"))      //dep_list=[depenv1 depenv2]
				fmt.Fprintf(os.Stdout, "dep_string=%s\n", args.GetString("dep_string")) // dep_string=ee
				fmt.Fprintf(os.Stdout, "subnargs=%v\n", args.GetArray("subnargs"))      // subnargs=ww
			}
		}
	}
	return
}

func ExampleExtArgsParse_GetCmdKey() {
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
            "nojsonoption" : true,
            "nohelpoption" : true
        }`)
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

func ExampleExtArgsParse_GetCmdOpts() {
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
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{"prog" : "cmd1"}`))
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

func ExampleExtArgsParse_GetSubCommands() {
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
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{"prog" : "cmd1"}`))
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
