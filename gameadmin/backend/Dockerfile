FROM centos:7

ARG USER_ID
ARG GROUP_ID

RUN groupadd --gid $GROUP_ID zy 
RUN useradd zy -u $USER_ID -g $GROUP_ID -m -s /bin/bash 

ADD https://studygolang.com/dl/golang/go1.19.2.linux-amd64.tar.gz /home/zy
RUN cd /home/zy && tar -xzf ./go1.19.2.linux-amd64.tar.gz && chown -R zy:zy /home/zy/go

ENV PATH /home/zy/go/bin:$PATH

ENV GOPROXY https://goproxy.io,direct
ENV GOCACHE /gocache/cache
ENV GOMODCACHE /gocache/mod
ENV CGO_ENABLED 0

USER zy