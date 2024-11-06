package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

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

func (s *server) PingPongServerStream(req *pb.PingReqest, stream grpc.ServerStreamingServer[pb.PingResponse]) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&pb.PingResponse{
			Pong: fmt.Sprintf("Hello, [%d] ServerStreaming:%s", i, req.GetPing()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

func (s *server) PingPongClientStream(stream pb.PingPong_PingPongClientStreamServer) error {
	pingList := make([]string, 0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("recived pingpong list: ", pingList)
			pong := fmt.Sprintf("Hello, ClientStreaming!, %v", pingList)
			return stream.SendAndClose(&pb.PingResponse{
				Pong: pong,
			})
		}
		if err != nil {
			return err
		}
		pingList = append(pingList, req.GetPing())
	}
}

func (s *server) PingPongBiStreams(stream pb.PingPong_PingPongBiStreamsServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return err
		}
		if err != nil {
			return err
		}
		pong := fmt.Sprintf("Hello, BiStreaming!, %v", req.GetPing())
		if err := stream.Send(&pb.PingResponse{
			Pong: pong,
		}); err != nil {
			return err
		}

	}
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
