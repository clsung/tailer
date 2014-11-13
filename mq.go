package tailer

import (
	"sync/atomic"

	"github.com/apcera/nats"
)

type MessageQueue struct {
	messages    chan *nats.Msg
	length      int64
	maxMessages int64
}

func NewMessageQueue(maxMessages int64) *MessageQueue {
	return &MessageQueue{
		messages:    make(chan *nats.Msg, maxMessages),
		maxMessages: maxMessages,
	}
}

func (mq *MessageQueue) Push(m *nats.Msg) int64 {
	mq.messages <- m
	atomic.AddInt64(&mq.length, 1)
	return mq.length
}

func (mq *MessageQueue) Pop() (*nats.Msg, bool) {
	select {
	case m := <-mq.messages:
		atomic.AddInt64(&mq.length, -1)
		return m, true
	default:
		return nil, false
	}
}

func (mq *MessageQueue) Len() int64 {
	return atomic.LoadInt64(&mq.length)
}
