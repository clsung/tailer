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
func WatchDir(path string) {
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
						log.Println("TODO: add event to tailer:", ev)
					} else if ev.Op&fsnotify.Rename == fsnotify.Rename {
						log.Println("TODO: remove event from tailer:", ev)
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
