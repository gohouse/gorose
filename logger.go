package gorose

import (
	"errors"
	"io"
	"log/slog"
)

type ILogger interface {
	Log(sqls string, bindings ...any)
	Error(err error)
	LastSql() SqlItem
}

// logger implementation by default

type SqlItem struct {
	Sql      string
	Bindings []any
	Err      error
}

type logger struct {
	*slog.Logger
	lv      slog.Level
	lastSql SqlItem
}

func DefaultLogger(lv slog.Level, writer io.Writer) ILogger {
	opts := slog.HandlerOptions{
		Level: lv,
		//AddSource: true,
	}
	return &logger{Logger: slog.New(slog.NewTextHandler(writer, &opts)), lv: lv}
}

func (l *logger) LastSql() SqlItem {
	if l.lv > slog.LevelDebug {
		return SqlItem{Err: errors.New("only record when slog level in debug mod")}
	}
	return l.lastSql
}

func (l *logger) Log(sqls string, bindings ...any) {
	l.With("bindings", bindings).Debug(sqls)
	if l.lv <= slog.LevelDebug {
		l.lastSql = SqlItem{Sql: sqls, Bindings: bindings}
	}
}

func (l *logger) Error(err error) {
	l.Logger.Error(err.Error())
}
