package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"sync"

	"github.com/miles990/grpcpool/grpcpool"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func CallSayHello(wg *sync.WaitGroup) {
	defer wg.Done()

	conn := grpcpool.GetManager().DefaultConnPool().Get()   // get connection from pool
	defer grpcpool.GetManager().DefaultConnPool().Put(conn) // put connection back to pool

	// call grpc service
	result, err := pb.NewGreeterClient(conn).SayHello(context.Background(), &pb.HelloRequest{Name: "world"})
	if err != nil {
		slog.Error("example", "SayHello", err.Error())
		return
	}
	slog.Info("example", "SayHello", fmt.Sprintf("result:%v", result))
}

func main() {
	// new connection pool
	grpcpool.GetManager().NewConnPool(10, *addr)

	var wg sync.WaitGroup
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go CallSayHello(&wg)
	}

	wg.Wait()

	slog.Info("example", "SayHello", "all done")
	// release connection pool
	grpcpool.GetManager().DefaultConnPool().Close()
}
