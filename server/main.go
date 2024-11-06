package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/marnysan111/gRPC_Practice/pingpong"
	pb "github.com/marnysan111/gRPC_Practice/pingpong"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "Server Port")
)

type server struct {
	pb.UnimplementedPingPongServer
}

// contextってあるけど、いつ使うんだろ。ってかなんで引数として存在してるんだ？使ってないのに
func (s *server) PingPong(_ context.Context, in *pb.PingReqest) (*pb.PingResponse, error) {
	log.Println("Received: ", in.GetPing())
	return &pb.PingResponse{Pong: "Hello " + in.GetPing()}, nil
}

func (s *server) PingPongServerStream(req *pingpong.PingReqest, stream grpc.ServerStreamingServer[pingpong.PingResponse]) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&pingpong.PingResponse{
			Pong: fmt.Sprintf("Hello, [%d] ServerStreaming:%s", i, req.GetPing()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
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
