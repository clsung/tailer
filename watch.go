package tailer

import (
	log "github.com/Sirupsen/logrus"

	"gopkg.in/fsnotify.v1"
)

func (s *Tailer) isUnwantEvent(ev fsnotify.Event) bool {
	if s.visitedFile[ev.Name] {
		return true
	}
	if RegexNotWatch.MatchString(ev.Name) {
		s.visitedFile[ev.Name] = true
		return true
	}
	return false
}

// watchDir watches new files added to the dir, and start another tail for it
func (s *Tailer) watchDir(path string) {
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
				if !s.isUnwantEvent(ev) {
					if ev.Op&fsnotify.Create == fsnotify.Create {
						s.addToTail(ev.Name)
						log.Warningf("TODO: create event: %s", ev.Name)
					} else if ev.Op&fsnotify.Write == fsnotify.Write {
						//log.Warningf("TODO: write event: %s", ev.Name)
					} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Warningf("TODO: remove event: %s", ev.Name)
					}
				}
			case err := <-watcher.Errors:
				log.Errorf("error: %v", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
