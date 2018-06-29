package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
)

// LevelList 蚕宝宝等级列表
func LevelList() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,level,redeemitem,updatetime from swlevel order by id")
}

// EditLevel 编辑蚕宝宝等级
func EditLevel(redeemitem, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update swlevel set redeemitem = ?,updatetime = ? where id = ?", redeemitem, nowTime, id)
}
