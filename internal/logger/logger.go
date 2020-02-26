package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

/*
   Creation Time: 2019 - Mar - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	DefaultLevel  zap.AtomicLevel
	DefaultLogger *zapLogger
)

func init() {
	DefaultLevel = zap.NewAtomicLevel()
	consoleWriteSyncer := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "zapLogger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	DefaultLogger = NewZapLogger(
		zapcore.NewCore(consoleEncoder, consoleWriteSyncer, DefaultLevel),
		2,
	)

}

type Level = zapcore.Level
type Field = zapcore.Field
type CheckedEntry = zapcore.CheckedEntry

type Logger interface {
	Log(level Level, msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Check(Level, string) *CheckedEntry
	Sync() error
	SetLevel(level Level)
}

func log(l *zap.Logger, level Level, msg string, fields ...Field) {
	l.Check(level, msg).Write(fields...)
}

type zapLogger struct {
	*zap.Logger
	zap.AtomicLevel
}

func (l *zapLogger) Log(level Level, msg string, fields ...Field) {
	log(l.Logger, level, msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	log(l.Logger, DebugLevel, msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	log(l.Logger, InfoLevel, msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	log(l.Logger, WarnLevel, msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	log(l.Logger, ErrorLevel, msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	log(l.Logger, FatalLevel, msg, fields...)
}

func (l *zapLogger) Check(level Level, msg string) *CheckedEntry {
	return l.Logger.Check(level, msg)
}

func (l *zapLogger) SetLevel(level Level) {
	l.AtomicLevel.SetLevel(level)
}

func NewConsoleLogger() *zapLogger {
	l := new(zapLogger)
	l.AtomicLevel = zap.NewAtomicLevel()
	consoleWriteSyncer := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "zapLogger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	l.Logger = zap.New(
		zapcore.NewCore(consoleEncoder, consoleWriteSyncer, l.AtomicLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	)
	return l
}

func NewFileLogger(filename string) *zapLogger {
	l := new(zapLogger)

	l.AtomicLevel = zap.NewAtomicLevelAt(DebugLevel)
	fileLog, _ := os.Create(filename)

	syncer := zapcore.Lock(fileLog)
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	l.Logger = zap.New(
		zapcore.NewCore(encoder, syncer, l.AtomicLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	)
	return l
}

func NewZapLogger(core zapcore.Core, skip int) *zapLogger {
	l := new(zapLogger)
	l.Logger = zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(ErrorLevel),
		zap.AddCallerSkip(skip),
	)
	return l
}

func NewNop() *zapLogger {
	l := new(zapLogger)
	l.AtomicLevel = zap.NewAtomicLevel()
	l.Logger = zap.NewNop()
	return l
}

func InitLogger(logLevel zapcore.Level, sentryDSN string) Logger {
	level := zap.NewAtomicLevelAt(logLevel)
	consoleWriteSyncer := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriteSyncer, level),
	)

	zapLogger := NewZapLogger(core, 2)
	DefaultLogger = zapLogger
	return zapLogger
}

func Debug(msg string, fields ...Field) {
	log(DefaultLogger.Logger, DebugLevel, msg, fields...)
}

func Info(msg string, fields ...Field) {
	log(DefaultLogger.Logger, InfoLevel, msg, fields...)
}

func Warn(msg string, fields ...Field) {
	log(DefaultLogger.Logger, WarnLevel, msg, fields...)
}

func Error(msg string, fields ...Field) {
	log(DefaultLogger.Logger, ErrorLevel, msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	log(DefaultLogger.Logger, FatalLevel, msg, fields...)
}

func Check(level Level, msg string) *CheckedEntry {
	return DefaultLogger.Check(level, msg)
}

func SetLevel(level Level) {
	DefaultLevel.SetLevel(level)
}

func Enabled(level Level) bool {
	return DebugLevel.Enabled(level)
}
