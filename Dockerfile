FROM golang:latest

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && apt-get install -y \
    bash-completion \
    bind9-host \
    curl \
    iptraf-ng \
    iputils-tracepath \
    mtr \
    net-tools \
    telnet \
    traceroute \
    util-linux \
    wget

WORKDIR /app

COPY main.go .
COPY go.mod .
RUN go build -o app .

CMD ["./app"]
