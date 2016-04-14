package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/clsung/tailer"
	"github.com/jessevdk/go-flags"
)

var version = "v0.0.1"

type cmdOpts struct {
	OptVersion    bool   `short:"v" long:"version" description:"print the version and exit"`
	OptNats       bool   `long:"nats" description:"Using nats to publish"`
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
	configFile := opts.OptConfigfile
	if configFile == "" {
		configFile = os.Getenv("TAILER_CONFIG")
	}
	if configFile != "" {
		file, err := os.Open(configFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Using default config")
			} else {
				fmt.Printf("error: %v\n", err)
				exitCode = 1
				return
			}
		} else {
			fmt.Printf("Using config file: %s\n", configFile)
			decoder := json.NewDecoder(file)
			err = decoder.Decode(&config)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				exitCode = 1
				return
			}
		}
	}

	watchDirs := strings.Split(args[0], ",")

	tailer, err := tailer.NewTailer(opts.OptNats, config)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		exitCode = 1
		return
	}
	tailer.Serve(watchDirs, config.FileGlob)

	// TODO: exit if all watched files removed or closed
	<-done
	fmt.Println("Exit with 0")
}
