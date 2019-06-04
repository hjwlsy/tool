package tool

import (
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DirName(argv ...string) string {
	file := ""
	if len(argv) > 0 && argv[0] != "" {
		file = argv[0]
	} else {
		file, _ = exec.LookPath(os.Args[0])
	}
	path, _ := filepath.Abs(file)
	directory := filepath.Dir(path)
	return strings.Replace(directory, "\\", "/", -1)
}

func GetProPath() string {
	return DirName("root")
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func IsFile(f string) bool {
	if fi, err := os.Stat(f); err != nil {
		return false
	} else {
		return !fi.IsDir()
	}
}

func IsDir(f string) bool {
	if fi, err := os.Stat(f); err != nil {
		return false
	} else {
		return fi.IsDir()
	}
}

func MkDirAll(path string) {
	if err := os.MkdirAll(path, 0777); err != nil {
		logs.Error("MkDirAll失败：" + err.Error())
	}
}

func ReadFile(filename string) string {
	if !IsFile(filename) {
		logs.Error("文件不存在" + filename)
		return ""
	}
	if ret, err := ioutil.ReadFile(filename); err == nil {
		return Bytes2String(ret)
	} else {
		logs.Error("ReadFile错误" + err.Error())
		return ""
	}
}

func WriteFile(filename string, data string) {
	path := DirName(filename)
	if !IsDir(path) {
		MkDirAll(path)
	}
	if err := ioutil.WriteFile(filename, String2Bytes(data), 0777); err != nil {
		logs.Error("WriteFile错误" + err.Error())
	}
}

func AppendFile(filename string, data string) {
	if !IsFile(filename) {
		WriteFile(filename, data)
		return
	}
	if f, err := os.OpenFile(filename, os.O_WRONLY, 0644); err != nil {
		logs.Error("os.OpenFile错误" + err.Error())
	} else {
		if n, err := f.Seek(0, 2); err != nil {
			logs.Error("f.Seek错误" + err.Error())
		} else {
			if _, err = f.WriteAt(String2Bytes(data), n); err != nil {
				logs.Error("f.WriteAt错误" + err.Error())
			}
		}
		_ = f.Close()
	}
}
