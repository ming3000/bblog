package bblog

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type BBLog interface {
	io.Writer
	Close() error
}

type BaseLogger struct {
	manager *Manager
	option  *Option

	fileObj  *os.File
	filePath string
}

func (b *BaseLogger) ReOpen(file string) error {
	return nil
}

func (b *BaseLogger) Write(data []byte) (int, error) {
	fmt.Println(data)
	return 0, nil
}

func (b *BaseLogger) Close() error {
	return nil
}

type LockLogger struct {
	BaseLogger
	sync.Mutex
}

type BufferLogger struct {
	BaseLogger
	buffer *[]byte
	swap   int
}

func NewBBLog(opt *Option) (BBLog, error) {
	if opt.LogPath == "" || opt.FileName == "" {
		return nil, ErrInvalidOption
	}

	if err := os.MkdirAll(opt.LogPath, 0700); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(opt.LogFilePath(), DefaultFileFlag, DefaultFileMode)
	if err != nil {
		return nil, err
	}

}
