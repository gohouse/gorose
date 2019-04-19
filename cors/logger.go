package cors

import (
	"fmt"
	"os"
	"path"
	"time"
)

type LoggerHandler interface {
	Write(sql string, runtime string, datetime string)
}

type Logger struct {
	FilePath string
}

func (w *Logger) Write(sql string, runtime string, datetime string) {
	f := readFile(fmt.Sprintf("%s/%v.log", w.FilePath,
		time.Now().Format("2006-01-02")))

	defer f.Close()

	content := fmt.Sprintf("[%v] %v '%v'\n",
		datetime, runtime, sql)

	buf := []byte(content)
	f.Write(buf)
}

func readFile(filepath string) *os.File {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path.Dir(filepath),os.ModePerm)
		file, err = os.Create(filepath)
	}
	return file
}

func NewDefaultLogger(filePath ...string) *Logger {
	if len(filePath)>0 {
		return &Logger{FilePath:filePath[0]}
	}
	return &Logger{}
}
