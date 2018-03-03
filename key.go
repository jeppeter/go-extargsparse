package extargsparse

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	ATTR_SPLIT       = "split"
	ATTR_SPLIT_EQUAL = "split="
)

var helpExpr *regexp.Regexp = regexp.MustCompile(`(?i)##([^\#\!]+)##$`)
var cmdExpr *regexp.Regexp = regexp.MustCompile(`(?i)^([^\#\<\>\+\$\!]+)`)
var prefixExpr *regexp.Regexp = regexp.MustCompile(`(?i)\+([a-zA-Z]+[a-zA-Z_\-0-9]*)`)
var funcExpr *regexp.Regexp = regexp.MustCompile(`(?i)<([^\<\>\#\$\| \t\!]+)>`)
var flagExpr *regexp.Regexp = regexp.MustCompile(`(?i)^([a-zA-Z_\|\?\-]+[a-zA-Z_0-9\|\?\-]*)`)
var mustFlagExpr *regexp.Regexp = regexp.MustCompile(`(?i)^\$([a-zA-Z_\|\?]+[a-zA-Z_0-9\|\?\-]*)`)
var attrExpr *regexp.Regexp = regexp.MustCompile(`\!([^\<\>\$!\#\|]+)\!`)
var flagwords = []string{"flagname", "helpinfo", "shortflag", "nargs", "varname"}
var flagspecial = []string{"value", "prefix"}

func parseAttr(attr string) (kattr map[string]string, err error) {
	var splitchar string = ";"
	var splitstrings []string
	var splitexpr *regexp.Regexp
	var equalexpr *regexp.Regexp
	var vk []string
	var curs string

	kattr = nil
	err = nil
	if strings.HasPrefix(attr, ATTR_SPLIT_EQUAL) {
		splitchar = attr[len(ATTR_SPLIT_EQUAL):(len(ATTR_SPLIT_EQUAL) + 1)]
		switch splitchar {
		case "\\":
			splitchar = "\\\\"
		case ".":
			splitchar = "\\."
		case "/":
			splitchar = "/"
		case ":":
			splitchar = ":"
		case "+":
			splitchar = "\\+"
		case "@":
			splitchar = "@"
		default:
			return nil, fmt.Errorf(format_error("unknown splitchar [%s]", splitchar))
		}
	}
	splitexpr, err = regexp.Compile(splitchar)
	if err != nil {
		return
	}
	equalexpr, err = regexp.Compile("=")
	if err != nil {
		return
	}

	kattr = make(map[string]string)
	splitstrings = splitexpr.Split(attr, -1)
	for _, curs = range splitstrings {
		if strings.HasPrefix(curs, ATTR_SPLIT_EQUAL) || curs == "" {
			continue
		}
		vk = equalexpr.Split(curs, 2)
		if len(vk) < 2 {
			continue
		}
		kattr[vk[0]] = vk[1]
	}

	err = nil
	return
}

func setAttr(attr interface{}) (kattr map[string]string, err error) {
	var k string
	var v interface{}
	var vmap map[string]interface{}
	var vstr string
	kattr = make(map[string]string)

	switch attr.(type) {
	case map[string]interface{}:
		vmap = attr.(map[string]interface{})
		for k, v = range vmap {
			if strings.ToLower(k) == ATTR_SPLIT || k == "" {
				continue
			}
			switch v.(type) {
			case string:
				vstr = v.(string)
			default:
				vstr = fmt.Sprintf("%v", v)
			}
			kattr[k] = vstr
		}
	default:
		return kattr, fmt.Errorf(format_error("not valid type [%s]", reflect.TypeOf(attr)))
	}

	err = nil
	return
}

type ExtKeyParse struct {
	/*this are the inner member*/
	longPrefix  string
	shortPrefix string
	noChange    bool
	value       interface{}
	prefix      string
	flagName    string
	helpInfo    string
	shortFlag   string
	nargs       interface{}
	varName     string
	cmdName     string
	function    string
	origKey     string
	isCmd       bool
	isFlag      bool
	typeName    string
	attr        map[string]string
}

