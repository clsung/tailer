package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
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
	opts.MaxReconnect = 10
	opts.ReconnectWait = (1 * time.Second)
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
	nc.Opts.AsyncErrorCB = func(nc *nats.Conn, s *nats.Subscription, err error) {
		log.Printf("Got asyncerror %v, %v: %v!\n", nc.ConnectedUrl(), s, err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := nc.LastError()
			if err != nil {
				log.Printf("Error: %v", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()
	if len(os.Args) == 1 {
		log.Fatal("Need specify the topic")
	}
	topic := os.Args[1]
	nc.Subscribe(topic, func(m *nats.Msg) {
		fmt.Printf("[%s]: %s\n", m.Subject, string(m.Data))
	})

	wg.Wait()
}
