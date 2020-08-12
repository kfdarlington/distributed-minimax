package web

import (
	"fmt"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/engine/leader/pools"
	"io/ioutil"
	"net/http"
	"time"
)

// Create starts the server using mux
func Create(pools *pools.Pools, port int) *http.Server {
	lgger := logger.Init("Engine Leader Web", true, false, ioutil.Discard)
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      NewRouter(pools, lgger),
		ReadTimeout:  time.Duration(1000) * time.Millisecond, // TODO remove hardcoding
		WriteTimeout: time.Duration(1000) * time.Millisecond, // TODO remove hardcoding
	}
}
