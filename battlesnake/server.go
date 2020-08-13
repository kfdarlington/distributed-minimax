package battlesnake

import (
	"fmt"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/battlesnake/web"
	"github.com/kristian-d/distributed-minimax/engine/leader"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/kristian-d/distributed-minimax/battlesnake/game"
)

// Create starts the snake server using mux
func Create(engine *leader.Leader, port int) *http.Server {
	game.InitGames()
	lgger := logger.Init("BattleSnake Web", true, false, ioutil.Discard)
	var host string
	if os.Getenv("ENV") == "dev" {
		host = "localhost"
	} else {
		host = ""
	}
	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      web.NewRouter(engine, lgger),
		ReadTimeout:  time.Duration(1000) * time.Millisecond, // TODO remove hardcoding
		WriteTimeout: time.Duration(1000) * time.Millisecond, // TODO remove hardcoding
	}
}
