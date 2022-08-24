package logger

import (
	"fmt"
	"log"
	"os"
)

func New(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("%-25s", prefix), log.Ldate|log.Ltime)
}

func NewOnPort(prefix string, port int) *log.Logger {
	return New(fmt.Sprintf("%s  :%d", prefix, port))
}
