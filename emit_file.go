package tailer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/apcera/nats"
)

const fileLayout = "200601021504"

// FileEmitter print to file, separated by time
type FileEmitter struct {
	out    *os.File
	dir    string
	prefix string
	ticker *time.Ticker
	fLock  sync.Mutex
}

// NewFileEmitter return a new FileEmitter
func NewFileEmitter(directory, prefix string) (Emitter, error) {
	var err error
	directory, err = filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	f := &FileEmitter{dir: directory, prefix: prefix}
	return f, nil
}

// Start starts the fileEmitter
func (f *FileEmitter) Start() (err error) {
	fmt.Printf("Start fileEmitter")
	f.ticker = time.NewTicker(time.Minute)
	err = f.rotate()
	if err != nil {
		log.Errorf("Start error: %v", err)
		return err
	}
	go func() {
		for _ = range f.ticker.C {
			err = f.rotate()
			if err != nil {
				log.Errorf("Ticker error: %v", err)
				break
			}
		}
	}()
	return err
}

// Start stops the fileEmitter
func (f *FileEmitter) Stop() {
	fmt.Printf("Stop fileEmitter")
	if f.out != nil {
		f.out.Sync()
	}
	f.ticker.Stop()
}

// Emit call print the nats.Msg
func (f *FileEmitter) Emit(m *nats.Msg) (err error) {
	msg := fmt.Sprintf("[%s]: %s\n", m.Subject, m.Data)
	f.fLock.Lock()
	_, err = f.out.WriteString(msg)
	f.fLock.Unlock()
	if err != nil {
		log.Errorf("Emit error: %v", err)
		return err
	}
	return nil
}

func (f *FileEmitter) rotate() (err error) {
	filePath := fmt.Sprintf("%s/%s%s.log", f.dir, f.prefix, time.Now().UTC().Format(fileLayout))
	fmt.Printf("Write to %s\n", filePath)

	f.fLock.Lock()
	defer f.fLock.Unlock()
	if f.out != nil {
		f.out.Sync()
		err = f.out.Close()
		if err != nil {
			log.Errorf("Close File error: %v", err)
			return err
		}
	}
	f.out, err = os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("Create/Append File error: %v", err)
		return err
	}
	_, err = f.out.Seek(0, os.SEEK_END)
	if err != nil {
		log.Errorf("Seek File error: %v", err)
		return err
	}
	return nil
}
