package tailer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ActiveState/tail"
	"github.com/golang/glog"
)

// TailFile tail -f the file and emit with publisher
func TailFile(pub Publisher, filename string, done chan bool) {
	defer func() { done <- true }()
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
	for line := range t.Lines {
		err = pub.Publish([]byte(fmt.Sprintf("%s: %s", base, line.Text)))
		if err != nil {
			glog.Errorf("error: %v", err)
		}
	}
	err = t.Wait()
	if err != nil {
		glog.Errorf("error: %v", err)
	}
}
