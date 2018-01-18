package main

import (
	"flag"
	"fmt"
	"log"
	"net"

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
	fmt.Println("Server starting")
	lis, err := net.Listen("tcp", *serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPingPoseidonServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
