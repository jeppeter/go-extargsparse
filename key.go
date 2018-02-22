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
	helpinfo    string
	shortFlag   string
	nargs       interface{}
	varname     string
	cmdname     string
	function    string
	origkey     string
	iscmd       bool
	isflag      bool
	typename    string
	Attr        map[string]string
}

func (self *extKeyParse) parse(prefix string, key string, value interface{}, isflag bool, ishelp bool, isjsonfile bool, longprefix string, shortprefix string, nochange bool) error {
	var flagmode bool = false
	var cmdmode bool = false
	var flags string = ""
	var ok int
	var matchstrings []string
	var sarr []string
	var hexpr *regexp.Regexp
	self.origkey = key
	self.longPrefix = longprefix
	self.shortPrefix = shortprefix
	self.noChange = nochange

	/*now to test whether it is the flag one*/
	if strings.Contains(self.origkey, "$") {
		if self.origkey[0:1] != "$" {
			return fmt.Errorf("(%s) not right format for ($)", self.origkey)
		}
		ok = 1
		if strings.Contains(self.origkey[1:], "$") {
			ok = 0
		}
		if ok != 1 {
			return fmt.Errorf("(%s) has ($) more than one", self.origkey)
		}
	}

	if isflag || ishelp || isjsonfile {
		matchstrings = flagExpr.FindStringSubmatch(self.origkey)
		if len(matchstrings) > 1 {
			flags = matchstrings[1]
		}
		if len(flags) == 0 {
			matchstrings = mustFlagExpr.FindStringSubmatch(self.origkey)
			if len(matchstrings) > 1 {
				flags = matchstrings[1]
			}
		}

		if len(flags) == 0 && self.origkey[0:1] == "$" {
			self.flagName = "$"
			flagmode = true
		}

		if len(flags) > 0 {
			if strings.Contains(flags, "|") {
				hexpr = regexp.MustCompile(`\|`)
				sarr = hexpr.Split(flags, -1)
				if len(sarr) > 2 || len(sarr[1]) != 1 || len(sarr[0]) <= 1 {
					return fmt.Errorf("(%s) (%s)flag only accept (longop|l) format", self.origkey, flags)
				}
				self.flagName = sarr[0]
				self.shortFlag = sarr[1]
			} else {
				self.flagName = flags
			}
			flagmode = true
		}
	} else {
		matchstrings = mustFlagExpr.FindStringSubmatch(self.origkey)
		if len(matchstrings) > 1 {
			flags = matchstrings[1]
			if strings.Contains(flags, "|") {
				hexpr = regexp.MustCompile(`\|`)
				sarr = hexpr.Split(flags, -1)
				if len(sarr) > 2 || len(sarr[1]) != 1 || len(sarr[0]) <= 1 {
					return fmt.Errorf("(%s) (%s)flag only accept (longop|l) format", self.origkey, flags)
				}
				self.flagName = sarr[0]
				self.shortFlag = sarr[1]
			} else {
				self.flagName = flags
			}
			flagmode = true
		} else if self.origkey == "$" {
			self.flagName = "$"
			flagmode = true
		}

		matchstrings = cmdExpr.FindStringSubmatch(self.origkey)
		if len(matchstrings) > 1 {
		}
	}
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
