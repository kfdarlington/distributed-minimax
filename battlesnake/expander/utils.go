package expander

import (
	"github.com/kristian-d/distributed-minimax/battlesnake/game"
)

func prependHead(body []game.Coordinate, head game.Coordinate) []game.Coordinate { // theoretically faster than using straight append(), should be tested
	body = append(body, head)
	copy(body[1:], body)
	body[0] = head
	return body
}

func shiftBody(body []game.Coordinate, head game.Coordinate) []game.Coordinate {
	for i := len(body) - 1; i >= 1; i-- {
		body[i] = body[i - 1]
	}
	body[0] = head
	return body
}

func killSnake(b game.Board, snakeValue uint32) {
	grid := b.Grid
	snake := b.Snakes[snakeValue]
	for _, bodyPart := range snake.Body {
		grid.SetValue(bodyPart.Y, bodyPart.X, game.EMPTY)
	}
	delete(b.Snakes, snakeValue)
}

func turnComplete(b game.Board) bool {
	for _, snake := range b.Snakes {
		if !snake.Moved && snake.Value != game.ME {
			return false
		}
	}
	return true
}

func prelimaryCheck(b game.Board, snakeValue uint32, coord game.Coordinate) bool {
	// NOTE: out of bounds must be checked before calling this function (due to coords using uint32, which cannot be signed)
	value := b.Grid.GetValue(coord.Y, coord.X)
	if value == game.FOOD {
		return true // moving into a location with food, therefore death is not guaranteed
	}
	if b.Snakes[snakeValue].Health == 1 {
		return false // starvation is next turn and location does not contain food, therefore guaranteed death
	}
	if value == game.EMPTY {
		return true // moving into empty location, therefore death is not guaranteed
	}
	otherSnake := b.Snakes[value]
	if otherSnake.Moved {
		if coord.X == otherSnake.Body[0].X && coord.Y == otherSnake.Body[0].Y && len(b.Snakes[snakeValue].Body) > len(otherSnake.Body) {
			return true // moving onto a head value and has size advantage, therefore death is not guaranteed
		} else {
			return false // moving onto a body value, tail value, or head value without size advantage, therefore guaranteed death
		}
	} else {
		if coord.X == otherSnake.Body[len(otherSnake.Body)-1].X && coord.Y == otherSnake.Body[len(otherSnake.Body)-1].Y {
			return true // moving onto a tail value, therefore death is not guaranteed
		} else {
			return false // moving onto a body value, therefore guaranteed death
		}
	}
}
