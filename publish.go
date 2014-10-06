package tailer

import "fmt"

type Publisher interface {
	SetTopic(topic string)
	Publish(message []byte) error
	Close()
}

type SimplePublisher struct {
}

func (f *SimplePublisher) SetTopic(topic string) {}
func (f *SimplePublisher) Publish(msg []byte) error {
	_, err := fmt.Println(string(msg))
	return err
}

func (f *SimplePublisher) Close() {
}
