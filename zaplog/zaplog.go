package zaplog

import (
	"fmt"
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createHook(filePath string, maxAgeDay int) (*rotatelogs.RotateLogs, error) {
	filePathP := fmt.Sprintf("%s.%%Y%%m%%d", filePath)
	logf, err := rotatelogs.New(
		filePathP,
		rotatelogs.WithLinkName(filePath),
		rotatelogs.WithMaxAge(time.Duration(maxAgeDay)*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		log.Printf("failed to create rotatelogs:%s, %s\n", filePath, err)
		return nil, errors.Wrapf(err, "create rotatelogs:%s", filePath)
	}
	return logf, nil
}

func cutLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	lvlStr := l.CapitalString()
	if len(lvlStr) <= 6 {
		lvlStr = lvlStr[:3]
	}
	enc.AppendString(lvlStr)
}

type loggerOptions struct {
	errFilePath   string
	fileMaxAgeDay int
	stdout        bool
}

type LoggerOption func(*loggerOptions)

func WithErrFilePath(errFilePath string) LoggerOption {
	return func(opts *loggerOptions) {
		opts.errFilePath = errFilePath
	}
}
func WithFileMaxAgeDay(day int) LoggerOption {
	return func(opts *loggerOptions) {
		opts.fileMaxAgeDay = day
	}
}

func WithStdout(stdout bool) LoggerOption {
	return func(opts *loggerOptions) {
		opts.stdout = stdout
	}
}

var defaultLoggerOptions = loggerOptions{"", 10, true}

func NewLogger(logFilePath string, opts ...LoggerOption) *zap.Logger {
	options := defaultLoggerOptions
	for _, o := range opts {
		o(&options)
	}
	errFilePath := options.errFilePath
	fileMaxAgeDay := options.fileMaxAgeDay
	stdout := options.stdout

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "t",
		LevelKey:         "lv",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		ConsoleSeparator: " ",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      cutLevelEncoder,
		EncodeTime:       zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05.000]"),
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	var logHook *rotatelogs.RotateLogs
	if logFilePath != "" {
		logHook1, err := createHook(logFilePath, fileMaxAgeDay)
		if err != nil {
			log.Printf("Create Hook Error:%s, %s", logFilePath, err)
		}
		logHook = logHook1
	}

	logAtomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)

	var logWriteSyncer zapcore.WriteSyncer

	// logHook == nil的情况下, 强制stdout
	if logHook != nil {
		if stdout {
			logWriteSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(logHook))
		} else {
			logWriteSyncer = zapcore.AddSync(logHook)
		}
	} else {
		logWriteSyncer = zapcore.AddSync(os.Stdout)
	}

	logCore := zapcore.NewCore(encoder, logWriteSyncer, logAtomicLevel)

	var errCore zapcore.Core
	if errFilePath != "" {
		errFileHook, err := createHook(errFilePath, fileMaxAgeDay)
		if err != nil {
			log.Printf("Create Hook Error:%s, %s", errFilePath, err)
		}

		errAtomicLevel := zap.NewAtomicLevelAt(zap.ErrorLevel)

		if errFileHook != nil {
			errCore = zapcore.NewCore(encoder, zapcore.AddSync(errFileHook), errAtomicLevel)
		}
	}
	var coreTee zapcore.Core
	if errCore != nil {
		coreTee = zapcore.NewTee(logCore, errCore)
	} else {
		coreTee = logCore
	}

	logger := zap.New(coreTee)
	return logger
}

var Logger = NewLogger("")

func Info(logger *zap.Logger, format string, v ...interface{}) {
	logger.Info(fmt.Sprintf(format, v...))
}

func Error(logger *zap.Logger, format string, v ...interface{}) {
	logger.Error(fmt.Sprintf(format, v...))
}
