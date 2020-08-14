package leader

import (
	"context"
	"github.com/kristian-d/distributed-minimax/engine/pb"
)

func (l *Leader) evaluate(ctx context.Context, board *pb.Board) float32 {
	// find an available follower
	follower := l.pool.Activate(ctx)
	if follower == nil {
		l.logger.Info("follower was not made available before context expired")
		return 0
	}
	defer func() {
		if err := l.pool.MarkAsIdle(follower); err != nil {
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
	l.logger.Infof("received evaluation of %f", score)
	return score
}
