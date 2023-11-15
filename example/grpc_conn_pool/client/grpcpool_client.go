package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/miles990/grpcpool/grpcpool"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	grpcpool.GetManager().NewConnPool(10, *addr)

	conn := grpcpool.GetManager().DefaultConnPool().Get()
	defer grpcpool.GetManager().DefaultConnPool().Put(conn)

	result, err := pb.NewGreeterClient(conn).SayHello(context.Background(), &pb.HelloRequest{Name: "world"})
	if err != nil {
		slog.Error("example", "SayHello", err.Error())
		return
	}
	slog.Info("example", "SayHello", fmt.Sprintf("result:%v", result))
}
