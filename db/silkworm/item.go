package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
)

// ItemList 获得物品列表
func ItemList() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,exp,limitday,updatetime from item order by id")
}

// EditItem 编辑物品
func EditItem(exp, limitday, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update item set exp = ?,limitday = ?,updatetime = ? where id = ?", exp, limitday, nowTime, id)
}

// ItemTypeList 获得分类列表
func ItemTypeList(isLeaf bool) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	var types int
	if isLeaf {
		types = 0
	} else {
		types = 1
	}
	return mysqlConn.GetResults(mysql.Statement, "select id,name,exp,limitday,updatetime from item where types = ? order by id", types)
}

// ItemInfo 获得单个物品信息
func ItemInfo(id string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select id,name,exp,limitday,types,img from item where id = ?", id)
}
