package drivers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func BackQuotes(arg any) string {
	var tmp []string
	if v, ok := arg.(string); ok {
		split := strings.Split(v, " ")
		split2 := strings.Split(split[0], ".")
		if len(split2) > 1 {
			if split2[1] == "*" {
				tmp = append(tmp, fmt.Sprintf("`%s`.%s", split2[0], split2[1]))
			} else {
				tmp = append(tmp, fmt.Sprintf("`%s`.`%s`", split2[0], split2[1]))
			}
		} else {
			tmp = append(tmp, fmt.Sprintf("`%s`", split2[len(split2)-1]))
		}
		tmp = append(tmp, split[1:]...)
	}
	return strings.Join(tmp, " ")
}

func ToSlice(arg any) []any {
	ref := reflect.Indirect(reflect.ValueOf(arg))
	var res []any
	switch ref.Kind() {
	case reflect.Slice:
		l := ref.Len()
		v := ref.Slice(0, l)
		for i := 0; i < l; i++ {
			res = append(res, v.Index(i).Interface())
		}
	default:
		res = append(res, ref.Interface())
	}
	return res
}

func NamedSprintf(format string, a ...any) string {
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(fmt.Sprintf(regexp.MustCompile(`:\w+`).ReplaceAllString(format, "%s"), a...), " "))
}

func TrimPrefixAndOr(s string) string {
	return regexp.MustCompile(`(?i)^\s*(and|or)\s+`).ReplaceAllString(s, "")
}

func assertsEqual(t *testing.T, expect, real any) {
	marshal, err := json.Marshal(expect)
	assertsError(t, err)
	bytes, err := json.Marshal(real)
	assertsError(t, err)
	if string(marshal) != string(bytes) {
		methodName, file, line := getCallerInfo()
		t.Errorf("[%s] Error\n\t Trace - %s:%v\n\tExpect - %T %s\n\t   Got - %T %s\n------------------------------------------------------", methodName, file, line, expect, marshal, real, bytes)
	}
}

func assertsError(t *testing.T, err error) {
	if err != nil {
		methodName, file, line := getCallerInfo()
		t.Errorf("[%s] Error\n\t Trace - %s:%v\n\t%s\n------------------------------------------------------", methodName, file, line, err.Error())
	}
}

func getCallerInfo() (string, string, int) {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var i int
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		i++
		if i == 1 {
			continue
		}
		//fmt.Printf("Method: %s\nFile: %s\nLine: %d\n\n", frame.Function, frame.File, frame.Line)
		lastDotIndex := strings.LastIndex(frame.Function, ".")
		methodName := frame.Function[lastDotIndex+1:]
		//t.Logf("[%s] errors on file:line: \n\t\t -> %s:%v\n", methodName, frame.File, frame.Line)
		if i == 2 {
			return methodName, frame.File, frame.Line
		}
		//break
	}
	return "", "", 0
}

func jsonLog(t *testing.T, data any) {
	marshal, err := json.Marshal(data)
	assertsError(t, err)
	t.Logf("json data: %s", marshal)
}

func Map[Data any, Datas ~[]Data, Result any](datas Datas, mapper func(Data) Result) []Result {
	results := make([]Result, 0, len(datas))
	for _, data := range datas {
		results = append(results, mapper(data))
	}
	return results
}
