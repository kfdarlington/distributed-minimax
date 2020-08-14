package expander

import (
	"context"
	"github.com/kristian-d/distributed-minimax/battlesnake/game"
	"github.com/kristian-d/distributed-minimax/engine/pb"
	"math"
	"sync"
)

func boardBranchesBySnakeMove(b game.Board, snakeValue uint32, outChan chan<- game.Board) {
	defer close(outChan)

	// determine which coordinates are possible to move to given the board dimensions and current head coordinate
	head := b.Snakes[snakeValue].Body[0]
	coords := make([]game.Coordinate, 0, 4)
	if head.X > 0 {
		coords = coords[:len(coords) + 1]
		coords[len(coords) - 1] = game.Coordinate{X: head.X-1, Y: head.Y}
	}
	if head.Y > 0 {
		coords = coords[:len(coords) + 1]
		coords[len(coords) - 1] = game.Coordinate{X: head.X, Y: head.Y-1}
	}
	if head.X < b.Grid.GetWidth() - 1 {
		coords = coords[:len(coords) + 1]
		coords[len(coords) - 1] = game.Coordinate{X: head.X+1, Y: head.Y}
	}
	if head.Y < b.Grid.GetHeight() - 1 {
		coords = coords[:len(coords) + 1]
		coords[len(coords) - 1] = game.Coordinate{X: head.X, Y: head.Y+1}
	}

	// of the possible move options, do some more checking and, if move still possible, move snake to new coordinate
	moveOptions := 0
	for _, coord := range coords {
		if prelimaryCheck(b, snakeValue, coord) {
			moveOptions++
			newBoard := b.Copy()
			moveSnake(newBoard, snakeValue, coord)
			outChan <- newBoard
		}
	}

	// if the snake had no move options, then it is forced to die (otherwise we can assume the snake will not kill itself)
	if moveOptions == 0 {
		newBoard := b.Copy()
		killSnake(newBoard, snakeValue)
		outChan <- newBoard
	}
}

func boardBranches(b game.Board, outChan chan<- game.Board) {
	valueSnakeMap := b.Snakes
	maxSize := 0
	var largestSnakeValue uint32
	for value, snake := range valueSnakeMap {
		if !snake.Moved && len(snake.Body) > maxSize && snake.Value != game.ME {
			maxSize = len(snake.Body)
			largestSnakeValue = value
		}
	}
	boardBranchesBySnakeMove(b, largestSnakeValue, outChan)
}

func Expand(ctx context.Context, pb *pb.Board, maximizingPlayer bool, outChan chan<- *pb.Board) {
	defer close(outChan)
	board := game.BoardFromProtobuf(pb)
	if maximizingPlayer {
		branchChan := make(chan game.Board, 4)
		go boardBranchesBySnakeMove(board, game.ME, branchChan)
		for {
			select {
			case branch, ok := <-branchChan:
				if !ok {
					return
				} else {
					var terminalState bool
					if _, ok := branch.Snakes[game.ME]; !ok || len(branch.Snakes) == 1 { // board state is terminal if I am dead or the only snake left
						terminalState = true
					} else {
						terminalState = false
					}
					outChan <- branch.ToProtobuf(terminalState)
				}
			case <-ctx.Done():
				return
			}
		}
	} else {
		// buffer channels to the maximum possible number of outputs so that there are no blocks
		maxOutputs := int64(math.Pow(3, float64(len(board.Snakes)-1)))
		branchChan := make(chan game.Board, maxOutputs)
		var wg sync.WaitGroup
		wg.Add(1)
		branchChan <- board
		go func() {
			wg.Wait()
			close(branchChan)
		}()
		for {
			select {
			case branch, ok := <-branchChan:
				if !ok {
					return
				} else if turnComplete(branch) {
					var terminalState bool
					if _, ok := branch.Snakes[game.ME]; !ok || len(branch.Snakes) == 1 { // board state is terminal if I am dead or the only snake left
						terminalState = true
					} else {
						terminalState = false
					}
					outChan <- branch.ToProtobuf(terminalState)
					wg.Done()
				} else { // expand further (move another snake that has yet to be moved)
					newBranchChan := make(chan game.Board, 4)
					go func(c <-chan game.Board) {
						defer wg.Done()
						for branch := range c {
							wg.Add(1)
							branchChan <- branch
						}
					}(newBranchChan)
					go boardBranches(branch, newBranchChan)
				}
			case <-ctx.Done():
				return
			}
		}
	}
}
