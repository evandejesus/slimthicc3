package main

import (
	"log/slog"

	"github.com/notnil/chess"
)

var pieceValues = map[chess.PieceType]int{chess.Pawn: 100, chess.Knight: 280, chess.Bishop: 320, chess.Rook: 479, chess.Queen: 929, chess.King: 0}

var knightPST = []float64{
	-10, -10, -10, -10, -10, -10, -10, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, -30, -10, -10, -10, -10, -30, -10,
}
var bishopPST = []float64{
	-10, -10, -10, -10, -10, -10, -10, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 5, 5, 5, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, -10, -20, -10, -10, -20, -10, -10,
}

var kingPST = []float64{
	-40, -40, -40, -40, -40, -40, -40, -40,
	-40, -40, -40, -40, -40, -40, -40, -40,
	-40, -40, -40, -40, -40, -40, -40, -40,
	-40, -40, -40, -40, -40, -40, -40, -40,
	-40, -40, -40, -40, -40, -40, -40, -40,
	-40, -40, -40, -40, -40, -40, -40, -40,
	-20, -20, -20, -20, -20, -20, -20, -20,
	0, 20, 40, -20, 0, -20, 40, 20,
}

var kingEndgamePST = []float64{
	0, 10, 20, 30, 30, 20, 10, 0,
	10, 20, 30, 40, 40, 30, 20, 10,
	20, 30, 40, 50, 50, 40, 30, 20,
	30, 40, 50, 60, 60, 50, 40, 30,
	30, 40, 50, 60, 60, 50, 40, 30,
	20, 30, 40, 50, 50, 40, 30, 20,
	10, 20, 30, 40, 40, 30, 20, 10,
	0, 10, 20, 30, 30, 20, 10, 0,
}

var pawnPST = []float64{
	0, 0, 0, 0, 0, 0, 0, 0,
	5, 10, 15, 20, 20, 15, 10, 5,
	4, 8, 12, 16, 16, 12, 8, 4,
	3, 6, 9, 12, 12, 9, 6, 3,
	2, 4, 6, 8, 8, 6, 4, 2,
	1, 2, 3, -10, -10, 3, 2, 1,
	0, 0, 0, -40, -40, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var flipPST = []int{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

func simpleBestMove(game *chess.Game) *chess.Move {
	bestValue := -9999.0
	var bestMove = &chess.Move{}
	for _, m := range game.ValidMoves() {
		tmpGame := game.Clone()
		tmpGame.Move(m)
		eval := -evaluateBoard(tmpGame.Position())
		slog.Debug("evaluateBoard2", "move", m.String(), "eval", eval)

		if eval > bestValue {
			bestValue = eval
			bestMove = m
		}
	}
	return bestMove
}

func evaluateBoard(pos *chess.Position) float64 {
	pawnRank := make(map[chess.Color][]int)
	pawnRank[chess.White] = make([]int, 10)
	pawnRank[chess.Black] = make([]int, 10)
	pieceMaterial := make(map[chess.Color]float64)
	pawnMaterial := 0.0
	score := 0.0

	for i := 0; i < 10; i++ {
		pawnRank[chess.White][i] = 0
		pawnRank[chess.Black][i] = 7
	}

	// pass 1
	for sqr, piece := range pos.Board().SquareMap() {
		switch piece.Type() {
		case chess.NoPieceType:
			continue
		case chess.Pawn:
			if pawnRank[piece.Color()][sqr.File()] < int(sqr.Rank()) {
				pawnRank[piece.Color()][sqr.File()] = int(sqr.Rank())
			}
			if piece.Color() == chess.White {
				pawnMaterial += float64(pieceValues[chess.Pawn])
			} else {
				pawnMaterial -= float64(pieceValues[chess.Pawn])
			}
		default:
			pieceMaterial[piece.Color()] += float64(pieceValues[piece.Type()])
		}

	}

	// pass 2
	score = (pieceMaterial[chess.White] - pieceMaterial[chess.Black]) + pawnMaterial
	for sqr, piece := range pos.Board().SquareMap() {
		switch piece.Type() {
		case chess.Pawn:
			if piece.Color() == chess.White {
				score += pawnPST[int(sqr)]
			} else {
				score -= pawnPST[flipPST[int(sqr)]]
			}
		case chess.Knight:
			if piece.Color() == chess.White {
				score += knightPST[int(sqr)]
			} else {
				score -= knightPST[flipPST[int(sqr)]]
			}
		case chess.Bishop:
			if piece.Color() == chess.White {
				score += bishopPST[int(sqr)]
			} else {
				score -= bishopPST[flipPST[int(sqr)]]
			}
		case chess.King:
			if piece.Color() == chess.White {
				if pieceMaterial[chess.Black] < 1200 {
					score += kingEndgamePST[int(sqr)]
				} else {
					score += kingPST[int(sqr)]
				}
			} else {
				if pieceMaterial[chess.White] < 1200 {

					score -= kingEndgamePST[flipPST[int(sqr)]]
				} else {
					score -= kingPST[flipPST[int(sqr)]]
				}
			}
		}
	}

	if pos.Turn() == chess.White {
		return score
	}
	return -score
}
