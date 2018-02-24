package extargsparse

import (
	"encoding/json"
	"testing"
)

func Test_key_A001(t *testing.T) {
	flags, err := newExtKeyParse("", "$flag|f+type", "string", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "flag")
	check_equal(t, flags.Longopt(), "--type-flag")
	check_equal(t, flags.Shortopt(), "-f")
	check_equal(t, flags.Optdest(), "type_flag")
	check_equal(t, flags.Value(), "string")
	check_equal(t, flags.ShortFlag(), "f")
	check_equal(t, flags.Prefix(), "type")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.VarName(), "type_flag")
}

func Test_key_A002(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : []}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("", "$flag|f+type", vmap["code"], true)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "flag")
	check_equal(t, flags.ShortFlag(), "f")
	check_equal(t, flags.Longopt(), "--type-flag")
	check_equal(t, flags.Shortopt(), "-f")
	check_equal(t, flags.Optdest(), "type_flag")
	check_equal(t, flags.Value(), vmap["code"])
	check_equal(t, flags.TypeName(), "list")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.VarName(), "type_flag")
}

func Test_key_A003(t *testing.T) {
	var flags *extKeyParse
	var err error
	flags, err = newExtKeyParse("", "flag|f", false, false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "flag")
	check_equal(t, flags.ShortFlag(), "f")
	check_equal(t, flags.Longopt(), "--flag")
	check_equal(t, flags.Shortopt(), "-f")
	check_equal(t, flags.Optdest(), "flag")
	check_equal(t, flags.Value(), false)
	check_equal(t, flags.TypeName(), "bool")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.VarName(), "flag")
}

func Test_key_A004(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("newtype", "flag<flag.main>##help for flag##", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.CmdName(), "flag")
	check_equal(t, flags.Function(), "flag.main")
	check_equal(t, flags.TypeName(), "command")
	check_equal(t, flags.Prefix(), "newtype")
	check_equal(t, flags.HelpInfo(), "help for flag")
	check_equal(t, flags.FlagName(), "")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.Value(), vmap["code"])
	check_equal(t, flags.IsFlag(), false)
	check_equal(t, flags.IsCmd(), true)
	check_equal(t, flags.VarName(), "")
}

func Test_key_A005(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("", "flag<flag.main>##help for flag##", "", true)
	check_equal(t, err, nil)
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.TypeName(), "string")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.FlagName(), "flag")
	check_equal(t, flags.HelpInfo(), "help for flag")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.Value(), "")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.VarName(), "flag.main")
	check_equal(t, flags.Longopt(), "--flag")
	check_equal(t, flags.Shortopt(), "")
	check_equal(t, flags.Optdest(), "flag")
}

func Test_key_A006(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {"new": false}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("", "flag+type<flag.main>##main", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.CmdName(), "flag")
	check_equal(t, flags.Prefix(), "type")
	check_equal(t, flags.Function(), "flag.main")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.FlagName(), "")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.IsFlag(), false)
	check_equal(t, flags.IsCmd(), true)
	check_equal(t, flags.TypeName(), "command")
	check_equal(t, flags.Value(), vmap["code"])
	check_equal(t, flags.VarName(), "")
}

func Test_key_A007(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("", "+flag", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.Prefix(), "flag")
	check_equal(t, flags.Value(), vmap["code"])
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.FlagName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.TypeName(), "prefix")
	check_equal(t, flags.VarName(), "")
}

func Test_key_A008(t *testing.T) {
	var err error
	_, err = newExtKeyParse("", "+flag## help ##", nil, false)
	check_not_equal(t, err, nil)
}

func Test_key_A009(t *testing.T) {
	var err error
	_, err = newExtKeyParse("", "+flag<flag.main>", nil, false)
	check_not_equal(t, err, nil)
}

func Test_key_A010(t *testing.T) {
	var err error
	_, err = newExtKeyParse("", "flag|f2", "", false)
	check_not_equal(t, err, nil)
}

func Test_key_A011(t *testing.T) {
	var err error
	_, err = newExtKeyParse("", "f|f2", "", false)
	check_not_equal(t, err, nil)
}

