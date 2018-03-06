package extargsparse_test

import (
	extargsparse "."
	"fmt"
	"io/ioutil"
	"os"
)

func ExampleNewExtArgsOptions() {
	var options *extargsparse.ExtArgsOptions
	var err error
	var confstr = fmt.Sprintf(`{
		"%s" : 90.0		
		}`, extargsparse.OPT_SCREEN_WIDTH)
	options, err = extargsparse.NewExtArgsOptions(confstr)
	if err == nil {
		fmt.Fprintf(os.Stdout, "screenwidth=%d\n", options.GetInt(extargsparse.OPT_SCREEN_WIDTH)) // screenwidth=90
	}
	return
}

func ExampleExtArgsOptions_SetValue() {
	var options *extargsparse.ExtArgsOptions
	var err error
	options, err = extargsparse.NewExtArgsOptions(`{}`)
	if err == nil {
		options.SetValue(extargsparse.OPT_SCREEN_WIDTH, float64(100.0))
		fmt.Fprintf(os.Stdout, "screenwidth=%d\n", options.GetInt(extargsparse.OPT_SCREEN_WIDTH)) //screenwidth=100
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
	var confstr = fmt.Sprintf(`{"%s" : true}`, extargsparse.OPT_NO_HELP_OPTION)
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
            "%s" : true,
            "%s" : true
        }`, extargsparse.OPT_NO_JSON_OPTION, extargsparse.OPT_NO_HELP_OPTION)
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

func ExampleExtKeyParse_CmdName() {
	var loads = `{
		"dep" : {

		},
		"rdep": {

		}
	}`
	var err error
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var flag *extargsparse.ExtKeyParse
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{}`))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				flag, _ = parser.GetCmdKey("")
				fmt.Fprintf(os.Stdout, "cmdname=%s\n", flag.CmdName())
				flag, _ = parser.GetCmdKey("dep")
				fmt.Fprintf(os.Stdout, "cmdname=%s\n", flag.CmdName())
				flag, _ = parser.GetCmdKey("rdep")
				fmt.Fprintf(os.Stdout, "cmdname=%s\n", flag.CmdName())
				/*
					Output:
					cmdname=main
					cmdname=dep
					cmdname=ip
				*/

			}
		}
	}
}

func ExampleExtKeyParse_FlagName() {
	var loads = `{
		"verbose|v" : "+",
		"dep" : {
			"cc|c" : ""
		},
		"rdep": {
			"dd|C" : ""
		}
	}`
	var err error
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var flag *extargsparse.ExtKeyParse
	var opts []*extargsparse.ExtKeyParse
	var finded bool
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{}`))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				opts, _ = parser.GetCmdOpts("")
				for _, flag = range opts {
					if flag.FlagName() == "verbose" {
						fmt.Fprintf(os.Stdout, "flagname=%s\n", flag.FlagName())
					}
				}

				opts, _ = parser.GetCmdOpts("dep")
				for _, flag = range opts {
					if flag.FlagName() == "cc" {
						fmt.Fprintf(os.Stdout, "flagname=%s\n", flag.FlagName())
					}
				}
				opts, _ = parser.GetCmdOpts("rdep")
				for _, flag = range opts {
					if flag.FlagName() == "dd" {
						fmt.Fprintf(os.Stdout, "flagname=%s\n", flag.FlagName())
					}
				}

				opts, _ = parser.GetCmdOpts("rdep")
				finded = false
				for _, flag = range opts {
					if flag.FlagName() == "cc" {
						fmt.Fprintf(os.Stdout, "flagname=%s\n", flag.FlagName())
						finded = true
					}
				}

				if !finded {
					fmt.Fprintf(os.Stdout, "can not found cc for rdep cmd\n")
				}

			}
		}
	}
	/*
		Output:
		flagname=verbose
		flagname=cc
		flagname=dd
		can not found cc for rdep cmd
	*/
}

func ExampleExtKeyParse_Function() {
	var loads = `{
		"verbose|v" : "+",
		"dep<dep_handler>" : {
			"cc|c" : ""
		},
		"rdep<rdep_handler>": {
			"dd|C" : ""
		}
	}`
	var err error
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var flag *extargsparse.ExtKeyParse
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{}`))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				flag, _ = parser.GetCmdKey("")
				fmt.Fprintf(os.Stdout, "main funcion:%s\n", flag.Function())

				flag, _ = parser.GetCmdKey("dep")
				fmt.Fprintf(os.Stdout, "dep funcion:%s\n", flag.Function())

				flag, _ = parser.GetCmdKey("rdep")
				fmt.Fprintf(os.Stdout, "rdep funcion:%s\n", flag.Function())
			}
		}
	}
	/*
		Notice:
			if options not set fmt.Sprintf(`{"%s" : false}`,extargsparse.OPT_FUNC_UPPER_CASE)
			the real function is Dep_handler and Rdep_handler
		Output:
			main funcion:
			dep funcion:dep_handler
			rdep funcion:rdep_handler
	*/
}

