package main

import (
	_ "embed"
	"math/rand"
	"strings"

	"github.com/notnil/chess"
)

//go:embed book.txt
var bookContents []byte

func bookMatch(history []*chess.Move, bookMoves []string) bool {
	for i, m := range history {
		if bookMoves[i] != m.String() {
			return false
		}
	}
	return true
}

func getBookMove(game *chess.Game) *chess.Move {
	bookLines := strings.Split(string(bookContents), "\n")

	var potentialMoves []*chess.Move

	for _, line := range bookLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		bookMoves := strings.Split(line, " ")
		history := game.Moves()

		if len(history) < len(bookMoves) && bookMatch(history, bookMoves) {
			notation := chess.UCINotation{}
			move, err := notation.Decode(game.Position(), bookMoves[len(history)])
			check(err)

			if !contains(potentialMoves, move) && contains(game.ValidMoves(), move) {
				potentialMoves = append(potentialMoves, move)
			}
		}
	}

	if len(potentialMoves) == 0 {
		return nil
	}
	return potentialMoves[rand.Intn(len(potentialMoves))]
}

// position startpos moves e2e4 g7g6
// position startpos moves e2e4 d7d6 d2d4 g7g6 g1f3
// position startpos moves e2e4 e7e6 d2d4 d7d5 e4e5 c7c5 d4c5

func contains(s []*chess.Move, e *chess.Move) bool {
	for _, a := range s {
		if a.String() == e.String() {
			return true
		}
	}
	return false
}
