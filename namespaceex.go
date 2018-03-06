package extargsparse

import (
	"fmt"
	"sort"
	"strconv"
)

// NameSpaceEx structure to store the parse command line result
//    use multifunction to get the coding
type NameSpaceEx struct {
	logger *logObject
	obj    map[string]interface{}
}

func newNameSpaceEx() *NameSpaceEx {
	self := &NameSpaceEx{logger: newLogObject("extargsparse")}
	self.obj = make(map[string]interface{})
	return self
}

// SetValue to set value
func (self *NameSpaceEx) SetValue(k string, v interface{}) {
	self.obj[k] = v
	return
}

// GetValue return interface{} if not set return nil
func (self *NameSpaceEx) GetValue(k string) interface{} {
	var v interface{} = nil
	var ok bool
	v, ok = self.obj[k]
	if !ok {
		return nil
	}
	return v
}

// IsAccessed return whether the key has been setted
func (self *NameSpaceEx) IsAccessed(k string) bool {
	var ok bool
	_, ok = self.obj[k]
	return ok
}

// GetBool return true/false value ,no key or not the type return false
func (self *NameSpaceEx) GetBool(k string) bool {
	var v bool = false
	var ok bool
	var val interface{}
	val, ok = self.obj[k]
	if !ok {
		return v
	}

	switch val.(type) {
	case bool:
		v = val.(bool)
	}
	return v
}

// GetString return string type ,default "" on no key or not type of string
func (self *NameSpaceEx) GetString(k string) string {
	var v interface{}
	v = self.GetValue(k)
	if v == nil {
		return ""
	}

	switch v.(type) {
	case string:
		return v.(string)
	}

	return fmt.Sprintf("%v", v)
}

// GetInt return int type , default 0 on nokey or not type of int
func (self *NameSpaceEx) GetInt(k string) int {
	var v interface{}
	var err error
	var vstr string
	var vint int
	v = self.GetValue(k)
	if v == nil {
		return 0
	}

	switch v.(type) {
	case int:
		return v.(int)
	case uint32:
		return int(v.(uint32))
	case uint64:
		return int(v.(uint64))
	case int32:
		return int(v.(int32))
	case int64:
		return int(v.(int64))
	case float32:
		return int(v.(float32))
	case float64:
		return int(v.(float64))
	}

	vstr = fmt.Sprintf("%v", v)
	vint, err = strconv.Atoi(vstr)
	if err != nil {
		return 0
	}
	return vint
}

// GetFloat return type of float64 ,default 0.0 on no key set or not type of float64
func (self *NameSpaceEx) GetFloat(k string) float64 {
	var v interface{}
	var err error
	var vstr string
	var vint float64
	v = self.GetValue(k)
	if v == nil {
		return 0.0
	}

	switch v.(type) {
	case float64:
		return v.(float64)
	case uint32:
		return float64(v.(uint32))
	case uint64:
		return float64(v.(uint64))
	case int32:
		return float64(v.(int32))
	case int64:
		return float64(v.(int64))
	case float32:
		return float64(v.(float32))
	case int:
		return float64(v.(float64))
	}

	vstr = fmt.Sprintf("%v", v)
	vint, err = strconv.ParseFloat(vstr, 64)
	if err != nil {
		return 0.0
	}
	return vint
}

// GetArray return []string, default []string{} on no key set or not type of []string
func (self *NameSpaceEx) GetArray(k string) []string {
	var v interface{}
	var ve string
	var va []string
	var varr []string
	var vstr string
	v = self.GetValue(k)
	varr = make([]string, 0)
	if v == nil {
		return varr
	}

	switch v.(type) {
	case []string:
		va = v.([]string)
		for _, ve = range va {
			vstr = fmt.Sprintf("%s", ve)
			varr = append(varr, vstr)
		}
		return varr
	}
	return varr
}

// GetKeys get all available keys for setted sorted by alphabet value
func (self *NameSpaceEx) GetKeys() []string {
	var keys []string
	keys = make([]string, 0)
	for k, _ := range self.obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Format return string with all keys set like format {key1=val1;key2=value2;...;keyn=valn}
func (self *NameSpaceEx) Format() string {
	var s string = ""
	var keys []string
	var cnt int = 0
	s += "{"
	keys = self.GetKeys()
	for _, k := range keys {
		if cnt > 0 {
			s += ";"
		}
		s += fmt.Sprintf("%s=%s", k, self.GetString(k))
		cnt++
	}
	s += "}"
	return s
}
