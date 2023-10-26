package extargsparse

import (
	"os"
	"syscall"
	"unsafe"
)

const MAX_MODULE_NAME32 = 255

type ModuleEntry32 struct {
	Size         uint32
	ModuleID     uint32
	ProcessID    uint32
	GlblcntUsage uint32
	ProccntUsage uint32
	ModBaseAddr  uintptr
	ModBaseSize  uint32
	ModuleHandle syscall.Handle
	Module       [MAX_MODULE_NAME32 + 1]uint16
	ExePath      [syscall.MAX_PATH]uint16
}

const SizeofModuleEntry32 = unsafe.Sizeof(ModuleEntry32{})

var modkernel32 *syscall.LazyDLL = nil
var procModule32FirstW *syscall.LazyProc = nil

func init() {
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	procModule32FirstW = modkernel32.NewProc("Module32FirstW")
}

var (
	errnoERROR_IO_PENDING       = syscall.Errno(997)
	errERROR_IO_PENDING   error = syscall.Errno(997)
	errERROR_EINVAL       error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

//sys	Module32First(snapshot syscall.Handle, moduleEntry *ModuleEntry32) (err error) = kernel32.Module32FirstW
//sys	Module32Next(snapshot syscall.Handle, moduleEntry *ModuleEntry32) (err error) = kernel32.Module32NextW

func Module32First(snapshot syscall.Handle, moduleEntry *ModuleEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procModule32FirstW.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(moduleEntry)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func get_current_process_exec_info() (startaddr uintptr, endaddr uintptr, err error) {
	var modinfo ModuleEntry32
	var pid int = os.Getpid()
	err = nil
	startaddr = 0
	endaddr = 0

	snaphd, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPMODULE|syscall.TH32CS_SNAPMODULE32, uint32(pid))
	if err != nil {
		return
	}

	defer syscall.CloseHandle(snaphd)
	modinfo.Size = uint32(SizeofModuleEntry32)
	err = Module32First(snaphd, &modinfo)
	if err != nil {
		return
	}

	startaddr = uintptr(modinfo.ModBaseAddr)
	endaddr = uintptr(modinfo.ModBaseAddr) + uintptr(modinfo.ModBaseSize)
	err = nil

	return
}
