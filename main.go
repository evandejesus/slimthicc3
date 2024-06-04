package main

import (
	"log"
	"os"
)

var uciInfo *log.Logger
var uciDebug *log.Logger
var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "", 0)
	uciInfo = log.New(os.Stdout, "info ", 0)
	uciDebug = log.New(os.Stdout, "debug ", 0)

	// run UCI loop
	UCI()
}
