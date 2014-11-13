FROM google/debian:wheezy

MAINTAINER Cheng-Lung Sung <clsung@gmail.com>

ADD tailer /tailer
ADD config.json /config.json

ENV NATS_CLUSTER nats://localhost:4222/
ENV TAILER_CONFIG
VOLUME /tmp
VOLUME /mnt/extend-disk/tmp

CMD ["/tailer", "--nats", "/tmp,/mnt/extend-disk/tmp"]
