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
	pb "github.com/samsung-cnct/sample-grpc-apiserver/api"
	"google.golang.org/grpc/reflection"
)

var (
  serverPort = flag.String("server_port", ":5300", "Server port")
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

	lis, err := net.Listen("tcp", *serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPingPoseidonServer(s, &server{})

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
	    s.Stop()
	    log.Print("service terminated")
	    os.Exit(0)
	}()


	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
