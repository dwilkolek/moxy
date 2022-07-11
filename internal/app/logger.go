package app

import (
	"log"
	"os"
)

var MoxyLogger = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
