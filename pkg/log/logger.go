// Package log provides context-aware and structured logging capabilities.
package logger

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/cakramediadata2022/chs_cloud_general/config"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gopkg.in/natefinch/lumberjack.v2"
	"moul.io/zapgorm2"
)

var jsonEncoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
var textEncoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
var consoleLogger = zapcore.NewCore(
	textEncoder,
	zapcore.AddSync(os.Stdout),
	zapcore.InfoLevel,
)

func zapOpts() []zap.Option {
	return []zap.Option{zap.AddCaller()}
}

func loggers(config config.Logger) []zapcore.Core {
	cores := make([]zapcore.Core, 0)
	cores = append(cores, consoleLogger)
	if config.LogFileEnabled {
		cores = append(cores, fileLogger(config))
	}
	return cores
}

// Logger is a logger that supports log levels, context and structured logging.
type Logger interface {
	// With returns a logger based off the root logger and decorates it with the given context and arguments.
	With(ctx context.Context, args ...interface{}) Logger

	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Error(args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	Debugf(format string, args ...interface{})
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...interface{})
}

type logger struct {
	*zap.SugaredLogger
}

type contextKey int

const (
	requestIDKey contextKey = iota
	correlationIDKey
)

func fileLogger(config config.Logger) zapcore.Core {
	return zapcore.NewCore(
		jsonEncoder,
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.LogFilename,
			MaxSize:    config.LogMaxSize,
			MaxBackups: config.LogMaxBackups,
			MaxAge:     config.LogMaxAge,
		}),
		zapcore.InfoLevel)
}

func logLevel(config config.Logger) zapcore.Level {
	if level, err := zapcore.ParseLevel(config.Level); err != nil {
		return zapcore.InfoLevel
	} else {
		return level
	}
}

func Init(config config.Logger) Logger {
	return NewWithZap(zap.New(zapcore.NewTee(loggers(config)...), zapOpts()...))
}

func InitOtelzap() *otelzap.Logger {
	log := otelzap.New(zap.NewExample())
	return log
}

// New creates a new logger using the default configuration.
func New() Logger {
	l, _ := zap.NewProduction()
	return NewWithZap(l)
}

// NewWithZap creates a new logger using the preconfigured zap logger.
func NewWithZap(l *zap.Logger) Logger {
	return &logger{l.Sugar()}
}

// NewForTest returns a new logger and the corresponding observed logs which can be used in unit tests to verify log entries.
func NewForTest() (Logger, *observer.ObservedLogs) {
	core, recorded := observer.New(zapcore.InfoLevel)
	return NewWithZap(zap.New(core)), recorded
}

// With returns a logger based off the root logger and decorates it with the given context and arguments.
//
// If the context contains request ID and/or correlation ID information (recorded via WithRequestID()
// and WithCorrelationID()), they will be added to every log message generated by the new logger.
//
// The arguments should be specified as a sequence of name, value pairs with names being strings.
// The arguments will also be added to every log message generated by the logger.
func (l *logger) With(ctx context.Context, args ...interface{}) Logger {
	if ctx != nil {
		if id, ok := ctx.Value(requestIDKey).(string); ok {
			args = append(args, zap.String("request_id", id))
		}
		if id, ok := ctx.Value(correlationIDKey).(string); ok {
			args = append(args, zap.String("correlation_id", id))
		}
	}
	if len(args) > 0 {
		return &logger{l.SugaredLogger.With(args...)}
	}
	return l
}

// WithRequest returns a context which knows the request ID and correlation ID in the given request.
func WithRequest(ctx context.Context, req *http.Request) context.Context {
	id := getRequestID(req)
	if id == "" {
		id = uuid.New().String()
	}
	ctx = context.WithValue(ctx, requestIDKey, id)
	if id := getCorrelationID(req); id != "" {
		ctx = context.WithValue(ctx, correlationIDKey, id)
	}
	return ctx
}

// getCorrelationID extracts the correlation ID from the HTTP request
func getCorrelationID(req *http.Request) string {
	return req.Header.Get("X-Correlation-ID")
}

// getRequestID extracts the correlation ID from the HTTP request
func getRequestID(req *http.Request) string {
	return req.Header.Get("X-Request-ID")
}

func GinzapWithConfig(config config.Logger) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(zap.New(zapcore.NewTee(loggers(config)...), zapOpts()...), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/health"},
	})
}

func ZapGorm(config config.Logger) zapgorm2.Logger {
	return zapgorm2.New(zap.New(zapcore.NewTee(loggers(config)...)))
}
