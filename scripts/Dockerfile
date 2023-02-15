FROM golang:1.19.3-alpine as builder

RUN apk --no-cache add git ca-certificates gcc g++

WORKDIR /go/src/github.com/go-sonic/

RUN git clone --recursive --depth 1 https://github.com/go-sonic/sonic.git

WORKDIR /go/src/github.com/go-sonic/sonic

RUN GOPROXY=https://goproxy.cn CGO_ENABLED=1 GOOS=linux go build -o sonic -ldflags="-s -w" -trimpath .

RUN mkdir -p /app/conf && \
    mkdir /app/resources && \
    cp -r /go/src/github.com/go-sonic/sonic/sonic /app/ && \
    cp -r /go/src/github.com/go-sonic/sonic/conf /app/ && \
    cp -r /go/src/github.com/go-sonic/sonic/resources /app/ && \
    cp /go/src/github.com/go-sonic/sonic/scripts/docker_init.sh /app/

FROM alpine:latest as prod

COPY --from=builder /app/ /app/

RUN apk update \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

VOLUME /sonic
EXPOSE 8080

WORKDIR /sonic
CMD /app/docker_init.sh && /app/sonic -config /sonic/conf/config.yaml