package leader

import (
	"context"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"math"
	"sync"
)

func (l *Leader) alphabeta(ctx context.Context, b *pb.Board, depth int, alpha float64, beta float64, maximizingPlayer bool, resultChan chan float64) {
	if depth == 0 {
		resultChan <- float64(l.evaluate(ctx, b))
	}

	// instantiate all that is needed for communication between go functions
	boardChan := make(chan *pb.Board)
	var mu sync.Mutex
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
				mu.Lock()
				expectedValueCount++
				go l.alphabeta(ctx2, newBoard, depth-1, alpha, beta, !maximizingPlayer, valueChan)
				mu.Unlock()
			}
		case newValue := <-valueChan:
			mu.Lock()
			expectedValueCount--
			value = compareFn(value, newValue)
			l.logger.Infof("value updated value=%f", value)
			*alphaOrBetaPtr = compareFn(*alphaOrBetaPtr, value)
			l.logger.Infof("%s updated %s=%f", alphaOrBetaDescriptor, alphaOrBetaDescriptor, *alphaOrBetaPtr)
			if beta <= alpha { // prune any sibling branches that have not run or are currently running -- "defer cancel()" ensures they will finish due to their context
				l.logger.Infof("pruning value=%f depth=%d alpha=%f beta=%f maximizingPlayer=%t", value, depth, alpha, beta, maximizingPlayer)
				resultChan <- value
				mu.Unlock()
				return
			} else if expectedValueCount == 0 && boardChan == nil { // we are not expecting and will never expect more values
				l.logger.Infof("exhausted branches, returning value=%f depth=%d alpha=%f beta=%f maximizingPlayer=%t", value, depth, alpha, beta, maximizingPlayer)
				resultChan <- value
				mu.Unlock()
				return
			} else { // continue handling values
				mu.Unlock()
			}
		case <-ctx2.Done():
			return
		}
	}
}
