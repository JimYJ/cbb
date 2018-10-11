package silkworm

import (
	"canbaobao/common"
	log "canbaobao/service/logs"
	"fmt"

	"github.com/JimYJ/easysql/mysql"
)

// CheckHatch 检查目前孵化中的数量
func CheckHatch(openid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from usersw left join user on uid = user.id where hatch = ? and openid = ?", 0, openid)
}

// Hatch 孵化蚕仔
func Hatch(uid, rucksackid string, swtype int, enabletime int64) bool {
	mysqlConn := common.GetMysqlConn()
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	_, err = tx.Insert("insert into usersw set uid = ?,swtype = ?,hatch = ?,exp = ?,name = ?,health = ?,level = ?,enabletime = ?,enable = ?", uid, swtype, 0, 0, "蚕仔", 100, 1, enabletime, 1) //enable改为1关闭孵化过程
	_, err2 := tx.Delete("delete from rucksack where id = ?", rucksackid)
	if err != nil || err2 != nil {
		log.Println(err, err2)
		tx.Rollback()
		return false
	}
	tx.Commit()
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
func ApplyPair(id, pairid, uid, pairuid string, pairtime int64) bool {
	mysqlConn := common.GetMysqlConn()
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	_, err = tx.Update("update usersw set pair = ?,pairtime = ?,pairsrc = ?,pairid = ?,pairuid = ? where id = ?", 1, pairtime, 1, pairid, pairuid, id)
	_, err2 := tx.Update("update usersw set pair = ?,pairtime = ?,pairsrc = ?,pairid = ?,pairuid = ? where id = ?", 1, pairtime, 0, id, uid, pairid)
	if err != nil || err2 != nil {
		log.Println(err, err2)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

// AllowPair 同意配对
func AllowPair(id, pairid string, pairtime int64) bool {
	mysqlConn := common.GetMysqlConn()
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	_, err = tx.Update("update usersw set pair = ?,pairtime = ? where id = ?", 2, pairtime, id)
	_, err2 := tx.Update("update usersw set pair = ?,pairtime = ? where id = ?", 2, pairtime, pairid)
	if err != nil || err2 != nil {
		log.Println(err, err2)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

// EndPair 结束/拒绝配对
func EndPair(id, pairid string) bool {
	mysqlConn := common.GetMysqlConn()
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	_, err = tx.Update("update usersw set pair = ?,pairtime = ?,pairsrc = ?,pairid = ?,pairuid = ? where id = ?", 0, 0, 0, 0, 0, id)
	_, err2 := tx.Update("update usersw set pair = ?,pairtime = ?,pairsrc = ?,pairid = ?,pairuid = ? where id = ?", 0, 0, 0, 0, 0, pairid)
	if err != nil || err2 != nil {
		log.Println(err, err2)
		tx.Rollback()
		return false
	}
	tx.Commit()
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

// GetSingleUserSWInfo 获取单只蚕宝宝的信息
func GetSingleUserSWInfo(id string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select id,uid,hatch,level,exp,swtype,enable,health from usersw where id = ?", id)
}

// UpExp 增加经验值
func UpExp(newExp, level, name, id, uid, rucksackid, keyTimes, keyDate, ip, nowTime, nowDate, itemid string, feedTimes, health int) bool {
	mysqlConn := common.GetMysqlConn()
	if itemid == "1" && health < 100 {
		health++
	}
	usersql := fmt.Sprintf("update user set loginip = ?,logintime = ?,%s = ?,%s = ? where id = ?", keyTimes, keyDate)
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	_, err = tx.Delete("delete from rucksack where uid = ? and id = ?", uid, rucksackid)
	_, err2 := tx.Update("update usersw set health = ?,exp = ?,level = ?,name = ? where id = ?", health, newExp, level, name, id)
	_, err3 := tx.Update(usersql, ip, nowTime, feedTimes, nowDate, uid)
	if err != nil || err2 != nil || err3 != nil {
		log.Println(err, err2, err3)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

// BeButterfly 化蝶
func BeButterfly(newExp, name, swid, id, uid, rucksackid, level, loginip, nowTime, nowDate, keyTimes, keyDate string, feedTimes int) bool {
	mysqlConn := common.GetMysqlConn()
	usersql := fmt.Sprintf("update user set level = ?,loginip = ?,logintime = ?,updatetime = ?,%s = ?,%s = ? where id = ?", keyTimes, keyDate)
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	_, err = tx.Delete("delete from rucksack where uid = ? and id = ?", uid, rucksackid)
	_, err2 := tx.Update("update usersw set exp = ?,level = ?,hatch = ?, name = ?,swid = ? where id = ?", newExp, 10, 1, name, swid, id)
	_, err3 := tx.Update(usersql, level, loginip, nowTime, nowTime, feedTimes, nowDate, uid)
	if err != nil || err2 != nil || err3 != nil {
		log.Println(err, err2, err3)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

// GetUserButterflyList 获得用户可兑换蝴蝶列表-按种类
func GetUserButterflyList(swid, uid, limit string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	sql := fmt.Sprintf("select id from usersw where swid = ? and uid = ? and hatch = ? and pair = ? limit %s", limit)
	return mysqlConn.GetResults(mysql.Statement, sql, swid, uid, 1, 0)
}

// GetUserButterflyAllList 获得用户蝴蝶列表-不限种类
func GetUserButterflyAllList(uid, limit string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	sql := fmt.Sprintf("select id from usersw where uid = ? and hatch = ? and pair = ? limit %s", limit)
	return mysqlConn.GetResults(mysql.Statement, sql, uid, 1, 0)
}

// GetSWlist 获取蚕宝宝
func GetSWlist() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select usersw.id,uid,health,user.name,leafday,sppday,mppday,lppday FROM usersw left join user on uid = user.id where hatch = 0")
}

// UpdateHealth 更新健康值
func UpdateHealth(updateHealth map[string]int) bool {
	mysqlConn := common.GetMysqlConn()
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	commit := true
	for id, health := range updateHealth {
		_, err := tx.Update("update usersw set health = ? where id = ?", health, id)
		if err != nil {
			log.Println(err)
			commit = false
			break
		}
	}
	if !commit {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return commit
}
