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

	"fmt"

	pb "github.com/samsung-cnct/sample-grpc-apiserver/api"
	"github.com/samsung-cnct/sample-grpc-apiserver/configs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Client starting")

	err := configs.InitEnvVars()
	if err != nil {
		log.Fatalf("failed to init config vars: %s", err)
	}

	_, port, _, address := configs.ParseGateWayEnvVars()
	serverAddr := fmt.Sprintf("%s:%d", address, port)
	name := configs.ParseClientEnvVars()

	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPingPoseidonClient(conn)

	r, err := c.HelloPoseidon(context.Background(), &pb.HelloPoseidonMsg{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	r, err = c.HelloPoseidonAgain(context.Background(), &pb.HelloPoseidonMsg{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
}
