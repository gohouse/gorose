package gorose

import "time"

type ILogger interface {
	Sql(sqlStr string, runtime time.Duration)
	Slow(sqlStr string, runtime time.Duration)
	Error(msg string)
	EnableSqlLog() bool
	EnableErrorLog() bool
	EnableSlowLog() float64
}
