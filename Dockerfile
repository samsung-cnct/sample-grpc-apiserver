# Dockerfile for gRPC Go
FROM golang:1.7

RUN apt-get update && apt-get -y install unzip && apt-get clean

# install protobuf
ENV PB_VER 3.1.0
ENV PB_URL https://github.com/google/protobuf/releases/download/v${PB_VER}/protoc-${PB_VER}-linux-x86_64.zip
RUN mkdir -p /tmp/protoc && \
    curl -L ${PB_URL} > /tmp/protoc/protoc.zip && \
    cd /tmp/protoc && \
    unzip protoc.zip && \
    cp /tmp/protoc/bin/protoc /usr/local/bin && \
    cp -R /tmp/protoc/include/* /usr/local/include && \
    chmod go+rx /usr/local/bin/protoc && \
    cd /tmp && \
    rm -r /tmp/protoc

# Get the source from GitHub
RUN go get google.golang.org/grpc
# Install protoc-gen-go
RUN go get github.com/golang/protobuf/protoc-gen-go

RUN go get github.com/samsung-cnct/sample-grpc-apiserver/api

RUN mkdir -p sample-grpc-apiserver/
ADD . sample-grpc-apiserver/
RUN cd sample-grpc-apiserver/ && \
    protoc -I api api/api.proto --go_out=plugins=grpc:api 

EXPOSE 5300
ENTRYPOINT go run sample-grpc-apiserver/server/main.go 