package main

import (
	_ "embed"
	"log/slog"
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
			slog.Debug("debug", "book_move", bookMoves[len(history)])
			notation := chess.UCINotation{}
			move, err := notation.Decode(game.Position(), bookMoves[len(history)])
			check(err)

			if !contains(potentialMoves, move) && contains(game.ValidMoves(), move) {
				potentialMoves = append(potentialMoves, move)
			}
		}
	}

	// if err := scanner.Err(); err != nil {
	// 	slog.Error(err.Error())
	// }

	if len(potentialMoves) == 0 {
		return nil
	}
	slog.Debug("debug", "potential_moves", potentialMoves)
	return potentialMoves[rand.Intn(len(potentialMoves))]
}

// position startpos moves e2e4 g7g6

func contains(s []*chess.Move, e *chess.Move) bool {
	for _, a := range s {
		if a.String() == e.String() {
			return true
		}
	}
	return false
}
