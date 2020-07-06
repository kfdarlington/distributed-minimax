package leader

import (
	"context"
	"github.com/kristian-d/distributed-minimax/engine/pb"
)

func (l *Leader) evaluate(ctx context.Context, board *pb.Board) float32 {
	// find an available follower
	follower, err := l.pools.Activate(ctx)
	if err != nil {
		l.logger.Errorf("error retrieving and activating idle follower err=%v", err)
		return 0
	}
	if follower == nil {
		l.logger.Info("follower was not made available before context expired")
		return 0
	}
	defer func() {
		if err := l.pools.MarkIdle(follower); err != nil {
			l.logger.Errorf("error resetting follower to idle follower=%v err=%v", *follower, err)
		}
	}()

	// begin request and wait for response
	client := *follower.GetClient()
	evaluateReply, err := client.GetEvaluation(ctx, &pb.EvaluateRequest{
		Board: board,
	})
	if err != nil {
		l.logger.Errorf("error requesting evaluation from follower client=%v err=%v", client, err)
		return 0
	}
	score := evaluateReply.GetScore()
	l.logger.Infof("calculated score of %f", score)
	return score
}
