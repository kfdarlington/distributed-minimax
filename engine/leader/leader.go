package leader

import (
	"context"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/battlesnake/game"
	"github.com/kristian-d/distributed-minimax/engine/leader/pools"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"io"
	"io/ioutil"
	"math"
	"time"
)

type Leader struct {
	logger *logger.Logger
	pools *pools.Pools
}

func alphabeta(n *expander.Node, depth int, alpha float64, beta float64, maximizingPlayer bool) float64 {
	if depth == 0 || n.Terminal {
		return evaluator.Evaluate(n.Game)
	}
	if maximizingPlayer {
		value := math.Inf(-1) // negative infinity
		for _, child := range n.Children {
			value = math.Max(value, alphabeta(child, depth-1, alpha, beta, false))
			alpha = math.Max(alpha, value)
			if beta <= alpha {
				break
			}
		}
		return value
	} else {
		value := math.Inf(1) // positive infinity
		for _, child := range n.Children {
			value = math.Min(value, alphabeta(child, depth-1, alpha, beta, true))
			beta = math.Min(beta, value)
			if beta <= alpha {
				break
			}
		}
		return value
	}
}

func (l *Leader) ComputeMove(board game.Board, deadline time.Duration) game.Move {
	follower, err := l.pools.Activate()
	if err != nil {
		l.logger.Errorf("error retrieving and activating idle follower err=%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	stream, err := (*follower.Client).GetExpansion(ctx, )
	if err != nil {
		l.logger.Errorf("error requesting expansion from follower err=%v", err)
	}
	for {
		expansion, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			l.logger.Errorf("error receiving expansion err=%v", err)
		}
		l.logger.Infof("received expansion expansion=%v", expansion)
	}
	return nil
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
