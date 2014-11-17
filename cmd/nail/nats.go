package main

import (
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/apcera/nats"
	"github.com/clsung/tailer"
)

var (
	logDir    = flag.String("nail_log_dir", "", "Directory to store logs.")
	regExp    = flag.String("nail_filter_pattern", "", "pattern to filter the regex")
	hasStdout = flag.Bool("nail_stdout", true, "send to stdout")
)

func main() {
	flag.Parse()
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
	if len(flag.Args()) < 1 {
		log.Fatal("Need specify the topic")
	}
	topic := flag.Arg(0)
	// TODO: make it flag
	mq := tailer.NewMessageQueue(1 << 20)
	nc.Subscribe(topic, func(m *nats.Msg) {
		mq.Push(m)
	})
	var matchLine *regexp.Regexp
	if *regExp != "" {
		log.Printf("Filter line by regex: %s", *regExp)

		matchLine, err = regexp.Compile(*regExp)
		if err != nil {
			log.Fatalf("regex %s error:%v", *regExp, err)
		}
	}
	emitters := make([]tailer.Emitter, 0)
	if *hasStdout {
		emitters = append(emitters, &tailer.StdoutEmitter{})
	}
	if *logDir != "" {
		fileEmitter, err := tailer.NewFileEmitter(*logDir)
		if err != nil {
			log.Fatalf("create file emitter error %v", err)
		}
		fileEmitter.Start()
		emitters = append(emitters, fileEmitter)
		defer fileEmitter.Stop()
	}

	go func() {
		for {
			m, ok := mq.Pop()
			if ok {
				// TODO: filter logic put here
				if matchLine != nil && !matchLine.Match(m.Data) {
					// skip unmatched
					continue
				}
				for _, emitter := range emitters {
					emitter.Emit(m)
				}
			}
		}
	}()

	wg.Wait()
}
