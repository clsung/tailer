#!/bin/bash
docker build -t clsung/tailer_build build/
docker run -d --name tailer_build clsung/tailer_build
docker cp tailer_build:/go/bin/tailer run/
docker stop tailer_build
docker rm tailer_build
docker rmi clsung/tailer_build
docker build -t clsung/tailer run/
