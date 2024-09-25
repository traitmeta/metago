package bzap

import (
	"context"
	"fmt"

	"github.com/traitmeta/metago/pkg/blog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	customMsg = "custom"
)

type ctxMarker struct{}

type CtxLogger struct {
	logger *zap.Logger
	fields []zapcore.Field
	ctx    context.Context
}

var (
	ctxMarkedKey = &ctxMarker{}
)

func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	l := &CtxLogger{
		logger: logger,
	}

	return context.WithValue(ctx, ctxMarkedKey, l)
}

func WithContext(ctx context.Context) *CtxLogger {
	l, ok := ctx.Value(ctxMarkedKey).(*CtxLogger)
	if !ok || l == nil {
		return NewContextLogger(ctx)
	}

	l.ctx = ctx

	return l
}

func NewContextLogger(ctx context.Context) *CtxLogger {
	return &CtxLogger{
		logger: GetZapLogger(),
		ctx:    ctx,
	}
}

func tagsToFields(ctx context.Context) []zapcore.Field {
	var fields []zapcore.Field
	tags := blog.Extract(ctx)
	for k, v := range tags.Values() {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}

func (l *CtxLogger) Extract() *zap.Logger {
	fields := tagsToFields(l.ctx)
	fields = append(fields, l.fields...)
	return l.logger.With(fields...)
}

func (l *CtxLogger) WithField(fields ...zap.Field) {
	l.fields = append(l.fields, fields...)
}

func (l *CtxLogger) Debug(msg string, fields ...zap.Field) {
	l.Extract().WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

func (l *CtxLogger) Info(msg string, fields ...zap.Field) {
	l.Extract().WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

func (l *CtxLogger) Warn(msg string, fields ...zap.Field) {
	l.Extract().WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

func (l *CtxLogger) Error(msg string, fields ...zap.Field) {
	l.Extract().WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

func (l *CtxLogger) Panic(msg string, fields ...zap.Field) {
	l.Extract().WithOptions(zap.AddCallerSkip(1)).Panic(msg, fields...)
}

func (l *CtxLogger) Debugf(format string, data ...interface{}) {
	l.Debug(customMsg, zap.String("content", fmt.Sprintf(format, data...)))
}

func (l *CtxLogger) Infof(format string, data ...interface{}) {
	l.Info(customMsg, zap.String("content", fmt.Sprintf(format, data...)))
}

func (l *CtxLogger) Warnf(format string, data ...interface{}) {
	l.Warn(customMsg, zap.String("content", fmt.Sprintf(format, data...)))
}

func (l *CtxLogger) Errorf(format string, data ...interface{}) {
	l.Error(customMsg, zap.String("content", fmt.Sprintf(format, data...)))
}

func (l *CtxLogger) Panicf(format string, data ...interface{}) {
	l.Error(customMsg, zap.String("content", fmt.Sprintf(format, data...)))
}
