package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
)

// TreeLevelList 桑树等级列表
func TreeLevelList() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,level,growthhours,maxhours,updatetime from tree order by id")
}

// EditTreeLevel 编辑桑树等级
func EditTreeLevel(growthhours, maxhours, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update tree set growthhours = ?,maxhours = ?,updatetime = ? where id = ?", growthhours, maxhours, nowTime, id)
}
