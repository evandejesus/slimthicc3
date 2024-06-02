package main

import (
	"math"

	"github.com/notnil/chess"
)

func Search(game *chess.Game, depth int) (float64, *chess.Move) {
	if st := game.Position().Status(); depth == 0 || st == chess.Checkmate || st == chess.Stalemate {
		return -EvaluateBoard(game.Position()), nil
	}

	var bestScore float64
	var bestMove *chess.Move
	if game.Position().Turn() == chess.White {
		// turn is White
		bestScore = math.Inf(-1)
		bestMove = nil
		for _, move := range game.ValidMoves() {
			newPos := game.Clone()
			err := newPos.Move(move)
			check(err)
			score, _ := Search(newPos, depth-1)
			if score > bestScore {
				uciInfo.Println("White", depth, score, move)
				bestScore = score
				bestMove = move
			}
		}
		return bestScore, bestMove
	} else {
		// turn is Black
		bestScore = math.Inf(1)
		bestMove = nil
		for _, move := range game.ValidMoves() {
			newPos := game.Clone()
			err := newPos.Move(move)
			check(err)
			score, _ := Search(newPos, depth-1)
			if score < bestScore {
				uciInfo.Println("Black", depth, score, move)
				bestScore = score
				bestMove = move
			}
		}
		return bestScore, bestMove
	}

}
