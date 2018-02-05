package main

import (
	"testing"

	pb "github.com/samsung-cnct/sample-grpc-apiserver/api"
	"golang.org/x/net/context"
)

var (
	helloName       = "Marcus"
	helloReply      = "Hello Marcus"
	helloReplyAgain = "Hello again Marcus"
)

func TestHelloPoseidon(t *testing.T) {
	s := server{}

	req := pb.HelloPoseidonMsg{Name: helloName}

	resp, err := s.HelloPoseidon(context.Background(), &req)
	if err != nil {
		t.Errorf("got an unexpected error: %s", err)
	}

	if resp.Message != helloReply {
		t.Errorf("got an unexpected reply: %s %s %s", resp.Message, "intead of ", helloReply)
	}

}

func TestHelloPoseidonAgain(t *testing.T) {
	s := server{}

	req := pb.HelloPoseidonMsg{Name: helloName}

	resp, err := s.HelloPoseidonAgain(context.Background(), &req)
	if err != nil {
		t.Errorf("got an unexpected error: %s", err)
	}

	if resp.Message != helloReplyAgain {
		t.Errorf("got an unexpected reply: %s %s %s", resp.Message, "intead of ", helloReplyAgain)
	}
}
