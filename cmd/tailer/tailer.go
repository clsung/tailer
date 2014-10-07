package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/clsung/tailer"
	"github.com/jessevdk/go-flags"
)

var version = "v0.0.1"

type cmdOpts struct {
	OptVersion    bool   `short:"v" long:"version" description:"print the version and exit"`
	OptNats       bool   `long:"nats" description:"Using nats to publish" default:"false"`
	OptConfigfile string `long:"config" description:"config file" optional:"yes"`
}

func main() {
	var err error
	var exitCode int

	defer func() { os.Exit(exitCode) }()
	if envvar := os.Getenv("GOMAXPROCS"); envvar == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	done := make(chan bool)

	opts := &cmdOpts{}
	p := flags.NewParser(opts, flags.Default)
	p.Usage = "[OPTIONS] DIR1[,DIR2...]"

	args, err := p.Parse()

	if opts.OptVersion {
		fmt.Fprintf(os.Stderr, "tailer: %s\n", version)
		return
	}

	if err != nil || len(args) == 0 {
		p.WriteHelp(os.Stderr)
		exitCode = 1
		return
	}

	config := tailer.Config{FileGlob: "*-*"} // assume abc-def.log
	if opts.OptConfigfile != "" {
		file, err := os.Open(opts.OptConfigfile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Using default config")
			} else {
				fmt.Println("error:", err)
				exitCode = 1
				return
			}
		} else {
			decoder := json.NewDecoder(file)
			err = decoder.Decode(&config)
			if err != nil {
				fmt.Println("error:", err)
				exitCode = 1
				return
			}
		}
	}

	watchDirs := strings.Split(args[0], ",")
	filesToTail := []string{}

	var pub tailer.Publisher
	if opts.OptNats {
		natsURL := os.Getenv("NATS_CLUSTER")
		if natsURL == "" {
			natsURL = "nats://localhost:4222"
		}
		pub, err = tailer.NewNatsPublisher(natsURL)
		if err != nil {
			exitCode = 1
			return
		}
	} else {
		pub = &tailer.SimplePublisher{}
	}

	// observe file and add to tailer
	addToTail := func(filePath string) error {
		tailer.TailFile(pub, filePath, done)
		return nil
	}
	fMap := map[string]interface{}{
		"onCreate": addToTail,
		"onWrite":  addToTail,
	}

	// examine the input dir and select how many files to watch and publish
	for _, dir := range watchDirs {
		fileGlobPattern := fmt.Sprintf("%s/%s", dir, config.FileGlob)
		files, _ := filepath.Glob(fileGlobPattern)
		filesToTail = append(filesToTail, files...)
		log.Printf("Files to watch now: %v", filesToTail)
		go tailer.WatchDir(dir, fMap)
	}

	for _, filePath := range filesToTail {
		tailer.TailFile(pub, filePath, done)
	}

	// TODO: exit if all watched files removed or closed
	<-done
}
