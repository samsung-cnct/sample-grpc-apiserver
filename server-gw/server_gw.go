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
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/samsung-cnct/sample-grpc-apiserver/api"
	c "github.com/samsung-cnct/sample-grpc-apiserver/configs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func startServerGW(addr string, gracefulStop chan os.Signal) error {
	var err error
	ctx := context.Background()
	log.Print("starting rest-gateway server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger/", serveSwagger)

	muxGateway := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterPingPoseidonHandlerFromEndpoint(ctx, muxGateway, addr, opts)
	if err != nil {
		return err
	}

	mux.Handle("/", muxGateway)

	// Chance here to gracefully handle being stopped.
	go func() {
		sig := <-gracefulStop
		log.Printf("caught sig: %+v", sig)
		log.Println("waiting for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		log.Print("rest-gateway server terminated")
		os.Exit(0)
	}()

	log.Printf("attempting to start rest-gateway server in address: %s", addr)

	return http.ListenAndServe(addr, mux)

}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	log.Printf("call to find swagger resource.... %s", r.URL.Path)
	if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
		log.Printf("Not a swagger file %s, missing suffix .swagger.json", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	gwSwaggerDir := c.ParseGWSwaggerEnvVars()

	log.Printf("Serving %s", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join(gwSwaggerDir, p)
	http.ServeFile(w, r, p)
}

func main() {
	ctx := context.Background()
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	err := c.InitEnvVars()
	if err != nil {
		log.Fatalf("failed to init config vars: %s", err)
	}

	gwPort, _, gwAddress, _ := c.ParseGateWayEnvVars()

	//  Get notified that server is being asked to stop
	// Handle SIGINT and SIGTERM.
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)

	// Server Gateway Code
	err = startServerGW(fmt.Sprintf("%s:%d", gwAddress, gwPort), gracefulStop)
	if err != nil {
		log.Fatalf("failed to start rest-gateway server: %s", err)
	}

}
