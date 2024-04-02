package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// DefaultLogger logger
	DefaultLogger *KLogger
)

// 初始化日志
func init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(getLogWriter())),
		zapcore.InfoLevel,
	)

	logger := zap.New(core)
	DefaultLogger = &KLogger{
		SugaredLogger: *logger.Sugar(),
	}
}

// 获取日志写入器
func getLogWriter() zapcore.WriteSyncer {
	// 按日期格式化日志文件名
	logFileName := time.Now().Format("2006-01-02") + ".log"

	// 创建 RollingFileOutput
	writer, _ := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFileName,
			MaxSize:    100, // 每个日志文件最大尺寸，单位：MB
			MaxBackups: 10,  // 最多保留的日志文件数量
			MaxAge:     30,  // 保留日志文件的最大天数
			Compress:   true,  // 是否压缩日志文件
		}),
	)

	return writer
}

// StandardLogger provides API compatibility with standard printf loggers
// eg. go-logging
type StandardLogger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
}

// Logger retrieves an event logger by name
func logger(system string) *KLogger {
	if len(system) == 0 {
		panic("missing name")
	}

	logger := getLogger(system)
	return &KLogger{
		system:        system,
		SugaredLogger: *logger,
	}
}

// KLogger implements the StandardLogger interface
type KLogger struct {
	zap.SugaredLogger
	system string
}

// FormatRFC3339 returns the given time in UTC with RFC3999Nano format.
func FormatRFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

// Logger retrieves an event logger by name
func Logger(system string) *KLogger {
	return logger(system)
}

func (h *KLogger) WithError(err error) *KLogger {
	h.SugaredLogger = *h.SugaredLogger.With("error", err)
	return h
}

func (h *KLogger) WithFields(args ...interface{}) *KLogger {
	h.SugaredLogger = *h.SugaredLogger.With(args...)
	return h
}

func (h *KLogger) WithCallerSkip(skip int) *KLogger {
	h.SugaredLogger = *h.SugaredLogger.Desugar().WithOptions(zap.AddCallerSkip(skip)).Sugar()
	return h
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	DefaultLogger.Infof(template, args...)
}

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	DefaultLogger.Debugf(template, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	DefaultLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	DefaultLogger.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Panic(args...)
}

func Fatalf(template string, args ...interface{}) {
	DefaultLogger.Panicf(template, args...)
}
