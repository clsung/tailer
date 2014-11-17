package tailer

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/apcera/nats"
	"github.com/golang/glog"
)

const fileLayout = "200601021504"

// FileEmitter print to file, seperated by time
type FileEmitter struct {
	out    *os.File
	dir    string
	ticker *time.Ticker
	fLock  sync.Mutex
}

func NewFileEmitter(directory string) (Emitter, error) {
	var err error
	directory, err = filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	f := &FileEmitter{dir: directory}
	return f, nil
}

func (f *FileEmitter) Start() {
	fmt.Printf("Start fileEmitter")
	f.ticker = time.NewTicker(time.Minute)
	f.rotate()
	go func() {
		for _ = range f.ticker.C {
			f.rotate()
		}
	}()
}

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
		glog.Errorf("Emit error: %v", err)
		return err
	}
	return nil
}

func (f *FileEmitter) rotate() (err error) {
	filePath := fmt.Sprintf("%s/%s.log", f.dir, time.Now().UTC().Format(fileLayout))

	f.fLock.Lock()
	defer f.fLock.Unlock()
	if f.out != nil {
		f.out.Sync()
		err = f.out.Close()
		if err != nil {
			glog.Errorf("Close File error: %v", err)
			return err
		}
	}
	f.out, err = os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		glog.Errorf("Create/Append File error: %v", err)
		return err
	}
	return nil
}
