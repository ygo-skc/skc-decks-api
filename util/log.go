package util

import (
	"context"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

func init() {
	slogOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, slogOpts)))
}

func NewRequestSetup(ctx context.Context, operation string, customAttributes ...slog.Attr) (*slog.Logger, context.Context) {
	defaults := []any{slog.String("requestID", uuid.New().String()), slog.String("operation", operation)}

	for _, customAttribute := range customAttributes {
		defaults = append(defaults, customAttribute)
	}

	l := slog.With(defaults...)
	return l, context.WithValue(ctx, loggerKey, l)
}

func AddAttribute(ctx context.Context, customAttributes ...slog.Attr) (*slog.Logger, context.Context) {
	newAttributes := []any{}

	for _, customAttribute := range customAttributes {
		newAttributes = append(newAttributes, customAttribute)
	}

	l := LoggerFromContext(ctx)
	l = l.With(newAttributes...)
	return l, context.WithValue(ctx, loggerKey, l)
}
