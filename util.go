package gorose

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
func GetRandomInt(num int) int {
	return rand.Intn(num)
}
func GetRandomWeightedIndex(weights []int) int {
	if len(weights) == 0 {
		return 0
	}
	if len(weights) == 1 {
		return 0
	}
	totalWeight := 0
	for _, w := range weights {
		totalWeight += w
	}
	if totalWeight == 0 {
		return rand.Intn(len(weights))
	}

	rnd := rand.Intn(totalWeight)

	currentWeight := 0
	for i, w := range weights {
		currentWeight += w
		if rnd < currentWeight {
			return i
		}
	}
	return -1 // 如果权重都为 0，或者总权重为 0，则返回 -1
}

//////////// struct field ptr 4 orm helpers ////////////

// PtrBool helper
func PtrBool(arg bool) *bool {
	return &arg
}

// PtrString helper
func PtrString(arg string) *string {
	return &arg
}

// PtrInt helper
func PtrInt(arg int) *int {
	return &arg
}

// PtrInt8 helper
func PtrInt8(arg int8) *int8 {
	return &arg
}

// PtrInt16 helper
func PtrInt16(arg int16) *int16 {
	return &arg
}

// PtrInt64 helper
func PtrInt64(arg int64) *int64 {
	return &arg
}

// PtrFloat64 helper
func PtrFloat64(arg float64) *float64 {
	return &arg
}

// PtrTime helper
func PtrTime(arg time.Time) *time.Time {
	return &arg
}

//////////// sql.Null* type helpers ////////////

// NullInt64From helper
func NullInt64From(arg int64) sql.NullInt64 { return sql.NullInt64{Int64: arg, Valid: true} }

// NullInt32From helper
func NullInt32From(arg int32) sql.NullInt32 { return sql.NullInt32{Int32: arg, Valid: true} }

// NullInt16From helper
func NullInt16From(arg int16) sql.NullInt16 { return sql.NullInt16{Int16: arg, Valid: true} }

// NullByteFrom helper
func NullByteFrom(arg byte) sql.NullByte { return sql.NullByte{Byte: arg, Valid: true} }

// NullFloat64From helper
func NullFloat64From(arg float64) sql.NullFloat64 { return sql.NullFloat64{Float64: arg, Valid: true} }

// NullBoolFrom helper
func NullBoolFrom(arg bool) sql.NullBool { return sql.NullBool{Bool: arg, Valid: true} }

// NullTimeFrom helper
func NullTimeFrom(arg time.Time) sql.NullTime { return sql.NullTime{Time: arg, Valid: true} }

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
func ToSliceAddressable(arg any) []any {
	ref := reflect.Indirect(reflect.ValueOf(arg))
	var res []any
	switch ref.Kind() {
	case reflect.Slice:
		l := ref.Len()
		v := ref.Slice(0, l)
		for i := 0; i < l; i++ {
			res = append(res, v.Index(i).Addr().Interface())
		}
	default:
		res = append(res, ref.Addr().Interface())
	}
	return res
}
func SliceContains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
func Map[Data any, Datas ~[]Data, Result any](datas Datas, mapper func(Data) Result) []Result {
	results := make([]Result, 0, len(datas))
	for _, data := range datas {
		results = append(results, mapper(data))
	}
	return results
}

func NamedSprintf(format string, a ...any) string {
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(fmt.Sprintf(regexp.MustCompile(`:\w+`).ReplaceAllString(format, "%s"), a...), " "))
}

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

func SortedMapKeys(data any) (cols []string) {
	// 从 map 中获取所有的键，并转换为切片
	keys := reflect.ValueOf(data).MapKeys()

	// 对切片进行排序
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	// 输出排序后的结果
	for _, key := range keys {
		cols = append(cols, key.String())
	}
	return
}

func IsExpression(obj any) (b bool) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	if rfv.Kind() == reflect.String && strings.Contains(rfv.String(), "?") {
		b = true
	}
	return
}
