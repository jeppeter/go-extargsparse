package extargsparse

import (
	"fmt"
	"strings"
)

type parseState struct {
	logObject
	cmdpaths      []*parserCompat
	curidx        int
	curcharidx    int
	shortcharargs int
	longargs      int
	keyidx        int
	validx        int
	args          []string
	ended         int
	longprefix    string
	shortprefix   string
	bundlemode    bool
	parseall      bool
	leftargs      []string
}

func newParseState(args []string, maincmd *parserCompat, optattr *ExtArgsOptions) *parseState {
	var err error
	assert_test(maincmd != nil, "maincmd can not set nil")
	self := &parseState{logObject: *newLogObject("extargsparse")}
	if optattr == nil {
		optattr, err = NewExtArgsOptions("{}")
		assert_test(err == nil, "parser option error [%s]", err.Error())
	}

	self.cmdpaths = make([]*parserCompat, 0)
	self.cmdpaths = append(self.cmdpaths, maincmd)
	self.curidx = 0
	self.curcharidx = -1
	self.shortcharargs = -1
	self.longargs = -1
	self.keyidx = -1
	self.validx = -1
	self.args = args
	self.ended = 0
	self.longprefix = optattr.GetString("longprefix")
	self.shortprefix = optattr.GetString("shortprefix")
	if len(self.shortprefix) == 0 || len(self.longprefix) == 0 ||
		self.shortprefix == self.longprefix {
		self.bundlemode = true
	} else {
		self.bundlemode = false
	}
	self.parseall = optattr.GetBool("parseall")
	self.leftargs = make([]string, 0)
	return self
}

func (self *parseState) FormatCmdnamePath(curparser []*parserCompat) string {
	var c *parserCompat
	var cmdname string = ""
	if len(curparser) == 0 {
		curparser = self.cmdpaths
	}

	for _, c = range curparser {
		if len(cmdname) > 0 {
			cmdname += "."
		}
		cmdname += c.CmdName
	}
	return cmdname
}

func (self *parseState) find_sub_command(name string) *ExtKeyParse {
	var cmdparent *parserCompat
	var c *parserCompat
	cmdparent = self.cmdpaths[len(self.cmdpaths)-1]
	for _, c = range cmdparent.SubCommands {
		if c.CmdName == name {
			self.cmdpaths = append(self.cmdpaths, c)
			return c.KeyCls
		}
	}
	return nil
}

func (self *parseState) AddParseArgs(nargs int) error {
	if self.curcharidx >= 0 {
		if nargs > 0 && self.shortcharargs > 0 {
			return fmt.Errorf("%s", format_error("[%s] already set args", self.args[self.curidx]))
		}
		if self.shortcharargs < 0 {
			self.shortcharargs = 0
		}
		self.shortcharargs += nargs
	} else {
		if self.longargs > 0 {
			return fmt.Errorf("%s", format_error("[%s] not handled ", self.args[self.curidx]))
		}
		if self.longargs < 0 {
			self.longargs = 0
		}
		self.longargs += nargs
	}
	return nil
}

