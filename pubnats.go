package tailer

import (
	"errors"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/apcera/nats"
)

// NatsPublisher will Publish messages via nats
type NatsPublisher struct {
	URL   string
	topic string
	nc    *nats.Conn
}

var (
	ErrNatsExceedMaxReconnects = errors.New("pubnats: exceed max reconnects")
)

// NewNatsPublisher return a Publisher using nats
func NewNatsPublisher(url string) (Publisher, error) {
	opts := nats.DefaultOptions
	opts.Servers = strings.Split(url, ",")
	opts.MaxReconnect = 10
	opts.ReconnectWait = (1 * time.Second)

	nc, err := opts.Connect()
	if err != nil {
		log.Errorf("error to connect to %s: %v", url, err)
		return nil, err
	}
	nc.Opts.DisconnectedCB = func(nc *nats.Conn) {
		log.Warningf("Got disconnected! Reconnects: %d", nc.Reconnects)
	}
	nc.Opts.ReconnectedCB = func(nc *nats.Conn) {
		log.Warningf("Got reconnected to %v!", nc.ConnectedUrl())
	}
	nc.Opts.ClosedCB = func(_ *nats.Conn) {
		log.Fatal("Got closed")
	}
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}
	return &NatsPublisher{
		URL:   url,
		topic: strings.Replace(hostname, "-", ".", -1),
		nc:    nc,
	}, nil
}

// SetTopic sets the publish topic
func (n *NatsPublisher) SetTopic(topic string) {
	n.topic = topic
}

func (n *NatsPublisher) isExceedMaxReconnects() bool {
	if uint64(n.nc.Opts.MaxReconnect) < n.nc.Reconnects {
		return true
	}
	return false
}

// Publish publish the message to server
func (n *NatsPublisher) Publish(msg []byte) error {
	log.Debugf("publish %s with topic %s", msg, n.topic)
	err := n.nc.Publish(n.topic, msg)
	if err != nil {
		if n.isExceedMaxReconnects() {
			return ErrNatsExceedMaxReconnects
		}
		return err
	}
	return nil
}

// Close close the channel
func (n *NatsPublisher) Close() {
	n.nc.Close()
}
