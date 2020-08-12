package leader

import (
	"context"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/battlesnake/game"
	"github.com/kristian-d/distributed-minimax/engine/leader/pools"
	"github.com/kristian-d/distributed-minimax/engine/leader/web"
	"io/ioutil"
	"math"
	"time"
)

type Leader struct {
	logger *logger.Logger
	pools *pools.Pools
}

func (l *Leader) ComputeMove(b game.Board, deadline time.Duration) game.Move {
	ctx, cancel := context.WithTimeout(context.Background(), deadline*time.Millisecond) // process the move for x ms, leaving (500 - x) ms for the network
	defer cancel()
	// absoluteDeadline := time.Now().UnixNano()/int64(time.Millisecond) + int64(deadline)
	root := b.Protobuf()
	// latestMove := game.UP // default move is some arbitrary direction for now
	depth := 2 // TODO change
	valueChan := make(chan float64)
	go l.alphabeta(ctx, root, depth, math.Inf(-1), math.Inf(1), true, valueChan)
	select {
	case value, ok := <-valueChan:
		l.logger.Infof("value=%f ok=%t", value, ok)
		// TODO something
	case <-ctx.Done():
		// TODO something
	}
	return game.UP
}

func (l *Leader) CloseConnections() {
	l.pools.DestroyConnections()
}

func CreateLeader() *Leader {
	p := pools.CreatePools()
	web.Create(p, 3001)
	return &Leader{
		logger: logger.Init("Leader", true, false, ioutil.Discard),
		pools: p,
	}
}
