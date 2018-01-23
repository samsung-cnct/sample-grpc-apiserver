# Setting Up Your Environment
Keep in mind that gRPC requires golang v1.6 or greater, to install run

```
go get -u google.golang.org/grpc
```

and also run

```
go get -u github.com/golang/protobuf/protoc-gen-go
```

Finally install latest version of protobuf from release page [https://github.com/google/protobuf/releases](https://github.com/google/protobuf/releases) which
contains the `protoc` binary that must be copied/moved to your `PATH`, e.g., `/usr/local/bin`.



# Generating Golang API
To build the api (generate go code from proto file) you will want to run a command like
```
protoc -I api/ api/api.proto -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:api
```

where in general you would run

```
protoc -I <input_directory> <path_to_file> --go_out=plugins=grpc:<output_directory>
```

# Generating Golang Reverse Proxy REST API
```
protoc -I api/ api/api.proto -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:api
```

# Generated Documentation
We can follow along [https://github.com/pseudomuto/protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc). You can
follow along the docker commands (recommended) or run locally (assuming you have golang set up properly)
via

```
go get -u github.com/pseudomuto/protoc-gen-doc/cmd/...
```

and, if you did the above, you can simply run

```
protoc --doc_out docs --doc_opt=markdown,api.md  api/hello.proto
```

where what the above comes from is

```
protoc --doc_out <output_folder> --doc_opt=<format, file_output> <api_proto_input_files>
```
