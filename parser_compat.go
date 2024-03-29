package extargsparse

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type parserCompat struct {
	logObject
	KeyCls        *ExtKeyParse
	CmdName       string
	CmdOpts       []*ExtKeyParse
	SubCommands   []*parserCompat
	HelpInfo      string
	CallFunction  string
	ScreenWidth   int
	Epilog        string
	Description   string
	Prog          string
	Usage         string
	Version       string
	FunctionUpper bool
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
		self.SubCommands = make([]*parserCompat, 0)
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
		self.SubCommands = make([]*parserCompat, 0)
		self.HelpInfo = ""
		self.CallFunction = ""
	}
	self.ScreenWidth = 80
	if opt != nil && opt.GetValue(OPT_SCREEN_WIDTH) != nil {
		self.ScreenWidth = opt.GetValue(OPT_SCREEN_WIDTH).(int)
	}

	self.FunctionUpper = true
	if opt != nil {
		self.FunctionUpper = opt.GetBool(OPT_FUNC_UPPER_CASE)
	}

	if self.ScreenWidth < 40 {
		self.ScreenWidth = 40
	}
	if opt != nil {
		self.Epilog = opt.GetString(OPT_EPILOG)
		self.Description = opt.GetString(OPT_DESCRIPTION)
		self.Prog = opt.GetString(OPT_PROG)
		self.Usage = opt.GetString(OPT_USAGE)
		self.Version = opt.GetString(OPT_VERSION)
	} else {
		self.Epilog = ""
		self.Description = ""
		self.Prog = ""
		self.Usage = ""
		self.Version = "0.0.1"
	}
	return self
}