func ExampleExtKeyParse_HelpInfo() {
	var loads = `{
		"verbose|v##we used verbose##" : "+",
		"$##this is args help##" : "*",
		"dep##dep help set##" : {
			"cc|c##cc sss##" : "",
			"$##this is dep subnargs help##" : "*"
		},
		"rdep##can not set rdep help##": {
			"dd|C##capital C##" : "",
			"$##this is rdep subnargs help##" : "*"
		}
	}`
	var err error
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var flag *extargsparse.ExtKeyParse
	var opts []*extargsparse.ExtKeyParse
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{}`))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				flag, _ = parser.GetCmdKey("")
				fmt.Fprintf(os.Stdout, "main helpinfo:%s\n", flag.HelpInfo())

				flag, _ = parser.GetCmdKey("dep")
				fmt.Fprintf(os.Stdout, "dep helpinfo:%s\n", flag.HelpInfo())

				flag, _ = parser.GetCmdKey("rdep")
				fmt.Fprintf(os.Stdout, "rdep helpinfo:%s\n", flag.HelpInfo())

				opts, _ = parser.GetCmdOpts("")
				for _, flag = range opts {
					if flag.TypeName() == "args" {
						fmt.Fprintf(os.Stdout, "main.args.HelpInfo=%s\n", flag.HelpInfo())
					} else if flag.FlagName() == "verbose" {
						fmt.Fprintf(os.Stdout, "verbose.HelpInfo=%s\n", flag.HelpInfo())
					}
				}
				opts, _ = parser.GetCmdOpts("dep")
				for _, flag = range opts {
					if flag.TypeName() == "args" {
						fmt.Fprintf(os.Stdout, "dep.subnargs.HelpInfo=%s\n", flag.HelpInfo())
					} else if flag.FlagName() == "cc" {
						fmt.Fprintf(os.Stdout, "dep.cc.HelpInfo=%s\n", flag.HelpInfo())
					}
				}
				opts, _ = parser.GetCmdOpts("rdep")
				for _, flag = range opts {
					if flag.TypeName() == "args" {
						fmt.Fprintf(os.Stdout, "rdep.subnargs.HelpInfo=%s\n", flag.HelpInfo())
					} else if flag.FlagName() == "dd" {
						fmt.Fprintf(os.Stdout, "rdep.dd.HelpInfo=%s\n", flag.HelpInfo())
					}
				}

			}
		}
	}
	/*
		Notice:
			HelpInfo is part of between ##(information)##
		Output:
			main helpinfo:
			dep helpinfo:dep help set
			rdep helpinfo:can not set rdep help
			main.args.HelpInfo=this is args help
			verbose.HelpInfo=we used verbose
			dep.subnargs.HelpInfo=this is dep subnargs help
			dep.cc.HelpInfo=cc sss
			rdep.subnargs.HelpInfo=this is rdep subnargs help
			rdep.dd.HelpInfo=capital C
	*/
}

func ExampleExtKeyParse_Longopt() {
	var loads = `{
		"verbose|v" : "+",
		"dep<dep_handler>" : {
			"cc|c" : ""
		},
		"rdep<rdep_handler>": {
			"dd|C" : ""
		}
	}`
	var confstr string
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var flag *extargsparse.ExtKeyParse
	var opts []*extargsparse.ExtKeyParse
	confstr = `{}`
	options, _ = extargsparse.NewExtArgsOptions(confstr)
	parser, _ = extargsparse.NewExtArgsParse(options, nil)
	parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	opts, _ = parser.GetCmdOpts("")
	for _, flag = range opts {
		if flag.TypeName() == "count" && flag.FlagName() == "verbose" {
			fmt.Fprintf(os.Stdout, "longprefix=%s\n", flag.LongPrefix())
			fmt.Fprintf(os.Stdout, "longopt=%s\n", flag.Longopt())
			fmt.Fprintf(os.Stdout, "shortopt=%s\n", flag.Shortopt())
		}
	}

	confstr = fmt.Sprintf(`{ "%s" : "++", "%s" : "+"}`, extargsparse.OPT_LONG_PREFIX, extargsparse.OPT_SHORT_PREFIX)
	options, _ = extargsparse.NewExtArgsOptions(confstr)
	parser, _ = extargsparse.NewExtArgsParse(options, nil)
	parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	opts, _ = parser.GetCmdOpts("")
	for _, flag = range opts {
		if flag.TypeName() == "count" && flag.FlagName() == "verbose" {
			fmt.Fprintf(os.Stdout, "longprefix=%s\n", flag.LongPrefix())
			fmt.Fprintf(os.Stdout, "longopt=%s\n", flag.Longopt())
			fmt.Fprintf(os.Stdout, "shortopt=%s\n", flag.Shortopt())
		}
	}

	/*
		Output:
			longprefix=--
			longopt=--verbose
			shortopt=-v
			longprefix=++
			longopt=++verbose
			shortopt=+v
	*/
}

func ExampleExtKeyParse_Nargs() {
	var loads = `{
		"verbose|v" : "+",
		"dep<dep_handler>" : {
			"cc|c" : "",
			"$" : "+"
		},
		"rdep<rdep_handler>": {
			"dd|C" : "",
			"$" : "?"
		},
		"$port" : {
			"nargs" : 1,
			"type" : "int",
			"value" : 9000
		}
	}`
	var confstr string
	var parser *extargsparse.ExtArgsParse
	var options *extargsparse.ExtArgsOptions
	var flag *extargsparse.ExtKeyParse
	var opts []*extargsparse.ExtKeyParse
	confstr = `{}`
	options, _ = extargsparse.NewExtArgsOptions(confstr)
	parser, _ = extargsparse.NewExtArgsParse(options, nil)
	parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	opts, _ = parser.GetCmdOpts("")
	for _, flag = range opts {
		if flag.TypeName() == "args" {
			fmt.Fprintf(os.Stdout, "args.nargs=%v\n", flag.Nargs())
		} else if flag.FlagName() == "port" {
			fmt.Fprintf(os.Stdout, "port.nargs=%d\n", flag.Nargs().(int))
		}
	}

	opts, _ = parser.GetCmdOpts("dep")
	for _, flag = range opts {
		if flag.TypeName() == "args" {
			fmt.Fprintf(os.Stdout, "dep.args.nargs=%v\n", flag.Nargs())
		}
	}
	opts, _ = parser.GetCmdOpts("rdep")
	for _, flag = range opts {
		if flag.TypeName() == "args" {
			fmt.Fprintf(os.Stdout, "rdep.args.nargs=%v\n", flag.Nargs())
		}
	}
	/*
		Output:
			args.nargs=*
			port.nargs=1
			dep.args.nargs=+
			rdep.args.nargs=?
	*/
}
