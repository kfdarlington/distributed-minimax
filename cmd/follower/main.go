package main

import (
	"flag"
	"fmt"
	"github.com/google/logger"
	minimax "github.com/kristian-d/distributed-minimax/engine/follower"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
)

func main() {
	lgger := logger.Init("MainFollowerLogger", true, false, ioutil.Discard)

	var host = flag.String("host", "", "host address of the follower")
	var port = flag.Int("port", 3000, "port of the follower")
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		lgger.Fatalf("failed to listen port=%d err=%v", *port, err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMinimaxServer(grpcServer, minimax.NewFollower())
	lgger.Infof("follower listening port=%d", *port)
	lgger.Fatal(grpcServer.Serve(listener))
}
