package logger

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/zgsm-ai/gatewayctl/internal/pkg/config"
)

const ctxLoggerKey = "zapLogger"

var (
	logger *Logger
	once   sync.Once
)

const LogLevelDebug = "debug"

type Logger struct {
	*zap.SugaredLogger
}

// InitLogger 初始化日志
func InitLogger(opts *Opts) {
	once.Do(
		func() {
			logger = NewLog(opts)
		},
	)
}

type Opts struct {
	FileName   string
	Level      string
	MaxSize    int
	MaxBackUps int
	MaxAge     int
	Compress   bool
	Encoding   string
	Env        string
}

func NewLog(opts *Opts) *Logger {
	log := NewZapLogger(opts)
	return &Logger{log.Sugar()}
}

func NewZapLogger(opts *Opts) *zap.Logger {
	// log address "out.log" User-defined
	lp := opts.FileName
	lv := opts.Level
	var level zapcore.LevelEnabler
	//debug<info<warn<error<fatal<panic
	switch lv {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	hook := lumberjack.Logger{
		Filename:   lp,              // Log file path
		MaxSize:    opts.MaxSize,    // Maximum size unit for each log file: M
		MaxBackups: opts.MaxBackUps, // The maximum number of backups that can be saved for log files
		MaxAge:     opts.MaxAge,     // Maximum number of days the file can be saved
		Compress:   opts.Compress,   // Compression or not
	}

	var encoder zapcore.Encoder
	if opts.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     timeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		})
	} else {
		encoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // Print to console and file
		level,
	)

	var log *zap.Logger
	if opts.Env != "prod" {
		log = zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	}
	return log
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//enc.AppendString(t.Format("2006-01-02 15:04:05"))
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000000"))
}

// WithValue Adds a field to the specified context
func WithValue(ctx context.Context, fields ...interface{}) context.Context {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
		c.Request = c.Request.WithContext(context.WithValue(ctx, ctxLoggerKey, WithContext(ctx).With(fields...)))
		return c
	}
	return context.WithValue(ctx, ctxLoggerKey, WithContext(ctx).With(fields...))
}

// WithContext Returns a zap instance from the specified context
func WithContext(ctx context.Context) *Logger {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	zl := ctx.Value(ctxLoggerKey)
	ctxLogger, ok := zl.(*zap.SugaredLogger)
	if ok {
		return &Logger{ctxLogger}
	}
	return logger
}

func Sync() error {
	return logger.Sync()
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Infof(template string, args ...interface{}) {
	logger.Info(template, args)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func NewOptsFromConfig() *Opts {
	return &Opts{
		FileName:   config.App.Log.LogFileName,
		Level:      config.App.Log.LogLevel,
		MaxBackUps: config.App.Log.MaxBackups,
		MaxAge:     config.App.Log.MaxAge,
		Compress:   config.App.Log.Compress,
		Encoding:   config.App.Log.Encoding,
		Env:        config.App.Env,
	}
}
