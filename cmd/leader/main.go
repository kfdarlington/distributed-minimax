package main

import (
	"flag"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/battlesnake"
	minimax "github.com/kristian-d/distributed-minimax/engine/leader"
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
	lgger := logger.Init("Main", true, false, ioutil.Discard)
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
	flag.Var(&addresses, "addr", "addresses of worker servers to connect to")
	var port = flag.Int("port", 3000, "port of the server")
	flag.Parse()

	engine := minimax.CreateLeader()
	defer engine.CloseConnections()

	srv := battlesnake.Create(engine, *port)
	lgger.Infof("server listening on port %d\n", *port)
	lgger.Fatal(srv.ListenAndServe())
}


