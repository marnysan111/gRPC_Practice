package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/marnysan111/gRPC_Practice/pingpong"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultPing = "pingpong"
	defaultType = "unary"
)

var (
	addr   = flag.String("addr", "localhost:50051", "the address to connect to")
	ping   = flag.String("ping", defaultPing, "Send PingPong")
	stream = flag.String("stream", defaultType, "gRPC Type")
)

func main() {
	flag.Parse()
	// ２つ目の引数なんだ？オプションっぽいけど
	// grpc.WithTransportCredentialsはコネクションでSSL/TLSを使用しないオプション
	// 前まではgrpc.Dialっていう関数でコネクションを作っていたが、今は非推奨
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("failed to connect: ", err)
	}
	defer conn.Close()
	c := pb.NewPingPongClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	fmt.Println(*stream)

	switch *stream {
	// unary RPC
	case "unary":
		req, err := c.PingPong(ctx, &pb.PingReqest{Ping: *ping})
		if err != nil {
			log.Fatalln("could not pingpong: ", err)
		}
		log.Println("PingPong!: ", req.GetPong())

	// Server streaming RPC
	case "server":
		stream, err := c.PingPongServerStream(context.Background(), &pb.PingReqest{Ping: *ping})
		if err != nil {
			log.Fatalln("could not pingpong: ", err)
		}
		for {
			// Recv = receive?
			req, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				log.Println("all the responses have already received.")
				break
			}
			if err != nil {
				log.Fatalln("server stream error: ", err)
			}
			fmt.Println(req)
		}
	default:
		fmt.Println("No Stream Type")
	}
}
