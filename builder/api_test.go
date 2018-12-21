package builder

import (
	"fmt"
	"testing"
)

func TestConfigParser_BuildQuery(test *testing.T) {
	dr, err := BuildQuery("mysql")
	if err != nil {
		test.Error("FAIL: driver failed.", err)
		return
	}

	test.Log(fmt.Sprintf("PASS: driver: %v", dr))
}
