package main

import (
	"log"
	"math/rand"

	"github.com/notnil/chess"
)

var pieceValues = map[chess.PieceType]int{chess.Pawn: 100, chess.Knight: 280, chess.Bishop: 320, chess.Rook: 479, chess.Queen: 929, chess.King: 60000}

func randomMove(game *chess.Game) *chess.Move {
	valid := game.ValidMoves()
	return valid[rand.Intn(len(valid))]
}

func simpleBestMove(game *chess.Game) *chess.Move {
	bestValue := -9999
	var bestMove = &chess.Move{}
	for _, m := range game.ValidMoves() {
		tmpGame := game.Clone()
		tmpGame.Move(m)
		eval := -evaluateBoard(tmpGame.Position())
		log.Printf("info move: %s eval: %d\n", m.String(), eval)

		if eval > bestValue {
			bestValue = eval
			bestMove = m
		}
	}
	return bestMove
}

func evaluateBoard(pos *chess.Position) int {
	eval := 0
	color := pos.Turn()
	squares := pos.Board().SquareMap()
	for _, v := range squares {
		if v.Color() == color {
			eval += pieceValues[v.Type()]
		} else {
			eval -= pieceValues[v.Type()]
		}
	}
	return eval
}
