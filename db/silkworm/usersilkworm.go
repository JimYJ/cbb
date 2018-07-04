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
	return mysqlConn.GetResults(mysql.Statement, "select id,swtype,hatch,exp,name,level,swid,health,enable,enabletime,pair,pairtime,pairid,pairsrc,pairuid from usersw where uid = ?", id)
}

// Enable 结束成长时间
func Enable(id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update usersw set enable = ?,enabletime = ? where id = ?", 1, 0, id)
}

// ApplyPair 申请配对
func ApplyPair(id, pairid, uid, pairuid string) bool {
	mysqlConn := common.GetMysqlConn()
	mysqlConn.TxBegin()
	_, err := mysqlConn.TxUpdate(mysql.Statement, "update usersw set pair = ?,pairtime = ?,pairsrc = ?,pairid = ?,pairuid = ? where id = ?", 1, 0, 1, pairid, pairuid, id)
	_, err2 := mysqlConn.TxUpdate(mysql.Statement, "update usersw set pair = ?,pairtime = ?,pairsrc = ?,pairid = ?,pairuid = ? where id = ?", 1, 0, 0, id, uid, pairid)
	if err != nil || err2 != nil {
		log.Println(err, err2)
		mysqlConn.TxRollback()
		return false
	}
	mysqlConn.TxCommit()
	return true
}

// AllowPair 同意配对
func AllowPair(id, pairid string, pairtime int64) bool {
	mysqlConn := common.GetMysqlConn()
	mysqlConn.TxBegin()
	_, err := mysqlConn.TxUpdate(mysql.Statement, "update usersw set pair = ?,pairtime = ? where id = ?", 2, pairtime, id)
	_, err2 := mysqlConn.TxUpdate(mysql.Statement, "update usersw set pair = ?,pairtime = ? where id = ?", 2, pairtime, pairid)
	if err != nil || err2 != nil {
		log.Println(err, err2)
		mysqlConn.TxRollback()
		return false
	}
	mysqlConn.TxCommit()
	return true
}

// EndPair 结束/拒绝配对
func EndPair(id, pairid string) bool {
	mysqlConn := common.GetMysqlConn()
	mysqlConn.TxBegin()
	_, err := mysqlConn.TxUpdate(mysql.Statement, "update usersw set pair = ?,pairtime = ?,pairsrc = ?,pairid = ?,pairuid = ? where id = ?", 0, 0, 0, 0, 0, id)
	_, err2 := mysqlConn.TxUpdate(mysql.Statement, "update usersw set pair = ?,pairtime = ?,pairsrc = ?,pairid = ?,pairuid = ? where id = ?", 0, 0, 0, 0, 0, pairid)
	if err != nil || err2 != nil {
		log.Println(err, err2)
		mysqlConn.TxRollback()
		return false
	}
	mysqlConn.TxCommit()
	return true
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

// CheckPairCondition 检测是否符合配对条件
func CheckPairCondition(userswid string) (string, string, string) {
	mysqlConn := common.GetMysqlConn()
	info, err := mysqlConn.GetRow(mysql.Statement, "select hatch,pair,uid from usersw where id = ?", userswid)
	if err != nil {
		log.Println(err)
		return "", "", ""
	}
	return info["hatch"], info["pair"], info["uid"]
}
