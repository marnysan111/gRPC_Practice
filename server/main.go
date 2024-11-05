package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/marnysan111/gRPC_Practice/pingpong"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "Server Port")
)

type server struct {
	pb.UnimplementedPingPongServer
}

func (s *server) PingPong(_ context.Context, in *pb.PingReqest) (*pb.PingReply, error) {
	log.Println("Received: ", in.GetPing())
	return &pb.PingReply{Pong: "Hello " + in.GetPing()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln("failed to listes: ", err)
	}
	s := grpc.NewServer()
	pb.RegisterPingPongServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalln("failed to serve: ", err)
	}
}
