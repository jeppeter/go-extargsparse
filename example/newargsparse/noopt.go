package main

import (
	"github.com/jeppeter/go-extargsparse"
)

func main() {
	var parser *extargsparse.ExtArgsParse
	var err error
	var loads = `{}`
	parser, err = extargsparse.NewExtArgsParse(nil, nil)
	if err == nil {
		parser.LoadCommandLineString(loads)
		parser.ParseCommandLine([]string{"-h"}, nil)
		/*
			Output:
		*/
	}
}