func (self *ExtKeyParse) getType(value interface{}) string {
	var bret bool

	switch value.(type) {
	case string:
		return "string"
	case bool:
		return "bool"
	case float32:
		_, bret = isFloatToInt(float64(value.(float32)))
		if bret {
			return "int"
		}
		return "float"
	case float64:
		_, bret = isFloatToInt(value.(float64))
		if bret {
			return "int"
		}
		return "float"
	case int:
		return "int"
	case int64:
		return "int"
	case int32:
		return "int"
	case uint:
		return "int"
	case uint32:
		return "int"
	case uint64:
		return "int"
	case map[string]interface{}:
		return "dict"
	case []interface{}:
		return "list"
	}

	if value == nil {
		return "string"
	}

	s := fmt.Sprintf("not valid type [%v]", value)
	panic(s)
	return ""
}

func (self *ExtKeyParse) getValue(value interface{}) interface{} {
	var iv int
	var bret bool
	if value == nil {
		return nil
	}
	switch value.(type) {
	case float64:
		iv, bret = isFloatToInt(value.(float64))
		if bret {
			return iv
		}
	case float32:
		iv, bret = isFloatToInt(float64(value.(float32)))
		if bret {
			return iv
		}
		return float64(value.(float32))
	case int16:
		return int(value.(int16))
	case int32:
		return int(value.(int32))
	case int64:
		return int(value.(int64))
	case uint16:
		return int(value.(uint16))
	case uint32:
		return int(value.(uint32))
	case uint64:
		return int(value.(uint64))
	}
	return value
}

func (self *ExtKeyParse) setFlag(prefix, key string, value interface{}) error {
	var vmap map[string]interface{}
	var findvalue bool = false
	var v interface{}
	var k, k2, vstr string
	var isflagwords bool
	var isflagspecial bool
	var newprefix string
	var err error
	self.isFlag = true
	self.isCmd = false
	self.origKey = key
	vmap = value.(map[string]interface{})
	for k, _ = range vmap {
		if k == "value" {
			findvalue = true
			break
		}
	}

	if !findvalue {
		self.typeName = "string"
		self.value = nil
	}

	for k, v = range vmap {
		isflagwords = false
		isflagspecial = false
		for _, k2 = range flagwords {
			if k2 == k {
				isflagwords = true
				break
			}
		}

		if isflagwords {
			switch v.(type) {
			case string:
				vstr = v.(string)
			default:
				if k != "nargs" {
					if v != nil {
						return fmt.Errorf(format_error("not value type"))
					} else {
						vstr = ""
					}
				}
			}
			switch k {
			case "flagname":
				self.flagName = vstr
			case "helpinfo":
				self.helpInfo = vstr
			case "shortflag":
				self.shortFlag = vstr
			case "nargs":
				switch v.(type) {
				case string:
					self.nargs = v
				default:
					vstr = fmt.Sprintf("%v", v)
					self.nargs, err = strconv.Atoi(vstr)
				}
			case "varname":
				self.varName = vstr
			default:
				return fmt.Errorf(format_error("[%s] not recognize", k))
			}

		} else {
			for _, k2 = range flagspecial {
				if k2 == k {
					isflagspecial = true
					break
				}
			}

			if isflagspecial {
				switch k {
				case "prefix":
					newprefix = ""
					vstr = ""
					switch v.(type) {
					case string:
						vstr = v.(string)
					}
					if len(prefix) > 0 && len(vstr) > 0 {
						newprefix = fmt.Sprintf("%s_%s", prefix, vstr)
					} else if len(vstr) > 0 {
						newprefix = fmt.Sprintf("%s", vstr)
					}
					self.prefix = newprefix
				case "value":
					self.value = self.getValue(v)
					self.typeName = self.getType(v)
				default:
					return fmt.Errorf(format_error("[%s] not valid specialword", k))
				}
			} else if k == "attr" {
				self.attr, err = setAttr(v)
				if err != nil {
					self.attr = make(map[string]string)
					return err
				}
			}
		}
	}

	if len(self.prefix) == 0 && len(prefix) > 0 {
		self.prefix = prefix
	}
	return nil
}

