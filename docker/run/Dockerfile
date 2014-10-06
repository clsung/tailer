FROM google/debian:wheezy

MAINTAINER Cheng-Lung Sung <clsung@gmail.com>

ADD tailer /tailer

ENV NATS_CLUSTER nats://localhost:4222/
VOLUME /tmp
VOLUME /mnt/extend-disk/tmp

CMD ["/tailer", "--nats", "/tmp,/mnt/extend-disk/tmp"]
