package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	 pb "../api"
	"google.golang.org/grpc/reflection"
	"strings"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"crypto/tls"
	"crypto/x509"
	"../certs"
	"google.golang.org/grpc/credentials"
	"net/http"
	"fmt"
)

var (
  	serverAddr = flag.String("server_addr", "127.0.0.1:8155", "The server address in the format of host:port")
	demoKeyPair  *tls.Certificate
	demoCertPool *x509.CertPool
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			fmt.Printf("\n\ntrying to serve http 1.1 request \nreponse writer: %+v \nrequest: %+v", w, r)
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func setupCerts() {
	pair, err := tls.X509KeyPair([]byte(certs.Cert), []byte(certs.Key))
	if err != nil {
		panic(err)
	}
	demoKeyPair = &pair
	demoCertPool = x509.NewCertPool()
	ok := demoCertPool.AppendCertsFromPEM([]byte(certs.Cert))
	if !ok {
		panic("bad certs")
	}
}


// GetPoseidon implements helloworld.GreeterServer
func (s *server) GetPoseidon(ctx context.Context, in *pb.HelloPoseidon) (*pb.PoseidonReply, error) {
	return &pb.PoseidonReply{Message: "Hello " + in.Name}, nil
}

func (s *server) GetPoseidonAgain(ctx context.Context, in *pb.HelloPoseidon) (*pb.PoseidonReply, error) {
        return &pb.PoseidonReply{Message: "Hello again " + in.Name}, nil
}

func main() {

	setupCerts()

	log.Println("Server starting")

	opts := []grpc.ServerOption{grpc.Creds(credentials.NewClientTLSFromCert(demoCertPool, *serverAddr))}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPingPoseidonServer(grpcServer, &server{})
	ctx := context.Background()

	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: *serverAddr,
		RootCAs:    demoCertPool,
	})

	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux()
	err := pb.RegisterPingPoseidonHandlerFromEndpoint(ctx, gwmux, *serverAddr, dopts)
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", gwmux)


	//  Get notified that server is being asked to stop
	// Handle SIGINT and SIGTERM.
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)

	// Chance here to gracefully handle being stopped.
	go func() {
	    sig := <-gracefulStop
	    log.Printf("caught sig: %+v", sig)
	    log.Println("Wait for 2 second to finish processing")
	    time.Sleep(2*time.Second)
		grpcServer.Stop()
	    log.Print("service terminated")
	    os.Exit(0)
	}()

	conn, err := net.Listen("tcp", *serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := &http.Server{
		Addr: *serverAddr,
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*demoKeyPair},
			NextProtos:   []string{"h2"},
		},
	}

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	if err := srv.Serve(tls.NewListener(conn, srv.TLSConfig)); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
