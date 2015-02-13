package tailer

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/apcera/nats"
	"github.com/golang/glog"
)

// NatsPublisher will Publish messages via nats
type NatsPublisher struct {
	URL   string
	topic string
	nc    *nats.Conn
}

var (
	ErrNatsConnectionClosed    = nats.ErrConnectionClosed
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
		glog.Errorf("error to connect to %s: %v", url, err)
		return nil, err
	}
	nc.Opts.DisconnectedCB = func(nc *nats.Conn) {
		glog.Warningf("Got disconnected! Reconnects: %d\n", nc.Reconnects)
	}
	nc.Opts.ReconnectedCB = func(nc *nats.Conn) {
		glog.Warningf("Got reconnected to %v!\n", nc.ConnectedUrl())
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
	glog.V(2).Infof("publish %s with topic %s", msg, n.topic)
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
