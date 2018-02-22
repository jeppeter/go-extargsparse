package extargsparse

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
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
var flagwords = []string{"flagname","helpinfo","shortflag","nargs","varname"}

func formatMap(kattr map[string]string) string {
	var s string = ""
	for k, v := range kattr {
		s += fmt.Sprintf("[%s]=[%s]\n", k, v)
	}
	return s
}

func parseAttr(attr string) (kattr map[string]string, err error) {
	var lattr string
	var splitchar string = ";"
	var splitstrings []string
	var splitexpr *regexp.Regexp
	var equalexpr *regexp.Regexp
	var vk []string
	var curs string

	kattr = nil
	err = nil
	lattr = strings.ToLower(attr)
	if strings.HasPrefix(lattr, ATTR_SPLIT_EQUAL) {
		splitchar = lattr[len(ATTR_SPLIT_EQUAL):(len(ATTR_SPLIT_EQUAL) + 1)]
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
			return nil, fmt.Errorf("unknown splitchar [%s]", splitchar)
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
	splitstrings = splitexpr.Split(lattr, -1)
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
		return kattr, fmt.Errorf("not valid type [%s]", reflect.TypeOf(attr))
	}

	err = nil
	return
}

type extKeyParse struct {
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
	Attr        map[string]string
}

func (self *extKeyParse) getType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case bool:
		return "bool"
	case float32:
		return "float"
	case float64:
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

	panic("not valid type [%v]", value)
	return ""
}

func (self *extKeyParse) setFlag(prefix,key string,value interface{}) error {
	var vmap map[string]interface{}
	var findvalue bool = false
	var v interface{}
	var k string
	self.isFlag = true
	self.isCmd = false
	self.origKey = key
	vmap = value.(map[string]interface{})
	for k,_ = range vmap {
		if k == "value" {
			findvalue = true
			break
		}
	}

	if ! findvalue {
		self.typeName = "string"
		self.value = nil
	}

	for k,v = range vmap {

	}

	return nil
}

func (self *extKeyParse) validate() error {
	return nil
}

func (self *extKeyParse) parse(prefix string, key string, value interface{}, isflag bool, ishelp bool, isjsonfile bool, longprefix string, shortprefix string, nochange bool) error {
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
			return fmt.Errorf("(%s) not right format for ($)", self.origKey)
		}
		ok = 1
		if strings.Contains(self.origKey[1:], "$") {
			ok = 0
		}
		if ok != 1 {
			return fmt.Errorf("(%s) has ($) more than one", self.origKey)
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
			flagmode = true
		}

		if len(flags) > 0 {
			if strings.Contains(flags, "|") {
				hexpr = regexp.MustCompile(`\|`)
				sarr = hexpr.Split(flags, -1)
				if len(sarr) > 2 || len(sarr[1]) != 1 || len(sarr[0]) <= 1 {
					return fmt.Errorf("(%s) (%s)flag only accept (longop|l) format", self.origKey, flags)
				}
				self.flagName = sarr[0]
				self.shortFlag = sarr[1]
			} else {
				self.flagName = flags
			}
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
					return fmt.Errorf("(%s) (%s)flag only accept (longop|l) format", self.origKey, flags)
				}
				self.flagName = sarr[0]
				self.shortFlag = sarr[1]
			} else {
				self.flagName = flags
			}
			flagmode = true
		} else if self.origKey == "$" {
			self.flagName = "$"
			flagmode = true
		}

		matchstrings = cmdExpr.FindStringSubmatch(self.origKey)
		if len(matchstrings) > 1 {
			if flagmode {
				panic("flagmode set")
			}
			if strings.Contains(matchstrings[1], "|") {
				flags = matchstrings[1]
				hexpr = regexp.MustCompile(`\|`)
				sarr = hexpr.Split(flags, -1)
				if len(sarr) > 2 || len(sarr[1]) != 1 || len(sarr[0]) <= 1 {
					return fmt.Errorf("(%s) (%s)flag only accept (longop|l) format", self.origKey, flags)
				}
				self.flagName = sarr[0]
				self.shortFlag = sarr[1]
				flagmode = true
			} else {
				self.cmdName = flags
				cmdmode = true
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

	self.value = value
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
		return fmt.Errorf("help type must be value None")
	}

	if cmdmode && self.typeName != "dict" {
		flagmode = true
		cmdmode = false
		self.isFlag = true
		self.isCmd = false
		self.flagName = self.cmdName
		self.cmdName = ""
	}

	if self.isFlag && self.typeName == "string" && self.value.(string) == "+" && self.flagName != "$" {
		self.typeName = "count"
		self.value = 0
		self.nargs = 0
	}

	if self.isFlag && self.flagName == "$" && self.typeName != "dict" {
		if !((self.typeName == "string" && strings.Contains("?+*",self.value.(string))) or self.typeName == "int" ) {
			return fmt.Errorf("(%s)(%s)(%s) for $ should option dict set opt or +?* specialcase or type int", prefix,self.origKey,fmt.Sprintf("%v", self.value))
		} else {
			self.nargs = self.value
			self.value = nil
			self.typeName = "string"
		}
	}

	if self.isFlag && self.typeName == "dict" && len(self.flagName) > 0 {
		err = self.setFlag(prefix,key,value)
		if err != nil {
			return err
		}
	}

	matchstrings = attrExpr.FindStringSubmatch(self.origKey)
	if len(matchstrings) > 1 {
		self.Attr = parseAttr(matchstrings[1])
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

func NewExtKeyParse_long(prefix string, key string, value interface{}, isflag bool, ishelp bool, isjsonfile bool, longprefix string, shortprefix string, nochange bool) (k *extKeyParse, err error) {
	p := &extKeyParse{}
	err = p.parse(prefix, key, value, isflag, ishelp, isjsonfile, longprefix, shortprefix, nochange)
	if err != nil {
		k = nil
		return
	}
	k = p
	err = nil
	return
}

func NewExtKeyParse_short(prefix string, key string, value interface{}, isflag bool) (k *extKeyParse, err error) {
	return NewExtKeyParse_long(prefix, key, value, isflag, false, false, "--", "-", false)
}

func NewExtKeyParse(prefix string, key string, value interface{}) (k *extKeyParse, err error) {
	return NewExtKeyParse_long(prefix, key, value, false, false, false, "--", "-", false)
}