func Test_key_A012(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("", "$flag|f<flag.main>", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.Value(), nil)
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.ShortFlag(), "f")
	check_equal(t, flags.FlagName(), "flag")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.TypeName(), "string")
	check_equal(t, flags.VarName(), "flag.main")
	check_equal(t, flags.Longopt(), "--flag")
	check_equal(t, flags.Shortopt(), "-f")
	check_equal(t, flags.Optdest(), "flag")
}

func Test_key_A013(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("", "$flag|f+cc<flag.main>", nil, false)
	check_equal(t, err, nil)
	check_equal(t, flags.Prefix(), "cc")
	check_equal(t, flags.Value(), nil)
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.ShortFlag(), "f")
	check_equal(t, flags.FlagName(), "flag")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.TypeName(), "string")
	check_equal(t, flags.VarName(), "flag.main")
	check_equal(t, flags.Longopt(), "--cc-flag")
	check_equal(t, flags.Shortopt(), "-f")
	check_equal(t, flags.Optdest(), "cc_flag")
}

func Test_key_A014(t *testing.T) {
	var err error
	_, err = newExtKeyParse("", "c$", "", false)
	check_not_equal(t, err, nil)
}

func Test_key_A015(t *testing.T) {
	var err error
	_, err = newExtKeyParse("", "$$", "", false)
	check_not_equal(t, err, nil)
}

func Test_key_A016(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {"nargs" : "+"}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("", "$", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "$")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.TypeName(), "args")
	check_equal(t, flags.VarName(), "args")
	check_equal(t, flags.Nargs().(string), "+")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
}

func Test_key_A017(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : 3.3}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("type", "flag+app## flag help ##", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "flag")
	check_equal(t, flags.Prefix(), "type_app")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.TypeName(), "float")
	check_equal(t, flags.Value(), 3.3)
	check_equal(t, flags.Longopt(), "--type-app-flag")
	check_equal(t, flags.Shortopt(), "")
	check_equal(t, flags.Optdest(), "type_app_flag")
	check_equal(t, flags.HelpInfo(), " flag help ")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.VarName(), "type_app_flag")
}

func Test_key_A018(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("", "flag+app<flag.main>## flag help ##", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "")
	check_equal(t, flags.Prefix(), "app")
	check_equal(t, flags.CmdName(), "flag")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.VarName(), "")
	check_equal(t, flags.TypeName(), "command")
	check_equal(t, flags.Value(), vmap["code"])
	check_equal(t, flags.Function(), "flag.main")
	check_equal(t, flags.HelpInfo(), " flag help ")
	check_equal(t, flags.IsFlag(), false)
	check_equal(t, flags.IsCmd(), true)
}

func Test_key_A019(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : { "prefix" : "good", "value": false}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("", "$flag## flag help ##", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "flag")
	check_equal(t, flags.Prefix(), "good")
	check_equal(t, flags.Value(), false)
	check_equal(t, flags.TypeName(), "bool")
	check_equal(t, flags.HelpInfo(), " flag help ")
	check_equal(t, flags.Nargs(), 0)
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.VarName(), "good_flag")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.Longopt(), "--good-flag")
	check_equal(t, flags.Shortopt(), "")
	check_equal(t, flags.Optdest(), "good_flag")
}

func Test_key_A020(t *testing.T) {
	var err error
	_, err = newExtKeyParse("", "$", nil, false)
	check_not_equal(t, err, nil)
}

func Test_key_A021(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : { "nargs" : "?", "value": null}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("command", "$## self define ##", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.Prefix(), "command")
	check_equal(t, flags.VarName(), "subnargs")
	check_equal(t, flags.FlagName(), "$")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.Value(), nil)
	check_equal(t, flags.TypeName(), "args")
	check_equal(t, flags.Nargs(), "?")
	check_equal(t, flags.HelpInfo(), " self define ")
}

func Test_key_A022(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("command", "+flag", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.Prefix(), "command_flag")
	check_equal(t, flags.Value(), vmap["code"])
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.FlagName(), "")
	check_equal(t, flags.VarName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.TypeName(), "prefix")
}

