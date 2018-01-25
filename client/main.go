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
	"flag"
	"log"
	"os"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "../api"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:5300", "The server address in the format of host:port")
	defaultName = flag.String("default_name", "trident", "Default service name")
)


func main() {
	log.Println("Client starting")

	// Set up a connection to the server.
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPingPoseidonClient(conn)

	// Contact the server and print out its response.
	name := *defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	r, err := c.GetPoseidon(context.Background(), &pb.HelloPoseidon{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	r, err = c.GetPoseidonAgain(context.Background(), &pb.HelloPoseidon{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}