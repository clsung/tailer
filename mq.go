package tailer

import (
	"sync/atomic"

	"github.com/apcera/nats"
)

// MessageQueue handles the file events
type MessageQueue struct {
	messages    chan *nats.Msg
	length      int64
	maxMessages int64
}

// NewMessageQueue returns a new MessageQueue
func NewMessageQueue(maxMessages int64) *MessageQueue {
	return &MessageQueue{
		messages:    make(chan *nats.Msg, maxMessages),
		maxMessages: maxMessages,
	}
}

// Push pushs the nats.Msg into queue
func (mq *MessageQueue) Push(m *nats.Msg) int64 {
	mq.messages <- m
	atomic.AddInt64(&mq.length, 1)
	return mq.length
}

// Pop pop out a nats.Msg from queue
func (mq *MessageQueue) Pop() (*nats.Msg, bool) {
	select {
	case m := <-mq.messages:
		atomic.AddInt64(&mq.length, -1)
		return m, true
	default:
		return nil, false
	}
}

// Len returns current MessageQueue length
func (mq *MessageQueue) Len() int64 {
	return atomic.LoadInt64(&mq.length)
}
