package leader

import (
	"context"
	"github.com/kristian-d/distributed-minimax/battlesnake/game"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"math"
)

// this function begins the alpha-beta pruning algorithm at the root of the minimax tree; it returns a move instead of
// a value based on an evaluation (the move from the root of the tree is all that we care about)
func (l *Leader) startalphabeta(ctx context.Context, b *pb.Board, depth int) game.Move {
	if depth == 0 {
		return game.DEFAULT_MOVE
	}

	// instantiate all that is needed for communication between go functions
	boardChan := make(chan *pb.Board)
	expectedValueCount := 0
	moveChan := make(chan struct{
		move game.Move
		evaluation float64
	})

	// ensures that all child processes of this alpha-beta call will exit if their branch is pruned or the function is otherwise done
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	evaluation := math.Inf(-1) // negative infinity
	alpha := math.Inf(-1) // negative infinity
	beta := math.Inf(1) // positive infinity
	bestMove := game.DEFAULT_MOVE // arbitrary move

	go l.expand(ctx2, b, boardChan)
	for {
		select {
		case newBoard, ok := <-boardChan:
			if !ok {
				boardChan = nil // disables this case in select
			} else {
				expectedValueCount++
				go func(board *pb.Board) {
					valueChan := make(chan float64)
					go l.alphabeta(ctx2, board, depth-1, alpha, beta, false, valueChan)
					select {
					case evaluation := <-valueChan:
						if move, err := game.GetOriginatingMove(b, board); err != nil {
							l.logger.Errorf("error getting originating move err=%v", err)
						} else {
							moveChan <- struct{
								move game.Move
								evaluation float64
							}{
								move,
								evaluation,
							}
						}
					case <-ctx2.Done():
					}
				}(newBoard)
			}
		case newMove := <-moveChan:
			expectedValueCount--
			evaluation = math.Max(evaluation, newMove.evaluation)
			l.logger.Infof("value updated value=%f", evaluation)
			if evaluation == newMove.evaluation { // if value was updated, updated the move too
				bestMove = newMove.move
				l.logger.Infof("move updated move=%s", bestMove)
			}
			alpha = math.Max(alpha, evaluation)
			l.logger.Infof("alpha updated alpha=%f", alpha)
			if beta <= alpha { // prune any sibling branches that have not run or are currently running -- "defer cancel()" ensures they will finish due to their context
				l.logger.Infof("pruning value=%s depth=%d alpha=%f beta=%f maximizingPlayer=true", bestMove, depth, alpha, beta)
				return bestMove
			} else if expectedValueCount == 0 && boardChan == nil { // we are not expecting and will never expect more values
				l.logger.Infof("exhausted branches, returning move=%s depth=%d alpha=%f beta=%f maximizingPlayer=true", bestMove, depth, alpha, beta)
				return bestMove
			}
		case <-ctx2.Done():
			return bestMove
		}
	}
}

func (l *Leader) alphabeta(ctx context.Context, b *pb.Board, depth int, alpha float64, beta float64, maximizingPlayer bool, resultChan chan float64) {
	if depth == 0 {
		resultChan <- float64(l.evaluate(ctx, b))
		return
	}

	// instantiate all that is needed for communication between go functions
	boardChan := make(chan *pb.Board)
	valueChan := make(chan float64)
	expectedValueCount := 0

	// ensures that all child processes of this alpha-beta call will exit if their branch is pruned or the function is otherwise done
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	// to reduce code redundancy, instantiate the values that change based on maximizing or minimizing player
	var value float64
	var compareFn func(x float64, y float64) float64
	var alphaOrBetaPtr *float64
	var alphaOrBetaDescriptor string
	if maximizingPlayer {
		value = math.Inf(-1) // negative infinity
		compareFn = math.Max
		alphaOrBetaPtr = &alpha
		alphaOrBetaDescriptor = "alpha"
	} else {
		value = math.Inf(1) // positive infinity
		compareFn = math.Min
		alphaOrBetaPtr = &beta
		alphaOrBetaDescriptor = "beta"
	}

	go l.expand(ctx2, b, boardChan)
	for {
		select {
		case newBoard, ok := <-boardChan:
			if !ok {
				boardChan = nil // disables this case in select
			} else {
				expectedValueCount++
				go l.alphabeta(ctx2, newBoard, depth-1, alpha, beta, !maximizingPlayer, valueChan)
			}
		case newValue := <-valueChan:
			expectedValueCount--
			value = compareFn(value, newValue)
			l.logger.Infof("value updated value=%f", value)
			*alphaOrBetaPtr = compareFn(*alphaOrBetaPtr, value)
			l.logger.Infof("%s updated %s=%f", alphaOrBetaDescriptor, alphaOrBetaDescriptor, *alphaOrBetaPtr)
			if beta <= alpha { // prune any sibling branches that have not run or are currently running -- "defer cancel()" ensures they will finish due to their context
				l.logger.Infof("pruning value=%f depth=%d alpha=%f beta=%f maximizingPlayer=%t", value, depth, alpha, beta, maximizingPlayer)
				resultChan <- value
				return
			} else if expectedValueCount == 0 && boardChan == nil { // we are not expecting and will never expect more values
				l.logger.Infof("exhausted branches, returning value=%f depth=%d alpha=%f beta=%f maximizingPlayer=%t", value, depth, alpha, beta, maximizingPlayer)
				resultChan <- value
				return
			}
		case <-ctx2.Done():
			return
		}
	}
}
