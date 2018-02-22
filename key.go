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
