package main

import (
	"github.com/notnil/chess"
)

const (
	pawnPhase   int = 0
	knightPhase int = 1
	bishopPhase int = 1
	rookPhase   int = 2
	queenPhase  int = 4
	totalPhase  int = pawnPhase*16 + knightPhase*4 + bishopPhase*4 + rookPhase*4 + queenPhase*2
	inf         int = 2147483647

	scaleFactor  int = 16
	tempoBonusMG int = 14
)

var PhaseValues = map[chess.PieceType]int{
	chess.Pawn:   pawnPhase,
	chess.Knight: knightPhase,
	chess.Bishop: bishopPhase,
	chess.Rook:   rookPhase,
	chess.Queen:  queenPhase,
}

var PieceValueMG = map[chess.PieceType]int{chess.Pawn: 84, chess.Knight: 333, chess.Bishop: 346, chess.Rook: 441, chess.Queen: 921}
var PieceValueEG = map[chess.PieceType]int{chess.Pawn: 106, chess.Knight: 244, chess.Bishop: 268, chess.Rook: 478, chess.Queen: 886}

type Eval struct {
	MGScores [2]int // middlegame scores
	EGScores [2]int // endgame scores
	phase    int
}

var MG_PST = [6][64]int{
	{
		// MG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		45, 52, 42, 43, 28, 34, 19, 9,
		-14, -3, 7, 14, 35, 50, 15, -6,
		-27, -6, -8, 13, 16, 4, -3, -25,
		-32, -28, -7, 5, 7, -1, -15, -30,
		-29, -25, -12, -12, -1, -5, 6, -17,
		-34, -23, -27, -18, -14, 10, 13, -22,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// MG Knight PST
		-43, -11, -8, -5, 1, -20, -4, -22,
		-31, -22, 19, 7, 5, 13, -8, -11,
		-21, 21, 8, 16, 36, 33, 19, 6,
		-6, 2, 0, 23, 8, 27, 4, 14,
		-3, 10, 12, 8, 16, 10, 19, 1,
		-19, -4, 3, 7, 22, 12, 15, -11,
		-21, -20, -9, 8, 9, 11, -5, 0,
		-19, -13, -20, -14, -2, 3, -11, -8,
	},
	{
		// MG Bishop PST
		-13, 0, -17, -8, -7, -5, -2, -3,
		-21, 0, -16, -10, 4, 1, -6, -41,
		-23, 6, 10, 8, 8, 26, 0, -10,
		-15, -4, 2, 22, 9, 10, -1, -16,
		0, 10, -2, 15, 17, -7, -1, 13,
		-2, 16, 13, 0, 5, 16, 14, 0,
		8, 11, 12, 3, 11, 23, 27, 3,
		-26, 3, -3, -1, 10, -5, -7, -15,
	},
	{
		// MG Rook PST
		3, 1, 0, 7, 7, -1, 0, 0,
		-6, -9, 7, 7, 7, 5, -4, -1,
		-12, 11, 0, 17, -2, 12, 23, -1,
		-17, -9, 4, 0, 3, 15, -1, -2,
		-24, -16, -16, -4, -1, -14, 2, -20,
		-30, -15, -6, -3, 0, 2, 2, -15,
		-25, -6, -6, 5, 8, 6, 8, -46,
		-3, 1, 6, 15, 17, 14, -13, -2,
	},
	{
		// MG Queen PST
		-10, 0, 0, 0, 10, 9, 5, 7,
		-19, -35, -5, 2, -9, 7, 1, 15,
		-10, -7, -4, -9, 15, 29, 24, 22,
		-14, -14, -15, -11, -1, -5, 3, -6,
		-8, -20, -8, -5, -4, -2, 2, -2,
		-13, 5, 2, 1, -1, 8, 4, 2,
		-20, 0, 10, 16, 16, 16, -6, 6,
		-3, -1, 7, 19, 5, -10, -9, -17,
	},
	{
		// MG King PST
		-3, 0, 2, 0, 0, 0, 1, -1,
		1, 4, 0, 7, 4, 2, 3, -2,
		2, 4, 7, 4, 4, 14, 12, 0,
		0, 2, 6, 0, 0, 2, 6, -9,
		-8, 5, 0, -8, -10, -10, -9, -23,
		-3, 5, 1, -8, -12, -12, 8, -24,
		6, 13, 0, -40, -23, -1, 25, 19,
		-28, 29, 17, -53, 2, -25, 34, 15,
	},
}