func (self *parseState) find_key_cls() (retkey *ExtKeyParse, err error) {
	var oldcharidx, oldidx, idx int
	var c, curarg string
	var curch string
	var cmd *parserCompat
	var opt *ExtKeyParse
	var keycls *ExtKeyParse
	retkey = nil
	err = nil

	if self.ended > 0 {
		return
	}

	if self.longargs >= 0 {
		assert_test(self.curcharidx < 0, "curcharidx[%d]", self.curcharidx)
		self.curidx += self.longargs
		assert_test(len(self.args) >= self.curidx, "len[%d] < [%d]", len(self.args), self.curidx)
		self.longargs = -1
		self.validx = -1
		self.keyidx = -1
	}

	oldcharidx = self.curcharidx
	oldidx = self.curidx
	self.Trace("oldcharidx [%d] oldidx [%d]", oldcharidx, oldidx)

	if oldidx >= len(self.args) {
		self.curidx = oldidx
		self.curcharidx = -1
		self.shortcharargs = -1
		self.longargs = -1
		self.keyidx = -1
		self.validx = -1
		self.ended = 1
		return
	}

	if oldcharidx >= 0 {
		c = self.args[oldidx]
		if len(c) <= oldcharidx {
			oldidx += 1
			self.Trace("oldidx [%d] [%s] [%d]", oldidx, c, oldcharidx)
			if self.shortcharargs > 0 {
				oldidx += self.shortcharargs
			}
			self.Trace("oldidx [%d] __shortcharargs [%d]", oldidx, self.shortcharargs)
			self.curidx = oldidx
			self.curcharidx = -1
			self.shortcharargs = -1
			self.keyidx = -1
			self.validx = -1
			self.longargs = -1
			return self.find_key_cls()
		}
		curch = c[oldcharidx:(oldcharidx + 1)]
		self.Trace("argv[%d][%d] %s", oldidx, oldcharidx, curch)
		idx = len(self.cmdpaths) - 1
		for idx >= 0 {
			cmd = self.cmdpaths[idx]
			for _, opt = range cmd.CmdOpts {
				if !opt.IsFlag() {
					continue
				}

				if opt.FlagName() == "$" {
					continue
				}

				if len(opt.ShortFlag()) != 0 {
					if opt.ShortFlag() == curch {
						self.keyidx = oldidx
						self.validx = (oldidx + 1)
						self.curidx = oldidx
						self.curcharidx = (oldcharidx + 1)
						self.Info("%s validx [%s]", opt.Format(), self.validx)
						retkey = opt
						err = nil
						return
					}
				}
			}
			idx -= 1
		}
		retkey = nil
		err = fmt.Errorf("%s", format_error("can not parse (%s)", self.args[oldidx]))
		return
	} else {
		if !self.bundlemode {
			curarg = self.args[oldidx]
			if self.longprefix != "" && strings.HasPrefix(curarg, self.longprefix) {
				if curarg == self.longprefix {
					/*the end of the list of args*/
					self.keyidx = -1
					self.curidx = oldidx + 1
					self.curcharidx = -1
					self.validx = (oldidx + 1)
					self.shortcharargs = -1
					self.longargs = -1
					self.ended = 1
					for idx = self.curidx; idx < len(self.args); idx++ {
						self.leftargs = append(self.leftargs, self.args[idx])
					}
					retkey = nil
					err = nil
					return
				}
				idx = len(self.cmdpaths) - 1
				for idx >= 0 {
					cmd = self.cmdpaths[idx]
					for _, opt = range cmd.CmdOpts {
						if !opt.IsFlag() {
							continue
						}
						if opt.FlagName() == "$" {
							continue
						}
						self.Info("[%d]longopt %s curarg %s", idx, opt.Longopt(), curarg)
						if opt.Longopt() == curarg {
							self.keyidx = oldidx
							oldidx += 1
							self.validx = oldidx
							self.shortcharargs = -1
							self.longargs = -1
							self.Info("oldidx %d (len %d)", oldidx, len(self.args))
							self.curidx = oldidx
							self.curcharidx = -1
							retkey = opt
							err = nil
							return
						}
					}
					idx -= 1
				}
				retkey = nil
				err = fmt.Errorf("%s", format_error("can not parse (%s)", self.args[oldidx]))
				return
			} else if len(self.shortprefix) > 0 && strings.HasPrefix(curarg, self.shortprefix) {
				if curarg == self.shortprefix {
					if self.parseall {
						self.leftargs = append(self.leftargs, curarg)
						oldidx += 1
						self.curidx = oldidx
						self.curcharidx = -1
						self.longargs = -1
						self.shortcharargs = -1
						self.keyidx = -1
						self.validx = -1
						return self.find_key_cls()
					} else {
						self.ended = 1
						for idx = oldidx; idx < len(self.args); idx += 1 {
							self.leftargs = append(self.leftargs, self.args[idx])
						}
						self.validx = oldidx
						self.keyidx = -1
						self.curidx = oldidx
						self.curcharidx = -1
						self.shortcharargs = -1
						self.longargs = -1
						retkey = nil
						err = nil
						return
					}
				}
				oldcharidx = len(self.shortprefix)
				self.curidx = oldidx
				self.curcharidx = oldcharidx
				return self.find_key_cls()
			}
		} else {
			/*
				not bundle mode ,it means that the long prefix and short prefix are the same
				so we should test one by one
				first to check for the long opt
			*/
			idx = len(self.cmdpaths) - 1
			curarg = self.args[oldidx]
			for idx >= 0 {
				cmd = self.cmdpaths[idx]
				for _, opt = range cmd.CmdOpts {
					if !opt.IsFlag() {
						continue
					}
					if opt.FlagName() == "$" {
						continue
					}
					self.Info("[%d](%s) curarg [%s]", idx, opt.Longopt(), curarg)
					if opt.Longopt() == curarg {
						self.keyidx = oldidx
						self.validx = oldidx + 1
						self.shortcharargs = -1
						self.longargs = -1
						self.Info("oldidx %d (len %d)", oldidx, len(self.args))
						self.curidx = (oldidx + 1)
						self.curcharidx = -1
						retkey = opt
						err = nil
						return
					}
				}
				idx -= 1
			}
			idx = len(self.cmdpaths) - 1
			for idx >= 0 {
				cmd = self.cmdpaths[idx]
				for _, opt = range cmd.CmdOpts {
					if !opt.IsFlag() {
						continue
					}
					if opt.FlagName() == "$" {
						continue
					}
					self.Info("[%d](%s) curarg [%s]", idx, opt.Shortopt(), curarg)
					if len(opt.Shortopt()) > 0 && opt.Shortopt() == curarg {
						self.keyidx = oldidx
						self.validx = (oldidx + 1)
						self.shortcharargs = -1
						self.longargs = -1
						self.curidx = oldidx
						self.curcharidx = len(opt.Shortopt())
						retkey = opt
						err = nil
						return
					}
				}
				idx -= 1
			}
		}
	}

	keycls = self.find_sub_command(self.args[oldidx])
	if keycls != nil {
		self.Info("find %s", self.args[oldidx])
		self.keyidx = oldidx
		self.curidx = (oldidx + 1)
		self.validx = (oldidx + 1)
		self.curcharidx = -1
		self.shortcharargs = -1
		self.longargs = -1
		retkey = keycls
		err = nil
		return
	}

	if self.parseall {
		self.leftargs = append(self.leftargs, self.args[oldidx])
		oldidx += 1
		self.keyidx = -1
		self.validx = oldidx
		self.curidx = oldidx
		self.curcharidx = -1
		self.shortcharargs = -1
		self.longargs = -1
		return self.find_key_cls()
	} else {
		self.ended = 1
		for idx = oldidx; idx < len(self.args); idx++ {
			self.leftargs = append(self.leftargs, self.args[idx])
		}
		self.keyidx = -1
		self.curidx = oldidx
		self.curcharidx = -1
		self.shortcharargs = -1
		self.longargs = -1
		retkey = nil
		err = nil
		return
	}

	err = nil
	retkey = nil
	return
}

func (self *parseState) StepOne() (validx int, optval interface{}, keycls *ExtKeyParse, err error) {
	if self.ended != 0 {
		validx = self.curidx
		optval = self.leftargs
		keycls = nil
		err = nil
		return
	}

	keycls, err = self.find_key_cls()
	if err != nil {
		return
	}

	if keycls == nil {
		validx = self.curidx
		optval = self.leftargs
		err = nil
		return
	}

	if !keycls.IsCmd() {
		optval = keycls.Optdest()
	} else if keycls.IsCmd() {
		optval = self.FormatCmdnamePath(self.cmdpaths)
	}
	validx = self.validx
	err = nil
	return
}

func (self *parseState) GetCmdPaths() []*parserCompat {
	return self.cmdpaths
}