func Test_key_A023(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {"prefix": "good","value": 3.9, "nargs": 1}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("", "$flag## flag help ##", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "flag")
	check_equal(t, flags.Prefix(), "good")
	check_equal(t, flags.Value(), 3.9)
	check_equal(t, flags.TypeName(), "float")
	check_equal(t, flags.HelpInfo(), " flag help ")
	check_equal(t, flags.Nargs(), 1)
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.Longopt(), "--good-flag")
	check_equal(t, flags.Shortopt(), "")
	check_equal(t, flags.Optdest(), "good_flag")
	check_equal(t, flags.VarName(), "good_flag")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
}

func Test_key_A024(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	js = `{"code" : {"prefix": "good","value": false, "nargs": 2}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	_, err = newExtKeyParse("", "$flag## flag help ##", vmap["code"], false)
	check_not_equal(t, err, nil)
}

func Test_key_A027(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("dep", "verbose|v", "+", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "verbose")
	check_equal(t, flags.ShortFlag(), "v")
	check_equal(t, flags.Prefix(), "dep")
	check_equal(t, flags.TypeName(), "count")
	check_equal(t, flags.Value(), 0)
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Nargs(), 0)
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.VarName(), "dep_verbose")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.Longopt(), "--dep-verbose")
	check_equal(t, flags.Shortopt(), "-v")
	check_equal(t, flags.Optdest(), "dep_verbose")
}

func Test_key_A028(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("", "verbose|v## new help info ##", "+", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "verbose")
	check_equal(t, flags.ShortFlag(), "v")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.TypeName(), "count")
	check_equal(t, flags.Value(), 0)
	check_equal(t, flags.HelpInfo(), " new help info ")
	check_equal(t, flags.Nargs(), 0)
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.Optdest(), "verbose")
	check_equal(t, flags.VarName(), "verbose")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.Longopt(), "--verbose")
	check_equal(t, flags.Shortopt(), "-v")
}

func Test_key_A029(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("", "rollback|R## rollback not set ##", true, false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "rollback")
	check_equal(t, flags.ShortFlag(), "R")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.TypeName(), "bool")
	check_equal(t, flags.Value(), true)
	check_equal(t, flags.HelpInfo(), " rollback not set ")
	check_equal(t, flags.Nargs(), 0)
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.Optdest(), "rollback")
	check_equal(t, flags.VarName(), "rollback")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.Longopt(), "--no-rollback")
	check_equal(t, flags.Shortopt(), "-R")
}

func Test_key_A030(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("", "maxval|m##max value set ##", 0xffffffff, false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "maxval")
	check_equal(t, flags.ShortFlag(), "m")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.TypeName(), "int")
	check_equal(t, flags.Value(), 0xffffffff)
	check_equal(t, flags.HelpInfo(), "max value set ")
	check_equal(t, flags.Nargs(), 1)
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.Optdest(), "maxval")
	check_equal(t, flags.VarName(), "maxval")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
	check_equal(t, flags.Longopt(), "--maxval")
	check_equal(t, flags.Shortopt(), "-m")
}

func Test_key_A031(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : ["maxval"]}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("", "maxval|m", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "maxval")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.Value(), vmap["code"])
	check_equal(t, flags.TypeName(), "list")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Nargs(), 1)
	check_equal(t, flags.ShortFlag(), "m")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.Longopt(), "--maxval")
	check_equal(t, flags.Shortopt(), "-m")
	check_equal(t, flags.Optdest(), "maxval")
	check_equal(t, flags.VarName(), "maxval")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.IsCmd(), false)
}

func Test_key_A032(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("", "$<numargs>", "+", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "$")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.Value(), nil)
	check_equal(t, flags.TypeName(), "args")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Nargs(), "+")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.VarName(), "numargs")
}

func Test_key_A033(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("", "$", "+", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "$")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.Value(), nil)
	check_equal(t, flags.TypeName(), "args")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Nargs(), "+")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.VarName(), "args")
}

func Test_key_A034(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("prefix", "$", "+", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "$")
	check_equal(t, flags.Prefix(), "prefix")
	check_equal(t, flags.Value(), nil)
	check_equal(t, flags.TypeName(), "args")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Nargs(), "+")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.VarName(), "subnargs")
}

func Test_key_A035(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("prefix", "$<newargs>", "+", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "$")
	check_equal(t, flags.Prefix(), "prefix")
	check_equal(t, flags.Value(), nil)
	check_equal(t, flags.TypeName(), "args")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Nargs(), "+")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.VarName(), "newargs")
}

func Test_key_A036(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse("prefix", "$<newargs>!func=args_opt_func;wait=cc!", "+", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "$")
	check_equal(t, flags.Prefix(), "prefix")
	check_equal(t, flags.Value(), nil)
	check_equal(t, flags.TypeName(), "args")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Nargs(), "+")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.VarName(), "newargs")
	check_equal(t, flags.Attr("func"), "args_opt_func")
	check_equal(t, flags.Attr("wait"), "cc")
}

func Test_key_A037(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = newExtKeyParse_long("prefix", "help|h!func=args_opt_func;wait=cc!", nil, false, true, false, "--", "-", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "help")
	check_equal(t, flags.Prefix(), "prefix")
	check_equal(t, flags.Value(), nil)
	check_equal(t, flags.TypeName(), "help")
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Nargs(), 0)
	check_equal(t, flags.ShortFlag(), "h")
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.VarName(), "prefix_help")
	check_equal(t, flags.Attr("func"), "args_opt_func")
	check_equal(t, flags.Attr("wait"), "cc")
	check_equal(t, flags.Longopt(), "--help")
	check_equal(t, flags.Shortopt(), "-h")
	check_equal(t, flags.Optdest(), "prefix_help")
}

func Test_key_A038(t *testing.T) {
	var err error
	var flag1, flag2 *extKeyParse
	flag1, err = newExtKeyParse_long("prefix", "help|h!func=args_opt_func;wait=cc!", nil, false, true, false, "--", "-", false)
	check_equal(t, err, nil)
	flag2, err = newExtKeyParse_long("prefix", "help|h!func=args_opt_func;wait=cc!", nil, false, false, false, "--", "-", false)
	check_equal(t, err, nil)
	check_equal(t, flag1.Equal(flag2), false)
	flag1, err = newExtKeyParse_long("prefix", "help|h!func=args_opt_func;wait=cc!", nil, false, true, false, "--", "-", false)
	check_equal(t, err, nil)
	flag2, err = newExtKeyParse_long("prefix", "help|h!func=args_opt_func;wait=cc!", nil, false, true, false, "--", "-", false)
	check_equal(t, err, nil)
	check_equal(t, flag1.Equal(flag2), true)
	flag1, err = newExtKeyParse_long("prefix", "help|h!func=args_opt_func!", nil, false, true, false, "--", "-", false)
	check_equal(t, err, nil)
	flag2, err = newExtKeyParse_long("prefix", "help|h!func=args_opt_func;wait=cc!", nil, false, true, false, "--", "-", false)
	check_equal(t, err, nil)
	check_equal(t, flag1.Equal(flag2), false)
}

func Test_key_A039(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {"modules": [],"$<NARGS>" : "+"}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("rdep", "ip", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.IsCmd(), true)
	check_equal(t, flags.CmdName(), "ip")
	check_equal(t, flags.Prefix(), "rdep")
	js = `{"code" : []}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("rdep_ip", "modules", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.Value(), vmap["code"])
	check_equal(t, flags.Prefix(), "rdep_ip")
	check_equal(t, flags.Longopt(), "--rdep-ip-modules")
	check_equal(t, flags.Shortopt(), "")
	check_equal(t, flags.Optdest(), "rdep_ip_modules")
	check_equal(t, flags.VarName(), "rdep_ip_modules")
}

