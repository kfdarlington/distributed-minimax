package game

import "fmt"

type GridValue uint8
const (
	EMPTY GridValue = iota
	FOOD
	ME
)

type Grid [][]GridValue

func createGrid(state Update, snakesMap SnakeByValue) Grid { // currently generates a new grid every update for simplicity
	height   := state.Board.Height
	width    := state.Board.Width
	grid     := make(Grid, height)
	contents := make([]GridValue, height*width)
	for i := range grid {
		start := i*width
		end   := start+width
		grid[i] = contents[start:end:end]
	}
	for _, snake := range snakesMap {
		for _, coordinate := range snake.Body {
			grid[coordinate.Y][coordinate.X] = snake.Value
		}
	}
	for _, coordinate := range state.Board.Food {
		grid[coordinate.Y][coordinate.X] = FOOD
	}
	return grid
}

func copyGrid(grid Grid) Grid {
	height    := len(grid)
	width     := len(grid[0])
	gridCopy  := make(Grid, height)
	contents  := make([]GridValue, height*width)
	for i := range grid {
		start := i*width
		end   := start+width
		gridCopy[i] = contents[start:end:end]
		copy(gridCopy[i], grid[i])
	}
	return gridCopy
}

func PrintGrid(grid Grid) {
	for _, row := range grid {
		fmt.Println(row)
	}
}
