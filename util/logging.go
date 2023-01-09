package util

import (
	"context"
	stdlog "log"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var gLogger Logger

func init() {
	InitGLogger()
}

func GLogger() Logger {
	return gLogger
}

func InitGLogger(options ...loggerOption) {
	logger, err := NewLoggerZapImplWithOption(options...)
	if err != nil {
		logger = NewLoggerStdImplWithOption()
		logger.Warn(err, "logger initialization failed. fall back logger.")
	}
	gLogger = logger
}

type LoggerLevel int

const (
	LoggerLevelDebug LoggerLevel = iota - 1
	LoggerLevelInfo
	LoggerLevelWarn
	LoggerLevelError
)

func (l LoggerLevel) String() string {
	switch l {
	case LoggerLevelDebug:
		return "debug"
	case LoggerLevelInfo:
		return "info"
	case LoggerLevelWarn:
		return "warn"
	case LoggerLevelError:
		return "error"
	default:
		return ""
	}
}

type LoggerFormat int

const (
	LoggerFormatJSON LoggerFormat = iota
	LoggerFormatConsole
)

func ParseLoggerLevelStr(v string) LoggerLevel {
	switch v {
	default:
		return LoggerLevelInfo
	case LoggerLevelDebug.String():
		return LoggerLevelDebug
	case LoggerLevelWarn.String():
		return LoggerLevelDebug
	case LoggerLevelError.String():
		return LoggerLevelDebug
	}
}

func ParseLoggerFormatStr(v string) LoggerFormat {
	switch v {
	default:
		return LoggerFormatJSON
	case LoggerFormatConsole.String():
		return LoggerFormatConsole
	}
}

func (f LoggerFormat) String() string {
	switch f {
	case LoggerFormatConsole:
		return "console"
	default:
		return "json"
	}
}

type loggerOption interface {
	apply(*loggerConf)
}

type loggerConf struct {
	Level  LoggerLevel
	Format LoggerFormat
}

type LoggerLevelOption int

func (o LoggerLevelOption) apply(c *loggerConf) {
	level := LoggerLevel(o)
	if level >= LoggerLevelDebug && level <= LoggerLevelError {
		c.Level = LoggerLevel(o)
	}
}

func OptionLoggerLevel(v LoggerLevel) LoggerLevelOption {
	return LoggerLevelOption(v)
}

type LoggerFormatOption int

func (o LoggerFormatOption) apply(c *loggerConf) {
	format := LoggerFormat(o)
	if format >= LoggerFormatJSON && format <= LoggerFormatConsole {
		c.Format = LoggerFormat(o)
	}
}

func OptionLoggerFormat(v LoggerFormat) LoggerFormatOption {
	return LoggerFormatOption(v)
}

type Logger interface {
	Enabled() bool
	WithValues(keysAndValues ...interface{}) Logger
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(err error, msg string, keysAndValues ...interface{})
	Error(err error, msg string, keysAndValues ...interface{})
	WithCallDepth(depth int) Logger
}

type loggerContextKey struct{}

func FromContext(ctx context.Context) Logger {
	if v, ok := ctx.Value(loggerContextKey{}).(Logger); ok {
		return v
	}

	return gLogger
}

func NewContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

type LoggerZapImpl struct {
	logger logr.Logger
}

func NewLoggerZapImplWithOption(options ...loggerOption) (Logger, error) {
	conf := &loggerConf{
		Level: LoggerLevelInfo,
	}
	for _, opt := range options {
		opt.apply(conf)
	}
	zc := zap.NewProductionConfig()
	zc.Level = zap.NewAtomicLevelAt(zapcore.Level(conf.Level))
	zc.EncoderConfig.TimeKey = "time"
	zc.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zc.Encoding = conf.Format.String()
	z, err := zc.Build()
	if err != nil {
		return nil, err
	}

	return &LoggerZapImpl{
		logger: zapr.NewLogger(z),
	}, nil
}

func (l *LoggerZapImpl) Enabled() bool {
	return l.logger.Enabled()
}

func (l *LoggerZapImpl) WithValues(keysAndValues ...interface{}) Logger {
	return &LoggerZapImpl{
		logger: l.logger.WithValues(keysAndValues...),
	}
}

func (l *LoggerZapImpl) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.GetSink().Info(int(LoggerLevelDebug), msg)
}

func (l *LoggerZapImpl) Info(msg string, keysAndValues ...interface{}) {
	l.logger.WithCallDepth(1).Info(msg, keysAndValues...)
}

func (l *LoggerZapImpl) Warn(err error, msg string, keysAndValues ...interface{}) {
	l.logger.GetSink().Info(int(LoggerLevelWarn), msg, "error", err)
}

func (l *LoggerZapImpl) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.WithCallDepth(1).Error(err, msg, keysAndValues...)
}

func (l *LoggerZapImpl) WithCallDepth(depth int) Logger {
	return &LoggerZapImpl{
		logger: l.logger.WithCallDepth(depth),
	}
}

type LoggerStdImpl struct {
	logger logr.Logger
}

func NewLoggerStdImplWithOption() Logger {
	return &LoggerZapImpl{
		logger: stdr.NewWithOptions(stdlog.New(os.Stderr, "", stdlog.LstdFlags), stdr.Options{LogCaller: stdr.All}),
	}
}

func (l *LoggerStdImpl) Enabled() bool {
	return l.logger.Enabled()
}

func (l *LoggerStdImpl) WithValues(keysAndValues ...interface{}) Logger {
	return &LoggerStdImpl{
		logger: l.logger.WithValues(keysAndValues...),
	}
}

func (l *LoggerStdImpl) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.WithCallDepth(1).V(1).Info(msg, keysAndValues...)
}

func (l *LoggerStdImpl) Info(msg string, keysAndValues ...interface{}) {
	l.logger.WithCallDepth(1).Info(msg, keysAndValues...)
}

func (l *LoggerStdImpl) Warn(err error, msg string, keysAndValues ...interface{}) {
	l.logger.GetSink().Info(-int(LoggerLevelWarn), msg, "error", err)
}

func (l *LoggerStdImpl) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.WithCallDepth(1).Error(err, msg, keysAndValues...)
}

func (l *LoggerStdImpl) WithCallDepth(depth int) Logger {
	return &LoggerStdImpl{
		logger: l.logger.WithCallDepth(depth),
	}
}
