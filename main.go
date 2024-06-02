package main

import (
	"log"
	"os"
)

var uciInfo *log.Logger
var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "", 0)
	uciInfo = log.New(os.Stdout, "info ", 0)

	// run UCI loop
	UCI()
}
