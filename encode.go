package tool

import "github.com/axgle/mahonia"

func convert(str string, src string, tag string) string {
	str = mahonia.NewDecoder(src).ConvertString(str)
	_, bytes, _ := mahonia.NewDecoder(tag).Translate([]byte(str), true)
	return string(bytes)
}
