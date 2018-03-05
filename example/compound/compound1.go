package main

import (
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"
)

func beforeParser() {
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
				strings.HasPrefix(k, "HTTP_") ||
				strings.HasPrefix(k, "SSL_") ||
				strings.HasPrefix(k, "EXTARGSPARSE_JSON") ||
				strings.HasPrefix(k, "EXTARGSPARSE_JSONFILE") {
				err = os.Unsetenv(sarr[0])
				if err == nil {
					delone = true
					break
				}
			}
		}
	}
	return
}

func makeWriteTempFile(s string) string {
	var fname string
	var f *os.File
	var err error

	err = fmt.Errorf("new error")
	for err != nil {
		f, err = ioutil.TempFile("", "newfile")
	}

	fname = f.Name()
	f.Close()

	err = fmt.Errorf("new error")
	for err != nil {
		err = ioutil.WriteFile(fname, []byte(s), 0600)
	}
	return fname
}

func safeRemoveFile(fname string) {
	if len(fname) > 0 {
		os.RemoveAll(fname)
	}
	return
}

func format_out_stack(level int) string {
	_, f, l, _ := runtime.Caller(level)
	return fmt.Sprintf("[%s:%d]", f, l)
}

func check_equal(orig, check interface{}) {
	var s string
	if !reflect.DeepEqual(orig, check) {
		s = fmt.Sprintf("%s orig [%v] != check[%v]", format_out_stack(2), orig, check)
		panic(s)
	}
}

func main() {
	var err error
	var parser *extargsparse.ExtArgsParse
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
	var depjson string = ""
	var rdepjson string = ""
	var rdepipjson string = ""
	var jsonfile string = ""
	var params []string
	var args *extargsparse.NameSpaceEx
	beforeParser()
	depjson = makeWriteTempFile(`{"float3":33.221}`)
	defer safeRemoveFile(depjson)
	rdepjson = makeWriteTempFile(`{"ip" : { "float4" : 40.3}}`)
	defer safeRemoveFile(rdepjson)
	jsonfile = makeWriteTempFile(`{"verbose": 30,"float3": 77.1}`)
	defer safeRemoveFile(jsonfile)
	rdepipjson = makeWriteTempFile(`{"float7" : 11.22,"float4" : 779.2}`)
	defer safeRemoveFile(rdepipjson)
	os.Setenv("EXTARGSPARSE_JSON", jsonfile)
	os.Setenv("DEP_JSON", depjson)
	os.Setenv("RDEP_JSON", rdepjson)
	os.Setenv("DEP_FLOAT3", fmt.Sprintf("33.52"))
	os.Setenv("RDEP_IP_FLOAT7", fmt.Sprintf("99.3"))
	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	check_equal(err, nil)
	err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
	check_equal(err, nil)
	params = []string{"-vvfvv", "33.21", "rdep", "ip", "--json", jsonfile, "--rdep-ip-json", rdepipjson}
	args, err = parser.ParseCommandLine(params, nil)
	check_equal(err, nil)
	check_equal(args.GetArray("subnargs"), []string{})
	check_equal(args.GetString("subcommand"), "rdep.ip")
	check_equal(args.GetInt("verbose"), 4)
	check_equal(args.GetFloat("float1"), 33.21)
	check_equal(args.GetFloat("dep_float3"), 33.52)
	check_equal(args.GetFloat("float2"), 6422.22)
	check_equal(args.GetFloat("float3"), 77.1)
	check_equal(args.GetFloat("rdep_ip_float4"), 779.2)
	check_equal(args.GetFloat("rdep_ip_float6"), 33.22)
	check_equal(args.GetFloat("rdep_ip_float7"), 11.22)
	return
}
