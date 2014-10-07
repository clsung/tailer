package main

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/apcera/nats"
)

func main() {
	opts := nats.DefaultOptions
	natsURL := os.Getenv("NATS_CLUSTER")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	opts.Servers = strings.Split(natsURL, ",")
	opts.MaxReconnect = 5
	opts.ReconnectWait = (20 * time.Second)
	nc, err := opts.Connect()
	log.SetFlags(0)

	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	nc.Opts.DisconnectedCB = func(_ *nats.Conn) {
		log.Println("Got disconnected!")
	}
	nc.Opts.ReconnectedCB = func(nc *nats.Conn) {
		log.Printf("Got reconnected to %v!\n", nc.ConnectedUrl())
	}

	done := make(chan struct{})
	regexNotNeed, err := regexp.Compile("(?:queue/queue.go|\\[Polling\\])")
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) == 1 {
		log.Fatal("Need specify the topic")
	}
	topic := os.Args[1]
	nc.Subscribe(topic, func(m *nats.Msg) {
		if !regexNotNeed.Match(m.Data) {
			log.Printf("[%s]: %s\n", m.Subject, string(m.Data))
		}
	})

	<-done
}
