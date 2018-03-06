package extargsparse

import (
	"fmt"
	"os"
)

// IoWriter is the interface for PrintHelp
//    extended example see https://github.com/jeppeter/go-extargsparse/blob/master/example/helpfunc/helpstr1.go
type IoWriter interface {
	Write(data []byte) (int, error)
	WriteString(s string) (int, error)
}

type FileIoWriter struct {
	IoWriter
	file *os.File
}

// it is called by PrintHelp
//    example see https://github.com/jeppeter/go-extargsparse/blob/master/example/helpfunc/filehelp.go
func NewFileWriter(f *os.File) *FileIoWriter {
	self := &FileIoWriter{}
	self.file = f
	return self
}

func (self *FileIoWriter) Write(data []byte) (int, error) {
	if self.file != nil {
		return self.file.Write(data)
	}
	return 0, fmt.Errorf("%s", format_error("no file assign"))

}

func (self *FileIoWriter) WriteString(s string) (int, error) {
	if self.file != nil {
		return self.file.WriteString(s)
	}
	return 0, fmt.Errorf("%s", format_error("no file assign"))
}
