package extargsparse

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.Name() == "." || entry.Name() == ".." {
			continue
		}
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func check_equal(t *testing.T, orig, check interface{}) {
	if !reflect.DeepEqual(orig, check) {
		t.Fatalf("%s[%s] orig [%v] != check[%v]", format_out_stack(2), t.Name(), orig, check)
	}
}

func check_not_equal(t *testing.T, orig, check interface{}) {
	if reflect.DeepEqual(orig, check) {
		t.Fatalf("%s[%s] orig [%v] == check[%v]", format_out_stack(2), t.Name(), orig, check)
	}
}

func makeWriteTempFileInner(s string) (fname string, err error) {
	var f *os.File
	f, err = ioutil.TempFile("", "tmpfile")
	if err != nil {

		return "", err
	}
	defer f.Close()
	_, err = f.WriteString(s)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func makeWriteTempFile(s string) string {
	var err error
	var fname string
	for {
		fname, err = makeWriteTempFileInner(s)
		if err == nil {
			return fname
		}
	}
	return ""
}

type stringIO struct {
	IoWriter
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

func formLine(tabs int, fmtstr string, a ...interface{}) string {
	var s string
	var i int
	s = ""
	for i = 0; i < tabs; i++ {
		s += fmt.Sprintf("    ")
	}
	s += fmt.Sprintf(fmtstr, a...)
	s += "\n"
	return s
}

type compileExec struct {
	logObject
	fname     string
	dname     string
	origfname string
	exename   string
	addedvars map[string]string
	modvars   map[string]string
	outsarr   []string
	errsarr   []string
}

func newComileExec() *compileExec {
	self := &compileExec{logObject: *newLogObject("extargsparse")}
	self.fname = ""
	self.dname = ""
	self.origfname = ""
	self.exename = ""
	self.addedvars = make(map[string]string)
	self.modvars = make(map[string]string)
	self.outsarr = make([]string, 0)
	self.errsarr = make([]string, 0)
	return self
}

func (self *compileExec) resetVars() {
	var k, v string
	for k, _ = range self.addedvars {
		os.Unsetenv(k)
	}
	self.addedvars = make(map[string]string)

	for k, v = range self.modvars {
		os.Setenv(k, v)
	}
	self.modvars = make(map[string]string)
	return
}

func (self *compileExec) Release(ok bool) {
	self.resetVars()
	self.outsarr = make([]string, 0)
	self.errsarr = make([]string, 0)

	if len(self.exename) > 0 {
		if ok {
			os.Remove(self.exename)
		} else {
			self.Error("exe [%s]", self.exename)
		}
		self.exename = ""
	}

	if len(self.origfname) > 0 {
		if ok {
			os.Remove(self.origfname)
		} else {
			self.Error("origname [%s]", self.origfname)
		}
		self.origfname = ""
	}

	if len(self.fname) > 0 {
		if ok {
			os.Remove(self.fname)
		} else {
			self.Error("fname [%s]", self.fname)
		}
		self.fname = ""
	}
	if len(self.dname) > 0 {
		if ok {
			os.RemoveAll(self.dname)
		} else {
			self.Error("dname [%s]", self.dname)
		}
	}
	return
}

func (self *compileExec) getParserStruct(tabs int, parser *ExtArgsParse, conf *ExtArgsOptions, cmdname string) (string, error) {
	var subcmds []string
	var err error
	var curcmd string
	var curs string
	var curname string
	var s string
	var opts []*ExtKeyParse
	var opt *ExtKeyParse
	var typename string
	var optname string
	s = ""
	subcmds, err = parser.GetSubCommands(cmdname)
	if err != nil {
		return "", err
	}

	for _, curcmd = range subcmds {
		curname = ""
		curname += cmdname
		if len(curname) > 0 {
			curname += "."
		}
		curname += curcmd
		if conf.GetBool(VAR_UPPER_CASE) {
			s += formLine(tabs, "%s struct {", ucFirst(curcmd))
		} else {
			s += formLine(tabs, "%s struct {", curcmd)
		}

		curs, err = self.getParserStruct(tabs+1, parser, conf, curname)
		if err != nil {
			return "", err
		}
		s += curs
		s += formLine(tabs, "}")
	}

	opts, err = parser.GetCmdOpts(cmdname)
	if err != nil {
		return "", err
	}
	for _, opt = range opts {
		if opt.IsFlag() && opt.TypeName() != "help" && opt.TypeName() != "jsonfile" && opt.TypeName() != "prefix" {
			switch opt.TypeName() {
			case "list":
				typename = "[]string"
			case "string":
				typename = "string"
			case "int":
				typename = "int"
			case "float":
				typename = "float64"
			case "count":
				typename = "int"
			case "bool":
				typename = "bool"
			case "args":
				typename = "[]string"
			default:
				return "", fmt.Errorf("%s", format_error("can not find type [%s]", opt.TypeName()))
			}
			if opt.TypeName() != "args" {
				optname = opt.FlagName()
			} else {
				if len(cmdname) > 0 {
					optname = "subnargs"
				} else {
					optname = "args"
				}
			}
			if conf.GetBool(VAR_UPPER_CASE) {
				s += formLine(tabs, "%s %s", ucFirst(optname), typename)
			} else {
				s += formLine(tabs, "%s %s", optname, typename)
			}
		}
	}
	return s, nil
}

func (self *compileExec) formPriority(priority interface{}) (string, error) {
	var pr []int
	var p int
	var s string
	var i int
	if priority == nil {
		return "nil", nil
	}
	pr = priority.([]int)
	s = "[]int{"
	for i, p = range pr {
		if i > 0 {
			s += ","
		}
		switch p {
		case SUB_COMMAND_JSON_SET:
			s += "extargsparse.SUB_COMMAND_JSON_SET"
		case COMMAND_JSON_SET:
			s += "extargsparse.COMMAND_JSON_SET"
		case ENVIRONMENT_SET:
			s += "extargsparse.ENVIRONMENT_SET"
		case ENV_SUB_COMMAND_JSON_SET:
			s += "extargsparse.ENV_SUB_COMMAND_JSON_SET"
		case ENV_COMMAND_JSON_SET:
			s += "extargsparse.ENV_COMMAND_JSON_SET"
		default:
			return "", fmt.Errorf("%s", format_error("can not set [%d]", p))
		}
	}
	s += "}"
	return s, nil
}

func (self *compileExec) formNsName(tabs int, key string, nsname string, actname string) string {
	return formLine(tabs, `fmt.Fprintf(os.Stdout,"%s=%%v\n", %s.%s("%s"))`, key, nsname, actname, key)
}

func (self *compileExec) formSName(tabs int, key string, sname string) string {
	return formLine(tabs, `fmt.Fprintf(os.Stdout,"%s.%s=%%v\n", %s.%s)`, sname, key, sname, key)
}

func (self *compileExec) getStructMemberName(options *ExtArgsOptions, cmdname string, flagname string) string {
	var sarr []string
	var i int
	var s string
	if options.GetBool(VAR_UPPER_CASE) && len(cmdname) > 0 {
		sarr = strings.Split(cmdname, ".")
		for i = 0; i < len(sarr); i++ {
			sarr[i] = ucFirst(sarr[i])
		}
		cmdname = strings.Join(sarr, ".")
	}

	s = cmdname
	if len(s) > 0 {
		s += "."
	}
	if options.GetBool(VAR_UPPER_CASE) {
		s += ucFirst(flagname)
	} else {
		s += flagname
	}
	return s
}

func (self *compileExec) formPrintoutInner(tabs int, parser *ExtArgsParse, options *ExtArgsOptions, nsname string, sname string, cmdname string) (string, error) {
	var opts []*ExtKeyParse
	var opt *ExtKeyParse
	var subcmds []string
	var curcmd string
	var curname string
	var curs string
	var err error
	var s string
	var flagname string
	var actname string

	s = ""
	subcmds, err = parser.GetSubCommands(cmdname)
	if err != nil {
		return "", nil
	}

	for _, curcmd = range subcmds {
		curname = ""
		curname += cmdname
		if len(curname) > 0 {
			curname += "."
		}
		curname += curcmd
		curs, err = self.formPrintoutInner(tabs, parser, options, nsname, sname, curname)
		if err != nil {
			return "", err
		}
		s += curs
	}

	opts, err = parser.GetCmdOpts(cmdname)
	if err != nil {
		return "", err
	}
	for _, opt = range opts {
		if opt.IsFlag() && opt.TypeName() != "help" && opt.TypeName() != "jsonfile" && opt.TypeName() != "prefix" {
			if opt.TypeName() != "args" {
				flagname = self.getStructMemberName(options, cmdname, opt.FlagName())
			} else {
				if len(cmdname) > 0 {
					flagname = self.getStructMemberName(options, cmdname, "subnargs")
				} else {
					flagname = self.getStructMemberName(options, cmdname, "args")
				}
			}

			s += self.formSName(tabs, flagname, sname)
		}
	}

	for _, opt = range opts {
		if opt.IsFlag() && opt.TypeName() != "help" && opt.TypeName() != "jsonfile" && opt.TypeName() != "prefix" && opt.TypeName() != "args" {
			switch opt.TypeName() {
			case "string":
				actname = "GetString"
			case "int":
				actname = "GetInt"
			case "bool":
				actname = "GetBool"
			case "count":
				actname = "GetInt"
			case "float":
				actname = "GetFloat"
			case "list":
				actname = "GetArray"
			default:
				return "", fmt.Errorf("%s", format_error("unknown type [%s] [%s]", opt.TypeName(), opt.Format()))
			}
			s += self.formNsName(tabs, opt.Optdest(), nsname, actname)
		}
	}

	return s, nil
}

func (self *compileExec) formPrintout(tabs int, parser *ExtArgsParse, options *ExtArgsOptions, nsname string, sname string) (string, error) {
	var err error
	var s string
	var curs string

	s = ""

	s += self.formNsName(tabs, "subcommand", nsname, "GetString")
	curs, err = self.formPrintoutInner(tabs, parser, options, nsname, sname, "")
	if err != nil {
		return "", err
	}
	s += curs
	s += formLine(tabs, "if len(%s.GetString(\"subcommand\")) > 0 {", nsname)
	s += self.formNsName(tabs+1, "subnargs", nsname, "GetArray")
	s += formLine(tabs, "} else {")
	s += self.formNsName(tabs+1, "args", nsname, "GetArray")
	s += formLine(tabs, "}")
	return s, nil
}

func (self *compileExec) writeScript(options string, commandline string, priority interface{}, printout bool, nsname, sname string) (string, error) {
	var s string
	var parser *ExtArgsParse
	var err error
	var conf *ExtArgsOptions
	var curs string
	var prstr string
	s = ""
	s += formLine(0, "package main")
	s += formLine(0, "")
	s += formLine(0, "import (")
	s += formLine(1, `"fmt"`)
	s += formLine(1, `"go-extargsparse"`)
	s += formLine(1, `"os"`)
	s += formLine(0, ")")

	if printout {
		conf = nil
		if len(options) > 0 {
			conf, err = NewExtArgsOptions(options)
			if err != nil {
				return "", err
			}
		}
		parser, err = NewExtArgsParse(conf, priority)
		if err != nil {
			return "", err
		}

		err = parser.LoadCommandLineString(commandline)
		if err != nil {
			return "", err
		}
		s += formLine(0, "")
		s += formLine(0, "type CommandArgs struct {")
		curs, err = self.getParserStruct(1, parser, conf, "")
		if err != nil {
			return "", err
		}
		s += curs
		s += formLine(0, "}")
	}

	s += formLine(0, "")
	s += formLine(0, "func main() {")

	s += formLine(1, "var parser *extargsparse.ExtArgsParse")
	s += formLine(1, "var options *extargsparse.ExtArgsOptions")
	s += formLine(1, "var err error")
	s += formLine(1, "var commandline = `%s`", commandline)
	if printout {
		s += formLine(1, "var %s *extargsparse.NameSpaceEx", nsname)
		s += formLine(1, "var %s *CommandArgs", sname)
		s += formLine(1, "%s = &CommandArgs{}", sname)
	}
	s += formLine(0, "")
	prstr, err = self.formPriority(priority)
	if err != nil {
		return "", err
	}
	if len(options) > 0 {
		s += formLine(1, "options, err = extargsparse.NewExtArgsOptions(`%s`)", options)
		s += formLine(1, "if err != nil {")
		s += formLine(2, `fmt.Fprintf(os.Stderr,"can not parse [%s] error[%%s]\n", err.Error())`, options)
		s += formLine(2, "os.Exit(3)")
		s += formLine(1, "}")
	} else {
		s += formLine(1, "options = nil")
	}
	s += formLine(1, "parser ,err = extargsparse.NewExtArgsParse(options,%s)", prstr)
	s += formLine(1, "if err != nil {")
	s += formLine(2, `fmt.Fprintf(os.Stderr,"new args error [%%s]\n", err.Error())`)
	s += formLine(2, "os.Exit(3)")
	s += formLine(1, "}")
	// now we should give the coding
	s += formLine(1, "err = parser.LoadCommandLineString(commandline)")
	s += formLine(1, "if err != nil {")
	s += formLine(2, `fmt.Fprintf(os.Stderr,"load commandline[%%s] error [%%s]\n", commandline, err.Error())`)
	s += formLine(2, "os.Exit(3)")
	s += formLine(1, "}")
	if printout {
		s += formLine(1, "%s, err = parser.ParseCommandLineEx(nil,parser,%s,nil)", nsname, sname)
		s += formLine(1, "if err != nil {")
		s += formLine(2, `fmt.Fprintf(os.Stderr,"parse command error [%%s]\n", err.Error())`)
		s += formLine(2, "os.Exit(3)")
		s += formLine(1, "}")

		curs, err = self.formPrintout(1, parser, conf, nsname, sname)
		if err != nil {
			return "", err
		}
		s += curs
	} else {
		s += formLine(1, "_, err = parser.ParseCommandLine(nil,nil)")
		s += formLine(1, "if err != nil {")
		s += formLine(2, `fmt.Fprintf(os.Stderr,"parse command error [%%s]\n", err.Error())`)
		s += formLine(2, "os.Exit(3)")
		s += formLine(1, "}")
	}

	s += formLine(0, "}")

	return s, nil
}

func (self *compileExec) makeSrcDir(copyfrom string) error {
	var dname string
	var err error
	var f *os.File
	if len(self.dname) > 0 || len(self.fname) > 0 || len(self.origfname) > 0 {
		return fmt.Errorf("%s", format_error("[%s] dir not delete", self.dname))
	}
	dname, err = ioutil.TempDir("", "gobuild")
	if err != nil {
		return err
	}
	self.dname = dname
	err = os.MkdirAll(filepath.Join(self.dname, "src"), os.ModeDir|0644)
	if err != nil {
		return fmt.Errorf("%s", format_error("can not mkdir [%s] err[%s]", filepath.Join(self.dname, "src"), err.Error()))
	}
	err = copyDir(copyfrom, filepath.Join(self.dname, "src", "go-extargsparse"))
	f, err = ioutil.TempFile("", "main")
	if err != nil {
		return fmt.Errorf("%s", format_error("can not make temp err[%s]", err.Error()))
	}
	defer f.Close()
	self.origfname = f.Name()
	self.fname = self.origfname
	self.fname += ".go"
	return nil
}

func (self *compileExec) WriteScript(options string, commandline string, priority interface{}, printout bool, nsname string, sname string) error {
	var s string
	var err error
	var fdir string
	self.Release(true)

	s, err = self.writeScript(options, commandline, priority, printout, nsname, sname)
	if err != nil {
		return err
	}

	fdir = getCallerFilename(1)
	if len(fdir) == 0 {
		return fmt.Errorf("%s", format_error("can not get caller"))
	}

	if len(fdir) > 0 {
		fdir, err = filepath.Abs(filepath.Dir(fdir))
	} else {
		fdir, err = filepath.Abs(".")
	}
	if err != nil {
		return fmt.Errorf("%s", format_error("get abs error[%s]", err.Error()))
	}

	err = self.makeSrcDir(fdir)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(self.fname, []byte(s), 0644)
	if err != nil {
		return fmt.Errorf("%s", format_error("write [%s] error[%s]", self.fname, err.Error()))
	}

	self.Info("write [%s]", self.fname)
	self.Info("dname [%s]", self.dname)
	return nil
}

func (self *compileExec) setVars(setvars map[string]string) error {
	var k, v string
	var oldv string
	for k, v = range setvars {
		oldv = os.Getenv(k)
		if len(oldv) > 0 {
			self.modvars[k] = oldv
		} else {
			self.addedvars[k] = ""
		}
		os.Setenv(k, v)
	}
	return nil
}

func (self *compileExec) removeVars(keys []string) error {
	var c, oldv string
	for _, c = range keys {
		oldv = os.Getenv(c)
		if len(oldv) > 0 {
			self.modvars[c] = oldv
			os.Unsetenv(c)
		}
	}
	return nil
}

func (self *compileExec) RunCmd(setvars map[string]string, params ...string) error {
	var cmdrun *exec.Cmd
	var gobin string
	var obuf, ebuf *bytes.Buffer
	var err error
	var c string

	if len(self.exename) > 0 {
		err = fmt.Errorf("%s", "already run [%s]", self.exename)
		return err
	}

	gobin = getExecutableName("go")
	self.exename = getExecutableName(self.origfname)
	if len(self.fname) == 0 {
		return fmt.Errorf("%s", format_error("not set fname yet"))
	}

	/*now to add*/
	setvars["GOPATH"] = fmt.Sprintf("%s%c%s", self.dname, os.PathListSeparator, os.Getenv("GOPATH"))

	err = self.setVars(setvars)
	if err != nil {
		return err
	}
	defer self.resetVars()

	err = self.removeVars([]string{"EXTARGSPARSE_LOGLEVEL"})
	if err != nil {
		return err
	}

	// now we should give the export name

	// now we should set the go build
	cmdrun = exec.Command(gobin, "build", "-o", self.exename, self.fname)
	obuf = bytes.NewBufferString("")
	ebuf = bytes.NewBufferString("")
	cmdrun.Stdout = obuf
	cmdrun.Stderr = ebuf
	err = cmdrun.Run()
	if err != nil {
		err = fmt.Errorf("%s", format_error("compile [%s] error [%s] output [%s]", self.fname, err.Error(), ebuf.String()))
		return err
	}

	cmdrun = nil
	obuf = nil
	ebuf = nil
	cmdrun = exec.Command(self.exename)
	cmdrun.Args = make([]string, 0)
	cmdrun.Args = append(cmdrun.Args, self.exename)
	for _, c = range params {
		cmdrun.Args = append(cmdrun.Args, c)
	}
	obuf = bytes.NewBufferString("")
	ebuf = bytes.NewBufferString("")
	cmdrun.Stdout = obuf
	cmdrun.Stderr = ebuf
	err = cmdrun.Run()
	if err != nil {
		err = fmt.Errorf("%s", format_error("run [%s] %v error [%s] output [%s]", self.exename, params, err.Error(), ebuf.String()))
		return err
	}
	self.Trace("obuf [%s]", obuf.String())
	self.outsarr = strings.Split(obuf.String(), "\n")
	self.errsarr = strings.Split(ebuf.String(), "\n")
	err = nil
	obuf = nil
	ebuf = nil
	cmdrun = nil
	return nil
}

func (self *compileExec) GetOut() []string {
	return self.outsarr
}

func (self *compileExec) GetErr() []string {
	return self.errsarr
}
