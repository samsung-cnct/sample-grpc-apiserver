# Sample GRPC Apiservice

## Setting Up Your Environment
Keep in mind that gRPC requires golang v1.6 or greater, to install run

```
go get -u google.golang.org/grpc
```

and also run these commands for our other required packages

```
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
```

Finally install latest version of protobuf from release page [https://github.com/google/protobuf/releases](https://github.com/google/protobuf/releases) which
contains the `protoc` binary that must be copied/moved to your `PATH`, e.g., `/usr/local/bin`.

## Generating API, etc...,
You may generate all things that can be generated, including documentation, by running
```
./bin/gen_apis.sh
```

which should generate api and restful-api stubs, docs, and swagger documentation. If you wish
to manually generate any of these, we recommend viewing the subsections that follow.

### Generating Golang API
To build the api (generate go code from proto file) you will want to run a command like
```
protoc -I api/ api/api.proto -I ./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:api
```

where in general you would run

```
protoc -I <input_directory> <path_to_file> --go_out=plugins=grpc:<output_directory>
```

### Generating Golang Reverse Proxy REST API
```
protoc -I api/ api/api.proto
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
--grpc-gateway_out=logtostderr=true:api
```

### Generated Documentation
We can follow along
[https://github.com/pseudomuto/protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc). You
can
follow along the docker commands (recommended) or run locally (assuming you have golang set up
properly)
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

### Swagger docs via REST Gateway
Once you have created a rest-gateway server you can access swagger documentation by following a url such as
```
http://localhost:6300/swagger/api.swagger.json
```
 
## Testing
You can run all the suite of tests, including analysis tools
via
```$shell
./bin/test.sh
```

If the above fails, it will exit with a code of 1, indicating failure.
Please run the above and make sure all tests succeed before 
pushing your PR.

You can also run individual static analysis tests via:
```$shell
./bin/verify.sh
```

### Cleaning
If `goftm` fails, feel free to run
```$shell
./bin/clean/gofmt-clean.sh
./bin/clean/goimports-clean.sh
```

and that should clean up any formatting issues.
