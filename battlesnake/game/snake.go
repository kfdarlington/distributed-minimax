package game

import (
	"fmt"
	"github.com/kristian-d/distributed-minimax/engine/pb"
)

type Coordinate struct {
	X uint32 `json:"x"`
	Y uint32 `json:"y"`
}

type snakeRaw struct {
	Id     string       `json:"id"`
	Name   string       `json:"name"`
	Health uint32       `json:"health"`
	Body   []Coordinate `json:"body"`
	Shout  string       `json:"shout"`
}

type SnakeByValue map[uint32]*Snake

type Snake struct {
	Body   []Coordinate
	Health uint32
	Moved  bool
	Value  uint32
}

func (s *Snake) copy() *Snake {
	body := make([]Coordinate, len(s.Body))
	for i, location := range s.Body {
		body[i] = location
	}
	return &Snake{
		Value:  s.Value,
		Body:   body,
		Health: s.Health,
		Moved:  s.Moved,
	}
}

func (sbv SnakeByValue) copy() SnakeByValue {
	newMap := make(SnakeByValue)
	for k, v := range sbv {
		newMap[k] = v.copy()
	}
	return newMap
}

func createSnakeMappings(rawSnakes []snakeRaw, myId string) SnakeByValue {
	snakesMapping := make(SnakeByValue)
	for i, rawSnake := range rawSnakes {
		var value uint32
		if rawSnake.Id == myId {
			value = ME
		} else {
			value = uint32(i + 1) + ME // ensures that values are unique
		}
		snakesMapping[value] = &Snake{
			Health:         rawSnake.Health,
			Body:           rawSnake.Body,
			Value:          value,
			Moved:          false,
		}
	}
	return snakesMapping
}

func (sbv SnakeByValue) ToProtobuf() map[uint32]*pb.Board_Snake {
	snakes := make(map[uint32]*pb.Board_Snake)
	for k, v := range sbv {
		snakes[k] = v.ToProtobuf()
	}
	return snakes
}

func (s *Snake) ToProtobuf() *pb.Board_Snake {
	body := make([]*pb.Board_Snake_Coordinate, len(s.Body), len(s.Body))
	for i, coordinate := range s.Body {
		body[i] = &pb.Board_Snake_Coordinate{
			X: coordinate.X,
			Y: coordinate.Y,
		}
	}
	return &pb.Board_Snake{
		Body: body,
		Health: s.Health,
		Value: s.Value,
	}
}

func SnakesFromProtobuf(pb map[uint32]*pb.Board_Snake) SnakeByValue {
	snakes := make(SnakeByValue)
	for k, v := range pb {
		snakes[k] = SnakeFromProtobuf(v)
	}
	return snakes
}

func SnakeFromProtobuf(pb *pb.Board_Snake) *Snake {
	body := make([]Coordinate, len(pb.GetBody()))
	for i, coordinate := range pb.GetBody() {
		body[i] = Coordinate{
			X: coordinate.GetX(),
			Y: coordinate.GetY(),
		}
	}
	moved := false
	// we can always be set to moved as we move first and on our own
	if pb.GetValue() == ME {
		moved = true
	}
	return &Snake{
		Body: body,
		Health: pb.GetHealth(),
		Value: pb.GetValue(),
		Moved: moved,
	}
}

func (s *Snake) Print() {
	fmt.Printf("{\n\tValue: %d\n\tHealth: %d\n\tSize: %d\n\tMoved: %d\n}\n", s.Value, s.Health, len(s.Body), s.Moved)
}
