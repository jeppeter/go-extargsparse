package extargsparse

import (
	"encoding/json"
	"fmt"
)

type parserCompat struct {
	logObject
	KeyCls       *extKeyParse
	CmdName      string
	CmdOpts      []*extKeyParse
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

func newParserCompat(keycls *extKeyParse, opt *ExtArgsOptions) *parserCompat {
	var vmap map[string]interface{}
	var js string
	var err error
	self := &parserCompat{logObject: *newLogObject("extargsparse")}
	if keycls != nil {
		assert_test(keycls.IsCmd(), "%s must be cmd", keycls.Format())
		self.KeyCls = keycls
		self.CmdName = keycls.CmdName()
		self.CmdOpts = make([]*extKeyParse, 0)
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
		self.CmdOpts = make([]*extKeyParse, 0)
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
