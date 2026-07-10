package logger

import "go.uber.org/zap"

type logger struct {
	*zap.Logger
}

type Logger interface {
}
