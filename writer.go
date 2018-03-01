package extargsparse

import (
	"fmt"
	"os"
)

type IoWriter interface {
	Write(data []byte) (int, error)
	WriteString(s string) (int, error)
}

type FileIoWriter struct {
	IoWriter
	file *os.File
}

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