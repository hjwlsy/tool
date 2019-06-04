package tool

import (
	"fmt"
	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/astaxie/beego/logs"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"sort"
	"time"
)

func ImgResize(filename string, size int) (newfn string) {
	if !IsFile(filename) {
		logs.Error("文件不存在" + filename)
		return
	}
	if reader, err := os.Open(filename); err == nil {
		if m, _, err := image.Decode(reader); err == nil {
			bounds := m.Bounds()
			w, h := bounds.Dx(), bounds.Dy()
			if w > size || h > size {
				newx, newy := 0, 0
				if w >= h {
					newx, newy = size, size*h/w
				} else {
					newx, newy = size*w/h, size
				}
				dst := image.NewRGBA(image.Rect(0, 0, newx, newy))
				if err = graphics.Scale(dst, m); err == nil {
					path := GetProPath() + "/static/resize/" + Int2String(size) + "/"
					MkDirAll(path)
					file := path + Int64ToString(time.Now().UnixNano()) + ".jpg"
					if dstreader, err := os.Create(file); err == nil {
						if err = jpeg.Encode(dstreader, dst, nil); err == nil {
							newfn = file
						} else {
							logs.Error("jpeg.Encode错误" + err.Error())
						}
					} else {
						logs.Error("os.Create错误" + err.Error())
					}
				} else {
					logs.Error("graphics.Scale错误" + err.Error())
				}
			} else {
				newfn = filename
			}
		} else {
			logs.Error("image.Decode错误" + err.Error())
		}
		_ = reader.Close()
	} else {
		logs.Error("os.Open错误" + err.Error())
	}
	return
}

func Rgba2BW(filename string) (bw map[int]map[int]bool, w int, h int) {
	if !IsFile(filename) {
		logs.Error("文件不存在" + filename)
		return
	}
	if reader, err := os.Open(filename); err == nil {
		if m, _, err := image.Decode(reader); err == nil {
			bounds := m.Bounds()
			w, h = bounds.Dx(), bounds.Dy()
			i := 0
			u8s := make([]byte, w*h)
			var u8smin, u8smax byte
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, _ := m.At(x, y).RGBA()
					u8s[i] = byte((r + g + b) >> 8 / 3)
					if i == 0 {
						u8smax = u8s[i]
						u8smin = u8s[i]
					} else {
						if u8s[i] > u8smax {
							u8smax = u8s[i]
						}
						if u8s[i] < u8smin {
							u8smin = u8s[i]
						}
					}
					i++
				}
			}
			t := u8smin>>1 + u8smax>>1
			bw = make(map[int]map[int]bool)
			for k, v := range u8s {
				x := k % w
				if bw[x] == nil {
					bw[x] = make(map[int]bool)
				}
				bw[x][k/w] = v < t
			}
		} else {
			logs.Error("image.Decode错误" + err.Error())
		}
		_ = reader.Close()
	} else {
		logs.Error("os.Open错误" + err.Error())
	}
	return
}

func ShowBW(bw map[int]map[int]bool, w int, h int) {
	fmt.Println("@@@@@@@@@@@@@@@开始@@@@@@@@@@@@@@@")
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if bw[x][y] {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("")
	}
	fmt.Println("@@@@@@@@@@@@@@@结束@@@@@@@@@@@@@@@")
}

func checkDamp(bw map[int]map[int]bool, x int, y int) bool {
	for _, sx := range [3]int{x - 1, x, x + 1} {
		for _, sy := range [3]int{y - 1, y, y + 1} {
			if x == sx && y == sy {
				continue
			}
			if v, ok := bw[sx][sy]; ok && !v {
				return true
			}
		}
	}
	return false
}

func DampBW(bw map[int]map[int]bool) {
	damp := make(map[int]map[int]bool)
	for x, row := range bw {
		for y, v := range row {
			if v && checkDamp(bw, x, y) {
				if damp[x] == nil {
					damp[x] = make(map[int]bool)
				}
				damp[x][y] = true
			}
		}
	}
	for x, row := range damp {
		for y := range row {
			bw[x][y] = false
		}
	}
}

func markGroup(x, y int, group map[int]map[int]int, bw map[int]map[int]bool) {
	for _, sx := range [3]int{x - 1, x, x + 1} {
		for _, sy := range [3]int{y - 1, y, y + 1} {
			if x == sx && y == sy {
				continue
			}
			if bw[sx][sy] && group[sx][sy] == 0 {
				if group[sx] == nil {
					group[sx] = make(map[int]int)
				}
				group[sx][sy] = group[x][y]
				markGroup(sx, sy, group, bw)
			}
		}
	}
}

type dot struct{ x, y int }

