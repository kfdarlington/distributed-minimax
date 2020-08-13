package web

import (
	"fmt"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/engine/leader/pools"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Create starts the server using mux
func Create(pools *pools.Pool, port int) *http.Server {
	lgger := logger.Init("Engine Leader Web", true, false, ioutil.Discard)
	var host string
	if os.Getenv("ENV") == "dev" {
		host = "localhost"
	} else {
		host = ""
	}
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      NewRouter(pools, lgger),
		ReadTimeout:  time.Duration(10000) * time.Millisecond, // TODO remove hardcoding
		WriteTimeout: time.Duration(10000) * time.Millisecond, // TODO remove hardcoding
	}
}
