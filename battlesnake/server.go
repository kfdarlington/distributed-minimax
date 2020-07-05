package battlesnake

import (
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/battlesnake/web"
	"github.com/kristian-d/distributed-minimax/engine/leader"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kristian-d/distributed-minimax/battlesnake/game"
)

// Start starts the snake server using mux
func Create(engine *leader.Leader, port int) *http.Server {
	game.InitGames()
	lgger := logger.Init("Web", true, false, ioutil.Discard)
	return &http.Server{
		Addr:         ":" + string(port),
		Handler:      web.NewRouter(engine, lgger),
		ReadTimeout:  time.Duration(1000) * time.Millisecond, // TODO remove hardcoding
		WriteTimeout: time.Duration(1000) * time.Millisecond, // TODO remove hardcoding
	}
}
