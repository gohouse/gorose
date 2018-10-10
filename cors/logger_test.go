package cors

import (
	"fmt"
	"testing"
	"time"
)

func TestFileParser_Json(test *testing.T) {

	t := time.Now()

	sql := "select * from users limit 1"

	confP := NewDefaultLogger()
	confP.Write(sql, time.Since(t).String(), t.Format("2006-01-02 15:04:05"))

	test.Log(fmt.Sprintf("PASS: sqllog %v", sql))
}

