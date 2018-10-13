package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
)

// DelVendor 删除店铺
func DelVendor(id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete(mysql.Statement, "delete from litemall_vendor where id = ?", id)
}

// AddVendor 新增店铺
func AddVendor(name, leader, leaderphone, nowTime, province, city, county string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into litemall_vendor set name = ?,leader = ?,leaderphone = ?,createtime = ?,updatetime = ?,province = ?,city = ?,county = ?", name, leader, leaderphone, nowTime, nowTime, province, city, county)
}

// EditVendor 编辑店铺
func EditVendor(name, leader, leaderphone, nowTime, id, province, city, county string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update litemall_vendor set name = ?,leader = ?,leaderphone = ?,updatetime = ?,province = ?,city = ?,county = ? where id = ?", name, leader, leaderphone, nowTime, province, city, county, id)
}

// GetVendor 获取店铺
func GetVendor() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,leader,leaderphone,createtime,updatetime,province,city,county from litemall_vendor ORDER BY id desc")
}

// GetVendorName 获取店铺名称
func GetVendorName(id string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select name from litemall_vendor where id = ?", id)
}

// GetVendorByArea 根据地域筛选获取店铺
func GetVendorByArea(province, city, county string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name from litemall_vendor where province = ? and city = ? and county = ? ORDER BY id desc", province, city, county)
}
