package tailer

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/apcera/nats"
)

const group = "workers"

// NatsSubscriber will Subscribe messages via nats
type NatsSubscriber struct {
	url string
	mq  *MessageQueue
	nc  *nats.Conn
}

// NewNatsSubscriber return a Subscriber using nats
func NewNatsSubscriber(url string) (Subscriber, error) {
	opts := nats.DefaultOptions
	opts.Servers = strings.Split(url, ",")
	opts.MaxReconnect = 10
	opts.ReconnectWait = (1 * time.Second)

	nc, err := opts.Connect()
	if err != nil {
		return nil, fmt.Errorf("error to connect to %s: %v", url, err)
	}
	nc.Opts.ClosedCB = func(_ *nats.Conn) {
		log.Print("Got closed")
	}
	nc.Opts.ReconnectedCB = func(nc *nats.Conn) {
		log.Printf("Got reconnected to %v!\n", nc.ConnectedUrl())
	}
	nc.Opts.AsyncErrorCB = func(nc *nats.Conn, s *nats.Subscription, err error) {
		log.Printf("Got asyncerror %v, %v: %v!\n", nc.ConnectedUrl(), s, err)
	}
	return &NatsSubscriber{
		url: url,
		nc:  nc,
	}, nil
}

// Subscribe the topic and push to the message queue
func (n *NatsSubscriber) Subscribe(topic string) error {
	if n.mq == nil {
		return errors.New("subnats: message queue not bind")
	}
	_, err := n.nc.QueueSubscribe(topic, group, func(m *nats.Msg) {
		n.mq.Push(m)
	})
	if err != nil {
		log.Printf("subscribe error: %v", err)
	}
	return err
}

// Bind to the specified message queue
func (n *NatsSubscriber) Bind(mq *MessageQueue) {
	n.mq = mq
	n.nc.Opts.DisconnectedCB = func(_ *nats.Conn) {
		log.Printf("Got disconnected! Queued %d messagse\n", n.mq.Len())
	}
}

// Close close the channel
func (n *NatsSubscriber) Close() {
	n.nc.Close()
}

// LastError stats the latest error
func (n *NatsSubscriber) LastError() (err error) {
	err = n.nc.LastError()
	if err != nil {
		return err
	}
	return
}
