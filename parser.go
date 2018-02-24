package extargsparse

import (
	"reflect"
)

var parser_reserver_args = []string{"subcommand", "subnargs", "nargs", "extargs", "args"}
var parser_priority_args = []int{SUB_COMMAND_JSON_SET, COMMAND_JSON_SET, ENVIRONMENT_SET, ENV_SUB_COMMAND_JSON_SET, ENV_COMMAND_JSON_SET}

type ExtArgsParse struct {
	logObject
	options           *ExtArgsOptions
	mainCmd           *parserCompat
	errorHandler      string
	helpHandler       string
	outputMode        []string
	ended             int
	longPrefix        string
	shortPrefix       string
	noHelpOption      bool
	noJsonOption      bool
	helpLong          string
	helpShort         string
	jsonLong          string
	cmdPrefixAdded    bool
	loadPriority      []int
	loadCommandMap    map[string]reflect.Value
	optParseHandleMap map[string]reflect.Value
	parseSetMap       map[int]reflect.Value
	setJsonValue      map[string]reflect.Value
}

func NewExtArgsParse(options *ExtArgsOptions, priority interface{}) (self *ExtArgsParse, err error) {
	err = nil
	self = nil
	return
}
