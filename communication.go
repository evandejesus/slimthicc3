package main

import (
	"bufio"
	"os"
	"reflect"
	"strings"

	"github.com/notnil/chess"
)

func UCI() {
	game := chess.NewGame()
	uciNotation := chess.UCINotation{}
	// searcher
	r := bufio.NewReader(os.Stdin)
	for {
		input, _ := r.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")
		switch {
		case args[0] == "quit":
			return
		case args[0] == "isready":
			logger.Println("readyok")
		case args[0] == "uci":
			logger.Println("id name slimthicc3")
			logger.Println("id author evan")
			logger.Println("uciok")
		case args[0] == "ucinewgame" || len(args) == 2 && reflect.DeepEqual(args[:2], []string{"position", "startpos"}):
			game = chess.NewGame()
		case args[0] == "position":
			// handle startpos
			switch {
			case args[1] == "startpos":
				game = chess.NewGame()
				if len(args) > 2 && args[2] == "moves" {
					for _, s := range args[3:] {
						m, err := uciNotation.Decode(game.Position(), s)
						check(err)
						game.Move(m)
					}
				}
			// handle fen
			case args[1] == "fen":
				// position fen rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
				f, err := chess.FEN(strings.Join(args[2:8], " "))
				check(err)
				game = chess.NewGame(f)
				if len(args) > 8 && args[8] == "moves" {
					for _, s := range args[9:] {
						m, err := uciNotation.Decode(game.Position(), s)
						check(err)

						game.Move(m)
					}
				}
			}

		case args[0] == "go":
			var m *chess.Move

			m = getBookMove(game)
			if m == nil {
				var eval int
				eval, m = Search(game, 3, -inf, inf)
				uciDebug.Println("eval", eval)
			}
			logger.Printf("bestmove %s", m)
		}
	}
}
