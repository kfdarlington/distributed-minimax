package game

import "fmt"

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type snakeRaw struct {
	Id     string       `json:"id"`
	Name   string       `json:"name"`
	Health int          `json:"health"`
	Body   []Coordinate `json:"body"`
	Shout  string       `json:"shout"`
}

type SnakeByValue map[GridValue]Snake

type Snake struct {
	Body   []Coordinate
	Health int
	Moved  bool
	Value  GridValue
}

func copySnake(snake Snake) Snake {
	body := make([]Coordinate, len(snake.Body))
	for i, location := range snake.Body {
		body[i] = location
	}
	return Snake{
		Value:  snake.Value,
		Body:   body,
		Health: snake.Health,
		Moved:  snake.Moved,
	}
}

func copySnakeByValues(snakeByValues SnakeByValue) SnakeByValue {
	newMap := make(SnakeByValue)
	for k, v := range snakeByValues {
		newMap[k] = copySnake(v)
	}
	return newMap
}

func createSnakeMappings(rawSnakes []snakeRaw, myId string) map[GridValue]Snake {
	snakesMapping := make(map[GridValue]Snake)
	for i, rawSnake := range rawSnakes {
		var value GridValue
		if rawSnake.Id == myId {
			value = ME
		} else {
			value = GridValue(i + 1) + ME // ensures that values are unique
		}
		snakesMapping[value] = Snake{
			Health:         rawSnake.Health,
			Body:           rawSnake.Body,
			Value:          value,
			Moved:          false,
		}
	}
	return snakesMapping
}

func PrintSnake(snake Snake) {
	fmt.Printf("Value: %d\nHealth: %d\n Size: %d\n Moved: %d\n\n", snake.Value, snake.Health, len(snake.Body), snake.Moved)
}
