package tool

import (
	"github.com/astaxie/beego/logs"
	"strings"
	"time"
)

var TimeMap = map[string]string{
	"Y": "2006", "m": "01", "d": "02",
	"y": "06", "n": "1", "j": "2",
	"H": "15", "i": "04", "s": "05",
}

func Date(argv ...interface{}) string {
	var format string
	var timestamp int64
	var date time.Time
	argc := len(argv)
	if argc > 0 {
		format = argv[0].(string)
		if argc > 1 {
			if v, ok := argv[1].(int); ok {
				timestamp = int64(v)
			} else if v, ok := argv[1].(uint); ok {
				timestamp = int64(v)
			} else if v, ok := argv[1].(string); ok {
				timestamp = String2Int64(v)
			} else {
				timestamp = argv[1].(int64)
			}
		}
	}
	if format == "" {
		format = "Y-m-d H:i:s"
	}
	if timestamp > 0 {
		date = time.Unix(timestamp, 0)
	} else {
		date = time.Now()
	}
	for k, v := range TimeMap {
		format = strings.Replace(format, k, v, -1)
	}
	return date.Format(format)
}

func StrToTime(format string, str string) int64 {
	if format == "" {
		format = "Y-m-d H:i:s"
	}
	for k, v := range TimeMap {
		format = strings.Replace(format, k, v, -1)
	}
	if t, err := time.ParseInLocation(format, str, time.Local); err == nil {
		return t.Unix()
	} else {
		logs.Error("StrToTime转换失败" + err.Error())
		return 0
	}
}

func Time2Int() int {
	return int(time.Now().Unix())
}

func Time2Uint() uint {
	return uint(time.Now().Unix())
}

func Time2String() string {
	return Int64ToString(time.Now().Unix())
}

func Microtime() float64 {
	return float64(time.Now().UnixNano()/1e6) / 1e3
}

func GetYmd(timestamp uint) uint {
	return StringToUint(Date("Ymd", timestamp))
}

func GetTimestamp(ymd uint) uint {
	if ymd < 1 {
		ymd = GetYmd(0)
	}
	timestamp := StrToTime("Ymd", Uint2String(ymd))
	if timestamp > 0 {
		return uint(timestamp)
	} else {
		return GetTimestamp(0)
	}
}

func GetYmdBeforeDay(day uint) uint {
	return GetYmd(Time2Uint() - 86400*day)
}

func GetYmdAfterDay(day uint) uint {
	return GetYmd(Time2Uint() + 86400*day)
}

func GetYesterdayYmd() uint {
	return GetYmdBeforeDay(1)
}

func GetYesterdayTime() uint {
	return GetTimestamp(GetYesterdayYmd())
}

func GetYmdSubDay(ymd, day uint) uint {
	return GetYmd(GetTimestamp(ymd) - 86400*day)
}

func GetYmdAddDay(ymd, day uint) uint {
	return GetYmd(GetTimestamp(ymd) + 86400*day)
}
