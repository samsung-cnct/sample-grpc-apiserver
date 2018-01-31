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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/samsung-cnct/sample-grpc-apiserver/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := configs.InitEnvVars()
	if err != nil {
		log.Fatalf("failed to init config vars: %s", err)
	}

	log.Println("Server starting")

	gwPort, port, gwAddress, address := c.ParseGateWayEnvVars()

	serverAddr := fmt.Sprintf("%s:%d", address, port)
	gwServerAddr := fmt.Sprintf("%s:%d", gwAddress, gwPort)

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

	log.Printf("Setting up GPRC to listen on : %s", serverAddr)

	conn, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	go grpcServer.Serve(conn)

	mux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterPingPoseidonHandlerFromEndpoint(ctx, mux, serverAddr, dialOpts)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Printf("Setting up and Serving REST on : %v", gwServerAddr)

	err = http.ListenAndServe(gwServerAddr, mux)
	if err != nil {
		log.Fatalf("failed to listen restful endpoint: %v", err)
	}
}