func BW2Dot(bw map[int]map[int]bool, w int, h int) []dot {
	group := make(map[int]map[int]int)
	id := 0
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if bw[x][y] && group[x][y] == 0 {
				if group[x] == nil {
					group[x] = make(map[int]int)
				}
				id++
				group[x][y] = id
				markGroup(x, y, group, bw)
			}
		}
	}
	group2 := make([]struct{ minX, maxX, minY, maxY, num int }, id)
	for x, row := range group {
		for y, id := range row {
			if group2[id-1].num > 0 {
				if x > group2[id-1].maxX {
					group2[id-1].maxX = x
				} else if x < group2[id-1].minX {
					group2[id-1].minX = x
				}
				if y > group2[id-1].maxY {
					group2[id-1].maxY = y
				} else if y < group2[id-1].minY {
					group2[id-1].minY = y
				}
			} else {
				group2[id-1].minX = x
				group2[id-1].maxX = x
				group2[id-1].minY = y
				group2[id-1].maxY = y
			}
			group2[id-1].num++
		}
	}
	dots := make([]dot, id)
	for k, v := range group2 {
		dots[k].x = (v.minX + v.maxX) >> 1
		dots[k].y = (v.minY + v.maxY) >> 1
	}
	return dots
}

type dotCom struct {
	min, max dot
	ll       int
	l        float64
}

func DotAnalyze(dots []dot) (com map[int]map[int]dotCom, start int, end int, num int) {
	count := len(dots)
	combination := make(map[int]map[int]dotCom)
	lmin, lmax, minll := 0, 0, 0
	for min := 0; min < count-1; min++ {
		combination[min] = make(map[int]dotCom)
		for max := min + 1; max < count; max++ {
			ll := PowerInt(dots[min].x-dots[max].x, 2) + PowerInt(dots[min].y-dots[max].y, 2)
			combination[min][max] = dotCom{min: dots[min], max: dots[max], ll: ll, l: math.Sqrt(float64(ll))}
			if minll == 0 || minll > ll {
				minll = ll
				lmin = min
				lmax = max
			}
		}
	}
	lmin2, lmax2, minll2 := 0, 0, 0
	for min, row := range combination {
		for max, v := range row {
			if min == lmin && max == lmax {
				continue
			}
			if min == lmin || min == lmax || max == lmin || max == lmax {
				if minll2 == 0 || minll2 > v.ll {
					minll2 = v.ll
					lmin2 = min
					lmax2 = max
				}
			}
		}
	}
	if lmin2 == lmax || lmax2 == lmax {
		lmin, lmax = lmax, lmin
	}
	return combination, lmin, lmax, count
}

type dotInfo struct {
	angle, length float64
}

func getDotAngle(com map[int]map[int]dotCom, start int, end int) dotInfo {
	min, max := start, end
	if min > max {
		min, max = max, min
	}
	dot := com[min][max]
	angle := math.Asin(float64(dot.max.y-dot.min.y) / dot.l) //-0.5pi~0.5pi
	if start > end {
		angle = -angle //-0.5pi~0.5pi
	}
	if (start < end && dot.max.x >= dot.min.x) || (start > end && dot.max.x < dot.min.x) {
		if angle < 0 {
			angle += 2 * math.Pi //x>=0,y<0--0~0.5pi,1.5pi~2pi
		}
	} else {
		angle = math.Pi - angle //x<0,y>=0--0.5pi~1.5pi
	}
	return dotInfo{360 - angle*180/math.Pi, dot.l}
}

type dotRelation struct {
	angle, length int
}

type dotRelations []dotRelation

func (a dotRelations) Len() int {
	return len(a)
}

func (a dotRelations) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a dotRelations) Less(i, j int) bool {
	if a[j].angle == a[i].angle {
		return j > i
	}
	return a[j].angle > a[i].angle
}

func DotSort(com map[int]map[int]dotCom, start int, end int, num int) []dotRelation {
	data := make([]dotInfo, num)
	for i := 0; i < num; i++ {
		if i == start {
			continue
		}
		data[i] = getDotAngle(com, start, i)
	}
	angles := make(dotRelations, MaxInt(num-2, 0))
	for i, j := 0, 0; i < num; i++ {
		if i == start || i == end {
			continue
		}
		a := data[i].angle - data[end].angle
		if a < 0 {
			a += 360
		}
		angles[j] = dotRelation{Round(a * 100), Round(data[i].length / data[end].length * 10000)}
		j++
	}
	sort.Sort(angles)
	return angles
}

func GetImgKey(filename string, damp bool) []dotRelation {
	if resize := ImgResize(filename, 1000); resize != "" {
		filename = resize
	}
	bw, w, h := Rgba2BW(filename)
	if damp {
		DampBW(bw)
	}
	dots := BW2Dot(bw, w, h)
	com, start, end, num := DotAnalyze(dots)
	return DotSort(com, start, end, num)
}

func Compare(s, c dotRelations, ta, tl int) float64 {
	if len(s) != len(c) || len(s) == 0 {
		return 0
	}
	ta = MaxInt(ta, 1) * 100 //ta=100--角度浮动1度
	tl = MaxInt(tl, 1) * 100 //tl=100--距离浮动1%
	sort.Sort(s)
	sort.Sort(c)
	score := 0.0
	for i := 0; i < len(s); i++ {
		diffa, diffl := AbsInt(s[i].angle-c[i].angle), AbsInt(s[i].length-c[i].length)
		if diffa > 2*ta || diffl > 2*tl {
			continue
		}
		if diffa > ta {
			score -= float64(diffa-ta) / float64(ta)
		}
		if diffl > tl {
			score -= float64(diffl-tl) / float64(tl)
		}
		score += 2
	}
	return score / float64(len(s)*2)
}
