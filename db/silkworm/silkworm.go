package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
)

// LevelList 蚕宝宝等级列表
func LevelList() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,level,exp,redeemitem,updatetime,img from swlevel order by id")
}

// EditLevel 编辑蚕宝宝等级
func EditLevel(redeemitem, exp, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update swlevel set redeemitem = ?,exp = ?,updatetime = ? where id = ?", redeemitem, exp, nowTime, id)
}

// GetLevel 蚕宝宝等级信息
func GetLevel(level string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select id,name,level,exp,redeemitem from swlevel where level = ?", level)
}
