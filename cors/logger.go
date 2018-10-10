package cors

import (
	"fmt"
	"os"
	"time"
)

type LoggerHandler interface {
	Write(sql string, runtime string, datetime string)
}

type Logger struct {
}

func (w *Logger) Write(sql string, runtime string, datetime string) {
	f := readFile(fmt.Sprintf("%v.log", time.Now().Format("2006-01-02")))

	defer f.Close()

	content := fmt.Sprintf("[Datetime: %v][Runtime: %v][Sql: %v]\n",
		datetime, runtime, sql)
	//fmt.Println(content)

	buf := []byte(content)
	f.Write(buf)
}

func readFile(filepath string) *os.File {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil && os.IsNotExist(err) {
		file, err = os.Create(filepath)
	}
	return file
}

func NewDefaultLogger() *Logger {
	return &Logger{}
}
