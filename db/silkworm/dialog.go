package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
	"log"
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

// GetRandomDialog 随机获取蚕宝宝对话
func GetRandomDialog() string {
	mysqlConn := common.GetMysqlConn()
	content, err := mysqlConn.GetVal(mysql.Statement, "SELECT content FROM `swdialog` AS t1 JOIN (SELECT ROUND(RAND() * ((SELECT MAX(id) FROM `swdialog`)-(SELECT MIN(id) FROM `swdialog`))+(SELECT MIN(id) FROM `swdialog`)) AS id) AS t2 WHERE t1.id >= t2.id ORDER BY t1.id LIMIT 1")
	if err != nil {
		log.Println(err)
	}
	return content
}
