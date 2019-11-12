package bblog

import "io"

type RollingWriter interface {
	io.Writer
	Close() error
}
