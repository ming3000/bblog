package bblog

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
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
	// default open file mode is rw-r--r--
	DefaultFileMode = os.FileMode(0644)
	DefaultFileFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND

	DefaultFileBytesStr = "512M"
	DefaultFileBytes    = 1024 * 1024 * 512

	DefaultFileTagFormat = "200601021504"
)

const (
	// the default buffer size is 1M
	BufferSize = 0x100000
	// the queue for async write
	QueueSize = 1024
	// the precision defined the precision about the reopen operation condition check duration within second
	Precision = 1
)

var (
	ErrInvalidArgument      = errors.New("invalid argument")
	ErrorWriteContextClosed = errors.New("write context closed")
	ErrOther                = errors.New("other error")
)

type Option struct {
	// the LogPath & FileName define the full path of the log.
	// the current log is located at [LogPath]/[FileName].log
	// the truncated log is located at [LogPath]/[FileName].log.[tag]
	LogPath  string `json:"log_path"`
	FileName string `json:"file_name"`

	RollingPolicy int `json:"rolling_policy"`
	// cron job like pattern
	RollingCronJobPattern string `json:"rolling_cron_job_pattern"`
	// file size to start rolling
	RollingFileBytes string `json:"rolling_file_bytes"`

	WriteMode  int `json:"write_mode"`
	BufferSize int `json:"buffer_size"`
}

func (o *Option) LogFilePath() string {
	return path.Join(o.LogPath, o.FileName)
}

func (o *Option) ComputeRollingFileSize() int64 {
	rollingFileSizeStr := strings.ToUpper(o.RollingFileBytes)
	rollingFileSizeByte := []byte(rollingFileSizeStr)

	var tempValue int
	var dstValue int64
	var err error

	switch {
	case strings.Contains(rollingFileSizeStr, "K"):
		tempValue, err = strconv.Atoi(string(rollingFileSizeByte[:len(rollingFileSizeByte)-1]))
		dstValue = int64(tempValue) * 1024
	case strings.Contains(rollingFileSizeStr, "KB"):
		tempValue, err = strconv.Atoi(string(rollingFileSizeByte[:len(rollingFileSizeByte)-2]))
		dstValue = int64(tempValue) * 1024
	case strings.Contains(rollingFileSizeStr, "M"):
		tempValue, err = strconv.Atoi(string(rollingFileSizeByte[:len(rollingFileSizeByte)-1]))
		dstValue = int64(tempValue) * 1024 * 1024
	case strings.Contains(rollingFileSizeStr, "MB"):
		tempValue, err = strconv.Atoi(string(rollingFileSizeByte[:len(rollingFileSizeByte)-2]))
		dstValue = int64(tempValue) * 1024 * 1024
	case strings.Contains(rollingFileSizeStr, "G"):
		tempValue, err = strconv.Atoi(string(rollingFileSizeByte[:len(rollingFileSizeByte)-1]))
		dstValue = int64(tempValue) * 1024 * 1024 * 1024
	case strings.Contains(rollingFileSizeStr, "GB"):
		tempValue, err = strconv.Atoi(string(rollingFileSizeByte[:len(rollingFileSizeByte)-2]))
		dstValue = int64(tempValue) * 1024 * 1024 * 1024
	default:
		err = fmt.Errorf("unit error")
	}

	if err != nil {
		return DefaultFileBytes
	} else {
		return dstValue
	}

}

func NewDefaultOption() Option {
	return Option{
		LogPath:  "./log",
		FileName: "log",

		RollingPolicy:         TimeRolling,
		RollingCronJobPattern: "0 0 0 * * *",
		RollingFileBytes:      DefaultFileBytesStr,

		WriteMode: WriteModeLock,
	}
}
