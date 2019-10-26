package tool

import (
	"os/exec"
	"strings"
)

func Command(cmd string) (msg string, err error) {
	arg := []string{"/c"}
	cmd = strings.ReplaceAll(cmd, "\t", " ")
	arr := strings.Split(cmd, " ")
	for _, v := range arr {
		v = strings.TrimSpace(v)
		if v != "" {
			arg = append(arg, v)
		}
	}
	gbk, err := exec.Command("cmd", arg...).CombinedOutput()
	if len(gbk) == 0 {
		return "", err
	}
	msg = convert(string(gbk), "gbk", "utf-8")
	return msg, err
}