func (self *ExtKeyParse) Optdest() string {
	var optdest string
	if !self.isFlag || len(self.flagName) == 0 || self.typeName == "args" {
		s := fmt.Sprintf("can not set (%s) optdest", self.origKey)
		panic(s)
	}

	optdest = ""
	if len(self.prefix) > 0 {
		optdest += fmt.Sprintf("%s_", self.prefix)
	}

	optdest += self.flagName

	if !self.noChange {
		optdest = strings.ToLower(optdest)
	}
	optdest = strings.Replace(optdest, "-", "_", -1)
	return optdest
}

func (self *ExtKeyParse) ShortFlag() string {
	return self.shortFlag
}

func (self *ExtKeyParse) IsFlag() bool {
	return self.isFlag
}

func (self *ExtKeyParse) IsCmd() bool {
	return self.isCmd
}

func (self *ExtKeyParse) TypeName() string {
	return self.typeName
}

func (self *ExtKeyParse) validate() error {
	if self.isFlag {
		assert_test(!self.isCmd, "cmdmode setted")
		if len(self.function) > 0 {
			return fmt.Errorf(format_error("(%s) can not accept function", self.origKey))
		}

		if self.typeName == "dict" && len(self.flagName) > 0 {
			return fmt.Errorf(format_error("(%s) flag can not accept dict", self.origKey))
		}

		if self.typeName != self.getType(self.value) && self.typeName != "count" && self.typeName != "help" && self.typeName != "jsonfile" {
			return fmt.Errorf(format_error("(%s) value (%v) not match type (%s)", self.origKey, self.value, self.typeName))
		}

		if len(self.flagName) == 0 {
			if len(self.prefix) == 0 {
				return fmt.Errorf(format_error("(%s) should at least for prefix", self.origKey))
			}
			self.typeName = "prefix"
			if self.getType(self.value) != "dict" {
				return fmt.Errorf(format_error("(%s) should used dict to make prefix", self.origKey))
			}
			if len(self.helpInfo) > 0 {
				return fmt.Errorf(format_error("(%s) should not have help info", self.origKey))
			}
			if len(self.shortFlag) > 0 {
				return fmt.Errorf(format_error("(%s) should not set shortflag", self.origKey))
			}
		} else if self.flagName == "$" {
			self.typeName = "args"
			if len(self.shortFlag) > 0 {
				return fmt.Errorf(format_error("(%s) can not set shortflag for args", self.origKey))
			}
		} else {
			if len(self.flagName) < 0 {
				return fmt.Errorf(format_error("(%s) can not accept (%s)short flag in flagname", self.origKey, self.flagName))
			}
		}

		if len(self.shortFlag) > 1 {
			return fmt.Errorf(format_error("(%s) can not accept (%s) for shortflag", self.origKey, self.shortFlag))
		}

		if self.typeName == "bool" {
			if self.nargs != nil && self.nargs.(int) != 0 {
				return fmt.Errorf(format_error("bool type (%s) can not accept not 0 nargs", self.origKey))
			}
			self.nargs = 0
		} else if self.typeName == "help" {
			if self.nargs != nil && self.nargs.(int) != 0 {
				return fmt.Errorf(format_error("help type (%s) can not accept not 0 nargs", self.origKey))
			}
			self.nargs = 0
		} else if self.typeName != "prefix" && self.flagName != "$" && self.typeName != "count" {
			if self.typeName != "$" && self.nargs != nil && self.nargs.(int) != 1 {
				return fmt.Errorf(format_error("(%s)only $ can accept nargs option", self.origKey))
			}
			self.nargs = 1
		} else {
			if self.flagName == "$" && self.nargs == nil {
				self.nargs = "*"
			}
		}
	} else {
		if len(self.cmdName) == 0 {
			return fmt.Errorf(format_error("(%s) not set cmdname", self.origKey))
		}

		if len(self.shortFlag) > 0 {
			return fmt.Errorf(format_error("(%s) has shortflag (%s)", self.origKey, self.shortFlag))
		}

		if self.nargs != nil {
			return fmt.Errorf(format_error("(%s) has nargs (%v)", self.origKey, self.nargs))
		}

		if self.typeName != "dict" {
			return fmt.Errorf(format_error("(%s) command must be dict", self.origKey))
		}

		if len(self.prefix) == 0 {
			self.prefix += self.cmdName
		}
		self.typeName = "command"
	}

	if self.isFlag && len(self.varName) == 0 && len(self.flagName) > 0 {
		if self.flagName != "$" {
			self.varName = self.Optdest()
		} else {
			if len(self.prefix) > 0 {
				self.varName = "subnargs"
			} else {
				self.varName = "args"
			}
		}
	}

	return nil
}

