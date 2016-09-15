FROM alpine:3.4

MAINTAINER Cheng-Lung Sung <clsung@gmail.com>

ENV GLIBC_VERSION 2.23-r3

# install glibc for executable
RUN apk add --update curl && \
    curl -Lo /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub && \
    curl -Lo glibc.apk "https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-${GLIBC_VERSION}.apk" && \
    apk add glibc.apk && \
    curl -Lo glibc-bin.apk "https://github.com/sgerrand/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-bin-${GLIBC_VERSION}.apk" && \
    apk add glibc-bin.apk && \
    /usr/glibc-compat/sbin/ldconfig /lib /usr/glibc-compat/lib && \
    echo 'hosts: files mdns4_minimal [NOTFOUND=return] dns mdns4' >> /etc/nsswitch.conf &&\
    apk del curl &&\
    rm -rf glibc.apk glibc-bin.apk /var/cache/apk/*


ADD tailer /tailer
ADD config.json /config.json

ENV NATS_CLUSTER nats://localhost:4222/
ENV TAILER_CONFIG ""
VOLUME /tmp
VOLUME /mnt/extend-disk/tmp
VOLUME /csi/logs

CMD ["/tailer", "--nats", "/tmp,/mnt/extend-disk/tmp,/csi/logs"]
