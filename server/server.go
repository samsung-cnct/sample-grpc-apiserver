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
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	pb "github.com/samsung-cnct/sample-grpc-apiserver/api"
	c "github.com/samsung-cnct/sample-grpc-apiserver/configs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// GetPoseidon implements helloworld.GreeterServer
func (s *server) HelloPoseidon(ctx context.Context, in *pb.HelloPoseidonMsg) (*pb.PoseidonReply, error) {
	return &pb.PoseidonReply{Message: "Hello " + in.Name}, nil
}

func (s *server) HelloPoseidonAgain(ctx context.Context, in *pb.HelloPoseidonMsg) (*pb.PoseidonReply, error) {
	return &pb.PoseidonReply{Message: "Hello again " + in.Name}, nil
}

func startServer(addr string, gracefulStop chan os.Signal) error {
	var err error
	log.Print("starting server")

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPingPoseidonServer(grpcServer, &server{})

	// Chance here to gracefully handle being stopped.
	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		log.Println("waiting for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		grpcServer.Stop()
		log.Print("server terminated")
		os.Exit(0)
	}()

	log.Printf("attempting to start server in address: %s", addr)

	return grpcServer.Serve(listener)
}

func main() {
	ctx := context.Background()
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	err := c.InitEnvVars()
	if err != nil {
		log.Fatalf("failed to init config vars: %s", err)
	}

	_, port, _, address := c.ParseGateWayEnvVars()

	//  Get notified that server is being asked to stop
	// Handle SIGINT and SIGTERM.
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)

	// Server Code
	err = startServer(fmt.Sprintf("%s:%d", address, port), gracefulStop)
	if err != nil {
		log.Fatalf("failed to start server: %s", err)
	}
}
