package worker

import (
	"context"
	"log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	log.Printf(format, v...)
}
