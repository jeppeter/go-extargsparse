package extargsparse

import (
	"encoding/json"
	"fmt"
	"strings"
)

type parserCompat struct {
	logObject
	KeyCls       *ExtKeyParse
	CmdName      string
	CmdOpts      []*ExtKeyParse
	SubCommands  []string
	HelpInfo     string
	CallFunction string
	ScreenWidth  int
	Epilog       string
	Description  string
	Prog         string
	Usage        string
	Version      string
}

func newParserCompat(keycls *ExtKeyParse, opt *ExtArgsOptions) *parserCompat {
	var vmap map[string]interface{}
	var js string
	var err error
	self := &parserCompat{logObject: *newLogObject("extargsparse")}
	if keycls != nil {
		assert_test(keycls.IsCmd(), "%s must be cmd", keycls.Format())
		self.KeyCls = keycls
		self.CmdName = keycls.CmdName()
		self.CmdOpts = make([]*ExtKeyParse, 0)
		self.SubCommands = make([]string, 0)
		self.HelpInfo = fmt.Sprintf("%s handler", self.CmdName)
		if len(keycls.HelpInfo()) > 0 {
			self.HelpInfo = keycls.HelpInfo()
		}
		self.CallFunction = ""
		if len(keycls.Function()) > 0 {
			self.CallFunction = keycls.Function()
		}
	} else {
		js = `{"code": {}}`
		err = json.Unmarshal([]byte(js), &vmap)
		assert_test(err == nil, "parse [%s] must be succ", js)
		self.KeyCls, err = newExtKeyParse("", "main", vmap["code"], false)
		assert_test(err == nil, "make main cmd must succ")
		self.CmdName = ""
		self.CmdOpts = make([]*ExtKeyParse, 0)
		self.SubCommands = make([]string, 0)
		self.HelpInfo = ""
		self.CallFunction = ""
	}
	self.ScreenWidth = 80
	if opt != nil && opt.GetValue("screenwidth") != nil {
		self.ScreenWidth = opt.GetValue("screenwidth").(int)
	}

	if self.ScreenWidth < 40 {
		self.ScreenWidth = 40
	}
	self.Epilog = ""
	self.Description = ""
	self.Prog = ""
	self.Usage = ""
	self.Version = ""
	return self
}

func (self *parserCompat) get_help_info(keycls *ExtKeyParse) string {
	var s string
	var err error
	var helpfunc func(kl *ExtKeyParse) string
	assert_test(keycls != nil, "must no be null")
	if keycls.Attr("opthelp") != "" {
		/*now it is the help function ,so we call this*/

		err = self.GetFuncPtr(keycls.Attr("opthelp"), &helpfunc)
		if err == nil {
			return helpfunc(keycls)
		}
		self.Warn("can not find function [%s] for opthelp", keycls.Attr("opthelp"))
	}

	/*ok we should make the */
	s = ""
	if keycls.TypeName() == "bool" {
		if keycls.Value() == true {
			s += fmt.Sprintf("%s set false default(True)", keycls.Optdest())
		} else {
			s += fmt.Sprintf("%s set true default(False)", keycls.Optdest())
		}
	} else if keycls.TypeName() == "string" && keycls.Value() != nil && keycls.Value() == "+" {
		if keycls.IsFlag() {
			s += fmt.Sprintf("%s inc", keycls.Optdest())
		} else {
			assert_test(false == true, "cmd(%s) can not set value(%v)", keycls.CmdName(), keycls.Value())
		}
	} else if keycls.TypeName() == "help" {
		s += "to display this help information"
	} else {
		if keycls.IsFlag() {
			s += fmt.Sprintf("%s set default(%s)", keycls.Optdest(), keycls.Value())
		} else {
			s += fmt.Sprintf("%s command exec", keycls.CmdName())
		}
	}

	if len(keycls.HelpInfo()) > 0 {
		s = keycls.HelpInfo()
	}
	return s
}

func (self *parserCompat) get_opt_help_optname(opt *ExtKeyParse) string {
	var s string = ""
	s += opt.Longopt()
	if len(opt.Shortopt()) > 0 {
		s += fmt.Sprintf("|%s", opt.Shortopt())
	}
	return s
}

func (self *parserCompat) get_opt_help_optexpr(opt *ExtKeyParse) string {
	var s string = ""
	if opt.TypeName() != "bool" && opt.TypeName() != "args" &&
		opt.TypeName() != "dict" && opt.TypeName() != "help" {
		s = opt.VarName()
		s = strings.Replace(s, "-", "_", -1)
	}
	return s
}

func (self *parserCompat) get_opt_help_opthelp(opt *ExtKeyParse) string {
	return self.get_help_info(opt)
}
