package leader

import (
	"context"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"math"
)

func (l *Leader) alphabeta(ctx context.Context, b *pb.Board, depth int, alpha float64, beta float64, maximizingPlayer bool) float32 {
	if depth == 0 {
		return l.evaluate(ctx, b)
	}
	boardChan := make(chan *pb.Board)
	// ensures that all child processes of this alpha-beta call will exit if their
	// branch is pruned or the function is otherwise done
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	if maximizingPlayer {
		value := math.Inf(-1) // negative infinity
		go l.expand(ctx2, b, boardChan)
		for board := range boardChan {
			value = math.Max(value, float64(l.alphabeta(ctx2, board, depth - 1, alpha, beta, false)))
			alpha = math.Max(alpha, value)
			if beta <= alpha {
				break
			}
		}
		return float32(value)
	} else {
		value := math.Inf(1) // positive infinity
		go l.expand(ctx2, b, boardChan)
		for board := range boardChan {
			value = math.Min(value, float64(l.alphabeta(ctx2, board, depth - 1, alpha, beta, true)))
			beta = math.Min(beta, value)
			if beta <= alpha {
				break
			}
		}
		return float32(value)
	}
}
