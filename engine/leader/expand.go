package leader

import (
	"context"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"io"
)

func (l *Leader) expand(ctx context.Context, board *pb.Board, boardChan chan<- *pb.Board) {
	defer close(boardChan)

	// find an available follower
	follower, err := l.pools.Activate(ctx)
	if err != nil {
		l.logger.Errorf("error retrieving and activating idle follower err=%v", err)
		return
	}
	if follower == nil {
		l.logger.Info("follower was not made available before context expired")
		return
	}
	defer func() {
		if err := l.pools.MarkIdle(follower); err != nil {
			l.logger.Errorf("error resetting follower to idle follower=%v err=%v", *follower, err)
		}
	}()

	// begin request and response streaming
	client := *follower.GetClient()
	stream, err := client.GetExpansion(ctx, &pb.ExpandRequest{
		Board: board,
	})
	if err != nil {
		l.logger.Errorf("error requesting expansion from follower client=%v err=%v", client, err)
	}

	// funnel streamed responses into the board channel
	for {
		expansion, err := stream.Recv()
		if err == io.EOF {
			l.logger.Info("stream reached EOF")
			break
		}
		if err != nil {
			l.logger.Errorf("error receiving expansion err=%v", err)
			break
		}
		boardChan <- expansion.GetBoard()
		l.logger.Infof("received expansion expansion=%v", expansion)
	}
}
