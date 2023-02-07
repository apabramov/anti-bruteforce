package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	Log *zap.Logger
}

func New(level string) (*Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	en := zapcore.NewJSONEncoder(config)

	al, err := zap.ParseAtomicLevel(level)

	core := zapcore.NewTee(
		zapcore.NewCore(en, zapcore.AddSync(os.Stdout), al),
	)

	l := zap.New(core)

	log := l.Named("anti-bruteforce")
	return &Logger{log}, err
}

func (l Logger) Info(msg string) {
	l.Log.Info(msg)
}

func (l Logger) Error(msg string) {
	l.Log.Error(msg)
}

func (l Logger) Warn(msg string) {
	l.Log.Warn(msg)
}

func (l Logger) Debug(msg string) {
	l.Log.Debug(msg)
}
