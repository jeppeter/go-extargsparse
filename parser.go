package extargsparse

var parser_reserver_args = []string{"subcommand", "subnargs", "nargs", "extargs", "args"}
var parser_priority_args = []int{SUB_COMMAND_JSON_SET, COMMAND_JSON_SET, ENVIRONMENT_SET, ENV_SUB_COMMAND_JSON_SET, ENV_COMMAND_JSON_SET}

type ExtArgsParse struct {
	logObject
}
