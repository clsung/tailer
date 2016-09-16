package tailer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hpcloud/tail"
)

// observe file and add to tailer
func (s *Tailer) addToTail(filePath string) error {
	go s.tailFile(filePath)
	return nil
}

// TailFile tail -f the file and emit with publisher
func (s *Tailer) tailFile(filename string) {
	defer func() {
		log.Warningf("Stop %s", filename)
		s.waitGroup.Done()
	}()
	s.waitGroup.Add(1)
	base := filepath.Base(filename)
	if RegexNotWatch.MatchString(base) {
		log.Warningf("Skip %s", filename)
		return
	}
	log.Warningf("Tail %s, %d are watched", filename, atomic.LoadInt64(&s.numOfTail))
	t, err := tail.TailFile(filename, tail.Config{
		Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END},
		Poll: s.polling,
	})
	if err != nil {
		log.Errorf("initial tail file error: %v", err)
		return
	}
	s.fileLock.Lock()
	atomic.AddInt64(&s.numOfTail, 1)
	s.fileLock.Unlock()
	go func() {
		time.Sleep(50 * time.Millisecond)
		for {
			// FIXME: if files being watched is more then xxx, then we do this
			// every 60 mins, we check if the file is not appended anymore, then unsubscribe it
			fi, err := os.Stat(filename)
			if err != nil {
				t.Killf("stat file %s error: %v", filename, err)
				break
			}
			if time.Now().After(fi.ModTime().Add(time.Hour)) {
				log.Warningf("unwatch modified time > 60 minutes: %s", filename)
				t.Kill(nil)
				break
			} else {
				time.Sleep(time.Hour)
			}
		}
	}()
	for line := range t.Lines {
		if s.matchLine != nil && !s.matchLine.MatchString(line.Text) {
			continue
		} else {
			err = s.publisher.Publish([]byte(fmt.Sprintf("%s: %s", base, line.Text)))
			if err != nil {
				if err == ErrNatsExceedMaxReconnects {
					log.Fatalf("publish error: %v", err)
				} else {
					log.Errorf("publish error: %v", err)
				}
			}
		}
	}
	err = t.Wait()
	if err != nil {
		log.Errorf("wait error: %v", err)
	}
	s.fileLock.Lock()
	atomic.AddInt64(&s.numOfTail, -1)
	s.fileLock.Unlock()
}
