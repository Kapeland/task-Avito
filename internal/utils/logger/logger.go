package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Kapeland/task-Avito/internal/utils/config"
	slogmulti "github.com/samber/slog-multi"
	slogsampling "github.com/samber/slog-sampling"
)

var logger *Logger

func GetLogger() Logger {
	if logger != nil {
		return *logger
	}

	panic("Create logger first.")
}

type Logger struct {
	Logger   *slog.Logger
	logLevel *slog.LevelVar
}

func (lgr Logger) Info(msg, tp, method, after string) {
	lgr.Logger.Info(msg, slog.String("type", tp), slog.String("method", method), slog.String("after", after))
}

func (lgr Logger) InfoMsg(msg string, args ...any) {
	lgr.Logger.Info(msg, args...)
}

func (lgr Logger) Debug(msg, tp, method, after string) {
	lgr.Logger.Debug(msg, slog.String("type", tp), slog.String("method", method), slog.String("after", after))
}

func (lgr Logger) DebugMsg(msg string, args ...any) {
	lgr.Logger.Debug(msg, args...)
}

func (lgr Logger) Warn(msg, tp, method, after string) {
	lgr.Logger.Warn(msg, slog.String("type", tp), slog.String("method", method), slog.String("after", after))
}

func (lgr Logger) WarnMsg(msg string, args ...any) {
	lgr.Logger.Warn(msg, args...)
}

func (lgr Logger) Error(msg, tp, method, after string) {
	lgr.Logger.Error(msg, slog.String("type", tp), slog.String("method", method), slog.String("after", after))
}

func (lgr Logger) ErrorMsg(msg string, args ...any) {
	lgr.Logger.Error(msg, args...)
}

func CreateLogger(cfg *config.Config) Logger {
	samplingOption := slogsampling.UniformSamplingOption{
		Rate: cfg.Logger.LogRate,
	}

	logLevel := &slog.LevelVar{} // INFO

	switch strings.ToLower(cfg.Logger.Lvl) {
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "warn":
		logLevel.Set(slog.LevelWarn)
	case "error":
		logLevel.Set(slog.LevelError)
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)

	if !cfg.Project.Debug {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	logger = &Logger{
		Logger: slog.New(
			slogmulti.
				Pipe(samplingOption.NewMiddleware()).
				Handler(handler),
		),
		logLevel: logLevel,
	}

	return *logger
}

// SetLogLvl allows to set log lvl in runtime
func SetLogLvl(lvl slog.Level) {
	if logger != nil {
		logger.logLevel.Set(lvl)
	}
}
