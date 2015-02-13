package tailer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/ActiveState/tail"
	"github.com/golang/glog"
)

// observe file and add to tailer
func (s *Tailer) addToTail(filePath string) error {
	go s.tailFile(filePath)
	return nil
}

// TailFile tail -f the file and emit with publisher
func (s *Tailer) tailFile(filename string) {
	defer func() {
		if glog.V(2) {
			glog.Warningf("Stop %s", filename)
		}
		s.waitGroup.Done()
	}()
	s.waitGroup.Add(1)
	base := filepath.Base(filename)
	if RegexNotWatch.MatchString(base) {
		if glog.V(2) {
			glog.Warningf("Skip %s", filename)
		}
		return
	}
	glog.Warningf("Tail %s, %d are watched", filename, atomic.LoadInt64(&s.numOfTail))
	t, err := tail.TailFile(filename, tail.Config{
		Follow: true, Location: &tail.SeekInfo{0, os.SEEK_END},
	})
	if err != nil {
		glog.Errorf("initial tail file error: %v", err)
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
				glog.Warningf("unwatch modified time > 60 minutes: %s", filename)
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
					glog.Fatalf("publish error: %v", err)
				} else {
					glog.Errorf("publish error: %v", err)
				}
			}
		}
	}
	err = t.Wait()
	if err != nil {
		glog.Errorf("wait error: %v", err)
	}
	s.fileLock.Lock()
	atomic.AddInt64(&s.numOfTail, -1)
	s.fileLock.Unlock()
}
