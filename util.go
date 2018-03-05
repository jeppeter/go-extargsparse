package extargsparse

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unsafe"
)

func assert_test(ischeck bool, fmtstr string, a ...interface{}) {
	if !ischeck {
		s := fmt.Sprintf(fmtstr, a...)
		panic(s)
	}
}

func format_out_stack(level int) string {
	_, f, l, _ := runtime.Caller(level)
	return fmt.Sprintf("[%s:%d]", f, l)
}

func format_error_ex(level int, fmtstr string, a ...interface{}) string {
	s := format_out_stack(level + 1)
	s += fmt.Sprintf(fmtstr, a...)
	return s
}

func format_error(fmtstr string, a ...interface{}) string {
	return format_error_ex(2, fmtstr, a...)
}

func formatMap(kattr map[string]string) string {
	var s string = ""
	var ks []string
	var k, v string
	ks = make([]string, 0)
	for k, _ = range kattr {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k = range ks {
		v = kattr[k]
		s += fmt.Sprintf("[%s]=[%s]\n", k, v)
	}
	return s
}

func check_in_array(sarr []string, s string) bool {
	for _, k := range sarr {
		if k == s {
			return true
		}
	}
	return false
}

func keyDebug(fmtstr string, a ...interface{}) {
	s := format_out_stack(2)
	s += fmt.Sprintf(fmtstr, a...)
	fmt.Fprintf(os.Stderr, "%s\n", s)
	return
}

func findInterfaceField(a *reflect.Value, key string) int {
	var maxfld int
	var pt reflect.Type
	var i int
	maxfld = a.NumField()
	pt = a.Type()
	for i = 0; i < maxfld; i++ {
		if pt.Field(i).Name == key {
			return i
		}
	}
	return -1
}

func setMemberValueInner(a *reflect.Value, name string, value interface{}) error {
	var sarr []string
	var idx int
	var rf, vrf reflect.Value
	var s, rs string
	sarr = strings.Split(name, ".")
	if len(sarr) == 0 {
		return fmt.Errorf("%s", format_error("not set name []"))
	}

	if len(sarr) == 1 {
		idx = findInterfaceField(a, name)
		if idx < 0 {
			return fmt.Errorf("%s", format_error("can not find [%s]", name))
		}
		rf = a.Field(idx)
		rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
		s = rf.Type().String()
		vrf = reflect.ValueOf(value)
		if vrf.Kind() == reflect.Interface || vrf.Kind() == reflect.Ptr {
			rs = vrf.Elem().Type().String()
		} else {
			rs = vrf.Type().String()
		}

		if s != rs {
			return fmt.Errorf("%s", format_error("rf type[%s] != value[%v] type[%s]", s, value, rs))
		}
		rf.Set(reflect.ValueOf(value))
		return nil
	}

	idx = findInterfaceField(a, sarr[0])
	if idx < 0 {
		return fmt.Errorf("%s", format_error("can not find [%s] for [%s]", sarr[0], name))
	}
	rf = a.Field(idx)
	rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	return setMemberValueInner(&rf, strings.Join(sarr[1:], "."), value)
}

func setMemberValue(a interface{}, name string, value interface{}) error {
	var rf reflect.Value
	rf = reflect.ValueOf(a).Elem()
	return setMemberValueInner(&rf, name, value)
}

func ucFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func lcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func getCallerPackage(skip int) string {
	/*to add skip this will give the caller package*/
	skip = skip + 1
	var pcs []uintptr
	var ok bool = false
	var size int = skip + 1
	var ret int
	var sarr []string
	var curname string
	var name string
	var retname string
	var i int
	for !ok {
		pcs = make([]uintptr, size)
		ret = runtime.Callers(0, pcs)
		if ret < size {
			ok = true
		} else {
			size <<= 1
		}
	}

	if ret < skip {
		return ""
	}

	name = runtime.FuncForPC(pcs[skip]).Name()
	sarr = strings.Split(name, ".")

	retname = ""
	for i, curname = range sarr {
		if i < (len(sarr) - 1) {
			if !strings.HasPrefix(curname, "(") {
				if len(retname) > 0 {
					retname += "."
				}
				retname += curname
			}
		}
	}
	return retname
}

func getCallerFilename(skip int) string {
	/*to add skip this will give the caller package*/
	var retname string
	var ok bool
	_, retname, _, ok = runtime.Caller(skip + 1)
	if !ok {
		return ""
	}
	return retname
}

func isFloatToInt(fv float64) (iv int, bret bool) {
	var mzeros *regexp.Regexp
	var cv *regexp.Regexp
	var vstr string
	var err error
	var dotstring string
	var numstr string
	var basenum string
	var cnt int
	var ibase int
	var li int
	vstr = fmt.Sprintf("%v", fv)
	mzeros = regexp.MustCompile(`^[0-9]+$`)
	cv = regexp.MustCompile(`^([0-9])\.([0-9]+)[eE][\+]([0-9]+)$`)
	iv = 0
	bret = false
	if mzeros.MatchString(vstr) {
		iv, err = strconv.Atoi(vstr)
		if err == nil {
			bret = true
		}
		return
	}
	if cv.MatchString(vstr) {
		matchstrings := cv.FindStringSubmatch(vstr)
		if len(matchstrings) >= 3 {
			basenum = matchstrings[1]
			dotstring = matchstrings[2]
			numstr = matchstrings[3]
			li, err = strconv.Atoi(numstr)
			if err == nil && li >= len(dotstring) {
				ibase, err = strconv.Atoi(basenum)
				if err == nil {
					iv = ibase
					for cnt = 0; cnt < li; cnt++ {
						iv *= 10
					}
					cnt, err = strconv.Atoi(dotstring)
					if err == nil {
						for len(dotstring) < li {
							cnt *= 10
							li--
						}
						iv += cnt
						bret = true
						return
					}
				}
			}
		}
	}
	return
}

func getExecutableName(name string) string {
	if strings.ToLower(runtime.GOOS) == "windows" {
		return fmt.Sprintf("%s.exe", name)
	}
	return name
}