func (self *ExtKeyParse) parse(prefix string, key string, value interface{}, isflag bool, ishelp bool, isjsonfile bool, longprefix string, shortprefix string, nochange bool) error {
	var flagmode bool = false
	var cmdmode bool = false
	var flags string = ""
	var ok int
	var matchstrings []string
	var sarr []string
	var hexpr *regexp.Regexp
	var newprefix string = ""
	var err error
	self.origKey = key
	self.longPrefix = longprefix
	self.shortPrefix = shortprefix
	self.noChange = nochange

	/*now to test whether it is the flag one*/
	if strings.Contains(self.origKey, "$") {
		if self.origKey[0:1] != "$" {
			return fmt.Errorf(format_error("(%s) not right format for ($)", self.origKey))
		}
		ok = 1
		if strings.Contains(self.origKey[1:], "$") {
			ok = 0
		}
		if ok != 1 {
			return fmt.Errorf(format_error("(%s) has ($) more than one", self.origKey))
		}
	}

	if isflag || ishelp || isjsonfile {
		matchstrings = flagExpr.FindStringSubmatch(self.origKey)
		if len(matchstrings) > 1 {
			flags = matchstrings[1]
		}
		if len(flags) == 0 {
			matchstrings = mustFlagExpr.FindStringSubmatch(self.origKey)
			if len(matchstrings) > 1 {
				flags = matchstrings[1]
			}
		}

		if len(flags) == 0 && self.origKey[0:1] == "$" {
			self.flagName = "$"
			self.isFlag = true
			flagmode = true
		}

		if len(flags) > 0 {
			if strings.Contains(flags, "|") {
				hexpr = regexp.MustCompile(`\|`)
				sarr = hexpr.Split(flags, -1)
				if len(sarr) > 2 || len(sarr[1]) != 1 || len(sarr[0]) <= 1 {
					return fmt.Errorf(format_error("(%s) (%s)flag only accept (longop|l) format", self.origKey, flags))
				}
				self.flagName = sarr[0]
				self.shortFlag = sarr[1]
			} else {
				self.flagName = flags
			}
			self.isFlag = true
			flagmode = true
		}
	} else {
		matchstrings = mustFlagExpr.FindStringSubmatch(self.origKey)
		if len(matchstrings) > 1 {
			flags = matchstrings[1]
			if strings.Contains(flags, "|") {
				hexpr = regexp.MustCompile(`\|`)
				sarr = hexpr.Split(flags, -1)
				if len(sarr) > 2 || len(sarr[1]) != 1 || len(sarr[0]) <= 1 {
					return fmt.Errorf(format_error("(%s) (%s)flag only accept (longop|l) format", self.origKey, flags))
				}
				self.flagName = sarr[0]
				self.shortFlag = sarr[1]
			} else {
				self.flagName = flags
			}
			self.isFlag = true
			flagmode = true
		} else if len(self.origKey) > 0 && self.origKey[0:1] == "$" {
			self.flagName = "$"
			self.isFlag = true
			flagmode = true
		} else {
			matchstrings = cmdExpr.FindStringSubmatch(self.origKey)
			if len(matchstrings) > 1 {
				assert_test(!flagmode, "flagmode set")
				if strings.Contains(matchstrings[1], "|") {
					flags = matchstrings[1]
					hexpr = regexp.MustCompile(`\|`)
					sarr = hexpr.Split(flags, -1)
					if len(sarr) > 2 || len(sarr[1]) != 1 || len(sarr[0]) <= 1 {
						return fmt.Errorf(format_error("(%s) (%s)flag only accept (longop|l) format", self.origKey, flags))
					}
					self.flagName = sarr[0]
					self.shortFlag = sarr[1]
					self.isFlag = true
					flagmode = true
				} else {
					flags = matchstrings[1]
					self.cmdName = flags
					cmdmode = true
					self.isCmd = true
				}
			}
		}
	}

	matchstrings = helpExpr.FindStringSubmatch(self.origKey)
	if len(matchstrings) > 1 {
		self.helpInfo = matchstrings[1]
	}

	newprefix = ""
	if len(prefix) > 0 {
		newprefix = prefix
	}

	matchstrings = prefixExpr.FindStringSubmatch(self.origKey)
	if len(matchstrings) > 1 {
		if len(newprefix) > 0 {
			newprefix += "_"
		}
		newprefix += matchstrings[1]
		self.prefix = newprefix
	} else {
		if len(newprefix) > 0 {
			self.prefix = newprefix
		}
	}

	if flagmode {
		self.isFlag = true
		self.isCmd = false
	}

	if cmdmode {
		self.isFlag = false
		self.isCmd = true
	}

	if !flagmode && !cmdmode {
		self.isFlag = true
		self.isCmd = false
	}

	self.value = self.getValue(value)
	if !isjsonfile && !ishelp {
		self.typeName = self.getType(value)
	} else if isjsonfile {
		self.typeName = "jsonfile"
		self.nargs = 1
	} else if ishelp {
		self.typeName = "help"
		self.nargs = 0
	}

	if self.typeName == "help" && value != nil {
		return fmt.Errorf(format_error("help type must be value None"))
	}

	if cmdmode && self.typeName != "dict" {
		flagmode = true
		cmdmode = false
		self.isFlag = true
		self.isCmd = false
		self.flagName = self.cmdName
		self.cmdName = ""
	}

	if self.isFlag && self.typeName == "string" && self.value != nil && self.value.(string) == "+" && self.flagName != "$" {
		self.typeName = "count"
		self.value = 0
		self.nargs = 0
	}

	if self.isFlag && self.flagName == "$" && self.typeName != "dict" {
		if !((self.typeName == "string" && self.value != nil && strings.Contains("?+*", self.value.(string))) || self.typeName == "int") {
			return fmt.Errorf(format_error("(%s)(%s)(%s) for $ should option dict set opt or +?* specialcase or type int", prefix, self.origKey, fmt.Sprintf("%v", self.value)))
		} else {
			self.nargs = self.value
			self.value = nil
			self.typeName = "string"
		}
	}

	if self.isFlag && self.typeName == "dict" && len(self.flagName) > 0 {
		err = self.setFlag(prefix, key, value)
		if err != nil {
			return err
		}
	}

	matchstrings = attrExpr.FindStringSubmatch(self.origKey)
	if len(matchstrings) > 1 {
		self.attr, err = parseAttr(matchstrings[1])
		if err != nil {
			return err
		}
	}

	matchstrings = funcExpr.FindStringSubmatch(self.origKey)
	if len(matchstrings) > 1 {
		if flagmode {
			self.varName = matchstrings[1]
		} else {
			self.function = matchstrings[1]
		}
	}
	return self.validate()
}

