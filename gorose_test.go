package gorose

import (
	"fmt"
	"testing"
)

func TestGorose_Open(test *testing.T) {
	conn, err = InitConn()
	if err != nil {
		test.Error("FAIL: open failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: open: %v", conn.Db.MasterDb))
}
