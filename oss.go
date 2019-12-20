package tool

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego"
	"strings"
)

func Size2string(size int64) (str string) {
	if size < 1<<10 {
		return Int64ToString(size) + "B"
	}
	if size < 1<<20 {
		size += 1 << 9
		return Int64ToString(size>>10) + "KB"
	}
	if size < 1<<30 {
		size += 1 << 19
		return Int64ToString(size>>20) + "MB"
	}
	size += 1 << 29
	return Int64ToString(size>>30) + "GB"
}

func OssUpload(tofile string) string {
	keyid := beego.AppConfig.String("accessKeyId")
	secret := beego.AppConfig.String("accessKeySecret")
	endpoint := beego.AppConfig.String("endpoint")
	bucketc := beego.AppConfig.String("bucket")
	if client, err := oss.New(endpoint, keyid, secret); err != nil {
		return ""
	} else if bucket, err := client.Bucket(bucketc); err != nil {
		return ""
	} else {
		obj := strings.TrimLeft(tofile, "static/")
		if err := bucket.PutObjectFromFile(obj, tofile); err != nil {
			return ""
		} else {
			return "http://" + bucketc + "." + endpoint + "/" + obj
		}
	}
}

var exts = [...]string{
	"gif", "jpg", "jpeg", "bmp", "png", "ico", "psd",
	"mp3", "wma", "wav", "amr",
	"rm", "rmvb", "wmv", "avi", "mpg", "mpeg", "mp4", "mov", "flv", "swf", "mkv", "ogg", "ogv", "webm", "mid",
	"txt", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "pps", "pdf", "chm", "md", "json", "sql",
	"rar", "zip", "7z", "tar", "gz", "bz2", "cab", "iso", "tar.gz", "mmap", "xmind", "md", "xml",
}

func IsExtLimit(ext string) bool {
	for _, s := range exts {
		if s == ext {
			return false
		}
	}
	return true
}
