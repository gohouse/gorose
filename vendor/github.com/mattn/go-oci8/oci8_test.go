package oci8

import (
	"reflect"
	"testing"
	"time"
)

func TestParseDSN(t *testing.T) {
	var (
		pacific *time.Location
		err     error
	)

	if pacific, err = time.LoadLocation("America/Los_Angeles"); err != nil {
		panic(err)
	}
	var dsnTests = []struct {
		dsnString   string
		expectedDSN *DSN
	}{
		{"oracle://xxmc:xxmc@107.20.30.169:1521/ORCL?loc=America%2FLos_Angeles", &DSN{Username: "xxmc", Password: "xxmc", Connect: "107.20.30.169:1521/ORCL", prefetch_rows: 10, Location: pacific}},
		{"xxmc/xxmc@107.20.30.169:1521/ORCL?loc=America%2FLos_Angeles", &DSN{Username: "xxmc", Password: "xxmc", Connect: "107.20.30.169:1521/ORCL", prefetch_rows: 10, Location: pacific}},
		{"sys/syspwd@107.20.30.169:1521/ORCL?loc=America%2FLos_Angeles&as=sysdba", &DSN{Username: "sys", Password: "syspwd", Connect: "107.20.30.169:1521/ORCL", prefetch_rows: 10, Location: pacific, operationMode: 0x00000002}}, // with operationMode: 0x00000002 = C.OCI_SYDBA
		{"xxmc/xxmc@107.20.30.169:1521/ORCL", &DSN{Username: "xxmc", Password: "xxmc", Connect: "107.20.30.169:1521/ORCL", prefetch_rows: 10, Location: time.Local}},
		{"xxmc/xxmc@107.20.30.169/ORCL", &DSN{Username: "xxmc", Password: "xxmc", Connect: "107.20.30.169/ORCL", prefetch_rows: 10, Location: time.Local}},
	}

	for _, tt := range dsnTests {
		actualDSN, err := ParseDSN(tt.dsnString)

		if err != nil {
			t.Errorf("ParseDSN(%s) got error: %+v", tt.dsnString, err)
		}

		if !reflect.DeepEqual(actualDSN, tt.expectedDSN) {
			t.Errorf("ParseDSN(%s): expected %+v, actual %+v", tt.dsnString, tt.expectedDSN, actualDSN)
		}
	}
}

func TestIsBadConn(t *testing.T) {
	var errorCode = "ORA-03114"
	if !isBadConnection(errorCode) {
		t.Errorf("TestIsBadConn: expected %+v, actual %+v", true, isBadConnection(errorCode))
	}
}