var EG_PST [6][64]int = [6][64]int{
	{
		// EG Pawn PST
		0, 0, 0, 0, 0, 0, 0, 0,
		77, 74, 63, 53, 59, 60, 72, 77,
		17, 11, 11, 11, 11, -6, 14, 8,
		-3, -14, -18, -31, -29, -25, -20, -18,
		-12, -14, -24, -31, -29, -28, -27, -28,
		-22, -20, -25, -20, -21, -24, -34, -34,
		-16, -22, -11, -19, -13, -23, -32, -34,
		0, 0, 0, 0, 0, 0, 0, 0,
	},
	{
		// EG Knight PST
		-36, -16, -7, -14, -4, -20, -20, -29,
		-17, 2, -7, 14, 2, -7, -9, -19,
		-13, -7, 14, 12, 4, 6, 0, -13,
		-5, 8, 24, 18, 22, 15, 11, -4,
		-3, 4, 20, 30, 22, 25, 15, -2,
		-7, 1, 3, 19, 10, -2, -4, -4,
		-10, -2, -1, 0, 6, -8, -3, -13,
		-12, -28, -8, 1, -5, -12, -27, -12,
	},
	{
		// EG Bishop PST
		-9, -5, -9, -5, -2, -4, -5, -8,
		0, 2, 8, -7, 1, 0, -2, -8,
		8, 0, 0, 1, 0, 1, 5, 6,
		0, 7, 7, 8, 3, 5, 2, 6,
		-1, 0, 12, 8, 0, 6, 0, -5,
		0, 0, 3, 6, 8, -1, 0, -1,
		-6, -12, -7, 0, 0, -8, -9, -13,
		-11, 0, -6, 0, -3, -4, -5, -9,
	},
	{
		// EG Rook PST
		8, 9, 11, 13, 13, 12, 13, 9,
		3, 5, 1, 0, -1, 0, 6, 2,
		9, 5, 7, 2, 2, 1, 0, 0,
		3, 3, 6, 0, 0, 0, 0, 4,
		5, 4, 9, 0, -3, -2, -6, -2,
		0, 0, -6, -5, -9, -14, -7, -12,
		-2, -5, -1, -7, -9, -11, -13, -1,
		-7, -3, 0, -8, -13, -12, -4, -24,
	},
	{
		// EG Queen PST
		-12, 4, 8, 4, 10, 9, 3, 6,
		-17, -7, -1, 7, 3, 6, 1, 0,
		-5, -1, -4, 12, 14, 20, 12, 14,
		-2, 2, 2, 9, 13, 7, 18, 22,
		-9, 3, 1, 15, 5, 10, 12, 10,
		-6, -20, 0, -15, 0, -1, 10, 7,
		-6, -14, -31, -27, -19, -12, -11, -4,
		-12, -22, -19, -30, -8, -13, -6, -15,
	},
	{
		// EG King PST
		-15, -11, -11, -6, -2, 3, 4, -9,
		-9, 14, 11, 13, 13, 28, 19, 1,
		-1, 18, 19, 15, 16, 35, 34, 4,
		-12, 14, 21, 25, 19, 25, 18, -5,
		-23, -6, 14, 21, 20, 18, 5, -16,
		-21, -6, 5, 13, 15, 9, -2, -12,
		-27, -10, 2, 9, 9, 1, -12, -26,
		-43, -34, -20, -5, -26, -9, -35, -55,
	},
}

var PassedPawnPST_MG = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	45, 52, 42, 43, 28, 34, 19, 9,
	48, 43, 43, 30, 24, 31, 12, 2,
	28, 17, 13, 10, 10, 19, 6, 1,
	14, 0, -9, -7, -13, -7, 9, 16,
	5, 3, -3, -14, -3, 10, 13, 19,
	8, 9, 2, -8, -3, 8, 16, 9,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var PassedPawnPST_EG = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	77, 74, 63, 53, 59, 60, 72, 77,
	91, 83, 66, 40, 30, 61, 67, 84,
	55, 52, 42, 35, 30, 34, 56, 52,
	29, 26, 21, 18, 17, 19, 34, 30,
	8, 6, 5, 1, 1, -1, 14, 7,
	2, 3, -4, 0, -2, -1, 7, 6,
	0, 0, 0, 0, 0, 0, 0, 0,
}

