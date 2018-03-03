package extargsparse

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	default_LOGNAME       = "extargsparse"
	log_FATAL_LEVEL       = 0
	log_ERROR_LEVEL       = 1
	log_WARN_LEVEL        = 2
	log_INFO_LEVEL        = 3
	log_DEBUG_LEVEL       = 4
	log_TRACE_LEVEL       = 5
	log_default_CALLSTACK = 2
)

type logObject struct {
	level int
}

func getEnvLevel(k string) string {
	return strings.ToUpper(fmt.Sprintf("%s_LOGLEVEL", k))
}

func newLogObject(name string) *logObject {
	var levelstr string
	var level int = log_ERROR_LEVEL
	var err error
	levelstr = os.Getenv(getEnvLevel(name))
	if levelstr != "" {
		level, err = strconv.Atoi(levelstr)
		if err != nil {
			level = log_ERROR_LEVEL
		}
	}

	return &logObject{level: level}
}

func (l *logObject) formatCallMsg(callstack int, fmtstr string, a ...interface{}) string {
	s := format_out_stack(callstack + 1)
	s += fmt.Sprintf(fmtstr, a...)
	return s
}

func (l *logObject) inner_call_msg(callstack int, needlevel int, fmtstr string, a ...interface{}) {
	if l.level >= needlevel {
		levelname := "ERROR"
		switch needlevel {
		case 0:
			levelname = "FATAL"
		case 1:
			levelname = "ERROR"
		case 2:
			levelname = "WARN"
		case 3:
			levelname = "INFO"
		case 4:
			levelname = "DEBUG"
		case 5:
			levelname = "TRACE"
		}
		s := l.formatCallMsg(callstack+1, "<%s> ", levelname)
		s += fmt.Sprintf(fmtstr, a...)
		fmt.Fprintf(os.Stderr, "%s\n", s)
	}
}

func (l *logObject) Fatal_long(callstack int, fmtstr string, a ...interface{}) {
	l.inner_call_msg(callstack+log_default_CALLSTACK, log_FATAL_LEVEL, fmtstr, a...)
	return
}

func (l *logObject) Fatal(fmtstr string, a ...interface{}) {
	l.inner_call_msg(log_default_CALLSTACK, log_FATAL_LEVEL, fmtstr, a...)
	return
}

func (l *logObject) Error_long(callstack int, fmtstr string, a ...interface{}) {
	l.inner_call_msg(callstack+log_default_CALLSTACK, log_ERROR_LEVEL, fmtstr, a...)
	return
}

func (l *logObject) Error(fmtstr string, a ...interface{}) {
	l.inner_call_msg(log_default_CALLSTACK, log_ERROR_LEVEL, fmtstr, a...)
	return
}

func (l *logObject) Warn_long(callstack int, fmtstr string, a ...interface{}) {
	l.inner_call_msg(callstack+log_default_CALLSTACK, log_WARN_LEVEL, fmtstr, a...)
	return
}

func (l *logObject) Warn(fmtstr string, a ...interface{}) {
	l.inner_call_msg(log_default_CALLSTACK, log_WARN_LEVEL, fmtstr, a...)
	return
}

func (l *logObject) Info_long(callstack int, fmtstr string, a ...interface{}) {
	l.inner_call_msg(callstack+log_default_CALLSTACK, log_INFO_LEVEL, fmtstr, a...)
	return
}

func (l *logObject) Info(fmtstr string, a ...interface{}) {
	l.inner_call_msg(log_default_CALLSTACK, log_INFO_LEVEL, fmtstr, a...)
	return
}

func (l *logObject) Debug_long(callstack int, fmtstr string, a ...interface{}) {
	l.inner_call_msg(callstack+log_default_CALLSTACK, log_DEBUG_LEVEL, fmtstr, a...)
	return
}

func (l *logObject) GetFuncPtr(uppercase bool, funcname string, outptr interface{}) error {
	var realname string
	var sarr []string
	if len(funcname) == 0 {
		return fmt.Errorf("%s", format_error("not valid funcname [%s]", funcname))
	}
	sarr = strings.Split(funcname, ".")
	if len(sarr) > 1 {
		realname = strings.Join(sarr[:(len(sarr)-1)], ".")
		realname += "."
		if uppercase {
			realname += ucFirst(sarr[len(sarr)-1])
		} else {
			realname += sarr[len(sarr)-1]
		}
	} else {
		if uppercase {
			realname = fmt.Sprintf("main.%s", ucFirst(funcname))
		} else {
			realname = fmt.Sprintf("main.%s", funcname)
		}
	}
	return getFunc(outptr, realname)
}

func (l *logObject) Debug(fmtstr string, a ...interface{}) {
	l.inner_call_msg(log_default_CALLSTACK, log_DEBUG_LEVEL, fmtstr, a...)
	return
}
func (l *logObject) Trace_long(callstack int, fmtstr string, a ...interface{}) {
	l.inner_call_msg(callstack+log_default_CALLSTACK, log_TRACE_LEVEL, fmtstr, a...)
	return
}

func (l *logObject) Trace(fmtstr string, a ...interface{}) {
	l.inner_call_msg(log_default_CALLSTACK, log_TRACE_LEVEL, fmtstr, a...)
	return
}
