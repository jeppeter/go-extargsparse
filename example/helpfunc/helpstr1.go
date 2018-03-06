package main

import (
	"bytes"
	"fmt"
	"github.com/jeppeter/go-extargsparse"
	"os"
)

type stringIO struct {
	extargsparse.IoWriter
	obuf *bytes.Buffer
}

func newStringIO() *stringIO {
	p := &stringIO{}
	p.obuf = bytes.NewBufferString("")
	return p
}

func (self *stringIO) Write(data []byte) (int, error) {
	return self.obuf.Write(data)
}

func (self *stringIO) WriteString(s string) (int, error) {
	return self.obuf.WriteString(s)
}

func (self *stringIO) String() string {
	return self.obuf.String()
}

func main() {
	var sio *stringIO
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
	options, err = extargsparse.NewExtArgsOptions(fmt.Sprintf(`{"%s" : "cmd1"}`, extargsparse.OPT_PROG))
	if err == nil {
		parser, err = extargsparse.NewExtArgsParse(options, nil)
		if err == nil {
			err = parser.LoadCommandLineString(fmt.Sprintf("%s", loads))
			if err == nil {
				sio = newStringIO()
				parser.PrintHelp(sio, "")
				fmt.Fprintf(os.Stdout, "main help:\n")
				fmt.Fprintf(os.Stdout, "%s", sio.String())

				sio = newStringIO()
				parser.PrintHelp(sio, "dep")
				fmt.Fprintf(os.Stdout, "dep help:\n")
				fmt.Fprintf(os.Stdout, "%s", sio.String())

				sio = newStringIO()
				parser.PrintHelp(sio, "rdep")
				fmt.Fprintf(os.Stdout, "rdep help:\n")
				fmt.Fprintf(os.Stdout, "%s", sio.String())
			}
		}
	}

}
