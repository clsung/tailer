package tailer

import (
	"log"
	"os"
	"strings"

	"github.com/apcera/nats"
)

type NatsPublisher struct {
	URL   string
	topic string
	//ec       *nats.EncodedConn
	nc *nats.Conn
}

func NewNatsPublisher(url string) (Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Printf("error to connect to %s: %v", url, err)
		return nil, err
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
