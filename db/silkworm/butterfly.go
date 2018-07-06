package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
)

// ButterflyList 蝴蝶种类列表
func ButterflyList() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,type,updatetime,img from butterfly order by id")
}

// EditButterfly 编辑蝴蝶种类
func EditButterfly(name, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update butterfly set name = ?,updatetime = ? where id = ?", name, nowTime, id)
}

// ButterflyName 获得蝴蝶名称
func ButterflyName(id int) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select name from butterfly where id = ?", id)
}
