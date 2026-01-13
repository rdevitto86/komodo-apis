package logger

import (
	"context"
	"komodo-forge-apis-go/http/services/redaction"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/lmittmann/tint"
)

var (
	slogger  *slog.Logger
	initOnce sync.Once
)

func Init(name string, level string, env string) {
	initOnce.Do(func() {
		parseLevel := func() slog.Level {
			switch strings.ToUpper(level) {
				case "DEBUG":  					return slog.LevelDebug
				case "WARN":  					return slog.LevelWarn
				case "ERROR", "FATAL": 	return slog.LevelError
				default:      					return slog.LevelInfo
			}
		}
		isLocal := func() bool {
			env = strings.ToLower(env)
			return env == "local" || env == "dev" || env == "development"
		}

		var maxLevel slog.LevelVar
		var handler slog.Handler

		maxLevel.Set(parseLevel())

		if isLocal() {
			handler = tint.NewHandler(os.Stdout, &tint.Options{Level: &maxLevel, TimeFormat: "15:04:05"})
		} else {
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: &maxLevel})
		}

		slogger = slog.New(&RedactingLogger{Handler: handler}).With(slog.String("app", name))
		slog.SetDefault(slogger)
	})
}

func Debug(msg string, args ...any) { slogger.Debug(msg, args...) }
func Info(msg string, args ...any) { slogger.Info(msg, args...) }
func Warn(msg string, args ...any) { slogger.Warn(msg, args...) }

func Error(msg string, err error, args ...any) {
	if err != nil { args = append(args, AttrError(err)) }
	slogger.Error(msg, args...)
}

func Fatal(msg string, err error, args ...any) {
	if err != nil { args = append(args, AttrError(err)) }
	slogger.Error(msg, args...)
	os.Exit(1)
}

type RedactingLogger struct { slog.Handler }

func (rl *RedactingLogger) Handle(ctx context.Context, rec slog.Record) error {
	clean := rec.Clone()
	rec.Attrs(func(attr slog.Attr) bool {
		clean.AddAttrs(slog.Any(attr.Key, redaction.RedactPair(attr.Key, attr.Value.Any())))
		return true
	})
	return rl.Handler.Handle(ctx, clean)
}
