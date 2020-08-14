package game

import (
	"errors"
)

type Update struct {
	Game struct {
		Id string `json:"id"`
	} `json:"game"`
	Turn  int `json:"turn"`
	Board struct {
		Height uint32        `json:"height"`
		Width  uint32        `json:"width"`
		Food   []Coordinate  `json:"food"`
		RawSnakes []snakeRaw `json:"snakes"`
	} `json:"board"`
	You snakeRaw `json:"you"`
}

func CreateGame(state Update) {
	snakesMap := createSnakeMappings(state.Board.RawSnakes, state.You.Id)
	grid      := createGrid(state, snakesMap)

	Games[state.Game.Id] = &Game{
		Id:    state.Game.Id,
		Board: Board{
			Grid:   grid,
			Snakes: snakesMap,
		},
		PreviousMaxDepth: 0,
	}
}

func UpdateGame(state Update) error {
	if game, ok := Games[state.Game.Id]; ok {
		game.Board.Snakes = createSnakeMappings(state.Board.RawSnakes, state.You.Id)
		game.Board.Grid   = createGrid(state, game.Board.Snakes)
		return nil
	} else {
		return errors.New("no game with given id for update")
	}
}

func DeleteGame(state Update) error {
	if _, ok := Games[state.Game.Id]; !ok {
		return errors.New("no game with given id for delete")
	}
	delete(Games, state.Game.Id) // garbage collector will do the rest
	return nil
}
