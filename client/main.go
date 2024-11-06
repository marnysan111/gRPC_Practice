package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/marnysan111/gRPC_Practice/pingpong"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultPing = "pingpong"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	ping = flag.String("ping", defaultPing, "Send PingPong")
)

func main() {
	flag.Parse()
	// ２つ目の引数なんだ？オプションっぽいけど
	// grpc.WithTransportCredentialsはコネクションでSSL/TLSを使用しないオプション
	// 前まではgrpc.Dialっていう関数でコネクションを作っていたが、今は非推奨
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalln("failed to connect: ", err)
	}
	defer conn.Close()
	c := pb.NewPingPongClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.PingPong(ctx, &pb.PingReqest{Ping: *ping})
	if err != nil {
		log.Fatalln("could not pingpong: ", err)
	}
	log.Println("PingPong!: ", r.GetPong())
}
