package bblog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"unsafe"
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

func (b *BaseLogger) ReOpen(fileName string) error {
	if err := os.Rename(b.filePath, fileName); err != nil {
		return err
	}

	newFileObj, err := os.OpenFile(b.filePath, DefaultFileFlag, DefaultFileMode)
	if err != nil {
		return err
	}
	oldFile := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&b.fileObj)), unsafe.Pointer(newFileObj))
	defer (*os.File)(oldFile).Close()
	return nil
}

func (b *BaseLogger) Write(data []byte) (int, error) {
	select {
	case fileName := <-b.manager.rollingChan:
		if err := b.ReOpen(fileName); err != nil {
			return 0, err
		}
	default:
	}
	fp := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&b.fileObj)))
	return (*os.File)(fp).Write(data)
}

func (b *BaseLogger) Close() error {
	return (*os.File)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&b.fileObj)))).Close()
}

type LockLogger struct {
	BaseLogger
	sync.Mutex
}

func (l *LockLogger) Write(data []byte) (int, error) {

}

func (l *LockLogger) Close() error {
	l.Lock()
	defer l.Unlock()
	return l.fileObj.Close()
}

type BufferLogger struct {
	BaseLogger
	buffer *[]byte
	swap   int
}

func (b *BufferLogger) Write(data []byte) (int, error) {
	fmt.Println(data)
	return 0, nil
}

func (b *BufferLogger) Close() error {
	_, _ = b.fileObj.Write(*b.buffer)
	return b.fileObj.Close()
}

func NewBBLog(opt *Option) (BBLog, error) {
	if opt.LogPath == "" || opt.FileName == "" {
		return nil, ErrInvalidOption
	}

	if err := os.MkdirAll(opt.LogPath, 0700); err != nil {
		return nil, err
	}

	theFile, err := os.OpenFile(opt.LogFilePath(), DefaultFileFlag, DefaultFileMode)
	if err != nil {
		return nil, err
	}

	theManager, err := NewManager(opt)
	if err != nil {
		return nil, err
	}

	var ret BBLog
	baseLogger := BaseLogger{
		option:   opt,
		manager:  theManager,
		filePath: opt.LogFilePath(),
		fileObj:  theFile,
	}
	switch opt.WriteMode {
	case WriteModeNone:
		ret = &baseLogger
	case WriteModeLock:
		ret = &LockLogger{BaseLogger: baseLogger}
	case WriteModeBuffered:
		theBuf := make([]byte, 0, opt.BufferSize)
		ret = &BufferLogger{BaseLogger: baseLogger, buffer: &theBuf, swap: 0}
	default:
		return nil, ErrInvalidOption
	}
	return ret, nil
}
