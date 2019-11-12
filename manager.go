package bblog

import (
	"github.com/robfig/cron"
	"strings"
	"sync"
	"time"
)

// Manager used to trigger rolling event.
type Manager struct {
	startAt time.Time
	cronJob *cron.Cron

	rollingFileSize int64

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
	m.cronJob.Stop()
}

func (m *Manager) ComputeRollingFileSize(cfg *Config) {
	upStr := strings.ToUpper(cfg.RollingFileSize)

	var ifValidFileSize bool
	ifValidFileSize = strings.Contains(upStr, "K")
	ifValidFileSize = ifValidFileSize || strings.Contains(upStr, "M")
	ifValidFileSize = ifValidFileSize || strings.Contains(upStr, "G")
	// byte unit
	ifValidFileSize = ifValidFileSize || strings.Contains(upStr, "KB")
	ifValidFileSize = ifValidFileSize || strings.Contains(upStr, "MB")
	ifValidFileSize = ifValidFileSize || strings.Contains(upStr, "GB")
	if !ifValidFileSize {
		// default rolling file size is 512MB
		m.rollingFileSize = 1024 * 1024 * 512
		return
	}

}

func (m *Manager) makeTheLogFileName(cfg *Config) (logFileName string) {
	m.lock.Lock()
	logFileName = cfg.LogFilePath() + ".log." + m.startAt.Format(cfg.RollingTimePattern)
	m.startAt = time.Now()
	m.lock.Unlock()
	return
}
