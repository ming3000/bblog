package bblog

import (
	"fmt"
	"github.com/robfig/cron"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Manager used to trigger rolling event.
type Manager struct {
	currentLogStartAt time.Time
	rollingCronJob    *cron.Cron
	rollingFileSize   int64

	fireChan    chan string
	contextChan chan int

	wg   sync.WaitGroup
	lock sync.Mutex
}

func (m *Manager) Fire() chan string {
	return m.fireChan
}

func (m *Manager) Close() {
	close(m.contextChan)
	m.rollingCronJob.Stop()
}

func (m *Manager) computeRollingFileSize(cfg *Option) {
	rollingFileSizeStr := strings.ToUpper(cfg.RollingFileSize)
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
		// default rolling file size is 512MB
		m.rollingFileSize = 1024 * 1024 * 512
	} else {
		m.rollingFileSize = dstValue
	}

}

func (m *Manager) makeTheLogFileName(cfg *Option) (logFileName string) {
	m.lock.Lock()
	logFileName = cfg.LogFilePath() + ".log." + m.currentLogStartAt.Format(cfg.RollingTimeTagFormat)
	m.currentLogStartAt = time.Now()
	m.lock.Unlock()
	return
}
