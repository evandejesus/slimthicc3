package main

import (
	"math/rand"

	"github.com/notnil/chess"
)

func randomMove(game *chess.Game) *chess.Move {
	valid := game.ValidMoves()
	return valid[rand.Intn(len(valid))]
}
