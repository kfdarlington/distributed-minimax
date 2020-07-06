package game

import "fmt"

const (
	EMPTY uint32 = iota
	FOOD
	ME
)

type Grid struct {
	height uint32
	width uint32
	values []uint32
}

func (g *Grid) GetHeight() uint32 {
	return g.height
}

func (g *Grid) GetWidth() uint32 {
	return g.width
}

func (g *Grid) GetValues() []uint32 {
	return g.values
}

func (g *Grid) GetValue(y uint32, x uint32) uint32 {
	return g.values[y*g.width + x]
}

func (g *Grid) SetValue(y uint32, x uint32, value uint32) {
	g.values[y*g.width + x] = value
}

func createGrid(state Update, snakesMap *SnakeByValue) *Grid { // currently generates a new grid every update for simplicity
	height   := state.Board.Height
	width    := state.Board.Width
	grid     := Grid{
		height: height,
		width: width,
		values: make([]uint32, height*width, height*width),
	}
	for _, snake := range *snakesMap {
		for _, coordinate := range snake.Body {
			grid.SetValue(coordinate.Y, coordinate.X, snake.Value)
		}
	}
	for _, coordinate := range state.Board.Food {
		grid.SetValue(coordinate.Y, coordinate.X, FOOD)
	}
	return &grid
}

func (g *Grid) Copy() *Grid {
	gridCopy  := Grid{
		height: g.height,
		width: g.width,
		values: make([]uint32, g.height*g.width, g.height*g.width),
	}
	copy(gridCopy.values, g.values)
	return &gridCopy
}

func (g *Grid) Print() {
	var start uint32
	var end uint32
	var i uint32 = 0
	for ; i < g.height; i++ {
		start = i*g.width
		end = start + g.width
		fmt.Println(g.values[start:end])
	}
}
