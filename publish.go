package tailer

import "fmt"

// Publisher will Publish messages
type Publisher interface {
	SetTopic(topic string)
	Publish(message []byte) error
	Close()
}

// SimplePublisher use fmt.Println to publish
type SimplePublisher struct {
}

// SetTopic NOOP for fmt
func (f *SimplePublisher) SetTopic(topic string) {}

// Publish call fmt.Println to print the msg
func (f *SimplePublisher) Publish(msg []byte) error {
	_, err := fmt.Println(string(msg))
	return err
}

// Close NOOP for fmt
func (f *SimplePublisher) Close() {
}
