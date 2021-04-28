FROM golang:1.13-alpine


FROM golang as build

ENV GOPROXY=https://goproxy.io

ADD . /12306

WORKDIR /12306

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api_server


FROM alpine:3.7

ENV GIN_MODE="release"
ENV PORT=3000

RUN echo "http://mirrors.aliyun.com/alpine/v3.7/main/" > /etc/apk/repositories && \
    apk update && \
    apk add ca-certificates && \
    echo "hosts: files dns" > /etc/nsswitch.conf && \
    mkdir -p /www/config

WORKDIR /www

COPY --from=build /12306/api_server /usr/bin/api_server
COPY --from=0 /usr/local/go/lib/time/zoneinfo.zip /opt/zoneinfo.zip
ENV ZONEINFO /opt/zoneinfo.zip

ADD ./config /www/config
RUN chmod +x /usr/bin/api_server

ENTRYPOINT ["api_server"]
