package bblog

import (
	"errors"
	"os"
	"path"
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
	WriteModeNone = iota
	WriteModeLock
	WriteModeAsync
	WriteModeBuffered
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
	// the truncated log is located at [LogPath]/[FileName].log.[RollingTimeTagFormat]
	// the compressed log is located at [LogPath]/[FileName].log.gz.[RollingTimeTagFormat]
	LogPath  string `json:"log_path"`
	FileName string `json:"file_name"`

	// postfix of truncated file
	RollingTimeTagFormat string `json:"rolling_time_tag_format"`
	// the mod will auto delete the rolling file, set 0 to disable auto clean
	MaxRollingRemain int `json:"max_rolling_remain"`

	RollingPolicy int `json:"rolling_policy"`
	// cron job like pattern
	RollingTimePattern string `json:"rolling_time_pattern"`
	RollingFileSize    string `json:"rolling_file_size"`

	WriteMode  int `json:"write_mode"`
	BufferSize int `json:"buffer_size"`
}

func (c *Config) LogFilePath() string {
	return path.Join(c.LogPath, c.FileName)
}

func NewDefaultConfig() Config {
	return Config{
		LogPath:  "./log",
		FileName: "log",

		RollingTimeTagFormat: "200601021504",
		MaxRollingRemain:     0,

		RollingPolicy:      TimeRolling,
		RollingTimePattern: "0 0 0 * * *",
		RollingFileSize:    "512M",

		WriteMode: WriteModeLock,
	}
}
