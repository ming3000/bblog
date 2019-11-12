package bblog

import (
	"errors"
	"os"
)

// 3 policy for rolling
const (
	// no rolling
	WithoutRolling = iota
	// rolling by time
	TimeRolling
	// rolling by file size
	FileSizeRolling
)

const (
	// the default buffer size is 1M
	BufferSize = 0x100000
	// the queue for async write
	QueueSize = 1024
	// the precision defined the precision about the reopen operation condition check duration within second
	Precision = 1
)

const (
	// default open file mode is rw-r--r--
	DefaultFileMode = os.FileMode(0644)
	DefaultFileFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND
)

var (
	ErrInvalidArgument      = errors.New("invalid argument")
	ErrorWriteContextClosed = errors.New("write context closed")
	ErrOther                = errors.New("other error")
)

type Config struct {
	// the LogPath & FileName define the full path of the log.
	// the current log is located at [LogPath]/[FileName].log
	// the truncated log is located at [LogPath]/[FileName].log.[TruncateRollingTag]
	// the compressed log is located at [LogPath]/[FileName].log.gz.[TruncateRollingTag]
	LogPath  string `json:"log_path"`
	FileName string `json:"file_name"`

	TruncateRollingTag string `json:"truncated_rolling_tag"`
	MaxRollingRemain   int    `json:"max_rolling_remain"`
}
