package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
)

// DelDialog 删除蚕宝宝对话
func DelDialog(id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete(mysql.Statement, "delete from swdialog where id = ?", id)
}

// AddDialog 新增蚕宝宝对话
func AddDialog(content, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into swdialog set content = ?,createtime = ?,updatetime = ?", content, nowTime, nowTime)
}

// EditDialog 编辑蚕宝宝对话
func EditDialog(content, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update swdialog set content = ?,updatetime = ? where id = ?", content, nowTime, id)
}

// GetDialog 获取蚕宝宝对话
func GetDialog() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,content,createtime,updatetime from swdialog ORDER BY id desc")
}
