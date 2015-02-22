package tailer

// Subscriber will Subscribe messages and do things
type Subscriber interface {
	Subscribe(topic string) error
	Bind(mq *MessageQueue)
	Close()
	LastError() error
}

// SimpleSubscriber use fmt.Println to handle the received msg
type SimpleSubscriber struct {
}

// Subscribe for specified topic
func (f *SimpleSubscriber) Subscribe(topic string) error { return nil }

// Bind to the specified message queue
func (f *SimpleSubscriber) Bind(mq *MessageQueue) {}

// Close NOOP for fmt
func (f *SimpleSubscriber) Close() {
}

// LastError stats the latest error
func (f *SimpleSubscriber) LastError() error { return nil }
