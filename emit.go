package tailer

import (
	"fmt"

	"github.com/apcera/nats"
)

// Emitter will emit Nats messages
type Emitter interface {
	Emit(msg *nats.Msg) error
	Start() error
	Stop()
}

// StdoutEmitter use fmt.Println to publish
type StdoutEmitter struct {
}

// Emit call fmt.Println to print the msg
func (f *StdoutEmitter) Emit(msg *nats.Msg) error {
	fmt.Printf("[%s]: %s\n", msg.Subject, string(msg.Data))
	return nil
}

func (f *StdoutEmitter) Start() error {
	return nil
}

func (f *StdoutEmitter) Stop() {
}
