package extargsparse

// constant for the options used
//   VAR_UPPER_CASE used in the extargsparse.ExtArgsOptions for variable for first character uppercase default is true
//   FUNC_UPPER_CASE used in the function for first character uppercase default is true
const (
	VAR_UPPER_CASE  = "varuppercase"
	FUNC_UPPER_CASE = "funcuppercase"
)

// constant for the priority in the NewExtArgsParse
//    COMMAND_SET for the command line input
//    SUB_COMMAND_JSON_SET  for the jsonfile specified in the subcommand
//    COMMAND_JSON_SET for the jsonfile specified in the top
//    ENVIRONMENT_SET  environment variable set
//    ENV_SUB_COMMAND_JSON_SET for the jsonfile specified by the evironment for subcommand
//    ENV_COMMAND_JSON_SET for the jsonfile specified by the environment for top
//    DEFAULT_SET  default value set by the json string
//    default priority is in the int order [COMMAND_SET,SUB_COMMAND_JSON_SET,COMMAND_JSON_SET,ENVIRONMENT_SET,ENV_SUB_COMMAND_JSON_SET,ENV_COMMAND_JSON_SET,DEFAULT_SET]
const (
	COMMAND_SET              = 10
	SUB_COMMAND_JSON_SET     = 20
	COMMAND_JSON_SET         = 30
	ENVIRONMENT_SET          = 40
	ENV_SUB_COMMAND_JSON_SET = 50
	ENV_COMMAND_JSON_SET     = 60
	DEFAULT_SET              = 70
)
