package extargsparse

import (
	"encoding/json"
	"testing"
)

func Test_A001(t *testing.T) {
	flags, err := NewExtKeyParse_short("", "$flag|f+type", "string", false)
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

func Test_A002(t *testing.T) {
	var v interface{}
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : []}`
	err = json.Unmarshal([]byte(js), &v)
	check_equal(t, err, nil)
	vmap = v.(map[string]interface{})
	flags, err = NewExtKeyParse_short("", "$flag|f+type", vmap["code"], true)
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

func Test_A003(t *testing.T) {
	var flags *extKeyParse
	var err error
	flags, err = NewExtKeyParse_short("", "flag|f", false, false)
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

func Test_A004(t *testing.T) {
	var v interface{}
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {}}`
	err = json.Unmarshal([]byte(js), &v)
	check_equal(t, err, nil)
	vmap = v.(map[string]interface{})
	flags, err = NewExtKeyParse_short("newtype", "flag<flag.main>##help for flag##", vmap["code"], false)
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

func Test_A005(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = NewExtKeyParse_short("", "flag<flag.main>##help for flag##", "", true)
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

func Test_A006(t *testing.T) {
	var v interface{}
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {"new": false}}`
	err = json.Unmarshal([]byte(js), &v)
	check_equal(t, err, nil)
	vmap = v.(map[string]interface{})
	flags, err = NewExtKeyParse_short("", "flag+type<flag.main>##main", vmap["code"], false)
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

func Test_A007(t *testing.T) {
	var v interface{}
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {}}`
	err = json.Unmarshal([]byte(js), &v)
	check_equal(t, err, nil)
	vmap = v.(map[string]interface{})
	flags, err = NewExtKeyParse_short("", "+flag", vmap["code"], false)
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

func Test_A008(t *testing.T) {
	var err error
	_, err = NewExtKeyParse_short("", "+flag## help ##", nil, false)
	check_not_equal(t, err, nil)
}

func Test_A009(t *testing.T) {
	var err error
	_, err = NewExtKeyParse_short("", "+flag<flag.main>", nil, false)
	check_not_equal(t, err, nil)
}

func Test_A010(t *testing.T) {
	var err error
	_, err = NewExtKeyParse_short("", "flag|f2", "", false)
	check_not_equal(t, err, nil)
}

func Test_A011(t *testing.T) {
	var err error
	_, err = NewExtKeyParse_short("", "f|f2", "", false)
	check_not_equal(t, err, nil)
}

func Test_A012(t *testing.T) {
	var v interface{}
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {}}`
	err = json.Unmarshal([]byte(js), &v)
	check_equal(t, err, nil)
	vmap = v.(map[string]interface{})
	flags, err = NewExtKeyParse_short("", "$flag|f<flag.main>", vmap["code"], false)
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

func Test_A013(t *testing.T) {
	var err error
	var flags *extKeyParse
	flags, err = NewExtKeyParse_short("", "$flag|f+cc<flag.main>", nil, false)
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

func Test_A014(t *testing.T) {
	var err error
	_, err = NewExtKeyParse_short("", "c$", "", false)
	check_not_equal(t, err, nil)
}

func Test_A015(t *testing.T) {
	var err error
	_, err = NewExtKeyParse_short("", "$$", "", false)
	check_not_equal(t, err, nil)
}

func Test_A016(t *testing.T) {
	var v interface{}
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {"nargs" : "+"}}`
	err = json.Unmarshal([]byte(js), &v)
	check_equal(t, err, nil)
	vmap = v.(map[string]interface{})
	flags, err = NewExtKeyParse_short("", "$", vmap["code"], false)
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

func Test_A017(t *testing.T) {
	var v interface{}
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : 3.3}`
	err = json.Unmarshal([]byte(js), &v)
	check_equal(t, err, nil)
	vmap = v.(map[string]interface{})
	flags, err = NewExtKeyParse_short("type", "flag+app## flag help ##", vmap["code"], false)
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

func Test_A018(t *testing.T) {
	var v interface{}
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : {}}`
	err = json.Unmarshal([]byte(js), &v)
	check_equal(t, err, nil)
	vmap = v.(map[string]interface{})
	flags, err = NewExtKeyParse_short("", "flag+app<flag.main>## flag help ##", vmap["code"], false)
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

func Test_A019(t *testing.T) {
	var v interface{}
	var vmap map[string]interface{}
	var err error
	var js string
	var flags *extKeyParse
	js = `{"code" : { "prefix" : "good", "value": false}}`
	err = json.Unmarshal([]byte(js), &v)
	check_equal(t, err, nil)
	vmap = v.(map[string]interface{})
	flags, err = NewExtKeyParse_short("", "$flag## flag help ##", vmap["code"], false)
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

func Test_A020(t *testing.T) {
	var err error
	_, err = NewExtKeyParse_short("", "$", nil, false)
	check_not_equal(t, err, nil)
}
