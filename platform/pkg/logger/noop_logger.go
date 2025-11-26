package logger

import (
	"context"

	"go.uber.org/zap"
)

type NoopLogger struct{}

func (l *NoopLogger) Info(ctx context.Context, msg string, field ...zap.Field)  {}
func (l *NoopLogger) Error(ctx context.Context, msg string, field ...zap.Field) {}
