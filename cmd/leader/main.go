package main

import (
	"flag"
	"github.com/google/logger"
	minimax "github.com/kristian-d/distributed-battlesnake/minimax/leader"
	"github.com/kristian-d/distributed-battlesnake/minimax/pb"
	"google.golang.org/grpc"
	"io/ioutil"
)

type addrFlag []string

func (addr *addrFlag) String() string {
	return "localhost:8980"
}

func (addr *addrFlag) Set(value string) error {
	*addr = append(*addr, value)
	return nil
}

func main() {
	lgger := logger.Init("MainLogger", true, false, ioutil.Discard)
	/*var cfg *config.Config
	env := os.Getenv("ENV")
	switch env {
	case "prod":
		var err error
		cfg, err = config.LoadEnv(); if err != nil {
			lgger.Fatalf("configuration file could not be loaded err=", err)
		}
	case "local":
		configPath := path.Join(projectpath.Root, "config/local.yml")
		var err error
		cfg, err = config.LoadFile(configPath); if err != nil {
			lgger.Fatalf("configuration file could not be loaded err=", err)
		}
	default:
		lgger.Infof("unknown environment or environment missing, defaulting to local env=%s", env)
		configPath := path.Join(projectpath.Root, "config/local.yml")
		var err error
		cfg, err = config.LoadFile(configPath); if err != nil {
			lgger.Fatalf("configuration file could not be loaded err=", err)
		}
	}*/
	var addresses addrFlag
	flag.Var(&addresses, "addr", "addresses of follower servers to connect to")
	flag.Parse()

	pool := make([]*pb.MinimaxClient, 0)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	for _, addr := range addresses {
		conn, err := grpc.Dial(addr, opts...)
		if err != nil {
			lgger.Errorf("failed to connect addr=%s err=%v", addr, err)
		} else {
			lgger.Infof("connected to follower addr=%s", addr)
			defer conn.Close()
			client := pb.NewMinimaxClient(conn)
			pool = append(pool, &client)
		}
	}

	leader, err := minimax.NewLeader(pool); if err != nil {
		lgger.Fatal(err)
	}
	lgger.Fatal(leader.Start())
}


