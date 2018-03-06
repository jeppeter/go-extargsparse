package extargsparse

// constant for the options used
//   OPT_PROG used for program name
//   OPT_USAGE used for help information
//   OPT_DESCRIPTION for description in help
//   OPT_EPILOG for help epilog
//   OPT_VERSION for version set
//   OPT_ERROR_HANDLER for error handler ,now is reserved
//   OPT_HELP_HANDLER for help option ,default is ""; "nohelp" for no help information print out
//   OPT_LONG_PREFIX for long prefix ,default is "--"
//   OPT_SHORT_PREFIX for short prefix, default is "-"
//   OPT_NO_HELP_OPTION for no help flag inserted into the opts default false
//   OPT_NO_JSON_OPTION for no json flag inserted into the opts default false
//   OPT_HELP_LONG for help flag flagname default "help"
//   OPT_HELP_SHORT for help flag shortflag default "h"
//   OPT_VAR_UPPER_CASE used in the extargsparse.ExtArgsOptions for variable for first character uppercase default is true
//   OPT_FUNC_UPPER_CASE used in the function for first character uppercase default is true
const (
	OPT_PROG            = "prog"
	OPT_USAGE           = "usage"
	OPT_DESCRIPTION     = "description"
	OPT_EPILOG          = "epilog"
	OPT_VERSION         = "version"
	OPT_ERROR_HANDLER   = "errorhandler"
	OPT_HELP_HANDLER    = "helphandler"
	OPT_LONG_PREFIX     = "longprefix"
	OPT_SHORT_PREFIX    = "shortprefix"
	OPT_NO_HELP_OPTION  = "nohelpoption"
	OPT_NO_JSON_OPTION  = "nojsonoption"
	OPT_HELP_LONG       = "helplong"
	OPT_HELP_SHORT      = "helpshort"
	OPT_VAR_UPPER_CASE  = "varuppercase"
	OPT_FUNC_UPPER_CASE = "funcuppercase"
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