func (self *parserCompat) get_help_info(keycls *ExtKeyParse) string {
	var s string
	var err error
	var helpfunc func(kl *ExtKeyParse) string
	assert_test(keycls != nil, "must no be null")
	if keycls.Attr("opthelp") != "" {
		/*now it is the help function ,so we call this*/
		self.Trace("call [%s] upper[%v]", keycls.Attr("opthelp"), self.FunctionUpper)
		err = self.GetFuncPtr(self.FunctionUpper, keycls.Attr("opthelp"), &helpfunc)
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
			s += fmt.Sprintf("%s set default(%v)", keycls.Optdest(), keycls.Value())
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
	if opt != nil {
		s += opt.Longopt()
		if len(opt.Shortopt()) > 0 {
			s += fmt.Sprintf("|%s", opt.Shortopt())
		}
	}
	return s
}

func (self *parserCompat) get_opt_help_optexpr(opt *ExtKeyParse) string {
	var s string = ""
	if opt != nil && opt.TypeName() != "bool" && opt.TypeName() != "args" &&
		opt.TypeName() != "dict" && opt.TypeName() != "help" {
		s = opt.VarName()
		s = strings.Replace(s, "-", "_", -1)
	}
	return s
}

func (self *parserCompat) get_opt_help_opthelp(opt *ExtKeyParse) string {
	return self.get_help_info(opt)
}

func (self *parserCompat) get_cmd_help_cmdname() string {
	var s string = ""
	if len(self.CmdName) > 0 {
		s = fmt.Sprintf("[%s]", self.CmdName)
	}
	return s
}

func (self *parserCompat) get_cmd_help_cmdhelp() string {
	var s string = ""
	if len(self.HelpInfo) > 0 {
		s = self.HelpInfo
	}
	return s
}

func (self *parserCompat) GetHelpSize(hs *helpSize, recursive int) *helpSize {
	var curopt *ExtKeyParse
	var chldparser *parserCompat
	if hs == nil {
		hs = newHelpSize()
	}

	hs.SetValue("cmdnamesize", len(self.get_cmd_help_cmdname())+1)
	hs.SetValue("cmdhelpsize", len(self.get_cmd_help_cmdhelp())+1)

	for _, curopt = range self.CmdOpts {
		if curopt.TypeName() == "args" {
			continue
		}
		hs.SetValue("optnamesize", len(self.get_opt_help_optname(curopt))+1)
		hs.SetValue("optexprsize", len(self.get_opt_help_optexpr(curopt))+1)
		hs.SetValue("opthelpsize", len(self.get_opt_help_opthelp(curopt))+1)
	}

	if recursive != 0 {
		for _, chldparser = range self.SubCommands {
			if recursive > 0 {
				hs = chldparser.GetHelpSize(hs, recursive-1)
			} else {
				hs = chldparser.GetHelpSize(hs, recursive)
			}
		}
	}

	for _, chldparser = range self.SubCommands {
		hs.SetValue("cmdnamesize", len(chldparser.get_cmd_help_cmdname())+1)
		hs.SetValue("cmdhelpsize", len(chldparser.get_cmd_help_cmdhelp())+1)
	}

	return hs
}

func (self *parserCompat) get_indent_string(s string, indentsize int, maxsize int) string {
	var rets string = ""
	var curs string = ""
	var ncurs string
	var i int
	var j int
	i = 0
	curs = ""
	for i = 0; i < indentsize; i++ {
		curs += " "
	}

	for j = 0; j < len(s); j++ {
		if (s[j:(j+1)] == " " || s[j:(j+1)] == "\t") && len(curs) >= maxsize {
			rets += curs + "\n"
			curs = ""
			for i = 0; i < indentsize; i++ {
				curs += " "
			}
			continue
		}
		curs += s[j:(j + 1)]
	}

	ncurs = strings.Trim(curs, "\t ")
	if len(ncurs) > 0 {
		rets += strings.TrimRight(curs, "\t ") + "\n"
	}
	return rets
}

func (self *parserCompat) GetHelpInfo(hs *helpSize, parentcmds []*parserCompat) string {
	var rets string = ""
	var rootcmds *parserCompat = nil
	var curcmd *parserCompat
	var curopt *ExtKeyParse
	var cmdname string
	var cmdhelp string
	var curs string
	var curint int
	var optname string
	var optexpr string
	var opthelp string
	if hs == nil {
		hs = self.GetHelpSize(hs, 0)
	}
	if len(self.Usage) > 0 {
		rets += fmt.Sprintf("%s", self.Usage)
	} else {
		rootcmds = self
		curcmd = self
		if len(parentcmds) > 0 {
			rootcmds = parentcmds[0]
		}

		if len(rootcmds.Prog) > 0 {
			rets += fmt.Sprintf("%s", rootcmds.Prog)
		} else {
			rets += fmt.Sprintf("%s", os.Args[0])
		}

		if len(parentcmds) > 0 {
			for _, curcmd = range parentcmds {
				rets += fmt.Sprintf(" %s", curcmd.CmdName)
			}
		}
		rets += fmt.Sprintf(" %s", self.CmdName)

		if len(self.HelpInfo) > 0 {
			rets += fmt.Sprintf(" %s", self.HelpInfo)
		} else {
			if len(self.CmdOpts) > 0 {
				rets += fmt.Sprintf(" [OPTIONS]")
			}

			if len(self.SubCommands) > 0 {
				rets += fmt.Sprintf(" [SUBCOMMANDS]")
			}

			for _, curopt = range self.CmdOpts {
				if curopt.TypeName() == "args" {
					switch curopt.Nargs().(type) {
					case string:
						curs = curopt.Nargs().(string)
						if curs == "+" {
							rets += fmt.Sprintf(" args...")
						} else if curs == "*" {
							rets += fmt.Sprintf(" [args...]")

						} else if curs == "?" {
							rets += fmt.Sprintf(" arg")
						}
					case int:
						curint = curopt.Nargs().(int)
						if curint > 1 {
							rets += " args..."
						} else if curint == 1 {
							rets += " arg"
						} else {
							rets += ""
						}
					default:
						assert_test(false == true, "%s nargs not valid", curopt.Format())
					}
				}
			}
		}

		rets += "\n"
	}

	if len(self.Description) > 0 {
		rets += fmt.Sprintf("%s\n", self.Description)
	}

	rets += "\n"

	if len(self.CmdOpts) > 0 {
		rets += "[OPTIONS]\n"

		for _, curopt = range self.CmdOpts {
			if curopt.TypeName() == "args" {
				continue
			}
			curs = ""
			curs += "    "
			optname = self.get_opt_help_optname(curopt)
			optexpr = self.get_opt_help_optexpr(curopt)
			opthelp = self.get_opt_help_opthelp(curopt)
			curs += fmt.Sprintf("%-*s %-*s %-*s\n", hs.GetValue("optnamesize"), optname, hs.GetValue("optexprsize"), optexpr, hs.GetValue("opthelpsize"), opthelp)
			if len(curs) < self.ScreenWidth {
				curs = ""
				curs += "    "
				curs += fmt.Sprintf("%-*s %-*s %-*s", hs.GetValue("optnamesize"), optname, hs.GetValue("optexprsize"), optexpr, hs.GetValue("opthelpsize"), opthelp)
				curs = strings.TrimRight(curs, " \t")
				curs += "\n"
				rets += curs
			} else {
				curs = ""
				curs += "    "
				curs += fmt.Sprintf("%-*s %-*s", hs.GetValue("optnamesize"), optname, hs.GetValue("optexprsize"), optexpr)
				rets += curs + "\n"
				if self.ScreenWidth > 60 {
					rets += self.get_indent_string(opthelp, 20, self.ScreenWidth)
				} else {
					rets += self.get_indent_string(opthelp, 15, self.ScreenWidth)
				}
			}
		}
	}

	if len(self.SubCommands) > 0 {
		rets += "\n"
		rets += "[SUBCOMMANDS]\n"
		for _, curcmd = range self.SubCommands {
			cmdname = curcmd.get_cmd_help_cmdname()
			cmdhelp = curcmd.get_cmd_help_cmdhelp()
			curs = ""
			curs += "    "
			curs += fmt.Sprintf("%-*s %-*s", hs.GetValue("cmdnamesize"), cmdname, hs.GetValue("cmdhelpsize"), cmdhelp)
			if len(curs) < self.ScreenWidth {
				rets += curs
				rets += "\n"
			} else {
				curs = ""
				curs += "    "
				curs += fmt.Sprintf("%-*s", hs.GetValue("cmdnamesize"), cmdname)
				rets += fmt.Sprintf("%s\n", curs)
				if self.ScreenWidth > 60 {
					rets += self.get_indent_string(cmdhelp, 20, self.ScreenWidth)
				} else {
					rets += self.get_indent_string(cmdhelp, 15, self.ScreenWidth)
				}
			}
		}
	}

	if len(self.Epilog) > 0 {
		rets += "\n"
		rets += fmt.Sprintf("\n%s\n", self.Epilog)
	}

	self.Trace("%s", rets)
	return rets
}

func (self *parserCompat) Format() string {
	var s string
	var i int = 0
	var curcmd *parserCompat
	var curopt *ExtKeyParse
	s = fmt.Sprintf("@%s|", self.CmdName)
	if self.KeyCls != nil {
		s += fmt.Sprintf("%s|", self.KeyCls.Format())
	} else {
		s += "nil|"
	}
	if len(self.SubCommands) > 0 {
		s += fmt.Sprintf("subcommands[%d]<", len(self.SubCommands))
		for i, curcmd = range self.SubCommands {
			if i > 0 {
				s += "|"
			}
			s += fmt.Sprintf("%s", curcmd.CmdName)
		}
		s += ">"
	}
	if len(self.CmdOpts) > 0 {
		s += fmt.Sprintf("cmdopts[%d]<", len(self.CmdOpts))
		for _, curopt = range self.CmdOpts {
			s += fmt.Sprintf("%s", curopt.Format())
		}
		s += ">"
	}
	return s
}
