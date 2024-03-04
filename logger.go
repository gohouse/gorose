package gorose

import (
	"log/slog"
	"os"
)

type ILogger interface {
	Log(sqls string, bindings ...any)
	Error(err error)
	LastSql() SqlItem
}

// logger implementation by default

type SqlItem struct {
	Sqls     string
	Bindings []any
}

type logger struct {
	*slog.Logger
	lastSql SqlItem
}

func DefaultLogger(lv slog.Level) *logger {
	opts := slog.HandlerOptions{
		Level: lv,
		//AddSource: true,
	}
	return &logger{Logger: slog.New(slog.NewTextHandler(os.Stdout, &opts))}
}

func (l *logger) LastSql() SqlItem {
	return l.lastSql
}

func (l *logger) Log(sqls string, bindings ...any) {
	l.lastSql = SqlItem{Sqls: sqls, Bindings: bindings}
	l.With("bindings", bindings).Debug(sqls)
}

func (l *logger) Error(err error) {
	l.Logger.Error(err.Error())
}
