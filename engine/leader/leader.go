package leader

import (
	"context"
	"github.com/google/logger"
	"github.com/kristian-d/distributed-minimax/battlesnake/game"
	"github.com/kristian-d/distributed-minimax/engine/leader/pools"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"io"
	"io/ioutil"
	"log"
	"math"
	"time"
)

type Leader struct {
	logger *logger.Logger
	pools *pools.Pools
}

func (l *Leader) Evaluate(board *pb.Board, deadline time.Duration) uint32 {
	follower, err := l.pools.Activate()
	if err != nil {
		l.logger.Errorf("error retrieving and activating idle follower err=%v", err)
		return 0
	}
	ctx, cancel := context.WithTimeout(context.Background(), deadline*time.Millisecond)
	defer cancel()
	client := *follower.GetClient()
	evaluateReply, err := client.GetEvaluation(ctx, &pb.EvaluateRequest{
		Board: board,
	})
	if err != nil {
		l.logger.Errorf("error requesting evaluation from follower client=%v err=%v", client, err)
		return 0
	}
	return evaluateReply.GetScore()
}

func (l *Leader) Expand(board *pb.Board, deadline time.Duration, boardChan chan<- *pb.Board) {
	defer close(boardChan)
	follower, err := l.pools.Activate()
	if err != nil {
		l.logger.Errorf("error retrieving and activating idle follower err=%v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), deadline*time.Millisecond)
	defer cancel()
	client := *follower.GetClient()
	stream, err := client.GetExpansion(ctx, &pb.ExpandRequest{
		Board: board,
	})
	if err != nil {
		l.logger.Errorf("error requesting expansion from follower client=%v err=%v", client, err)
	}
	for {
		expansion, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			l.logger.Errorf("error receiving expansion err=%v", err)
			continue
		}
		boardChan <- expansion.GetBoard()
		l.logger.Infof("received expansion expansion=%v", expansion)
	}
	if err := l.pools.MarkIdle(follower); err != nil {
		l.logger.Errorf("error resetting follower to idle err=%v", err)
	}
}

func (l *Leader) ComputeMove(board game.Board, deadline time.Duration) game.Move {
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
