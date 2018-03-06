package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var default_pkg_name []string = []string{"github.com", "jeppeter", "go-extargsparse"}
var default_length = 50

func run_test_case(pkgpath string) error {
	/*now we should */
	var cmd *exec.Cmd
	var gobin string
	var obuf, ebuf *bytes.Buffer
	var err error
	var curdir string
	fmt.Fprintf(os.Stdout, "%s", format_length(default_length, "test case ..."))
	os.Stdout.Sync()
	curdir, err = filepath.Abs("./")
	if err != nil {
		err = format_error("can not get current dir[%s]", err.Error())
		return err
	}
	defer os.Chdir(curdir)
	err = os.Chdir(pkgpath)
	if err != nil {
		err = format_error("can not chdir [%s] err[%s]", pkgpath, err.Error())
		return err
	}
	if strings.ToLower(runtime.GOOS) == "windows" {
		gobin = "go.exe"
	} else {
		gobin = "go"
	}
	cmd = exec.Command(gobin, "test", "-v")
	obuf = bytes.NewBufferString("")
	ebuf = bytes.NewBufferString("")
	cmd.Stdout = obuf
	cmd.Stderr = ebuf

	err = cmd.Start()
	if err != nil {
		err = format_error("can not run %v error [%s] errout\n%s", cmd.Args, err.Error(), ebuf.String())
		return err
	}

	fmt.Fprintf(os.Stdout, "OK\n")
	return nil
}

func add_gopath_env(builddir string) []string {
	var envs []string
	var newenvs []string
	var s string
	envs = os.Environ()
	newenvs = make([]string, 0)
	for _, s = range envs {
		if !strings.HasPrefix(s, "GOPATH=") {
			newenvs = append(newenvs, s)
		}
	}

	newenvs = append(newenvs, fmt.Sprintf("GOPATH=%s%c%s", builddir, filepath.ListSeparator, os.Getenv("GOPATH")))
	return newenvs
}

func copy_to_tempdir(srcdir string) (dstdir string, err error) {
	var realdst string
	dstdir, err = ioutil.TempDir("", "chkbuild")
	if err != nil {
		err = format_error("can not tempdir err[%s]", err.Error())
		return
	}

	defer func() {
		if err != nil && os.Getenv("NOT_REMOVE_FILE") == "" {
			os.RemoveAll(dstdir)
		}
	}()

	err = os.Chmod(dstdir, 0700)
	if err != nil {
		err = format_error("can not chmod 0700 for [%s] err[%s]", dstdir, err.Error())
		return
	}

	realdst = filepath.Join(dstdir, "src", default_pkg_name[0], default_pkg_name[1], default_pkg_name[2])

	err = os.MkdirAll(realdst, 0700)
	if err != nil {
		err = format_error("can not mkdir src [%s] err[%s]", dstdir, err.Error())
		return
	}
	err = copyDir(srcdir, realdst)
	if err != nil {
		err = format_error("can not copy [%s] => [%s] err[%s]", srcdir, realdst, err.Error())
		return
	}
	err = nil
	return
}

func compile_and_run(gofile string, envs []string) error {
	var cmd *exec.Cmd
	var gobin string
	var exebin string
	var obuf, ebuf *bytes.Buffer
	var err error

	fmt.Fprintf(os.Stdout, "%s", format_length(default_length, "%s ...", filepath.Base(gofile)))
	os.Stdout.Sync()
	exebin = gofile[:(len(gofile) - 3)]
	if strings.ToLower(runtime.GOOS) == "windows" {
		gobin = "go.exe"
		exebin = fmt.Sprintf("%s.exe", exebin)
	} else {
		gobin = "go"
	}
	cmd = exec.Command(gobin, "build", "-o", exebin, gofile)
	obuf = bytes.NewBufferString("")
	ebuf = bytes.NewBufferString("")
	cmd.Stdout = obuf
	cmd.Stderr = ebuf
	cmd.Env = envs
	err = cmd.Run()
	if err != nil {
		err = format_error("can not run %v err[%s]\n%s", cmd.Args, err.Error(), ebuf.String())
		return err
	}

	defer func() {
		if err == nil && os.Getenv("NOT_REMOVE_FILE") == "" {
			os.RemoveAll(exebin)
		} else {
			fmt.Fprintf(os.Stderr, "exebin [%s] for [%s]\n", exebin, gofile)
		}
	}()

	cmd = exec.Command(exebin)
	obuf = bytes.NewBufferString("")
	ebuf = bytes.NewBufferString("")
	cmd.Stdout = obuf
	cmd.Stderr = ebuf
	err = cmd.Run()
	if err != nil {
		err = format_error("can not run %v err[%s]\n%s", cmd.Args, err.Error(), ebuf.String())
		return err
	}
	fmt.Fprintf(os.Stdout, "OK\n")
	return nil
}

func scan_dir_build(curdir string, envs []string) error {
	var files []os.FileInfo
	var err error
	var curf os.FileInfo
	var curabs string
	files, err = ioutil.ReadDir(curdir)
	for _, curf = range files {
		if curf.IsDir() {
			if curf.Name() == "." || curf.Name() == ".." {
				continue
			}

			curabs = filepath.Join(curdir, curf.Name())
			err = scan_dir_build(curabs, envs)
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(curf.Name(), ".go") {
			curabs = filepath.Join(curdir, curf.Name())
			err = compile_and_run(curabs, envs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func run_example_case(pkgpath, examplepath string) error {
	var tempdir string = ""
	var err error
	var envs []string

	tempdir, err = copy_to_tempdir(pkgpath)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil && os.Getenv("NOT_REMOVE_FILE") == "" {
			os.RemoveAll(tempdir)
		} else {
			fmt.Fprintf(os.Stderr, "reserved tempdir [%s]\n", tempdir)
		}
	}()

	envs = add_gopath_env(tempdir)
	err = scan_dir_build(examplepath, envs)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var exename string = ""
	var err error
	var examplepath string = ""
	var pkgpath string = ""

	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(5)
		} else {
			fmt.Fprintf(os.Stdout, "run pkg[%s] example[%s] succ\n", pkgpath, examplepath)
		}
	}()

	exename, err = os.Executable()
	if err != nil {
		err = format_error("cann ot get executable [%s]", err.Error())
		return
	}

	pkgpath, err = filepath.Abs(filepath.Join(filepath.Dir(exename), ".."))
	if err != nil {
		err = format_error("can ot get abs for pkgpath[%s] err[%s]", exename, err.Error())
		return
	}

	examplepath, err = filepath.Abs(filepath.Join(filepath.Dir(exename), "..", "example"))
	if err != nil {
		err = format_error("can ot get abs for examplepath[%s] err[%s]", exename, err.Error())
		return
	}

	if len(os.Args) > 1 {
		examplepath = os.Args[2]
	}

	if len(os.Args) > 2 {
		pkgpath = os.Args[3]
	}

	err = run_test_case(pkgpath)
	if err != nil {
		return
	}

	err = run_example_case(pkgpath, examplepath)
	if err != nil {
		return
	}
	return

}
