FROM alpine:3.3

MAINTAINER Cheng-Lung Sung <clsung@gmail.com>

# install glibc for executable
RUN apk add --update curl && \
	cd /tmp && \
    curl -o glibc-2.21-r2.apk \
        "https://circle-artifacts.com/gh/andyshinn/alpine-pkg-glibc/6/artifacts/0/home/ubuntu/alpine-pkg-glibc/packages/x86_64/glibc-2.21-r2.apk" && \
    curl -o glibc-bin-2.21-r2.apk \
        "https://circle-artifacts.com/gh/andyshinn/alpine-pkg-glibc/6/artifacts/0/home/ubuntu/alpine-pkg-glibc/packages/x86_64/glibc-bin-2.21-r2.apk" && \
    apk add --allow-untrusted \
        glibc-2.21-r2.apk \
        glibc-bin-2.21-r2.apk && \
    /usr/glibc/usr/bin/ldconfig /lib /usr/glibc/usr/lib && \
	rm -f glibc-2.21-r2.apk glibc-bin-2.21-r2.apk && \
	apk del curl

ADD tailer /tailer
ADD config.json /config.json

ENV NATS_CLUSTER nats://localhost:4222/
ENV TAILER_CONFIG ""
VOLUME /tmp
VOLUME /mnt/extend-disk/tmp
VOLUME /csi/logs

CMD ["/tailer", "--nats", "/tmp,/mnt/extend-disk/tmp,/csi/logs"]
