package game

import "github.com/kristian-d/distributed-minimax/engine/pb"

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
	Grid           *Grid
	Snakes         *SnakeByValue
	MoveCoordinate MoveCoordinate // the move and coordinate pair that generated this board
}

func (b *Board) Copy() Board {
	return Board{
		Grid: b.Grid.Copy(),
		Snakes: b.Snakes.copy(),
		MoveCoordinate: b.MoveCoordinate,
	}
}

func (b *Board) Protobuf() *pb.Board {
	snakes := make(map[uint32]*pb.Board_Snake)
	for k, v := range *b.Snakes {
		body := make([]*pb.Board_Snake_Coordinate, len(v.Body), len(v.Body))
		for i, coordinate := range v.Body {
			body[i] = &pb.Board_Snake_Coordinate{
				X: coordinate.X,
				Y: coordinate.Y,
			}
		}
		snakes[k] = &pb.Board_Snake{
			Body: body,
			Health: v.Health,
			Value: v.Value,
		}
	}
	return &pb.Board{
		Grid: &pb.Board_Grid{
			Height: b.Grid.GetHeight(),
			Width: b.Grid.GetWidth(),
			Values: b.Grid.GetValues(),
		},
		Snakes: snakes,
	}
}

func InitGames() {
	Games = make(map[string]*Game)
}
