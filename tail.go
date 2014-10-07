package tailer

import (
	"fmt"
	"path/filepath"

	"github.com/ActiveState/tail"
	"github.com/golang/glog"
)

// TailFile tail -f the file and emit with publisher
func TailFile(pub Publisher, filename string, done chan bool) {
	defer func() { done <- true }()
	glog.Warningf("Tail %s", filename)
	t, err := tail.TailFile(filename, tail.Config{Follow: true})
	if err != nil {
		glog.Errorf("error: %v", err)
		return
	}
	base := filepath.Base(filename)
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
