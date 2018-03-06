package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"io/ioutil"
	"os"
)

func main() {
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
