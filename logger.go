package gorose

import (
	"fmt"
	"time"
)

// LogLevel ...
type LogLevel uint

const (
	// LOG_SQL ...
	LOG_SQL LogLevel = iota
	// LOG_SLOW ...
	LOG_SLOW
	// LOG_ERROR ...
	LOG_ERROR
)

// String ...
func (l LogLevel) String() string {
	switch l {
	case LOG_SQL:
		return "SQL"
	case LOG_SLOW:
		return "SLOW"
	case LOG_ERROR:
		return "ERROR"
	}
	return ""
}

// LogOption ...
type LogOption struct {
	FilePath     string
	EnableSqlLog bool
	// 是否记录慢查询, 默认0s, 不记录, 设置记录的时间阀值, 比如 1, 则表示超过1s的都记录
	EnableSlowLog  float64
	EnableErrorLog bool
}

// Logger ...
type Logger struct {
	filePath string
	sqlLog   bool
	slowLog  float64
	//infoLog  bool
	errLog bool
}

var _ ILogger = (*Logger)(nil)

//var onceLogger sync.Once
var logger *Logger

// NewLogger ...
func NewLogger(o *LogOption) *Logger {
	//onceLogger.Do(func() {
	logger = &Logger{filePath: "./"}
	if o.FilePath != "" {
		logger.filePath = o.FilePath
	}
	logger.sqlLog = o.EnableSqlLog
	logger.slowLog = o.EnableSlowLog
	logger.errLog = o.EnableErrorLog
	//})
	return logger
}

// DefaultLogger ...
func DefaultLogger() func(e *Engin) {
	return func(e *Engin) {
		e.logger = NewLogger(&LogOption{})
	}
}

// EnableSqlLog ...
func (l *Logger) EnableSqlLog() bool {
	return l.sqlLog
}

// EnableErrorLog ...
func (l *Logger) EnableErrorLog() bool {
	return l.errLog
}

// EnableSlowLog ...
func (l *Logger) EnableSlowLog() float64 {
	return l.slowLog
}

// Slow ...
func (l *Logger) Slow(sqlStr string, runtime time.Duration) {
	if l.EnableSlowLog() > 0 && runtime.Seconds() > l.EnableSlowLog() {
		logger.write(LOG_SLOW, "gorose_slow", sqlStr, runtime.String())
	}
}

// Sql ...
func (l *Logger) Sql(sqlStr string, runtime time.Duration) {
	if l.EnableSqlLog() {
		logger.write(LOG_SQL, "gorose_sql", sqlStr, runtime.String())
	}
}

// Error ...
func (l *Logger) Error(msg string) {
	if l.EnableErrorLog() {
		logger.write(LOG_ERROR, "gorose", msg, "0")
	}
}

func (l *Logger) write(ll LogLevel, filename string, msg string, runtime string) {
	now := time.Now()
	date := now.Format("20060102")
	datetime := now.Format("2006-01-02 15:04:05")
	f := readFile(fmt.Sprintf("%s/%v_%v.log", l.filePath, date, filename))
	content := fmt.Sprintf("[%v] [%v] %v --- %v\n", ll.String(), datetime, runtime, msg)
	withLockContext(func() {
		defer f.Close()
		buf := []byte(content)
		f.Write(buf)
	})
}
