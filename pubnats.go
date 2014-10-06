package tailer

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/apcera/nats"
)

type NatsPublisher struct {
	URL   string
	topic string
	nc    *nats.Conn
}

func NewNatsPublisher(url string) (Publisher, error) {
	opts := nats.DefaultOptions
	opts.Servers = strings.Split(url, ",")
	opts.MaxReconnect = 5
	opts.ReconnectWait = (2 * time.Second)

	nc, err := opts.Connect()
	if err != nil {
		log.Printf("error to connect to %s: %v", url, err)
		return nil, err
	}
	nc.Opts.DisconnectedCB = func(_ *nats.Conn) {
		log.Printf("Got disconnected!\n")
	}
	nc.Opts.ReconnectedCB = func(nc *nats.Conn) {
		log.Printf("Got reconnected to %v!\n", nc.ConnectedUrl())
	}
	hostname, _ := os.Hostname()
	return &NatsPublisher{
		URL:   url,
		topic: strings.Join([]string{"root", hostname}, "."),
		nc:    nc,
	}, nil
}

func (n *NatsPublisher) SetTopic(topic string) {
	n.topic = strings.Join([]string{"root", topic}, ".")
}

func (n *NatsPublisher) Publish(msg []byte) error {
	log.Printf("publish %s with topic %s", msg, n.topic)
	return n.nc.Publish(n.topic, msg)
}

func (n *NatsPublisher) Close() {
	n.nc.Close()
}
