package leader

import (
	"context"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/battlesnake/game"
	"github.com/kristian-d/distributed-minimax/engine/leader/pools"
	"github.com/kristian-d/distributed-minimax/engine/leader/web"
	"io/ioutil"
	"time"
)

type Leader struct {
	logger *logger.Logger
	pool *pools.Pool
}

func (l *Leader) ComputeMove(b game.Board, deadline time.Duration) string {
	ctx, cancel := context.WithTimeout(context.Background(), deadline*time.Millisecond) // process the move for x ms, leaving (500 - x) ms for the network (for battlesnake)
	defer cancel()
	// absoluteDeadline := time.Now().UnixNano()/int64(time.Millisecond) + int64(deadline)
	root := b.Protobuf()
	depth := 2 // TODO change
	move := l.startalphabeta(ctx, root, depth)
	return move
}

func (l *Leader) CloseConnections() {
	l.pool.DestroyConnections()
}

func CreateLeader() *Leader {
	lgger := logger.Init("Leader", true, false, ioutil.Discard)
	p := pools.CreatePool()
	server := web.Create(p, 3001)
	go func() {
		lgger.Fatal(server.ListenAndServe())
	}()
	lgger.Info("leader server listening on port 3001")
	return &Leader{
		logger: lgger,
		pool: p,
	}
}