// Flip white's perspective to black
var flipSq [2][64]int = [2][64]int{
	{
		0, 1, 2, 3, 4, 5, 6, 7,
		8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23,
		24, 25, 26, 27, 28, 29, 30, 31,
		32, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 44, 45, 46, 47,
		48, 49, 50, 51, 52, 53, 54, 55,
		56, 57, 58, 59, 60, 61, 62, 63,
	},

	{
		56, 57, 58, 59, 60, 61, 62, 63,
		48, 49, 50, 51, 52, 53, 54, 55,
		40, 41, 42, 43, 44, 45, 46, 47,
		32, 33, 34, 35, 36, 37, 38, 39,
		24, 25, 26, 27, 28, 29, 30, 31,
		16, 17, 18, 19, 20, 21, 22, 23,
		8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7,
	},
}

func EvaluateBoard(pos *chess.Position) int {

	eval := Eval{
		MGScores: [2]int{0, 0},
		EGScores: [2]int{0, 0},
		phase:    0,
	}

	// pass 1: determine base middlegame and endgame scores plus phase of game
	for sqr, piece := range pos.Board().SquareMap() {
		pieceType := piece.Type()
		switch pieceType {
		case chess.NoPieceType:
			continue
		default:
			eval.MGScores[piece.Color()-1] += PieceValueMG[pieceType] + MG_PST[pieceType-1][flipSq[piece.Color()-1][int(sqr)]]
			eval.EGScores[piece.Color()-1] += PieceValueEG[pieceType] + EG_PST[pieceType-1][flipSq[piece.Color()-1][int(sqr)]]
			eval.phase -= PhaseValues[pieceType]
		}
	}

	// pawnRank := make(map[chess.Color][]int)
	// pawnRank[chess.White] = make([]int, 10)
	// pawnRank[chess.Black] = make([]int, 10)
	// pieceMaterial := make(map[chess.Color]float64)
	// pawnMaterial := 0.0
	// score := 0.0

	// for i := 0; i < 10; i++ {
	// 	pawnRank[chess.White][i] = 0
	// 	pawnRank[chess.Black][i] = 7
	// }

	// // pass 1
	// for sqr, piece := range pos.Board().SquareMap() {
	// 	switch piece.Type() {
	// 	case chess.NoPieceType:
	// 		continue
	// 	case chess.Pawn:
	// 		if pawnRank[piece.Color()][sqr.File()] < int(sqr.Rank()) {
	// 			pawnRank[piece.Color()][sqr.File()] = int(sqr.Rank())
	// 		}
	// 		if piece.Color() == chess.White {
	// 			pawnMaterial += float64(pieceValues[chess.Pawn])
	// 		} else {
	// 			pawnMaterial -= float64(pieceValues[chess.Pawn])
	// 		}
	// 	default:
	// 		pieceMaterial[piece.Color()] += float64(pieceValues[piece.Type()])
	// 	}

	// }

	// // pass 2
	// score = (pieceMaterial[chess.White] - pieceMaterial[chess.Black]) + pawnMaterial
	// for sqr, piece := range pos.Board().SquareMap() {
	// 	switch piece.Type() {
	// 	case chess.Pawn:
	// 		if piece.Color() == chess.White {
	// 			score += pawnPST[int(sqr)]
	// 		} else {
	// 			score -= pawnPST[flipPST[int(sqr)]]
	// 		}
	// 	case chess.Knight:
	// 		if piece.Color() == chess.White {
	// 			score += knightPST[int(sqr)]
	// 		} else {
	// 			score -= knightPST[flipPST[int(sqr)]]
	// 		}
	// 	case chess.Bishop:
	// 		if piece.Color() == chess.White {
	// 			score += bishopPST[int(sqr)]
	// 		} else {
	// 			score -= bishopPST[flipPST[int(sqr)]]
	// 		}
	// 	case chess.King:
	// 		if piece.Color() == chess.White {
	// 			if pieceMaterial[chess.Black] < 1200 {
	// 				score += kingEndgamePST[int(sqr)]
	// 			} else {
	// 				score += kingPST[int(sqr)]
	// 			}
	// 		} else {
	// 			if pieceMaterial[chess.White] < 1200 {

	// 				score -= kingEndgamePST[flipPST[int(sqr)]]
	// 			} else {
	// 				score -= kingPST[flipPST[int(sqr)]]
	// 			}
	// 		}
	// 	}
	// }

	// return score
	eval.MGScores[pos.Turn()-1] += tempoBonusMG

	mgScore := eval.MGScores[pos.Turn()-1] - eval.MGScores[(pos.Turn())%2]
	egScore := eval.EGScores[pos.Turn()-1] - eval.EGScores[(pos.Turn())%2]

	phase := (eval.phase*256 + (totalPhase / 2)) / totalPhase
	score := int(((int(mgScore) * (int(256) - int(phase))) + (int(egScore) * int(phase))) / int(256))
	return score
}