func (self *ExtKeyParse) Format() string {
	var s string
	s = "{"
	s += fmt.Sprintf("<type:%s>", self.typeName)
	s += fmt.Sprintf("<origkey:%s>", self.origKey)
	if self.isCmd {
		s += fmt.Sprintf("<cmdname:%s>", self.cmdName)
		if len(self.function) > 0 {
			s += fmt.Sprintf("<function:%s>", self.function)
		}

		if len(self.helpInfo) > 0 {
			s += fmt.Sprintf("<helpinfo:%s>", self.helpInfo)
		}

		if len(self.prefix) > 0 {
			s += fmt.Sprintf("<prefix:%s>", self.prefix)
		}
	}

	if self.isFlag {
		if len(self.flagName) > 0 {
			s += fmt.Sprintf("<flagname:%s>", self.flagName)
		}

		if len(self.shortFlag) > 0 {
			s += fmt.Sprintf("<shortflag:%s>", self.shortFlag)
		}

		if len(self.prefix) > 0 {
			s += fmt.Sprintf("<prefix:%s>", self.prefix)
		}

		if self.nargs != nil {
			s += fmt.Sprintf("<nargs:%v>", self.nargs)
		}

		if len(self.varName) > 0 {
			s += fmt.Sprintf("<varname:%s>", self.varName)
		}

		if self.value != nil {
			s += fmt.Sprintf("<value:%v>", self.value)
		}

		s += fmt.Sprintf("<longprefix:%s>", self.longPrefix)
		s += fmt.Sprintf("<shortprefix:%s>", self.shortPrefix)
	}

	s += fmt.Sprintf("<attr:%s>", formatMap(self.attr))
	return s
}

