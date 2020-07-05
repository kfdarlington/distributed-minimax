package game

var Games map[string]*Game

type Game struct {
	Id               string
	Board            Board
	PreviousMaxDepth int
}

type Move string
const (
	UP    Move = "up"
	DOWN  Move = "down"
	LEFT  Move = "left"
	RIGHT Move = "right"
	NONE  Move = ""
)

type MoveCoordinate struct {
	Move Move
	Coordinate Coordinate
}

type Board struct {
	Grid           Grid
	Snakes         SnakeByValue
	MoveCoordinate MoveCoordinate // the move and coordinate pair that generated this board
}

func CopyBoard(board Board) Board {
	return Board{
		Grid: copyGrid(board.Grid),
		Snakes: copySnakeByValues(board.Snakes),
		MoveCoordinate: board.MoveCoordinate,
	}
}

func InitGames() {
	Games = make(map[string]*Game)
}
