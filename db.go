package tool

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"reflect"
	"strings"
)

func TableExists(table string) bool {
	o := orm.NewOrm()
	maps := make([]orm.Params, 0)
	if _, err := o.Raw("SHOW TABLES LIKE '" + table + "'").Values(&maps); err == nil && len(maps) == 1 {
		return true
	}
	return false
}

func GetMI(m interface{}) ([]string, []string, []interface{}) {
	columns := make([]string, 0)
	qmarks := make([]string, 0)
	values := make([]interface{}, 0)
	val := reflect.ValueOf(m).Elem()
	s := reflect.TypeOf(m).Elem()
	for index := 0; index < val.NumField(); index++ {
		tag := s.Field(index).Tag.Get("orm")
		if i := strings.Index(tag, "auto"); i >= 0 {
			continue
		}
		for _, v := range strings.Split(tag, ";") {
			if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 && strings.ToLower(v)[:i] == "column" {
				columns = append(columns, v[i+1:len(v)-1])
				qmarks = append(qmarks, "?")
				values = append(values, val.Field(index).Interface())
				break
			}
		}
	}
	return columns, qmarks, values
}

func CreateTable(date string) {
	if !TableExists("car_advplay_log" + date) {
		sql := "CREATE TABLE `car_advplay_log" + date + "` (" +
			"`id` int(11) NOT NULL AUTO_INCREMENT," +
			"`sn` varchar(20) NOT NULL DEFAULT '' COMMENT '设备sn'," +
			"`orderid` int(11) NOT NULL DEFAULT '0' COMMENT '父广告订单id'," +
			"`ordersubid` int(11) NOT NULL DEFAULT '0' COMMENT '子广告订单id'," +
			"`number` int(11) NOT NULL DEFAULT '0' COMMENT '播放次数'," +
			"`money` decimal(10,6) NOT NULL DEFAULT '0.000000' COMMENT '金额'," +
			"`lng` decimal(8,5) NOT NULL DEFAULT '0.00000' COMMENT '经度'," +
			"`lat` decimal(8,5) NOT NULL DEFAULT '0.00000' COMMENT '纬度'," +
			"`address` varchar(100) DEFAULT '' COMMENT '地址'," +
			"`createtime` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间'," +
			"PRIMARY KEY (`id`)" +
			") ENGINE=MyISAM DEFAULT CHARSET=utf8;"
		o := orm.NewOrm()
		if _, err := o.Raw(sql).Exec(); err != nil {
			logs.Error("CreateTable失败:" + err.Error())
		}
	}
}
