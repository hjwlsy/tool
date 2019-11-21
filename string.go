package tool

import "strconv"

func String2Int64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func String2Int(s string) (int, error) {
	return strconv.Atoi(s)
}

func StringToInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func StringToUint(s string) uint {
	n, _ := strconv.Atoi(s)
	return uint(n)
}

func String2Float64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func StringToFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func GetBytes(bytes interface{}) []byte {
	return bytes.([]byte)
}

func Bytes2String(bytes interface{}) string {
	return string(GetBytes(bytes))
}

func String2Bytes(s string) []byte {
	return []byte(s)
}

func Bytes2Float64(bytes interface{}) (float64, error) {
	return String2Float64(Bytes2String(bytes))
}

func Bytes2Int(bytes interface{}) (int, error) {
	return String2Int(Bytes2String(bytes))
}
