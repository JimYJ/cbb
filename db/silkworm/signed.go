package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
	"log"
)

// SignedList 获得签到
func SignedList() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,intro,dayitemid,weekitemid,updatetime from signed order by id")
}

// EditSigned 编辑签到
func EditSigned(intro, dayitemid, weekitemid, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update signed set intro = ?,dayitemid = ?,weekitemid = ?,updatetime = ? where id = ?", intro, dayitemid, weekitemid, nowTime, id)
}

// SignedItem 获得签到奖励物品ID
func SignedItem() (string, string) {
	mysqlConn := common.GetMysqlConn()
	rs, err := mysqlConn.GetRow(mysql.Statement, "select dayitemid,weekitemid from signed where id = ?", 1)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	return rs["dayitemid"], rs["weekitemid"]
}
