package game

import (
	"errors"
	"fmt"
	"github.com/kristian-d/distributed-minimax/engine/pb"
)

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
	DEFAULT_MOVE Move = "up"
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

func GetOriginatingMove(fromBoard *pb.Board, toBoard *pb.Board) (Move, error) {
	var fromHeadCoord *pb.Board_Snake_Coordinate
	if fromSnakes, ok := fromBoard.GetSnakes()[ME]; !ok {
		return NONE, errors.New("we do not exist in snake mapping from originating board")
	} else if len(fromSnakes.GetBody()) == 0 {
		return NONE, errors.New("we do not have a head (body length is 0) in snake body from originating board")
	} else {
		fromHeadCoord = fromSnakes.GetBody()[0]
	}

	var toHeadCoord *pb.Board_Snake_Coordinate
	if toSnakes, ok := toBoard.GetSnakes()[ME]; !ok {
		return NONE, errors.New("we do not exist in snake mapping from destination board")
	} else if len(toSnakes.GetBody()) == 0 {
		return NONE, errors.New("we do not have a head (body length is 0) in snake body from destination board")
	} else {
		toHeadCoord = toSnakes.GetBody()[0]
	}

	xDiff := int64(toHeadCoord.X) - int64(fromHeadCoord.X)
	yDiff := int64(toHeadCoord.Y) - int64(fromHeadCoord.Y)
	switch xDiff {
	case 0:
		switch yDiff {
		case 1:
			return DOWN, nil
		case -1:
			return UP, nil
		default:
			return NONE, errors.New(fmt.Sprintf("toBoard and fromBoard are not one move apart xDifference=%d yDifference=%d", xDiff, yDiff))
		}
	case 1:
		if yDiff == 0 {
			return RIGHT, nil
		} else {
			return NONE, errors.New(fmt.Sprintf("toBoard and fromBoard are not one move apart xDifference=%d yDifference=%d", xDiff, yDiff))
		}
	case -1:
		if yDiff == 0 {
			return LEFT, nil
		} else {
			return NONE, errors.New(fmt.Sprintf("toBoard and fromBoard are not one move apart xDifference=%d yDifference=%d", xDiff, yDiff))
		}
	default:
		return NONE, errors.New(fmt.Sprintf("toBoard and fromBoard are not one move apart xDifference=%d yDifference=%d", xDiff, yDiff))
	}
}

func InitGames() {
	Games = make(map[string]*Game)
}