func Test_key_A040(t *testing.T) {
	var err error
	var flag1, flag2 *extKeyParse
	flag1, err = newExtKeyParse_long("prefix", "json!func=args_opt_func;wait=cc!", nil, false, false, true, "--", "-", false)
	check_equal(t, err, nil)
	flag2, err = newExtKeyParse("prefix", "json!func=args_opt_func;wait=cc!", nil, false)
	check_equal(t, err, nil)
	check_equal(t, flag1.Equal(flag2), false)
	flag1, err = newExtKeyParse_long("prefix", "json!func=args_opt_func;wait=cc!", nil, false, false, true, "--", "-", false)
	check_equal(t, err, nil)
	flag2, err = newExtKeyParse_long("prefix", "json!func=args_opt_func;wait=cc!", nil, false, false, true, "--", "-", false)
	check_equal(t, err, nil)
	check_equal(t, flag1.Equal(flag2), true)
	check_equal(t, flag1.Optdest(), "prefix_json")
	check_equal(t, flag1.Longopt(), "--prefix-json")
}

func Test_key_A041(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {"nargs": 1,"attr" : {"func":"args_opt_func","wait": "cc"}}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("prefix", "$json", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.Prefix(), "prefix")
	check_equal(t, flags.IsFlag(), true)
	check_equal(t, flags.Attr("func"), "args_opt_func")
	check_equal(t, flags.Attr("wait"), "cc")
	check_equal(t, flags.FlagName(), "json")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.Longopt(), "--prefix-json")
	check_equal(t, flags.Shortopt(), "")
	check_equal(t, flags.Optdest(), "prefix_json")
	check_equal(t, flags.VarName(), "prefix_json")
}

