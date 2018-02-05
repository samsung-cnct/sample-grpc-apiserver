#!/bin/sh

set -e

apt-get install apt-utils -y
apt-get install unzip -y
apt-get install dh-autoreconf -y

# Install Protocol Buffers 3
git clone https://github.com/google/protobuf -b $PROTOBUF_TAG --depth 1
cd protobuf
./autogen.sh || exit 1
./configure --prefix=/usr || exit 1
make -j 3 || exit 1
make check || exit 1
make install || exit 1

cd ..
rm -rf protobuf


go get -u -v github.com/golang/protobuf/proto || exit 1
go get -u -v github.com/golang/protobuf/protoc-gen-go || exit 1
go get -u -v google.golang.org/grpc || exit 1
go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway || exit 1
go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger || exit 1
go get -u -v github.com/pseudomuto/protoc-gen-doc/cmd/... || exit 1



# static analysis tools
## goimports
go get golang.org/x/tools/cmd/goimports
## goconst
go get github.com/jgautheron/goconst/cmd/goconst
## gocyclo
go get github.com/fzipp/gocyclo
## gogas
go get github.com/GoASTScanner/gas/cmd/gas/...
## golint
go get -u github.com/golang/lint/golint
##
go get honnef.co/go/tools/cmd/gosimple
