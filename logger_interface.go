package gorose

import (
	"time"
	"io"
)

// ILogger ...
type ILogger interface {
	Sql(sqlStr string, runtime time.Duration)
	Slow(sqlStr string, runtime time.Duration)
	Error(msg string)
	EnableSqlLog() bool
	EnableErrorLog() bool
	EnableSlowLog() float64
}


// 暂时规划
type ilogger interface {
	// Persist 持久化
	Persist(w io.Writer)
	// Info 常规日志
	Info(args ...string)
	// Error 错误日志
	Error(args ...string)
	// Debug 调试日志
	Debug(args ...string)

	Infof(format string, args ...string)
	Errorf(format string, args ...string)
	Debugf(format string, args ...string)
	InfoWithCtx(ctx interface{}, args ...string)
	ErrorWithCtx(ctx interface{}, args ...string)
	DebugWithCtx(ctx interface{}, args ...string)
	InfofWithCtxf(ctx interface{}, format string, args ...string)
	ErrorfWithCtxf(ctx interface{}, format string, args ...string)
	DebugfWithCtxf(ctx interface{}, format string, args ...string)
}