package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
)

// GuideList 获得攻略图
func GuideList() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,imgurl,updatetime from guide order by id")
}

// EditGuide 编辑攻略图
func EditGuide(imgurl, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update guide set imgurl = ?,updatetime = ? where id = ?", imgurl, nowTime, id)
}
