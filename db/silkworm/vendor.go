package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
)

// DelVendor 删除店铺
func DelVendor(id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete(mysql.Statement, "delete from vendor where id = ?", id)
}

// AddVendor 新增店铺
func AddVendor(name, leader, leaderphone, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into vendor set name = ?,leader = ?,leaderphone = ?,createtime = ?,updatetime = ?", name, leader, leaderphone, nowTime, nowTime)
}

// EditVendor 编辑店铺
func EditVendor(name, leader, leaderphone, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update vendor set name = ?,leader = ?,leaderphone = ?,updatetime = ? where id = ?", name, leader, leaderphone, nowTime, id)
}

// GetVendor 获取店铺
func GetVendor() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,leader,leaderphone,createtime,updatetime from vendor ORDER BY id desc")
}

// GetVendorName 获取店铺名称
func GetVendorName(id string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select name from vendor where id = ?", id)
}
