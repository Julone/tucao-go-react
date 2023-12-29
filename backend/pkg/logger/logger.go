package logger

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path"
	"time"
)

/*
Logger exposes a logging framework to use in modules. It exposes level-specific logging functions and a set of common functions for compatibility.
*/
type Logger interface {
	IsDebug() bool
	SetLevel(level string)
	Debug(args ...interface{})
	Debugf(format string, v ...interface{})
	Info(args ...interface{})
	Infof(format string, v ...interface{})
	Warn(args ...interface{})
	Warnf(format string, v ...interface{})
	Error(args ...interface{})
	Errorf(format string, v ...interface{})
	Panic(args ...interface{})
	Panicf(format string, v ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, v ...interface{})
	WithField(key string, v interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	Sync()
}

type Option = zap.Option

var (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
)

var (
	WrapCore      = zap.WrapCore
	Hooks         = zap.Hooks
	Fields        = zap.Fields
	Development   = zap.Development
	AddCaller     = zap.AddCaller
	WithCaller    = zap.WithCaller
	AddCallerSkip = zap.AddCallerSkip
	AddStacktrace = zap.AddStacktrace
	IncreaseLevel = zap.IncreaseLevel
	WithFatalHook = zap.WithFatalHook
)

type localLogger struct {
	logger *zap.SugaredLogger
	lvl    *zap.AtomicLevel
}

func (l *localLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *localLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

func (l *localLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *localLogger) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

func (l *localLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *localLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warnf(format, v...)
}

func (l *localLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *localLogger) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

func (l *localLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *localLogger) Panicf(format string, v ...interface{}) {
	l.logger.Panicf(format, v...)
}

func (l *localLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *localLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
}

func (l *localLogger) Sync() {
	_ = l.logger.Sync()
}

func (l *localLogger) WithField(key string, v interface{}) Logger {
	return &localLogger{
		logger: l.logger.With(key, v),
	}
}

func (l *localLogger) WithFields(fields map[string]interface{}) Logger {
	f := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		f = append(f, k, v)
	}
	return &localLogger{
		logger: l.logger.With(f...),
	}
}

// SetLevel level contains debug,info,warn,error,dpanic,panic,fatal
func (l *localLogger) SetLevel(level string) {
	targetLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		targetLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	if l.lvl == nil {
		l.lvl = &targetLevel
	}
	l.lvl.SetLevel(targetLevel.Level())
}

func (l *localLogger) IsDebug() bool {
	return l.lvl.Level() == -1
}

// Replace the default Log instance
func Replace(new Logger) {
	Log = new
}

// Log default logger instance, will auto init
var Log = New("./", "rmp", "debug", 60*24*time.Hour, 24*time.Hour, 100, zap.AddCaller())

// New a Logger instance, output to stdout when logName is empty, maxSize(MB)
func New(logPath, logName, logLevel string, maxAge, rotate time.Duration, maxSize int64, opts ...Option) Logger {
	log := &localLogger{}
	log.SetLevel(logLevel)
	var ws zapcore.WriteSyncer
	if logName == "" {
		ws = zapcore.AddSync(os.Stdout)
	} else {
		writer := getWriter(logPath, logName, maxAge, rotate, maxSize)
		if writer != nil {
			ws = zapcore.AddSync(writer)
		} else {
			file, err := os.OpenFile(logPath+logName+"."+time.Now().Format(time.DateOnly)+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			ws = zapcore.AddSync(file)
		}
	}

	encodeCfg := zapcore.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		TimeKey:       "time",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder, //级别使用大写
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encodeCfg), ws, log.lvl)
	log.logger = zap.New(core, opts...).Sugar().WithOptions(zap.AddCallerSkip(1))
	return log
}

func getWriter(logPath, filename string, maxAge, rotate time.Duration, maxsize int64) io.Writer {
	var err error
	if len(logPath) == 0 {
		logPath, err = os.Getwd()
		if err != nil {
			return nil
		}
	}
	var baseLogPath string
	if len(filename) == 0 {
		baseLogPath = path.Join(logPath, "log")
	} else {
		baseLogPath = path.Join(logPath, filename)
	}
	options := []rotatelogs.Option{rotatelogs.WithLinkName(baseLogPath + "_symlink"), rotatelogs.WithMaxAge(maxAge), rotatelogs.WithRotationTime(rotate), rotatelogs.WithRotationSize(maxsize * 1024 * 1024)}
	hook, err := rotatelogs.New(baseLogPath+"."+time.Now().Format(time.DateOnly)+".log", options...)
	if err != nil {
		return nil
	}
	return hook
}