func (self *ExtKeyParse) Longopt() string {
	var longopt string
	if !self.isFlag || len(self.flagName) == 0 || self.typeName == "args" {
		s := fmt.Sprintf("can not set (%s) longopt", self.origKey)
		panic(s)
	}

	longopt = fmt.Sprintf("%s", self.longPrefix)
	if self.typeName == "bool" && self.value.(bool) {
		longopt += "no-"
	}

	if len(self.prefix) > 0 && self.typeName != "help" {
		longopt += fmt.Sprintf("%s_", self.prefix)
	}
	longopt += self.flagName

	if !self.noChange {
		longopt = strings.ToLower(longopt)
		longopt = strings.Replace(longopt, "_", "-", -1)
	}
	return longopt
}

func (self *ExtKeyParse) Shortopt() string {
	var shortopt string = ""
	if !self.isFlag || len(self.flagName) == 0 || self.typeName == "args" {
		s := fmt.Sprintf("can not set (%s) shortopt", self.origKey)
		panic(s)
	}

	if len(self.shortFlag) > 0 {
		shortopt = fmt.Sprintf("%s%s", self.shortPrefix, self.shortFlag)
	}
	return shortopt
}

func (self *ExtKeyParse) LongPrefix() string {
	return self.longPrefix
}

func (self *ExtKeyParse) ShortPrefix() string {
	return self.shortPrefix
}

func (self *ExtKeyParse) NeedArg() int {
	if !self.isFlag {
		return 0
	}
	if self.typeName == "int" || self.typeName == "list" ||
		self.typeName == "long" || self.typeName == "float" ||
		self.typeName == "string" || self.typeName == "jsonfile" {
		return 1
	}
	return 0
}

func (self *ExtKeyParse) Prefix() string {
	return self.prefix
}

func (self *ExtKeyParse) Value() interface{} {
	return self.value
}

func (self *ExtKeyParse) CmdName() string {
	return self.cmdName
}

func (self *ExtKeyParse) HelpInfo() string {
	return self.helpInfo
}

func (self *ExtKeyParse) Function() string {
	return self.function
}

func (self *ExtKeyParse) Attr(k string) string {
	if k == "" {
		return formatMap(self.attr)
	}

	v, ok := self.attr[k]
	if !ok {
		return ""
	}
	return v
}

func (self *ExtKeyParse) FlagName() string {
	return self.flagName
}

func (self *ExtKeyParse) VarName() string {
	return self.varName
}

func (self *ExtKeyParse) Nargs() interface{} {
	return self.nargs
}

func (self *ExtKeyParse) Equal(other *ExtKeyParse) bool {
	if self.Format() != other.Format() {
		return false
	}
	return true
}

func newExtKeyParse_long(prefix string, key string, value interface{}, isflag bool, ishelp bool, isjsonfile bool, longprefix string, shortprefix string, nochange bool) (k *ExtKeyParse, err error) {
	p := &ExtKeyParse{}
	p.nargs = nil
	p.value = nil
	p.isFlag = false
	p.isCmd = false
	err = p.parse(prefix, key, value, isflag, ishelp, isjsonfile, longprefix, shortprefix, nochange)
	if err != nil {
		k = nil
		return
	}
	k = p
	err = nil
	return
}

func newExtKeyParse(prefix string, key string, value interface{}, isflag bool) (k *ExtKeyParse, err error) {
	return newExtKeyParse_long(prefix, key, value, isflag, false, false, "--", "-", false)
}
