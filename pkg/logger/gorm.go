package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

const GormLoggerCallerSkip = 2

type GormLogger struct {
	ZapLogger     *zap.Logger
	slowThreshold time.Duration
}

// LogMode 实现 gorm logger 接口方法
func (g GormLogger) LogMode(gormLogLevel gormlogger.LogLevel) gormlogger.Interface {
	newlogger := g
	return &newlogger
}

// Info 实现 gorm logger 接口方法
func (g GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	fmtMsg := fmt.Sprintf(msg, data...)
	var fields []zap.Field
	allFields := addTraceFields(ctx, fields...)
	g.ZapLogger.Info(fmtMsg, allFields...)
}

// Warn 实现 gorm logger 接口方法
func (g GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	fmtMsg := fmt.Sprintf(msg, data...)
	var fields []zap.Field
	allFields := addTraceFields(ctx, fields...)
	g.ZapLogger.Warn(fmtMsg, allFields...)
}

// Error 实现 gorm logger 接口方法
func (g GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	fmtMsg := fmt.Sprintf(msg, data...)
	var fields []zap.Field
	allFields := addTraceFields(ctx, fields...)
	g.ZapLogger.Error(fmtMsg, allFields...)
}

// Trace 实现 gorm logger 接口方法
func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	var fields []zap.Field
	fields = append(fields, zap.Duration("elapsed", elapsed))
	fields = append(fields, zap.Int64("rows", rows))
	fields = append(fields, zap.String("sql", sql))

	switch {
	case err != nil:
		fields = append(fields, zap.String("error", err.Error()))
		allFields := addTraceFields(ctx, fields...)
		g.ZapLogger.Error("SqlErrorLog", allFields...)
	case g.slowThreshold != 0 && elapsed > g.slowThreshold:
		allFields := addTraceFields(ctx, fields...)
		g.ZapLogger.Warn("SqlSlowLog", allFields...)
	default:
		allFields := addTraceFields(ctx, fields...)
		g.ZapLogger.Info("SqlInfoLog", allFields...)
	}
}

func NewGormLogger() GormLogger {
	return GormLogger{
		ZapLogger:     logger.WithOptions(zap.AddCallerSkip(GormLoggerCallerSkip)),
		slowThreshold: 200 * time.Millisecond,
	}
}
