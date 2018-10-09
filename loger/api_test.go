package loger

import (
	"fmt"
	"testing"
	"time"
)

func TestFileParser_Json(test *testing.T) {

	t := time.Now()
	filepath := "sql.log"
	sql := "select ..."
	fmt.Println(time.Since(t).String())

	var confP = &Writer{}
	confP.File(filepath).Write(sql,time.Since(t).String())

	//if err != nil {
	//	test.Error("FAIL: json parser failed.", err)
	//	return
	//}

	test.Log(fmt.Sprintf("PASS: sqllog %v", sql))
}

