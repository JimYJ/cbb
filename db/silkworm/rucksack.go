package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
	"log"
	"strconv"
)

// ItemType
const (
	Silkworm = iota
	More
)

// AddRucksack 新增物品到背包
func AddRucksack(itemid, uid, nowTime string, swtype, take, itemtype int) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into rucksack set itemid = ?,uid = ?,itemtype = ?,swtype = ?,updatetime = ?,createtime = ?,take = ?",
		itemid, uid, itemtype, swtype, nowTime, nowTime, take)
}

// AddItemRucksack 普通物品进入背包
func AddItemRucksack(itemid, uid, nowTime string, itemtype int) (int64, error) {
	return AddRucksack(itemid, uid, nowTime, -1, 1, 1)
}

// AddLeafRucksack 生成桑叶等待拾取
func AddLeafRucksack(itemid, uid, nowTime string, itemtype int) (int64, error) {
	return AddRucksack(itemid, uid, nowTime, -1, 0, 1)
}

// AddSilkwormRucksack 生成蚕宝宝进入背包
func AddSilkwormRucksack(itemid, uid, swtype, img, nowTime string, itemtype int) (int64, error) {
	swtypeint, err := strconv.Atoi(swtype)
	if err != nil {
		return 0, err
	}
	return AddRucksack(itemid, uid, nowTime, swtypeint, 1, 0)
}

// UserRucksack 获得用户背包
func UserRucksack(openid string, isTake bool) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	var take int
	if isTake {
		take = 1
	} else {
		take = 0
	}
	return mysqlConn.GetResults(mysql.Statement, "select COUNT(*) as num,itemid from rucksack left join user on uid = user.id where openid = ? and take = ? GROUP BY itemid", openid, take)
}

// UserRucksackCount 用户背包物品数量
func UserRucksackCount(openid string, isTake bool) (string, error) {
	mysqlConn := common.GetMysqlConn()
	var take int
	if isTake {
		take = 1
	} else {
		take = 0
	}
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from rucksack left join user on uid = user.id where openid = ? and take = ?", openid, take)
}

// GetUserSWID 获得用户背包蚕仔ID
func GetUserSWID(openid string, swtype int) (string, string, error) {
	mysqlConn := common.GetMysqlConn()
	nums, err := mysqlConn.GetVal(mysql.Statement, "select count(*) as nums from rucksack left join user on uid = user.id where openid = ? and take = ? and swtype = ?", openid, 1, swtype)
	if err != nil || nums == "0" {
		return "0", "", err
	}
	id, _ := mysqlConn.GetVal(mysql.Statement, "select rucksack.id from rucksack left join user on uid = user.id where openid = ? and take = ? and swtype = ? limit 1", openid, 1, swtype)
	return "1", id, nil
}

// GetUserLeafUntake 获得用户未拾取桑叶
func GetUserLeafUntake(openid string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select rucksack.id from rucksack left join user on uid = user.id where openid = ? and take = ? and itemid = ? order by id", openid, 0, 1)
}

// TakeLeaf 收取桑叶
func TakeLeaf(openid, id, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	uinfo, err := GetUID(openid)
	if err != nil {
		return 0, err
	}
	uid := uinfo["id"]
	return mysqlConn.Update(mysql.Statement, "update rucksack set take = ?,createtime = ?,updatetime = ? where uid = ? and id = ?", 1, nowTime, nowTime, uid, id)
}

// GetUserLeafUntakeByID 获得用户未拾取桑叶
func GetUserLeafUntakeByID(id string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id from rucksack where uid = ? and take = ? and itemid = ? order by id", id, 0, 1)
}

// TakeLeafByID 收取桑叶
func TakeLeafByID(openid, loseUID, id, nowTime string) int {
	mysqlConn := common.GetMysqlConn()
	uinfo, err := GetUID(openid)
	if err != nil {
		return -1
	}
	takeUID := uinfo["id"]
	mysqlConn.TxBegin()
	var err2 error
	rs, err := mysqlConn.Delete(mysql.Statement, "delete from rucksack where uid = ? and id = ? and take = ? and itemid = ?", loseUID, id, 0, 1)
	if rs >= 1 {
		_, err2 = mysqlConn.Insert(mysql.Statement, "insert into rucksack set take = ?,uid = ?,itemid = ?,swtype = ?,createtime = ?,updatetime = ?", 1, takeUID, 1, -1, nowTime, nowTime)
	} else {
		mysqlConn.TxRollback()
		return -1
	}
	if err != nil || err2 != nil {
		log.Println(err, err2)
		mysqlConn.TxRollback()
		return -2
	}
	mysqlConn.TxCommit()
	return 1
}
