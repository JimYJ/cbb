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
func Hatch(uid, rucksackid string, swtype int, enabletime int64) bool {
	mysqlConn := common.GetMysqlConn()
	mysqlConn.TxBegin()
	_, err := mysqlConn.TxInsert(mysql.Statement, "insert into usersw set uid = ?,swtype = ?,hatch = ?,exp = ?,name = ?,health = ?,level = ?,enabletime = ?,enable = ?", uid, swtype, 0, 0, "蚕仔", 100, 1, enabletime, 0)
	_, err2 := mysqlConn.TxDelete(mysql.Statement, "delete from rucksack where id = ?", rucksackid)
	if err != nil || err2 != nil {
		log.Println(err, err2)
		mysqlConn.TxRollback()
		return false
	}
	mysqlConn.TxCommit()
	return true
}

// GetUserSilkworm 获得用户蚕宝宝列表
func GetUserSilkworm(id string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,swtype,hatch,exp,name,level,swid,health,enable,enabletime from usersw where uid = ?", id)
}

// Enable 结束成长时间
func Enable(id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update usersw set enable = ?,enabletime = ? where id = ?", 1, 0, id)
}

// GetUseButterflyCountByID 统计用户蝴蝶数量
func GetUseButterflyCountByID(id string) string {
	mysqlConn := common.GetMysqlConn()
	str, err := mysqlConn.GetVal(mysql.Statement, "select count(*) from usersw where uid = ? and hatch = ?", id, 1)
	if err != nil {
		log.Println(err)
	}
	return str
}

// GetUserButterflyCountByOpenID 统计用户蝴蝶数量
func GetUserButterflyCountByOpenID(openid string) string {
	mysqlConn := common.GetMysqlConn()
	str, err := mysqlConn.GetVal(mysql.Statement, "select count(*) from usersw left join user on user.id = uid where openid = ? and hatch = ?", openid, 1)
	if err != nil {
		log.Println(err)
	}
	return str
}
