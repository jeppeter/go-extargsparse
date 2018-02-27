package extargsparse

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

var parser_reserver_args = []string{"subcommand", "subnargs", "nargs", "extargs", "args"}
var parser_priority_args = []int{SUB_COMMAND_JSON_SET, COMMAND_JSON_SET, ENVIRONMENT_SET, ENV_SUB_COMMAND_JSON_SET, ENV_COMMAND_JSON_SET}

type ExtArgsParse struct {
	logObject
	options             *ExtArgsOptions
	mainCmd             *parserCompat
	errorHandler        string
	helpHandler         string
	outputMode          []string
	ended               int
	longPrefix          string
	shortPrefix         string
	noHelpOption        bool
	noJsonOption        bool
	helpLong            string
	helpShort           string
	jsonLong            string
	cmdPrefixAdded      bool
	loadPriority        []int
	loadCommandMap      map[string]reflect.Value
	optParseHandleMap   map[string]reflect.Value
	parsePrioritySetMap map[int]reflect.Value
	setJsonValueMap     map[string]reflect.Value
}

func is_valid_priority(k int) bool {
	for _, i := range parser_priority_args {
		if k == i {
			return true
		}
	}
	return false
}

func (self *ExtArgsParse) bindNameFunction(m map[string]reflect.Value, name string, fn interface{}) map[string]reflect.Value {
	var v reflect.Value
	v = reflect.ValueOf(fn)
	m[name] = v
	return m
}

func (self *ExtArgsParse) bindIntFunction(m map[int]reflect.Value, iv int, fn interface{}) map[int]reflect.Value {
	var v reflect.Value
	v = reflect.ValueOf(fn)
	m[iv] = v
	return m
}

func (self *ExtArgsParse) bindLoadCommandMap(name string, fn interface{}) {
	self.loadCommandMap = self.bindNameFunction(self.loadCommandMap, name, fn)
	return
}

func (self *ExtArgsParse) bindOptParseHandleMap(name string, fn interface{}) {
	self.optParseHandleMap = self.bindNameFunction(self.optParseHandleMap, name, fn)
	return
}

func (self *ExtArgsParse) bindSetJsonValueMap(name string, fn interface{}) {
	self.setJsonValueMap = self.bindNameFunction(self.setJsonValueMap, name, fn)
	return
}

func (self *ExtArgsParse) bindParsePrioritySetMap(iv int, fn interface{}) {
	self.parsePrioritySetMap = self.bindIntFunction(self.parsePrioritySetMap, iv, fn)
	return
}

func (self *ExtArgsParse) loadCommandLineBase(prefix string, keycls *ExtKeyParse, parsers []*parserCompat) error {
	if keycls.IsFlag() && keycls.FlagName() != "$" && check_in_array(parser_reserver_args, keycls.FlagName()) {
		return fmt.Errorf("%s", format_error("%s in the %v", keycls.FlagName(), parser_reserver_args))
	}
	return self.checkFlagInsert(keycls, parsers)
}

func (self *ExtArgsParse) loadCommandLineArgs(prefix string, keycls *ExtKeyParse, parsers []*parserCompat) error {
	return self.checkFlagInsert(keycls, parsers)
}

func (self *ExtArgsParse) formatCmdNamePath(parsers []*parserCompat) string {
	var cmdname string
	cmdname = ""
	for _, c := range parsers {
		if len(cmdname) > 0 {
			cmdname += "."
		}
		cmdname += c.CmdName
	}
	return cmdname
}

func (self *ExtArgsParse) findSubparserInner(name string, parentcmd *parserCompat) *parserCompat {
	var sarr []string
	var findcmd *parserCompat
	if parentcmd == nil {
		parentcmd = self.mainCmd
	}
	if len(name) == 0 {
		return parentcmd
	}

	sarr = strings.Split(name, ".")
	for _, c := range parentcmd.SubCommands {
		if c.CmdName == sarr[0] {
			findcmd = self.findSubparserInner(strings.Join(sarr[1:], "."), c)
			if findcmd != nil {
				return findcmd
			}
		}
	}
	return nil
}

func (self *ExtArgsParse) getSubparserInner(keycls *ExtKeyParse, parsers []*parserCompat) *parserCompat {
	var cmdname string
	var parentname string
	var cmdparser *parserCompat
	var curparser *parserCompat
	cmdname = ""
	parentname = self.formatCmdNamePath(parsers)
	cmdname += parentname
	if len(cmdname) > 0 {
		cmdname += "."
	}
	cmdname += keycls.CmdName()
	cmdparser = self.findSubparserInner(cmdname, nil)
	if cmdparser != nil {
		return cmdparser
	}
	cmdparser = newParserCompat(keycls, self.options)
	if len(parentname) > 0 {
		curparser = self.mainCmd
	} else {
		curparser = parsers[len(parsers)-1]
	}
	curparser.SubCommands = append(curparser.SubCommands, cmdparser)
	return cmdparser
}

func (self *ExtArgsParse) loadCommandSubparser(prefix string, keycls *ExtKeyParse, parsers []*parserCompat) error {
	var parser *parserCompat
	var nextparser []*parserCompat
	var newprefix string
	var err error
	var vmap map[string]interface{}
	if keycls.TypeName() != "dict" {
		return fmt.Errorf("%s", format_error("%s not valid dict", keycls.Format()))
	}
	vmap = keycls.Value().(map[string]interface{})
	if keycls.IsCmd() && check_in_array(parser_reserver_args, keycls.CmdName()) {
		return fmt.Errorf("%s", format_error("cmdname [%s] in [%v] reserved", keycls.CmdName(), parser_reserver_args))
	}
	parser = self.getSubparserInner(keycls, parsers)
	if parser == nil {
		return fmt.Errorf("%s", format_error("can not find [%s] ", keycls.Format()))
	}
	nextparser = make([]*parserCompat, 0)
	nextparser = append(nextparser, self.mainCmd)
	if len(parsers) > 0 {
		nextparser = parsers
	}
	nextparser = append(nextparser, parser)
	if self.cmdPrefixAdded {
		newprefix = prefix
		if len(newprefix) > 0 {
			newprefix += "_"
		}
		newprefix += keycls.CmdName()
	} else {
		newprefix = ""
	}
	err = self.loadCommandLineInner(newprefix, vmap, nextparser)
	nextparser = nextparser[:(len(nextparser) - 2)]
	return err
}

