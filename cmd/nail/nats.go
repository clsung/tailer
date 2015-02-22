package main

import (
	"flag"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/clsung/tailer"
)

var (
	logDir    = flag.String("nail_log_dir", "", "Directory to store logs.")
	logPrefix = flag.String("nail_log_prefix", "", "Log file prefix")
	regExp    = flag.String("nail_filter_pattern", "", "pattern to filter the regex")
	hasStdout = flag.Bool("nail_stdout", true, "send to stdout")
	workers   = flag.Int("nail_worker", 2, "queue workers")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	natsURL := os.Getenv("NATS_CLUSTER")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	if len(flag.Args()) < 1 {
		log.Fatal("Need specify the topic")
	}
	topic := flag.Arg(0)
	// TODO: make it flag
	mq := tailer.NewMessageQueue(1 << 20)
	var wg sync.WaitGroup
	var err error
	wg.Add(*workers)
	for i := 1; i <= *workers; i++ {
		subscriber, err := tailer.NewNatsSubscriber(natsURL)
		if err != nil {
			log.Fatal(err)
		}
		subscriber.Bind(mq)
		err = subscriber.Subscribe(topic)
		if err != nil {
			log.Fatal(err)
		}
		defer subscriber.Close()
		go func() {
			defer wg.Done()
			for {
				err := subscriber.LastError()
				if err != nil {
					log.Printf("Error: %v", err)
					break
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}
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
		fileEmitter, err := tailer.NewFileEmitter(*logDir, *logPrefix)
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
