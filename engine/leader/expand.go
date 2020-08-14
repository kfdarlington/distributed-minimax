package leader

import (
	"context"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"io"
)

func (l *Leader) expand(ctx context.Context, board *pb.Board, maximizingPlayer bool, boardChan chan<- *pb.Board) {
	defer close(boardChan)

	// find an available follower
	follower := l.pool.Activate(ctx)
	if follower == nil {
		l.logger.Info("follower was not made available before context expired")
		return
	}
	defer func() {
		if err := l.pool.MarkAsIdle(follower); err != nil {
			l.logger.Errorf("error resetting follower to idle follower=%v err=%v", *follower, err)
		}
	}()

	// begin request and response streaming
	client := *follower.GetClient()
	stream, err := client.GetExpansion(ctx, &pb.ExpandRequest{
		Board: board,
		IsMaximizerTurn: maximizingPlayer,
	})
	if err != nil {
		l.logger.Errorf("error requesting expansion from follower client=%v err=%v", client, err)
		return
	}

	// funnel streamed responses into the board channel
	for {
		expansion, err := stream.Recv()
		if err == io.EOF {
			l.logger.Info("stream reached EOF, exiting")
			return
		} else if err != nil {
			l.logger.Errorf("unknown error receiving expansion, exiting err=%v", err)
			return
		}
		l.logger.Info("received expansion")
		boardChan <- expansion.GetBoard()
	}
}