func (self *ExtArgsParse) loadCommandPrefix(prefix string, keycls *ExtKeyParse, parsers []*parserCompat) error {
	var vmap map[string]interface{}
	if len(prefix) > 0 && check_in_array(parser_reserver_args, prefix) {
		return fmt.Errorf("%s", format_error("prefix [%s] in [%v]", prefix, parser_reserver_args))
	}
	vmap = keycls.Value().(map[string]interface{})
	return self.loadCommandLineInner(prefix, vmap, parsers)
}

func (self *ExtArgsParse) stringAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	if validx >= len(params) {
		err = fmt.Errorf("%s", format_error("need args [%d] [%s] [%v]", validx, keycls.Format(), params))
		return 1, err
	}
	ns.SetValue(keycls.Optdest(), params[validx])
	return 1, nil
}

func (self *ExtArgsParse) boolAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	var b bool = false
	if keycls.Value() != nil {
		b = keycls.Value().(bool)
	}
	if b {
		ns.SetValue(keycls.Optdest(), false)
	} else {
		ns.SetValue(keycls.Optdest(), true)
	}
	return 0, nil
}

func (self *ExtArgsParse) intAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	var base int = 10
	var s string
	var i int64
	if validx >= len(params) {
		err = fmt.Errorf("%s", format_error("need args [%d] [%s] [%v]", validx, keycls.Format(), params))
		return 1, err
	}

	s = params[validx]
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
		base = 16
	} else if strings.HasPrefix(s, "x") || strings.HasPrefix(s, "X") {
		s = s[1:]
		base = 16
	}

	i, err = strconv.ParseInt(s, base, 64)
	if err != nil {
		err = fmt.Errorf("%s", format_error("parse [%s] error [%s]", params[validx], err.Error()))
		return 1, err
	}
	ns.SetValue(keycls.Optdest(), int(i))
	return 1, nil
}

func (self *ExtArgsParse) appendAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	var sarr []string
	if validx >= len(params) {
		err = fmt.Errorf("%s", format_error("need args [%d] [%s] [%v]", validx, keycls.Format(), params))
		return 1, err
	}
	sarr = ns.GetArray(keycls.Optdest())
	sarr = append(sarr, params[validx])
	ns.SetValue(keycls.Optdest(), sarr)
	return 1, nil
}

func (self *ExtArgsParse) printHelp(parsers []*parserCompat) string {
	var curcmd *parserCompat
	var cmdpaths []*parserCompat
	var i int
	if self.helpHandler == "nohelp" {
		return "no help information"
	}
	curcmd = self.mainCmd
	cmdpaths = make([]*parserCompat, 0)
	if len(parsers) > 0 {
		curcmd = parsers[len(parsers)-1]
		for i = 0; i < (len(parsers) - 1); i++ {
			cmdpaths = append(cmdpaths, parsers[i])
		}
	}
	return curcmd.GetHelpInfo(nil, cmdpaths)
}

func (self *ExtArgsParse) setCommandLineSelfArgsInner(paths []*parserCompat) error {
	var parentpaths []*parserCompat
	var curpaths []*parserCompat
	var err error
	var setted bool
	var cmdname string
	var prefix string
	parentpaths = make([]*parserCompat, 0)
	parentpaths = append(parentpaths, self.mainCmd)
	if len(paths) > 0 {
		parentpaths = paths
	}

	setted = false
	for _, opt := range parentpaths[len(parentpaths)-1].CmdOpts {
		if opt.IsFlag() && opt.FlagName() == "$" {
			setted = true
			break
		}
	}

	if !setted {
		cmdname = self.formatCmdFromCmdArray(parentpaths)
		prefix = strings.Replace(cmdname, ".", "_", -1)
		curkey, err := newExtKeyParse_long("", "$", "*", true, false, false, self.longPrefix, self.shortPrefix, self.options.GetBool("flagnochange"))
		if err != nil {
			return err
		}
		err = self.loadCommandLineArgs(prefix, curkey, parentpaths)
		if err != nil {
			return err
		}
	}

	for _, chld := range parentpaths[len(parentpaths)-1].SubCommands {
		curpaths = parentpaths
		curpaths = append(curpaths, chld)
		err = self.setCommandLineSelfArgsInner(curpaths)
		if err != nil {
			return err
		}
		curpaths = curpaths[:(len(curpaths) - 1)]
	}

	return nil
}

