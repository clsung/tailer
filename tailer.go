package tailer

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/golang/glog"
)

var (
	// RegexNotWatch sets the file extension to avoid watching
	RegexNotWatch = regexp.MustCompile("(?:^tailer\\.|^gobzip-|^\\..+\\.swp$|\\.gz$|\\.[0-9]$)")
)

// Tailer init the service functions
type Tailer struct {
	ch        chan bool
	waitGroup *sync.WaitGroup
	publisher Publisher
	matchLine *regexp.Regexp
	numOfTail int64
	fileLock  sync.Mutex
}

// Make a new Tailer
func NewTailer(publishToNats bool, config Config) (*Tailer, error) {
	var err error
	t := &Tailer{
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
	if len(config.Match) > 0 {
		glog.Warningf("Filter line by regex: %s", config.Match)
		t.matchLine, err = regexp.Compile(config.Match)
		if err != nil {
			return nil, err
		}
	}
	if publishToNats {
		natsURL := os.Getenv("NATS_CLUSTER")
		if natsURL == "" {
			natsURL = "nats://localhost:4222"
		}
		t.publisher, err = NewNatsPublisher(natsURL)
		if err != nil {
			return nil, err
		}
	} else {
		t.publisher = &SimplePublisher{}
	}
	return t, nil
}

func (s *Tailer) Serve(watchDirs []string, fileGlob string) {
	// examine the input dir and select how many files to watch and publish
	filesToTail := []string{}
	for _, dir := range watchDirs {
		fileGlobPattern := fmt.Sprintf("%s/%s", dir, fileGlob)
		files, _ := filepath.Glob(fileGlobPattern)
		filesToTail = append(filesToTail, files...)
		glog.Warningf("Files to watch now: %v", filesToTail)
		go s.watchDir(dir)
	}

	for _, filePath := range filesToTail {
		go s.tailFile(filePath)
	}

	s.waitGroup.Wait()
}
