package tool

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/skip2/go-qrcode"
	"math/rand"
	"net"
	"regexp"
	"strings"
	"time"
)

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write(String2Bytes(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func MyMD5(str string) string {
	bytes := []byte(str)
	length := len(bytes)
	mod := length & 0x3f
	mod = 56 - mod
	if mod <= 0 {
		mod += 64
	}
	patch := make([]byte, mod+8)
	patch[0] = 0x80
	length <<= 3
	for i := 0; i < 8; i++ {
		patch[mod+i] = byte(length & 0xff)
		length >>= 8
		if length <= 0 {
			break
		}
	}
	bytes = append(bytes, patch...)
	n := len(bytes) >> 6
	ti := [64]uint32{
		3614090360, 3905402710, 606105819, 3250441966, 4118548399, 1200080426, 2821735955, 4249261313,
		1770035416, 2336552879, 4294925233, 2304563134, 1804603682, 4254626195, 2792965006, 1236535329,
		4129170786, 3225465664, 643717713, 3921069994, 3593408605, 38016083, 3634488961, 3889429448,
		568446438, 3275163606, 4107603335, 1163531501, 2850285829, 4243563512, 1735328473, 2368359562,
		4294588738, 2272392833, 1839030562, 4259657740, 2763975236, 1272893353, 4139469664, 3200236656,
		681279174, 3936430074, 3572445317, 76029189, 3654602809, 3873151461, 530742520, 3299628645,
		4096336452, 1126891415, 2878612391, 4237533241, 1700485571, 2399980690, 4293915773, 2240044497,
		1873313359, 4264355552, 2734768916, 1309151649, 4149444226, 3174756917, 718787259, 3951481745,
	} //int((1 << 32) * math.Abs(math.Sin(float64(i)))) 1<=i<=64
	s := [64]uint8{
		7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
		5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
		4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
		6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
	}
	var A, B, C, D uint32
	A, B, C, D = 0x67452301, 0xefcdab89, 0x98badcfe, 0x10325476
	for i := 0; i < n; i++ {
		M := bytes[i<<6 : (i+1)<<6]
		m := make([]uint32, 16)
		for j := 0; j < 16; j++ {
			idx := j << 2
			m[j] = uint32(M[idx]) + uint32(M[idx+1])<<8 + uint32(M[idx+2])<<16 + uint32(M[idx+3])<<24
		}
		a, b, c, d := A, B, C, D
		for j := 0; j < 64; j++ {
			var fun, idx uint32
			if j < 16 {
				fun, idx = (b&c)|((^b)&d), uint32(j)
			} else {
				if j < 32 {
					fun, idx = (b&d)|(c&(^d)), uint32((1+5*j)%16)
				} else {
					if j < 48 {
						fun, idx = b^c^d, uint32((5+3*j)%16)
					} else {
						fun, idx = c^(b|(^d)), uint32((7*j)%16)
					}
				}
			}
			a = a + fun + ti[j] + m[idx]
			a = a<<s[j] + a>>(32-s[j]) + b
			a, b, c, d = d, a, b, c
		}
		A, B, C, D = A+a, B+b, C+c, D+d
	}
	temp := make([]uint8, 16)
	for i := 0; i < 4; i++ {
		temp[i] = uint8((A >> (i * 8)) & 0xff)
		temp[i+4] = uint8((B >> (i * 8)) & 0xff)
		temp[i+8] = uint8((C >> (i * 8)) & 0xff)
		temp[i+12] = uint8((D >> (i * 8)) & 0xff)
	}
	result := make([]byte, 32)
	h := [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	for k, v := range temp {
		result[k*2+1] = h[v&0x0f]
		result[k*2] = h[v&0xf0>>4]
	}
	return string(result)
}

func JsonDecode(s string, v interface{}) error {
	return json.Unmarshal(String2Bytes(s), v)
}

func JsonEncode(v interface{}) string {
	if bytes, err := json.Marshal(v); err == nil {
		return Bytes2String(bytes)
	} else {
		logs.Error("JsonEncode转换失败" + err.Error())
	}
	return ""
}

func QRCode(s string) string {
	smd5 := MD5(s)
	path := "static/qrcode/" + smd5[:2] + "/"
	MkDirAll(path)
	file := path + smd5[2:] + ".png"
	if !IsFile(file) {
		if err := qrcode.WriteFile(s, qrcode.Medium, 256, file); err != nil {
			logs.Error("生成二维码失败：" + err.Error())
		}
	}
	if IsFile(file) {
		return "/" + file
	}
	return ""
}

func Base64Encode(src string) string {
	return base64.StdEncoding.EncodeToString(String2Bytes(src))
}

func Base64Decode(src string) string {
	if ret, err := base64.StdEncoding.DecodeString(src); err == nil {
		return Bytes2String(ret)
	} else {
		logs.Error("base64Decode错误" + err.Error())
		return ""
	}
}

func IsMobile(useragent string) bool {
	useragent = strings.ToLower(useragent)
	keywords := []string{"android", "iphone", "ipod", "ipad", "windows phone", "mqqbrowser"}
	for _, v := range keywords {
		if strings.Contains(useragent, v) {
			return true
		}
	}
	pattern1 := `/(android|bb\d+|meego).+mobile|avantgo|bada/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|iris|kindle|lge |maemo|midp|mmp|mobile.+firefox|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i`
	if matched1, err := regexp.MatchString(pattern1, useragent); err == nil && matched1 {
		return true
	} else {
		pattern2 := `/1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i`
		if matched2, err := regexp.MatchString(pattern2, useragent[:4]); err == nil && matched2 {
			return true
		}
	}
	return false
}

func GetWid() uint {
	wid, err := beego.AppConfig.Int("wid")
	if err != nil || wid < 1 {
		panic("网站ID错误")
	}
	return uint(wid)
}

func IsMobileNumber(mobile string) bool {
	if ok, err := regexp.MatchString(`^1\d{10}$`, mobile); err == nil && ok {
		return true
	}
	return false
}

func FindInSet(s, sl string) bool {
	ss := strings.Split(sl, ",")
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func IsIP(ip string) bool {
	if i := strings.LastIndex(ip, ":"); i >= 0 {
		ip = ip[:i]
	}
	pattern := `^(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)$`
	if ok, err := regexp.MatchString(pattern, ip); err == nil && ok {
		return true
	}
	return false
}

func InArray(needle interface{}, haystack interface{}) (bool, error) {
	E1 := errors.New("Err:Type mismatch")
	E2 := errors.New("Err:Type unsupported")
	switch needle.(type) {
	case int:
		if list, ok := haystack.([]int); ok {
			nv := needle.(int)
			for _, v := range list {
				if v == nv {
					return true, nil
				}
			}
			return false, nil
		} else {
			return false, E1
		}
	case string:
		if list, ok := haystack.([]string); ok {
			nv := needle.(string)
			for _, v := range list {
				if v == nv {
					return true, nil
				}
			}
			return false, nil
		} else {
			return false, E1
		}
	default:
		return false, E2
	}
}

func GetMacList() (maclist string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, v := range interfaces {
		if v.Flags&net.FlagUp == 0 {
			continue
		}
		if v.Flags&net.FlagLoopback != 0 {
			continue
		}
		mac := strings.ToUpper(v.HardwareAddr.String())
		mac = strings.ReplaceAll(mac, ":", "")
		if len(mac) != 12 {
			continue
		}
		if maclist == "" {
			maclist = mac
		} else {
			maclist += "," + mac
		}
	}
	return
}

func Random(n int) (s string) {
	const l = "0123456789" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz"
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		idx := rand.Intn(62)
		s += string(l[idx])
	}
	return
}
