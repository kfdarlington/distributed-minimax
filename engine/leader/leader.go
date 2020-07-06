package leader

import (
	"context"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/battlesnake/game"
	"github.com/kristian-d/distributed-minimax/engine/leader/pools"
	"github.com/kristian-d/distributed-minimax/engine/pb"
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
	depth := 2
	l.alphabeta(ctx, root, depth, math.Inf(-1), math.Inf(1), true)
	return game.UP
}

func NewLeader(clients []*pb.MinimaxClient) (*Leader, error) {
	p, err := pools.CreateFollowerPools(clients)
	if err != nil {
		return nil, err
	}
	return &Leader{
		logger: logger.Init("Leader", true, false, ioutil.Discard),
		pools: p,
	}, nil
}
