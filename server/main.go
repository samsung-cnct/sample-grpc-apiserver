/*
 *
 * Copyright 2018 Samsung SDS Cloud Native Computing Team authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	pb "../api"
	"flag"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	serverAddrGRPC = flag.String("server_addr_grpc", "127.0.0.1:5300", "The server address in the format of host:port")
	serverAddrREST = flag.String("server_addr_rest", "127.0.0.1:6300", "The REST endpoint address, in the format of host:port")
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// GetPoseidon implements helloworld.GreeterServer
func (s *server) GetPoseidon(ctx context.Context, in *pb.HelloPoseidon) (*pb.PoseidonReply, error) {
	return &pb.PoseidonReply{Message: "Hello " + in.Name}, nil
}

func (s *server) GetPoseidonAgain(ctx context.Context, in *pb.HelloPoseidon) (*pb.PoseidonReply, error) {
	return &pb.PoseidonReply{Message: "Hello again " + in.Name}, nil
}

func main() {

	log.Println("Server starting")

	grpcServer := grpc.NewServer()
	pb.RegisterPingPoseidonServer(grpcServer, &server{})

	//  Get notified that server is being asked to stop
	// Handle SIGINT and SIGTERM.
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)

	// Chance here to gracefully handle being stopped.
	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		log.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		grpcServer.Stop()
		log.Print("service terminated")
		os.Exit(0)
	}()

	log.Printf("Setting up GPRC to listen on : %v", *serverAddrGRPC)

	conn, err := net.Listen("tcp", *serverAddrGRPC)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	go grpcServer.Serve(conn)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterPingPoseidonHandlerFromEndpoint(ctx, mux, *serverAddrGRPC, dialOpts)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Printf("Setting up and Serving REST on : %v", *serverAddrREST)

	http.ListenAndServe(*serverAddrREST, mux)
}