func (self *ExtArgsParse) checkVarNameInner(paths []*parserCompat, optchk *optCheck) error {
	var parentpaths []*parserCompat
	var curpaths []*parserCompat
	var opt *ExtKeyParse
	var c *parserCompat
	var copychk *optCheck
	var bval bool
	var err error

	if optchk == nil {
		optchk = newOptCheck()
	}
	parentpaths = make([]*parserCompat, 0)
	parentpaths = append(parentpaths, self.mainCmd)
	if len(paths) > 0 {
		parentpaths = paths
	}

	for _, opt = range parentpaths[len(parentpaths)-1].CmdOpts {
		if opt.IsFlag() {
			if opt.TypeName() == "help" || opt.TypeName() == "args" {
				continue
			}
			bval = optchk.AddAndCheck("varname", opt.VarName())
			if !bval {
				return fmt.Errorf("%s", format_error("opt varname[%s] is already", opt.VarName()))
			}

			bval = optchk.AddAndCheck("longopt", opt.Longopt())
			if !bval {
				return fmt.Errorf("%s", format_error("opt longopt[%s] is already", opt.Longopt()))
			}

			if len(opt.Shortopt()) > 0 {
				bval = optchk.AddAndCheck("shortopt", opt.Shortopt())
				if !bval {
					return fmt.Errorf("%s", format_error("opt shortopt[%s] is already", opt.Shortopt()))
				}
			}
		}
	}

	for _, c = range parentpaths[len(parentpaths)-1].SubCommands {
		curpaths = parentpaths
		curpaths = append(curpaths, c)
		copychk = newOptCheck()
		copychk.Copy(optchk)
		err = self.checkVarNameInner(curpaths, copychk)
		if err != nil {
			return err
		}
		curpaths = curpaths[:(len(curpaths) - 1)]
	}

	return nil
}

func (self *ExtArgsParse) setCommandLineSelfArgs() error {
	var paths []*parserCompat
	var err error
	if self.ended != 0 {
		return nil
	}
	paths = make([]*parserCompat, 0)
	err = self.setCommandLineSelfArgsInner(paths)
	if err != nil {
		return err
	}

	err = self.checkVarNameInner(paths, nil)
	if err != nil {
		return err
	}

	return nil
}

func (self *ExtArgsParse) findCommandInner(cmdname string, parsers []*parserCompat) *parserCompat {
	var sarr []string
	var curroot *parserCompat
	var nextparsers []*parserCompat
	sarr = strings.Split(cmdname, ".")
	curroot = self.mainCmd
	nextparsers = make([]*parserCompat, 0)
	if len(parsers) > 0 {
		nextparsers = parsers
		curroot = nextparsers[len(nextparsers)-1]
	}

	if len(sarr) > 1 {
		nextparsers = append(nextparsers, curroot)
		for _, c := range curroot.SubCommands {
			if c.CmdName == sarr[0] {
				nextparsers = make([]*parserCompat, 0)
				if len(parsers) > 0 {
					nextparsers = parsers
				}
				nextparsers = append(nextparsers, c)
				return self.findCommandInner(strings.Join(sarr[1:], "."), nextparsers)
			}
		}
	} else {
		for _, c := range curroot.SubCommands {
			if c.CmdName == sarr[0] {
				return c
			}
		}
	}
	return nil
}

func (self *ExtArgsParse) findCommandsInPath(cmdname string, parsers []*parserCompat) []*parserCompat {
	var commands []*parserCompat
	var i int
	var sarr []string
	var curcommand *parserCompat
	commands = make([]*parserCompat, 0)
	sarr = []string{""}
	if len(cmdname) > 0 {
		sarr = strings.Split(cmdname, ".")
	}
	if self.mainCmd != nil {
		commands = append(commands, self.mainCmd)
	}
	for i = 0; i < len(sarr) && len(cmdname) > 0; i++ {
		if i > 0 {
			curcommand = self.findCommandInner(sarr[i-1], commands)
			if curcommand == nil {
				break
			}
			commands = append(commands, curcommand)
		}
	}
	return commands
}

func (self *ExtArgsParse) PrintHelp(out *os.File, cmdname string) error {
	var err error
	var parsers []*parserCompat
	var s string
	var outs string
	err = self.setCommandLineSelfArgs()
	if err != nil {
		return err
	}

	parsers = make([]*parserCompat, 0)
	parsers = self.findCommandsInPath(cmdname, parsers)
	if len(parsers) == 0 {
		return fmt.Errorf("%s", format_error("can not find [%s] for help", cmdname))
	}

	s = self.printHelp(parsers)
	if len(self.outputMode) > 0 {
		if self.outputMode[len(self.outputMode)-1] == "bash" {
			outs = fmt.Sprintf("cat <<EOFMM\n%s\nEOFMM\nexit 0", s)
			os.Stdout.WriteString(outs)
			os.Exit(0)
		}
	}
	_, err = out.Write([]byte(s))
	return err
}

func (self *ExtArgsParse) helpAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	err = self.PrintHelp(os.Stdout, params[0])
	if err != nil {
		return 0, err
	}
	os.Exit(5)
	return 0, nil
}

func (self *ExtArgsParse) incAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	var i int
	i = ns.GetInt(keycls.Optdest())
	i++
	ns.SetValue(keycls.Optdest(), i)
	return 0, nil
}

func (self *ExtArgsParse) commandAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	return 0, nil
}

func (self *ExtArgsParse) floatAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	var f64 float64
	if validx >= len(params) {
		err = fmt.Errorf("%s", format_error("need args [%d] [%s] [%v]", validx, keycls.Format(), params))
		return 1, err
	}
	f64, err = strconv.ParseFloat(params[validx], 64)
	if err != nil {
		err = fmt.Errorf("%s", format_error("parse [%s] not float", params[validx]))
		return 1, err
	}
	ns.SetValue(keycls.Optdest(), f64)
	return 1, nil
}

func (self *ExtArgsParse) loadJsonFile(ns *NameSpaceEx, cmdname string, jsonfile string) error {
	return nil
}

