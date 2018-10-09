package loger

import (
	"fmt"
	"os"
	"time"
)

type Writer struct {
	FilePath string
}

func (w *Writer) File(filepath string) *Writer {
	w.FilePath = filepath
	// 检查文件是否存在, 不存在则创建之
	return w
}

func (w *Writer) Write(sql string, runtime string) {
	f := ReadFile(w.FilePath)
	AppendFileContent(f, sql, runtime)
}

type LogData struct {
	Datetime string
	Runtime string
	Sql string
}

func AppendFileContent(f *os.File,sql, runtime string)  {
	datetime:=time.Now().Format("2006-01-02 15:04:05")
	content:=fmt.Sprintf("[Datetime: %v][Runtime: %v][Sql: %v]\n",
		datetime,runtime,sql)
	//content := &LogData{
	//	datetime,
	//	runtime,
	//	sql,
	//}
	//buf,_:=json.Marshal(content)
	buf := []byte(content)
	f.Write(buf)
	f.Close()
}

func ReadFile(filepath string) *os.File {
	file,err:=os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
	if err!=nil && os.IsNotExist(err){
		file,err = os.Create(filepath)
	}
	return file
}
