package extargsparse

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func read_file_bytes(fname string) (rbytes []byte, err error) {
	var fp *os.File = os.Stdin
	err = nil
	rbytes = []byte{}
	defer func() {
		if fp != nil && fp != os.Stdin {
			fp.Close()
		}
		fp = nil
	}()

	if fname != "" {
		fp, err = os.Open(fname)
		if err != nil {
			err = fmt.Errorf("open [%s] error[%s]", fname, err.Error())
			return
		}
	}
	rbytes, err = ioutil.ReadAll(fp)
	return
}

func get_current_process_exec_info() (startaddr uintptr, endaddr uintptr, err error) {
	var exename string
	var exebyte []byte
	var exefile string
	var mapname string
	var mapbytes []byte
	var s string
	var lines []string
	var curline string
	var minaddr uintptr = 0xffffffffffffffff
	var maxaddr uintptr = 0
	var curstart uintptr
	var curend uintptr
	var reg *regexp.Regexp
	var sarr []string
	var hsarr []string
	var finded bool = false
	var ci int64
	var pid int = os.Getpid()
	startaddr = 0
	endaddr = 0
	err = fmt.Errorf("not valid")
	exename = fmt.Sprintf("/proc/%d/exe", pid)
	cmd := exec.Command("readlink", "-f", exename)
	exebyte, err = cmd.Output()
	if err != nil {
		return
	}
	exefile = string(exebyte)
	exefile = strings.TrimRight(exefile, "\n")
	exefile = strings.TrimRight(exefile, "\r")
	mapname = fmt.Sprintf("/proc/%d/maps", pid)
	mapbytes, err = read_file_bytes(mapname)
	if err != nil {
		return
	}
	s = string(mapbytes)
	lines = strings.Split(s, "\n")
	reg, err = regexp.Compile("\\s+")
	if err != nil {
		return
	}
	for _, curline = range lines {
		curline = strings.TrimRight(curline, "\r")
		if len(curline) == 0 {
			continue
		}

		if strings.HasSuffix(curline, exefile) {
			sarr = reg.Split(curline, -1)
			if len(sarr) > 2 {
				if strings.Contains(sarr[1], "x") {
					/*it means this area is executing one*/
					hsarr = strings.Split(sarr[0], "-")
					if len(hsarr) > 1 {
						ci, err = strconv.ParseInt(hsarr[0], 16, 64)
						if err != nil {
							return
						}
						curstart = uintptr(ci)
						ci, err = strconv.ParseInt(hsarr[1], 16, 64)
						if err != nil {
							return
						}
						curend = uintptr(ci)
						if curstart < minaddr {
							minaddr = curstart
						}

						if curend > maxaddr {
							maxaddr = curend
						}
						finded = true
					}
				}
			}
		}

	}

	if finded {
		err = nil
		startaddr = minaddr
		endaddr = maxaddr
	} else {
		err = fmt.Errorf("not valid")
	}

	return
}
