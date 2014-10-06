package tailer

import (
	"log"
	"strings"

	"gopkg.in/fsnotify.v1"
)

var (
	// WatchSuffix sets the file extension to watch
	WatchSuffix = ".log"
)

// WatchDir watches new files added to the dir, and start another tail for it
func WatchDir(onCreate func(string) error, onRemove func(string) error, path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if strings.HasSuffix(ev.Name, WatchSuffix) {
					if ev.Op&fsnotify.Create == fsnotify.Create {
						onCreate(ev.Name)
					} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
						onRemove(ev.Name)
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
