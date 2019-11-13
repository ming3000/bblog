package bblog

import (
	"io"
	"os"
)

type BBLogger interface {
	io.Writer
	Close() error
}

type BaseLogger struct {
	manager *Manager
	option  *Option

	fileObj  *os.File
	filePath string
}