func (self *ExtArgsParse) parseSubCommandJsonSet(ns *NameSpaceEx) error {
	var s string
	var cmds []*parserCompat
	var parsers []*parserCompat
	var idx int
	var subname string
	var prefix string
	var jsondst string
	var jsonfile string
	var err error
	s = ns.GetString("subcommand")
	if len(s) > 0 && !self.noJsonOption {
		parsers = make([]*parserCompat, 0)
		cmds = self.findCommandsInPath(s, parsers)
		idx = len(cmds)
		for idx = len(cmds); idx >= 2; idx-- {
			subname = self.formatCmdFromCmdArray(cmds[:idx])
			prefix = strings.Replace(subname, ".", "_", -1)
			jsondst = fmt.Sprintf("%s_%s", prefix, self.jsonLong)
			jsonfile = ns.GetString(jsondst)
			if len(jsonfile) > 0 {
				err = self.loadJsonFile(ns, subname, jsonfile)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (self *ExtArgsParse) parseCommandJsonSet(ns *NameSpaceEx) error {
	var jsonfile string
	if !self.noJsonOption && len(self.jsonLong) > 0 {
		jsonfile = ns.GetString(self.jsonLong)
		if len(jsonfile) > 0 {
			return self.loadJsonFile(ns, "", jsonfile)
		}
	}
	return nil
}

func (self *ExtArgsParse) setEnvironValue(ns *NameSpaceEx) error {
	return nil
}

func (self *ExtArgsParse) parseEnvironmentSet(ns *NameSpaceEx) error {
	return self.setEnvironValue(ns)
}

func (self *ExtArgsParse) parseEnvSubCommandJsonSet(ns *NameSpaceEx) error {
	var s string
	var cmds []*parserCompat
	var parsers []*parserCompat
	var prefix string
	var subname string
	var jsondst string
	var jsonfile string
	var err error
	var idx int
	s = ns.GetString("subcommand")
	if len(s) > 0 && !self.noJsonOption && len(self.jsonLong) > 0 {
		parsers = make([]*parserCompat, 0)
		cmds = self.findCommandsInPath(s, parsers)
		for idx = len(cmds); idx >= 2; idx-- {
			subname = self.formatCmdFromCmdArray(cmds[:idx])
			prefix = strings.Replace(subname, ".", "_", -1)
			prefix = fmt.Sprintf("%s_%s", prefix, self.jsonLong)
			jsondst = strings.ToUpper(prefix)
			jsonfile = os.Getenv(jsondst)
			if len(jsonfile) > 0 {
				err = self.loadJsonFile(ns, subname, jsonfile)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (self *ExtArgsParse) parseEnvCommandJsonSet(ns *NameSpaceEx) error {
	var jsonenv string
	var jsonfile string
	if !self.noJsonOption && len(self.jsonLong) > 0 {
		jsonenv = fmt.Sprintf("EXTARGSPARSE_%s", self.jsonLong)
		jsonenv = strings.Replace(jsonenv, ".", "_", -1)
		jsonenv = strings.Replace(jsonenv, "-", "_", -1)
		jsonenv = strings.ToUpper(jsonenv)
		jsonfile = os.Getenv(jsonenv)
		if len(jsonfile) > 0 {
			return self.loadJsonFile(ns, "", jsonfile)
		}
	}
	return nil
}

func (self *ExtArgsParse) jsonValueBase(ns *NameSpaceEx, keycls *ExtKeyParse, value interface{}) error {
	ns.SetValue(keycls.Optdest(), value)
	return nil
}

func (self *ExtArgsParse) jsonValueError(ns *NameSpaceEx, keycls *ExtKeyParse, value interface{}) error {
	return fmt.Errorf("%s", format_error("set [%s] error", keycls.Format()))
}

func NewExtArgsParse(options *ExtArgsOptions, priority interface{}) (self *ExtArgsParse, err error) {
	var pr []int
	var iv int
	if priority == nil {
		pr = parser_priority_args
	} else {
		switch priority.(type) {
		case []int:
			pr = make([]int, 0)
			for _, iv = range priority.([]int) {
				pr = append(pr, iv)
			}
		default:
			return nil, fmt.Errorf("%s", format_error("unknown type [%s] [%v]", reflect.ValueOf(priority).Type().Name(), priority))
		}
		for _, iv = range pr {
			if !is_valid_priority(iv) {
				return nil, fmt.Errorf("%s", format_error("not valid priority [%d]", iv))
			}
		}
	}

	if options == nil {
		options, err = NewExtArgsOptions("{}")
		if err != nil {
			return nil, err
		}
	}

	self = &ExtArgsParse{}

	self.options = options
	self.mainCmd = newParserCompat(nil, options)

	self.helpHandler = options.GetString("helphandler")
	self.outputMode = make([]string, 0)
	self.ended = 0
	self.longPrefix = options.GetString("longprefix")
	self.shortPrefix = options.GetString("shortprefix")
	self.noHelpOption = options.GetBool("nohelpoption")
	self.noJsonOption = options.GetBool("nojsonoption")
	self.helpLong = options.GetString("helplong")
	self.helpShort = options.GetString("helpshort")
	self.jsonLong = options.GetString("jsonlong")
	self.cmdPrefixAdded = options.GetBool("cmdprefixadded")

	self.loadCommandMap = make(map[string]reflect.Value)
	self.optParseHandleMap = make(map[string]reflect.Value)
	self.setJsonValueMap = make(map[string]reflect.Value)
	self.parsePrioritySetMap = make(map[int]reflect.Value)

	/*first to make loadCommandMap*/
	self.bindLoadCommandMap("string", self.loadCommandLineBase)
	self.bindLoadCommandMap("unicode", self.loadCommandLineBase)
	self.bindLoadCommandMap("int", self.loadCommandLineBase)
	self.bindLoadCommandMap("long", self.loadCommandLineBase)
	self.bindLoadCommandMap("float", self.loadCommandLineBase)
	self.bindLoadCommandMap("list", self.loadCommandLineBase)
	self.bindLoadCommandMap("bool", self.loadCommandLineBase)
	self.bindLoadCommandMap("args", self.loadCommandLineArgs)
	self.bindLoadCommandMap("command", self.loadCommandSubparser)
	self.bindLoadCommandMap("prefix", self.loadCommandPrefix)
	self.bindLoadCommandMap("count", self.loadCommandLineBase)
	self.bindLoadCommandMap("help", self.loadCommandLineBase)
	self.bindLoadCommandMap("jsonfile", self.loadCommandLineBase)

	/*optParsehandleMap*/
	self.bindOptParseHandleMap("string", self.stringAction)
	self.bindOptParseHandleMap("unicode", self.stringAction)
	self.bindOptParseHandleMap("bool", self.boolAction)
	self.bindOptParseHandleMap("int", self.intAction)
	self.bindOptParseHandleMap("long", self.intAction)
	self.bindOptParseHandleMap("list", self.appendAction)
	self.bindOptParseHandleMap("count", self.incAction)
	self.bindOptParseHandleMap("help", self.helpAction)
	self.bindOptParseHandleMap("jsonfile", self.stringAction)
	self.bindOptParseHandleMap("command", self.commandAction)
	self.bindOptParseHandleMap("float", self.floatAction)

	self.loadPriority = pr

	/*parsePrioritySetMap*/
	self.bindParsePrioritySetMap(SUB_COMMAND_JSON_SET, self.parseSubCommandJsonSet)
	self.bindParsePrioritySetMap(COMMAND_JSON_SET, self.parseCommandJsonSet)
	self.bindParsePrioritySetMap(ENVIRONMENT_SET, self.parseEnvironmentSet)
	self.bindParsePrioritySetMap(ENV_SUB_COMMAND_JSON_SET, self.parseEnvSubCommandJsonSet)
	self.bindParsePrioritySetMap(ENV_COMMAND_JSON_SET, self.parseEnvCommandJsonSet)

	/*setJsonValueMap*/
	self.bindSetJsonValueMap("string", self.jsonValueBase)
	self.bindSetJsonValueMap("unicode", self.jsonValueBase)
	self.bindSetJsonValueMap("bool", self.jsonValueBase)
	self.bindSetJsonValueMap("int", self.jsonValueBase)
	self.bindSetJsonValueMap("long", self.jsonValueBase)
	self.bindSetJsonValueMap("list", self.jsonValueBase)
	self.bindSetJsonValueMap("count", self.jsonValueBase)
	self.bindSetJsonValueMap("jsonfile", self.jsonValueBase)
	self.bindSetJsonValueMap("float", self.jsonValueBase)
	self.bindSetJsonValueMap("command", self.jsonValueError)
	self.bindSetJsonValueMap("help", self.jsonValueError)

	err = nil
	return
}

func (self *ExtArgsParse) checkFlagInsert(keycls *ExtKeyParse, parsers []*parserCompat) error {
	var lastparser *parserCompat
	if len(parsers) > 0 {
		lastparser = parsers[len(parsers)-1]
	} else {
		lastparser = self.mainCmd
	}
	for _, k := range lastparser.CmdOpts {
		if k.FlagName() != "$" && keycls.FlagName() != "$" {
			if k.TypeName() != "help" && keycls.TypeName() != "help" {
				if k.Optdest() == keycls.Optdest() {
					return fmt.Errorf("%s", format_error("[%s] already inserted", keycls.Optdest()))
				}
			} else if k.TypeName() == "help" && keycls.TypeName() == "help" {
				return fmt.Errorf("%s", format_error("help [%s] had already inserted", keycls.Format()))
			}
		} else if k.FlagName() == "$" && keycls.FlagName() == "$" {
			return fmt.Errorf("%s", format_error("args [%s] already inserted", keycls.Format()))
		}
	}
	lastparser.CmdOpts = append(lastparser.CmdOpts, keycls)
	return nil
}

func (self *ExtArgsParse) formatCmdFromCmdArray(parsers []*parserCompat) string {
	var cmdname string
	cmdname = ""
	for _, c := range parsers {
		if len(cmdname) > 0 {
			cmdname += "."
		}
		cmdname += c.CmdName
	}
	return cmdname
}

func (self *ExtArgsParse) loadCommandLineJsonFile(keycls *ExtKeyParse, parsers []*parserCompat) error {
	return self.checkFlagInsert(keycls, parsers)
}

func (self *ExtArgsParse) loadCommandLineJsonAdded(parsers []*parserCompat) error {
	var prefix string
	var key string
	var value interface{}
	var keycls *ExtKeyParse
	var err error
	prefix = ""
	key = fmt.Sprintf("%s##json input file to get the value set##", self.jsonLong)
	value = nil
	prefix = self.formatCmdFromCmdArray(parsers)
	prefix = strings.Replace(prefix, ".", "_", -1)
	keycls, err = newExtKeyParse_long(prefix, key, value, true, false, true, self.longPrefix, self.shortPrefix, false)
	assert_test(err == nil, "create json keycls error [%v]", err)
	return self.loadCommandLineJsonFile(keycls, parsers)
}

func (self *ExtArgsParse) loadCommandLineHelp(keycls *ExtKeyParse, parsers []*parserCompat) error {
	return self.checkFlagInsert(keycls, parsers)
}

func (self *ExtArgsParse) loadCommandLineHelpAdded(parsers []*parserCompat) error {
	var key string
	var keycls *ExtKeyParse
	var err error
	key = fmt.Sprintf("%s", self.helpLong)
	if len(self.helpShort) > 0 {
		key += fmt.Sprintf("|%s", self.helpShort)
	}
	keycls, err = newExtKeyParse_long("", key, nil, true, true, false, self.longPrefix, self.shortPrefix, false)
	assert_test(err == nil, "create help keycls error [%v]", err)
	return self.loadCommandLineHelp(keycls, parsers)
}

func (self *ExtArgsParse) callLoadCommandMapFunc(prefix string, keycls *ExtKeyParse, parsers []*parserCompat) error {
	var out []reflect.Value
	var in []reflect.Value
	in = make([]reflect.Value, 3)
	in[0] = reflect.ValueOf(prefix)
	in[1] = reflect.ValueOf(keycls)
	in[2] = reflect.ValueOf(parsers)
	out = self.loadCommandMap[keycls.TypeName()].Call(in)
	assert_test(len(out) == 1, format_error("out len [%d]", len(out)))
	return out[0].Interface().(error)
}

func (self *ExtArgsParse) loadCommandLineInner(prefix string, vmap map[string]interface{}, parsers []*parserCompat) error {
	var err error
	var parentpath []*parserCompat
	var k string
	var v interface{}
	var keycls *ExtKeyParse
	if !self.noJsonOption && len(self.jsonLong) > 0 {
		err = self.loadCommandLineJsonAdded(parsers)
		if err != nil {
			return err
		}
	}

	if !self.noHelpOption && len(self.helpLong) > 0 {
		err = self.loadCommandLineHelpAdded(parsers)
		if err != nil {
			return err
		}
	}

	parentpath = make([]*parserCompat, 0)
	parentpath = append(parentpath, self.mainCmd)
	if len(parsers) > 0 {
		parentpath = parsers
	}

	for k, v = range vmap {
		self.Info("%s , %s , %v , False", prefix, k, v)
		keycls, err = newExtKeyParse_long(prefix, k, v, false, false, false, self.longPrefix, self.shortPrefix, self.options.GetBool("flagnochange"))
		if err != nil {
			return err
		}

		err = self.callLoadCommandMapFunc(prefix, keycls, parsers)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *ExtArgsParse) loadCommandLine(vmap map[string]interface{}) error {
	var parsers []*parserCompat
	if self.ended != 0 {
		return fmt.Errorf("%s", format_error("you have call ParseCommandLine before call LoadCommandLineString"))
	}
	parsers = make([]*parserCompat, 0)

	return self.loadCommandLineInner("", vmap, parsers)
}

func (self *ExtArgsParse) LoadCommandLineString(s string) error {
	var vmap map[string]interface{}
	var err error
	err = json.Unmarshal([]byte(s), &vmap)
	if err != nil {
		return fmt.Errorf("%s", format_error("parse [%s] error [%s]", s, err.Error()))
	}
	return self.loadCommandLine(vmap)
}

func (self *ExtArgsParse) setArgs(ns *NameSpaceEx, cmdpaths []*parserCompat, vals interface{}) error {
	var params []string
	var argskeycls *ExtKeyParse
	var cmdname string
	var vstr string
	var vint int
	params = vals.([]string)
	cmdname = self.formatCmdNamePath(cmdpaths)
	argskeycls = nil
	for _, c := range cmdpaths[len(cmdpaths)-1].CmdOpts {
		if c.FlagName() == "$" {
			argskeycls = c
			break
		}
	}
	if argskeycls == nil {
		return fmt.Errorf("%s", format_error("can not find [%s]", cmdname))
	}

	vstr = ""
	switch argskeycls.Value().(type) {
	case string:
		vstr = argskeycls.Value().(string)
	case int:
		vint = argskeycls.Value().(int)
	default:
		return fmt.Errorf("%s", format_error("cmd [%s] [%v] unknown type[%s]", cmdname, argskeycls.Value(), reflect.ValueOf(argskeycls.Value()).Type().Name()))
	}

	if len(vstr) != 0 {
		switch vstr {
		case "*":
			break
		case "+":
			if len(params) < 1 {
				return fmt.Errorf("%s", format_error("[%s] args [%s] < 1", cmdname, vstr))
			}
		case "?":
			if len(params) > 1 {
				return fmt.Errorf("%s", format_error("[%s] args [%s] > 1", cmdname, vstr))
			}
		default:
			return fmt.Errorf("%s", format_error("[%s] args [%s] unknown", cmdname, vstr))
		}
	} else {
		if len(params) != vint {
			return fmt.Errorf("%s", format_error("[%s] args [%d] != %d", cmdname, len(params), vint))
		}
	}
	if len(cmdname) > 0 {
		ns.SetValue("subnargs", params)
		ns.SetValue("subcommand", cmdname)
	} else {
		ns.SetValue("args", params)
	}

	return nil
}

func (self *ExtArgsParse) callOptMethodFunc(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	var in []reflect.Value
	var out []reflect.Value
	in = make([]reflect.Value, 4)
	in[0] = reflect.ValueOf(ns)
	in[1] = reflect.ValueOf(validx)
	in[2] = reflect.ValueOf(keycls)
	in[3] = reflect.ValueOf(params)
	out = self.optParseHandleMap[keycls.TypeName()].Call(in)
	step = out[0].Interface().(int)
	err = out[1].Interface().(error)
	return
}

func (self *ExtArgsParse) callKeyOptMethodFunc(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	var callfunc func(ns *NameSpaceEx, valid int, keycls *ExtKeyParse, params []string) (step int, err error)
	err = self.GetFuncPtr(keycls.Attr("optparse"), &callfunc)
	if err != nil {
		return 0, err
	}
	return callfunc(ns, validx, keycls, params)
}

func (self *ExtArgsParse) callOptMethod(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	if keycls.Attr("optparse") != "" {

	}
	return self.callOptMethodFunc(ns, validx, keycls, params)
}

func (self *ExtArgsParse) parseArgs(params []string) (ns *NameSpaceEx, err error) {
	var pstate *parseState
	var validx int
	var optval interface{}
	var keycls *ExtKeyParse
	var cmdpaths []*parserCompat
	var helpcmdname string
	var step int
	var helpparams []string
	pstate = newParseState(params, self.mainCmd, self.options)
	ns = newNameSpaceEx()
	for {
		validx, optval, keycls, err = pstate.StepOne()
		if err != nil {
			return nil, err
		}
		if keycls == nil {
			cmdpaths = pstate.GetCmdPaths()
			err = self.setArgs(ns, cmdpaths, optval)
			if err != nil {
				return nil, err
			}
			break
		} else if keycls.TypeName() == "help" {
			cmdpaths = pstate.GetCmdPaths()
			helpcmdname = self.formatCmdFromCmdArray(cmdpaths)
			helpparams = []string{helpcmdname}
			step, err = self.callOptMethod(ns, validx, keycls, helpparams)
		} else {
			step, err = self.callOptMethod(ns, validx, keycls, params)
		}
		if err != nil {
			return nil, err
		}
		err = pstate.AddParseArgs(step)
		if err != nil {
			return nil, err
		}
	}
	return ns, nil
}

func (self *ExtArgsParse) callParseSetMapFunc(idx int, ns *NameSpaceEx) error {
	var in []reflect.Value
	var out []reflect.Value
	in = make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(ns)
	out = self.parsePrioritySetMap[idx].Call(in)
	if len(out) < 1 {
		return fmt.Errorf("%s", format_error("can not get error map return"))
	}
	if out[0].Interface() == nil {
		return nil
	}
	return out[0].Interface().(error)
}

func (self *ExtArgsParse) setFloatValue(ns *NameSpaceEx, opt *ExtKeyParse, fv float64) error {
	if opt.TypeName() != "float" {
		return fmt.Errorf("%s", format_error("[%s] not for [%v] set", opt.TypeName(), fv))
	}
	ns.SetValue(opt.Optdest(), fv)
	return nil
}

func (self *ExtArgsParse) setIntValue(ns *NameSpaceEx, opt *ExtKeyParse, iv int) error {
	if opt.TypeName() != "int" {
		return fmt.Errorf("%s", format_error("[%s] not for [%v] set", opt.TypeName(), iv))
	}
	ns.SetValue(opt.Optdest(), iv)
	return nil
}

func (self *ExtArgsParse) setJsonValueNotDefined(ns *NameSpaceEx, parser *parserCompat, dest string, value interface{}) error {
	var err error
	for _, c := range parser.SubCommands {
		err = self.setJsonValueNotDefined(ns, c, dest, value)
		if err != nil {
			return err
		}
	}

	for _, opt := range parser.CmdOpts {
		if opt.IsFlag() && opt.TypeName() != "prefix" && opt.TypeName() != "args" && opt.TypeName() != "help" {
			if opt.Optdest() == dest {
				if !ns.IsAccessed(dest) {
					if value != nil {
						switch value.(type) {
						case uint16:
							err = self.setIntValue(ns, opt, int(value.(uint16)))
						case uint32:
							err = self.setIntValue(ns, opt, int(value.(uint32)))
						case uint64:
							err = self.setIntValue(ns, opt, int(value.(uint64)))
						case int16:
							err = self.setIntValue(ns, opt, int(value.(int16)))
						case int32:
							err = self.setIntValue(ns, opt, int(value.(int32)))
						case int64:
							err = self.setIntValue(ns, opt, int(value.(int64)))
						case int:
							err = self.setIntValue(ns, opt, int(value.(int)))
						case float32:
							err = self.setFloatValue(ns, opt, float64(value.(float32)))
						case float64:
							err = self.setFloatValue(ns, opt, float64(value.(float64)))
						case string:
							if opt.TypeName() != "string" {
								return fmt.Errorf("%s", format_error("[%s] not for [%v] set", opt.TypeName(), value))
							}
							err = nil
						case []string:
							if opt.TypeName() != "list" {
								return fmt.Errorf("%s", format_error("[%s] not for [%v] set", opt.TypeName(), value))
							}
							ns.SetValue(opt.Optdest(), value)
							err = nil
						default:
							err = fmt.Errorf("%s", format_error("[%s] not for [%v] set", opt.TypeName(), value))
						}

						if err != nil {
							return err
						}
					} else {
						if opt.TypeName() != "string" {
							return fmt.Errorf("%s", format_error("[%s] not for nil set", opt.TypeName()))
						}
						ns.SetValue(opt.Optdest(), value)
					}
				}
			}
		}
	}

	return nil
}

func (self *ExtArgsParse) setParserDefaultValue(ns *NameSpaceEx, parser *parserCompat) error {
	var err error
	for _, c := range parser.SubCommands {
		err = self.setParserDefaultValue(ns, c)
		if err != nil {
			return err
		}
	}

	for _, opt := range parser.CmdOpts {
		if opt.IsFlag() && opt.TypeName() != "prefix" &&
			opt.TypeName() != "help" && opt.TypeName() != "args" {
			err = self.setJsonValueNotDefined(ns, parser, opt.Optdest(), opt.Value())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (self *ExtArgsParse) setDefaultValue(ns *NameSpaceEx) error {
	return self.setParserDefaultValue(ns, self.mainCmd)
}

func (self *ExtArgsParse) setStructPart(ns *NameSpaceEx, ostruct interface{}) error {
	return nil
}

func (self *ExtArgsParse) callbackFunc(funcname string, ns *NameSpaceEx, ostruct interface{}, Context interface{}) error {
	var callfunc func(ns *NameSpaceEx, ostruct interface{}, Context interface{}) error
	var err error
	err = self.GetFuncPtr(funcname, &callfunc)
	if err != nil {
		return err
	}
	return callfunc(ns, ostruct, Context)
}

func (self *ExtArgsParse) ParseCommandLine(params interface{}, Context interface{}, ostruct interface{}, mode interface{}) (ns *NameSpaceEx, err error) {
	var s string
	var realparams []string
	var subcmd string
	var cmds []*parserCompat
	var parsers []*parserCompat
	var funcname string
	var idx int
	if mode != nil {
		switch mode.(type) {
		case string:
			s = mode.(string)
		default:
			return nil, fmt.Errorf("%s", format_error("mode [%v] type error", mode))
		}
		self.outputMode = append(self.outputMode, s)
		defer func() {
			self.outputMode = self.outputMode[:(len(self.outputMode) - 1)]
		}()
	}
	err = self.setCommandLineSelfArgs()
	if err != nil {
		return nil, err
	}
	if params == nil {
		realparams = os.Args[1:]
	} else {
		switch params.(type) {
		case []string:
			realparams = params.([]string)
		default:
			return nil, fmt.Errorf("%s", format_error("params [%v] type error", params))
		}
	}

	ns, err = self.parseArgs(realparams)
	if err != nil {
		return nil, err
	}

	for _, idx = range self.loadPriority {
		err = self.callParseSetMapFunc(idx, ns)
		if err != nil {
			return nil, err
		}
	}

	err = self.setDefaultValue(ns)
	if err != nil {
		return nil, err
	}

	err = self.setStructPart(ns, ostruct)
	if err != nil {
		return nil, err
	}

	subcmd = ns.GetString("subcommand")
	if len(subcmd) > 0 {
		parsers = make([]*parserCompat, 0)
		cmds = self.findCommandsInPath(subcmd, parsers)
		if len(cmds) > 0 {
			funcname = cmds[len(cmds)-1].KeyCls.Function()
			if len(funcname) > 0 && (len(self.outputMode) == 0 || self.outputMode[len(self.outputMode)-1] == "") {
				err = self.callbackFunc(funcname, ns, ostruct, Context)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return ns, nil
}

func (self *ExtArgsParse) getSubCommands(name string, cmdpaths []*parserCompat) []string {
	var retnames []string
	var c *parserCompat
	var sarr []string
	retnames = make([]string, 0)
	if len(cmdpaths) == 0 {
		cmdpaths = append(cmdpaths, self.mainCmd)
	}
	if len(name) == 0 {
		for _, c = range cmdpaths[len(cmdpaths)-1].SubCommands {
			retnames = append(retnames, c.CmdName)
		}
		sort.Strings(retnames)
		return retnames
	}
	sarr = strings.Split(name, ".")
	for _, c = range cmdpaths[len(cmdpaths)-1].SubCommands {
		if c.CmdName == sarr[0] {
			cmdpaths = append(cmdpaths, c)
			return self.getSubCommands(strings.Join(sarr[1:], "."), cmdpaths)
		}
	}
	return retnames
}

func (self *ExtArgsParse) GetSubCommands(name string) ([]string, error) {
	var err error
	var retnames []string
	var cmdpaths []*parserCompat
	retnames = []string{}
	err = self.setCommandLineSelfArgs()
	if err != nil {
		return retnames, err
	}
	cmdpaths = make([]*parserCompat, 0)
	retnames = self.getSubCommands(name, cmdpaths)
	return retnames, nil
}

func (self *ExtArgsParse) getCmdKey(cmdname string, cmdpaths []*parserCompat) *ExtKeyParse {
	var retkey *ExtKeyParse = nil
	var sarr []string
	var c *parserCompat
	if len(cmdpaths) == 0 {
		cmdpaths = append(cmdpaths, self.mainCmd)
	}
	if len(cmdname) == 0 {
		retkey = cmdpaths[len(cmdpaths)-1].KeyCls
		return retkey
	}

	sarr = strings.Split(cmdname, ".")
	for _, c = range cmdpaths[len(cmdpaths)-1].SubCommands {
		if c.CmdName == sarr[0] {
			cmdpaths = append(cmdpaths, c)
			return self.getCmdKey(strings.Join(sarr[1:], "."), cmdpaths)
		}
	}

	return nil
}

func (self *ExtArgsParse) GetCmdKey(cmdname string) (*ExtKeyParse, error) {
	var err error
	var cmdpaths []*parserCompat
	err = self.setCommandLineSelfArgs()
	if err != nil {
		return nil, err
	}

	cmdpaths = make([]*parserCompat, 0)
	return self.getCmdKey(cmdname, cmdpaths), nil
}

func (self *ExtArgsParse) sortCmdOpts(opts []*ExtKeyParse) []*ExtKeyParse {
	sort.Slice(opts, func(i, j int) bool {
		return opts[i].CmdName() < opts[j].CmdName()
	})
	return opts
}

func (self *ExtArgsParse) getCmdOpts(cmdname string, cmdpaths []*parserCompat) []*ExtKeyParse {
	var retkeys []*ExtKeyParse
	var opt *ExtKeyParse
	var c *parserCompat
	var sarr []string
	retkeys = make([]*ExtKeyParse, 0)
	if len(cmdpaths) == 0 {
		cmdpaths = append(cmdpaths, self.mainCmd)
	}

	if len(cmdname) == 0 {
		for _, opt = range cmdpaths[len(cmdpaths)-1].CmdOpts {
			if opt.IsFlag() {
				retkeys = append(retkeys, opt)
			}
		}
		return self.sortCmdOpts(retkeys)
	}

	sarr = strings.Split(cmdname, ".")
	for _, c = range cmdpaths[len(cmdpaths)-1].SubCommands {
		if c.CmdName == sarr[0] {
			cmdpaths = append(cmdpaths, c)
			return self.getCmdOpts(strings.Join(sarr[1:], "."), cmdpaths)
		}
	}

	return retkeys
}

func (self *ExtArgsParse) GetCmdOpts(cmdname string) ([]*ExtKeyParse, error) {
	var retkeys []*ExtKeyParse
	var err error
	var cmdpaths []*parserCompat
	retkeys = make([]*ExtKeyParse, 0)
	cmdpaths = make([]*parserCompat, 0)
	err = self.setCommandLineSelfArgs()
	if err != nil {
		return retkeys, err
	}
	return self.getCmdOpts(cmdname, cmdpaths), nil
}
