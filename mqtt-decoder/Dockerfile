FROM golang:1.14.0-alpine

RUN apk --no-cache add mosquitto mosquitto-clients openrc bash

ARG PUID=1000
ARG PGID=1000
#
RUN addgroup -g ${PGID} -S influxdb \
    && adduser -u ${PUID} -G influxdb -h /var/lib/influxdb -D influxdb \
    && echo 'hosts: files dns' >> /etc/nsswitch.conf \
    && apk add -Uu --purge --no-cache \
        --repository http://dl-cdn.alpinelinux.org/alpine/edge/community \
        ca-certificates \
        influxdb \
    && chown -R influxdb:influxdb /var/lib/influxdb \
    && s6-setuidgid influxdb exec influxd -config /etc/influxdb.conf \
    && update-ca-certificates \
	&& rm -rf /var/cache/apk/* /tmp/*

VOLUME  ["/var/lib/influxdb"]

ENV GRAFANA_VERSION=v4.3.2

RUN apk add --no-cache --update --virtual .build-deps \
        build-base git libffi-dev linux-headers mercurial nodejs pcre-dev \
        postgresql-dev python3-dev tar \
    && runDeps="ca-certificates curl netcat-openbsd pcre postgresql postgresql-contrib" \
    && ln -s /bin/sh /bin/bash \
    && deluser xfs \
    && addgroup -g 33 www-data \
    && adduser -u 33 -D -G www-data -h /home/www www-data -s /bin/sh \
    && mkdir -p /home/www \
    && chown -R www-data.www-data /home/www \
    && find /usr/local \
       \( -type d -a -name test -o -name tests \) \
       -o \( -type f -a -name '*.pyc' -o -name '*.pyo' \) \
       -exec rm -rf '{}' + \
    && runDeps="$runDeps $( \
       scanelf --needed --nobanner --recursive /usr/local \
               | awk '{ gsub(/,/, "\nso:", $2); print "so:" $2 }' \
               | sort -u \
               | xargs -r apk info --installed \
               | sort -u \
   )" \
   && mkdir -p /go/src/github.com/grafana && cd /go/src/github.com/grafana \
   && git clone https://github.com/grafana/grafana.git -b ${GRAFANA_VERSION} \
   && cd grafana \
   && go run build.go setup \
   && go run build.go build  \
   && npm install \
   && npm install -g grunt-cli \
   && cd /go/src/github.com/grafana/grafana && grunt \
   && npm run build \
   && npm uninstall -g grunt-cli \
   && npm cache clear \
   && mkdir -p /etc/grafana/bin/ \
   && mkdir -p /etc/grafana/dashboard/ \
   && cp -a /go/src/github.com/grafana/grafana/bin /etc/grafana/ \
   && cp -ra /go/src/github.com/grafana/grafana/public_gen /etc/grafana/ \
   && mv /etc/grafana/public_gen /etc/grafana/public \
   && cp -ra /go/src/github.com/grafana/grafana/conf /etc/grafana/ \
   && go clean -i -r \
   && apk add --virtual .run-deps $runDeps \
   && apk del .build-deps \
   && rm -rf /go /tmp/* /var/cache/apk/* /root/.n* /etc/bin/phantomjs


VOLUME ["/var/log/grafana", "/etc/grafana"]

WORKDIR /etc/grafana/
COPY ./docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/etc/grafana/bin/grafana-server"]

EXPOSE 3000
EXPOSE 1883

