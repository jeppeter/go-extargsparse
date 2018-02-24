package extargsparse

import (
	"encoding/json"
	"fmt"
	"reflect"
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

func (self *ExtArgsParse) loadCommandLineBase(prefix string, keycls *ExtKeyParse, curparser []*parserCompat) error {
	return nil
}

func (self *ExtArgsParse) loadCommandLineArgs(prefix string, keycls *ExtKeyParse, curparser []*parserCompat) error {
	return nil
}

func (self *ExtArgsParse) loadCommandLineSubparser(prefix string, keycls *ExtKeyParse, curparser []*parserCompat) error {
	return nil
}

func (self *ExtArgsParse) loadCommandLinePrefix(prefix string, keycls *ExtKeyParse, curparser []*parserCompat) error {
	return nil
}

func (self *ExtArgsParse) stringAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	return 1, nil
}

func (self *ExtArgsParse) boolAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	return 0, nil
}

func (self *ExtArgsParse) intAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	return 1, nil
}

func (self *ExtArgsParse) appendAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	return 1, nil
}

func (self *ExtArgsParse) helpAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	return 0, nil
}

func (self *ExtArgsParse) incAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	return 0, nil
}

func (self *ExtArgsParse) commandAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	return 0, nil
}

func (self *ExtArgsParse) floatAction(ns *NameSpaceEx, validx int, keycls *ExtKeyParse, params []string) (step int, err error) {
	return 1, nil
}

func (self *ExtArgsParse) parseSubCommandJsonSet(ns *NameSpaceEx) error {
	return nil
}

func (self *ExtArgsParse) parseCommandJsonSet(ns *NameSpaceEx) error {
	return nil
}

func (self *ExtArgsParse) parseEnvironmentSet(ns *NameSpaceEx) error {
	return nil
}

func (self *ExtArgsParse) parseEnvSubCommandJsonSet(ns *NameSpaceEx) error {
	return nil
}

func (self *ExtArgsParse) parseEnvCommandJsonSet(ns *NameSpaceEx) error {
	return nil
}

func (self *ExtArgsParse) jsonValueBase(ns *NameSpaceEx, keycls *ExtKeyParse, value interface{}) error {
	return nil
}

func (self *ExtArgsParse) jsonValueError(ns *NameSpaceEx, keycls *ExtKeyParse, value interface{}) error {
	return nil
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
	self.bindLoadCommandMap("command", self.loadCommandLineSubparser)
	self.bindLoadCommandMap("prefix", self.loadCommandLinePrefix)
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

func (self *ExtArgsParse) loadCommandLine(vmap map[string]interface{}) error {
	if self.ended != 0 {
		return fmt.Errorf("%s", format_error("you have call ParseCommandLine before call LoadCommandLineString"))
	}
	return nil
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
