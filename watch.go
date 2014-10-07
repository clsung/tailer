package tailer

import (
	"strings"

	"github.com/golang/glog"

	"gopkg.in/fsnotify.v1"
)

var (
	// WatchSuffix sets the file extension to watch
	WatchSuffix = ".log"
)

// WatchDir watches new files added to the dir, and start another tail for it
func WatchDir(path string, fHandleMap map[string]interface{}) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if strings.HasSuffix(ev.Name, WatchSuffix) {
					if ev.Op&fsnotify.Create == fsnotify.Create {
						if f, ok := fHandleMap["onCreate"]; ok {
							f.(func(string) error)(ev.Name)
						} else {
							glog.Warningf("TODO: create event: %s", ev.Name)
						}
					} else if ev.Op&fsnotify.Write == fsnotify.Write {
						if f, ok := fHandleMap["onWrite"]; ok {
							f.(func(string) error)(ev.Name)
						} else {
							glog.Warningf("TODO: write event: %s", ev.Name)
						}
					} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
						if f, ok := fHandleMap["onRemove"]; ok {
							f.(func(string) error)(ev.Name)
						} else {
							glog.Warningf("TODO: remove event: %s", ev.Name)
						}
					}
				}
			case err := <-watcher.Errors:
				glog.Errorf("error: %v", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		glog.Fatal(err)
	}

	<-done
}
