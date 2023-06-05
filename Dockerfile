FROM golang:1.19.3-buster AS builder
MAINTAINER lory <cive148@126.com>

# this file is only for docker build testing

ENV SRC_DIR /gomicro

RUN go env -w GOPROXY=https://goproxy.io

# aliyun pub key
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 3B4FE6ACC0B21F32
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 871920D1991BC93C
# aliyun apt source
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse" > /etc/apt/sources.list
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse" >> /etc/apt/sources.list
RUN echo "deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse" >> /etc/apt/sources.list

# update apt
RUN apt-get clean && apt-get update
RUN apt-get install -y ca-certificates make

COPY . $SRC_DIR
RUN cd $SRC_DIR && export GIT_SSL_NO_VERIFY=true && git config --global http.sslVerify "false" && make

FROM ubuntu:20.04

RUN ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone && apt-get update && apt-get install -y tzdata
ENV TZ Asia/Shanghai
ENV SRC_DIR /gomicro

COPY --from=builder $SRC_DIR/client /usr/local/bin/client
COPY --from=builder $SRC_DIR/server /usr/local/bin/server
COPY --from=builder /etc/ssl/certs /etc/ssl/certs


ENV HOME_PATH /data

VOLUME $HOME_PATH

