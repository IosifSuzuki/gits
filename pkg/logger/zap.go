package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Level = zapcore.Level
type Field = zap.Field

const (
	InfoLevel   Level = zap.InfoLevel
	WarnLevel   Level = zap.WarnLevel
	ErrorLevel  Level = zap.ErrorLevel
	DPanicLevel Level = zap.DPanicLevel
	PanicLevel  Level = zap.PanicLevel
	FatalLevel  Level = zap.FatalLevel
	DebugLevel  Level = zap.DebugLevel
)

var (
	Error   = zap.Error
	String  = zap.String
	Float64 = zap.Float64
	Float32 = zap.Float32
	Int     = zap.Int
	Bool    = zap.Bool
	Skip    = zap.Skip
)

type logger struct {
	log *zap.Logger
}

func NewLogger(level Level) Logger {
	conf := zap.NewDevelopmentConfig()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(conf.EncoderConfig),
		os.Stdout,
		level,
	)
	return &logger{
		log: zap.New(core),
	}
}

func (l *logger) Debug(msg string, fields ...Field) {
	l.log.Debug(msg, fields...)
}

func (l *logger) Info(msg string, fields ...Field) {
	l.log.Debug(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...Field) {
	l.log.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...Field) {
	l.log.Error(msg, fields...)
}

func (l *logger) DPanic(msg string, fields ...Field) {
	l.log.DPanic(msg, fields...)
}

func (l *logger) Panic(msg string, fields ...Field) {
	l.log.Panic(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...Field) {
	l.log.Fatal(msg, fields...)
}