func Test_key_A042(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {}}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse("", "main", vmap["code"], false)
	check_equal(t, err, nil)
	check_equal(t, flags.Prefix(), "main")
	check_equal(t, flags.IsFlag(), false)
	check_equal(t, flags.IsCmd(), true)
	check_equal(t, flags.Attr(""), "")
	check_equal(t, flags.CmdName(), "main")
	check_equal(t, flags.Function(), "")
}

func Test_key_A043(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : true}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse_long("", "rollback|R## rollback not set ##", vmap["code"], false, false, false, "++", "+", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "rollback")
	check_equal(t, flags.ShortFlag(), "R")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.TypeName(), "bool")
	check_equal(t, flags.Value(), true)
	check_equal(t, flags.HelpInfo(), " rollback not set ")
	check_equal(t, flags.Nargs(), 0)
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.Optdest(), "rollback")
	check_equal(t, flags.VarName(), "rollback")
	check_equal(t, flags.Longopt(), "++no-rollback")
	check_equal(t, flags.Shortopt(), "+R")
}

func Test_key_A044(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : true}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse_long("", "rollback|R## rollback not set ##", vmap["code"], false, false, false, "++", "+", false)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "rollback")
	check_equal(t, flags.ShortFlag(), "R")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.TypeName(), "bool")
	check_equal(t, flags.Value(), true)
	check_equal(t, flags.HelpInfo(), " rollback not set ")
	check_equal(t, flags.Nargs(), 0)
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.Optdest(), "rollback")
	check_equal(t, flags.VarName(), "rollback")
	check_equal(t, flags.Longopt(), "++no-rollback")
	check_equal(t, flags.Shortopt(), "+R")
	check_equal(t, flags.LongPrefix(), "++")
	check_equal(t, flags.ShortPrefix(), "+")
}

func Test_key_A045(t *testing.T) {
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : false}`
	err = json.Unmarshal([]byte(js), &vmap)
	check_equal(t, err, nil)
	flags, err = newExtKeyParse_long("", "crl_CA_compromise", vmap["code"], false, false, false, "++", "+", true)
	check_equal(t, err, nil)
	check_equal(t, flags.FlagName(), "crl_CA_compromise")
	check_equal(t, flags.ShortFlag(), "")
	check_equal(t, flags.Prefix(), "")
	check_equal(t, flags.TypeName(), "bool")
	check_equal(t, flags.Value(), false)
	check_equal(t, flags.HelpInfo(), "")
	check_equal(t, flags.Nargs(), 0)
	check_equal(t, flags.CmdName(), "")
	check_equal(t, flags.Function(), "")
	check_equal(t, flags.Optdest(), "crl_CA_compromise")
	check_equal(t, flags.VarName(), "crl_CA_compromise")
	check_equal(t, flags.Longopt(), "++crl_CA_compromise")
	check_equal(t, flags.Shortopt(), "")
	check_equal(t, flags.LongPrefix(), "++")
	check_equal(t, flags.ShortPrefix(), "+")
}
