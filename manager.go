package bblog

import (
	"github.com/robfig/cron"
	"os"
	"sync"
	"time"
)

// Manager used to trigger rolling event.
type Manager struct {
	currentLogStartTime time.Time
	rollingCronJob      *cron.Cron

	rollingChan chan string
	contextChan chan int

	wg   sync.WaitGroup
	lock sync.Mutex
}

func (m *Manager) Rolling() chan string {
	return m.rollingChan
}

func (m *Manager) Close() {
	close(m.contextChan)
	m.rollingCronJob.Stop()
}

func (m *Manager) newLogFileName(opt *Option) (logFileName string) {
	m.lock.Lock()
	logFileName = opt.LogFilePath() + ".log." + m.currentLogStartTime.Format(DefaultFileTagFormat)
	m.currentLogStartTime = time.Now()
	m.lock.Unlock()
	return
}

func NewManager(opt *Option) (*Manager, error) {
	m := &Manager{
		currentLogStartTime: time.Now(),
		rollingCronJob:      cron.New(),

		wg:   sync.WaitGroup{},
		lock: sync.Mutex{},

		rollingChan: make(chan string),
		contextChan: make(chan int),
	}

	switch opt.RollingPolicy {
	case PolicyWithoutRolling:
		return m, nil
	case PolicyTimeRolling:
		err := m.rollingCronJob.AddFunc(opt.RollingCronJobPattern, func() {
			m.rollingChan <- m.newLogFileName(opt)
		})
		if err != nil {
			return nil, err
		} else {
			m.rollingCronJob.Start()
			return m, nil
		}
	case PolicyFileSizeRolling:
		m.wg.Add(1)
		defer m.wg.Wait()

		go func() {
			timer := time.Tick(time.Duration(DefaultFileSizeCheckDuration) * time.Second)
			logFilePath := opt.LogFilePath()
			m.wg.Done()

			for {
				select {
				// msg from writer, quit
				case <-m.contextChan:
					return
				case <-timer:
					if fileInfo, err := os.Stat(logFilePath); err != nil {
						continue
					} else {
						if fileInfo.Size() > opt.ComputeRollingFileSize() {
							m.rollingChan <- m.newLogFileName(opt)
						}
					} // else >>>>
				} // select >>>
			} // for >>
		}() // go func >

		return m, nil
	default:
		return m, nil
	}
}
