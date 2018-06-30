package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
	"log"
)

// CheckHatch 检查目前孵化中的数量
func CheckHatch(openid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from usersw left join user on uid = user.id where hatch = ? and openid = ?", 0, openid)
}

// Hatch 孵化蚕仔
func Hatch(uid, rucksackid string, swtype int) bool {
	mysqlConn := common.GetMysqlConn()
	mysqlConn.TxBegin()
	_, err := mysqlConn.TxInsert(mysql.Statement, "insert into usersw set uid = ?,swtype = ?,hatch = ?,exp = ?,name = ?,level = ?", uid, swtype, 0, 0, "蚕仔", 1)
	_, err2 := mysqlConn.TxDelete(mysql.Statement, "delete from rucksack where id = ?", rucksackid)
	if err != nil || err2 != nil {
		log.Println(err, err2)
		mysqlConn.TxRollback()
		return false
	}
	mysqlConn.TxCommit()
	return true
}
