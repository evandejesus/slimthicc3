package main

import (
	"github.com/notnil/chess"
)

func Search(game *chess.Game, depth int, alpha int, beta int) (int, *chess.Move) {
	if st := game.Position().Status(); depth == 0 || st == chess.Checkmate || st == chess.Stalemate {
		return EvaluateBoard(game.Position()), nil
	}

	var bestMove *chess.Move
	if game.Position().Turn() == chess.White {
		// turn is White
		bestMove = nil
		for _, move := range game.ValidMoves() {
			newPos := game.Clone()
			err := newPos.Move(move)
			check(err)
			score, _ := Search(newPos, depth-1, alpha, beta)
			if score > alpha {
				alpha = score
				bestMove = move
				if alpha >= beta {
					break
				}
			}
		}
		return alpha, bestMove
	} else {
		// turn is Black
		bestMove = nil
		for _, move := range game.ValidMoves() {
			newPos := game.Clone()
			err := newPos.Move(move)
			check(err)
			score, _ := Search(newPos, depth-1, alpha, beta)
			if score < beta {
				beta = score
				bestMove = move
				if alpha >= beta {
					break
				}
			}
		}
		return beta, bestMove
	}

}
