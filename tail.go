package tailer

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ActiveState/tail"
)

//type emitFunc interface{}

func TailFile(pub Publisher, filename string, done chan bool) {
	defer func() { done <- true }()
	log.Printf("Tail %s", filename)
	t, err := tail.TailFile(filename, tail.Config{Follow: true})
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	base := filepath.Base(filename)
	for line := range t.Lines {
		err = pub.Publish([]byte(fmt.Sprintf("%s: %s", base, line.Text)))
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
	err = t.Wait()
	if err != nil {
		log.Printf("error: %v", err)
	}
}
