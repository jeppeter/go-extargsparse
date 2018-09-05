# go-extargsparse 
> golang command line handle package inspired by [extargsparse](https://github.com/jeppeter/extargsparse)

### Release History
> Sep 5th 2018 Release 0.0.4 for fixup the usage format
> March 6th 2018 Release 0.0.2 for the first release

### simple example
```go
package main

import (
    "fmt"
    "github.com/jeppeter/go-extargsparse"
    "os"
)

func main() {
    var commandline = `{
        "verbose|v" : "+",
        "removed|R" : false,
        "floatv|f" : 3.3,
        "intv|i" : 5,
        "arrl|a" : [],
        "strv|s" : null,
        "$" : "+"
        }`
    var parser *extargsparse.ExtArgsParse
    var ns *extargsparse.NameSpaceEx
    var err error

    parser, err = extargsparse.NewExtArgsParse(nil, nil)
    if err != nil {
        fmt.Fprintf(os.Stderr, "can not init parser err [%s]\n", err.Error())
        os.Exit(5)
        return
    }

    err = parser.LoadCommandLineString(commandline)
    if err != nil {
        fmt.Fprintf(os.Stderr, "parse [%s] error[%s]\n", commandline, err.Error())
        os.Exit(5)
        return
    }

    if len(os.Args[1:]) == 0 {
        ns, err = parser.ParseCommandLine([]string{"-vvvv", "cc", "-f", "33.2", "--arrl", "wwwe", "-s", "3993"}, nil)
    } else {
        ns, err = parser.ParseCommandLine(nil, nil)
    }

    if err != nil {
        fmt.Fprintf(os.Stderr, "can not parser command err[%s]\n", err.Error())
        os.Exit(5)
        return
    }
    fmt.Fprintf(os.Stdout, "verbose=%d\n", ns.GetInt("verbose"))
    fmt.Fprintf(os.Stdout, "removed=%v\n", ns.GetBool("removed"))
    fmt.Fprintf(os.Stdout, "falotv=%f\n", ns.GetFloat("floatv"))
    fmt.Fprintf(os.Stdout, "intv=%d\n", ns.GetInt("intv"))
    fmt.Fprintf(os.Stdout, "arrl=%v\n", ns.GetArray("arrl"))
    fmt.Fprintf(os.Stdout, "strv=%s\n", ns.GetString("strv"))
    fmt.Fprintf(os.Stdout, "args=%v\n", ns.GetArray("args"))
    return
}
```

> if you run command 
```shell
./simple1 -R -i 30 ss ww
```

> output 
```shell
verbose=0
removed=true
falotv=3.300000
intv=30
arrl=[]
strv=
args=[ss ww]
```

### a little more complex example
```go
package main

import (
    "fmt"
    "github.com/jeppeter/go-extargsparse"
    "os"
)

type ArgStruct struct {
    Verbose int
    Removed bool
    Floatv  float64
    Intv    int
    Arrl    []string
    Strv    string
    Args    []string
}

func main() {
    var commandline = `{
        "verbose|v" : "+",
        "removed|R" : false,
        "floatv|f" : 3.3,
        "intv|i" : 5,
        "arrl|a" : [],
        "strv|s" : null,
        "$" : "+"
        }`
    var parser *extargsparse.ExtArgsParse
    var p *ArgStruct
    var err error

    p = &ArgStruct{}

    parser, err = extargsparse.NewExtArgsParse(nil, nil)
    if err != nil {
        fmt.Fprintf(os.Stderr, "can not init parser err [%s]\n", err.Error())
        os.Exit(5)
        return
    }

    err = parser.LoadCommandLineString(commandline)
    if err != nil {
        fmt.Fprintf(os.Stderr, "parse [%s] error[%s]\n", commandline, err.Error())
        os.Exit(5)
        return
    }

    if len(os.Args[1:]) == 0 {
        _, err = parser.ParseCommandLineEx([]string{"-vvvv", "cc", "-f", "33.2", "--arrl", "wwwe", "-s", "3993"}, nil, p, nil)
    } else {
        _, err = parser.ParseCommandLineEx(nil, nil, p, nil)
    }

    if err != nil {
        fmt.Fprintf(os.Stderr, "can not parser command err[%s]\n", err.Error())
        os.Exit(5)
        return
    }
    fmt.Fprintf(os.Stdout, "verbose=%d\n", p.Verbose)
    fmt.Fprintf(os.Stdout, "removed=%v\n", p.Removed)
    fmt.Fprintf(os.Stdout, "falotv=%f\n", p.Floatv)
    fmt.Fprintf(os.Stdout, "intv=%d\n", p.Intv)
    fmt.Fprintf(os.Stdout, "arrl=%v\n", p.Arrl)
    fmt.Fprintf(os.Stdout, "strv=%s\n", p.Strv)
    fmt.Fprintf(os.Stdout, "args=%v\n", p.Args)
    return
}
```

> run 
```shell
./simple2 -R -i 30 ss ww
```

> output 
```shell
verbose=0
removed=true
falotv=3.300000
intv=30
arrl=[]
strv=
args=[ss ww]
```

### call subcommand with handle
```go
package main

import (
    "fmt"
    "github.com/jeppeter/go-extargsparse"
    "os"
)

func dep_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
    if ns == nil {
        return nil
    }
    fmt.Fprintf(os.Stdout, "subcommand=%s\n", ns.GetString("subcommand"))
    fmt.Fprintf(os.Stdout, "verbose=%d\n", ns.GetInt("verbose"))
    fmt.Fprintf(os.Stdout, "dep_list=%v\n", ns.GetArray("dep_list"))
    fmt.Fprintf(os.Stdout, "dep_str=%s\n", ns.GetString("dep_str"))
    fmt.Fprintf(os.Stdout, "subnargs=%v\n", ns.GetArray("subnargs"))
    fmt.Fprintf(os.Stdout, "rdep_list=%v\n", ns.GetArray("rdep_list"))
    fmt.Fprintf(os.Stdout, "rdep_str=%s\n", ns.GetString("rdep_str"))
    os.Exit(0)
    return nil
}

func rdep_handler(ns *extargsparse.NameSpaceEx, ostruct interface{}, ctx interface{}) error {
    if ns == nil {
        return nil
    }
    fmt.Fprintf(os.Stdout, "subcommand=%s\n", ns.GetString("subcommand"))
    fmt.Fprintf(os.Stdout, "verbose=%d\n", ns.GetInt("verbose"))
    fmt.Fprintf(os.Stdout, "dep_list=%v\n", ns.GetArray("dep_list"))
    fmt.Fprintf(os.Stdout, "dep_str=%s\n", ns.GetString("dep_str"))
    fmt.Fprintf(os.Stdout, "subnargs=%v\n", ns.GetArray("subnargs"))
    fmt.Fprintf(os.Stdout, "rdep_list=%v\n", ns.GetArray("rdep_list"))
    fmt.Fprintf(os.Stdout, "rdep_str=%s\n", ns.GetString("rdep_str"))
    os.Exit(0)
    return nil
}

func init() {
    dep_handler(nil, nil, nil)
    rdep_handler(nil, nil, nil)
}

func main() {
    var commandline = `{
        "verbose|v" : "+",
        "dep<dep_handler>" : {
            "$" : "*",
            "list|L" :  [],
            "str|S" : ""
        },
        "rdep<rdep_handler>" : {
            "$" : "*",
            "list|l" : [],
            "str|s" : ""
        }
        }`
    var parser *extargsparse.ExtArgsParse
    var err error
    var options *extargsparse.ExtArgsOptions
    var confstr = fmt.Sprintf(`{ "%s" : false}`, extargsparse.OPT_FUNC_UPPER_CASE)
    options, err = extargsparse.NewExtArgsOptions(confstr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "can not make string [%s] err[%s]\n", confstr, err.Error())
        os.Exit(5)
        return
    }
    parser, err = extargsparse.NewExtArgsParse(options, nil)
    if err != nil {
        fmt.Fprintf(os.Stderr, "can not make parser err[%s]\n", err.Error())
        os.Exit(5)
        return
    }

    err = parser.LoadCommandLineString(commandline)
    if err != nil {
        fmt.Fprintf(os.Stderr, "can not load string [%s] err[%s]\n", commandline, err.Error())
        os.Exit(5)
        return
    }

    _, err = parser.ParseCommandLineEx(nil, nil, nil, nil)
    if err != nil {
        fmt.Fprintf(os.Stderr, "can not parse err[%s]\n", err.Error())
        os.Exit(5)
        return
    }
    return
}
```

> run
```shell
./subcmd1 -vvvv rdep --rdep-list css wwwwe
```

> output
```shell
subcommand=rdep
verbose=4
dep_list=[]
dep_str=
subnargs=[wwwwe]
rdep_list=[css]
rdep_str=
```

### [more documentation](https://godoc.org/github.com/jeppeter/go-extargsparse)