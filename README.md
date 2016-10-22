# tailer

`tailer` is a package for Go that tail -f the files in the watching
directory, and publish to a NATS server.

[![Go Report Card](https://goreportcard.com/badge/github.com/clsung/tailer)](https://goreportcard.com/report/github.com/clsung/tailer)
## Installation and Usage

Install using `go get github.com/clsung/tailer/cmd/tailer`.

Full documentation is available at
http://godoc.org/github.com/clsung/tailer

Below is an example of its usage which tail file to local nats server.

```bash
% export NATS_CLUSTER=nats://localhost:4222/
% tailer --nats /tmp
```

Also you can use docker image to run it. Say we want to
watch /mnt/extend-disk/tmp:
```bash
% docker run -e HOSTNAME=$HOSTNAME -e NATS_CLUSTER=$NATS_CLUSTER \
	-v /mnt/extend-disk/tmp:/tmp -d clsung/tailer
````
