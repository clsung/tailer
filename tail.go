package tailer

import (
	"fmt"
	"os"
	"path/filepath"

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
		s.waitGroup.Done()
	}()
	base := filepath.Base(filename)
	if RegexNotWatch.MatchString(base) {
		glog.Warningf("Skip %s", filename)
		return
	}
	glog.Warningf("Tail %s", filename)
	t, err := tail.TailFile(filename, tail.Config{
		// at least one line
		Follow: true, Location: &tail.SeekInfo{-1, os.SEEK_END},
	})
	if err != nil {
		glog.Errorf("error: %v", err)
		return
	}
	s.waitGroup.Add(1)
	for line := range t.Lines {
		err = s.publisher.Publish([]byte(fmt.Sprintf("%s: %s", base, line.Text)))
		if err != nil {
			glog.Errorf("error: %v", err)
		}
	}
	err = t.Wait()
	if err != nil {
		glog.Errorf("error: %v", err)
	}
}
