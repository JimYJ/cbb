package silkworm

import (
	"canbaobao/common"
	log "canbaobao/service/logs"
	"strconv"

	"github.com/JimYJ/easysql/mysql"
)

// ItemType
const (
	Silkworm = iota
	More
)

// AddRucksack 新增物品到背包
func AddRucksack(itemid, uid, nowTime string, swtype, take, itemtype int) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	rs, err := mysqlConn.Insert(mysql.Statement, "insert into rucksack set itemid = ?,uid = ?,itemtype = ?,swtype = ?,updatetime = ?,createtime = ?,take = ?",
		itemid, uid, itemtype, swtype, nowTime, nowTime, take)
	return rs, err
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
func AddSilkwormRucksack(itemid, uid, swtype, nowTime string, itemtype int) (int64, error) {
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

// UserRucksackItemCount 获得用户某项背包数量
func UserRucksackItemCount(openid, itemid string, isTake bool) (string, error) {
	mysqlConn := common.GetMysqlConn()
	var take int
	if isTake {
		take = 1
	} else {
		take = 0
	}
	return mysqlConn.GetVal(mysql.Statement, "select COUNT(*) as num from rucksack left join user on uid = user.id where itemid = ? AND openid = ? and take = ?", itemid, openid, take)
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

// GetUserLeafUntakeCountByID 获得用户未拾取桑叶数
func GetUserLeafUntakeCountByID(id string) int {
	mysqlConn := common.GetMysqlConn()
	rs, err := mysqlConn.GetVal(mysql.Statement, "select count(*) from rucksack where uid = ? and take = ? and itemid = ? order by id", id, 0, 1)
	if err != nil {
		return 0
	}
	var reInt int
	if len(rs) > 0 {
		reInt, _ = strconv.Atoi(rs)
	} else {
		reInt = 0
	}
	return reInt
}

// TakeLeafByID 偷桑叶
func TakeLeafByID(openid, loseUID, id, nowTime string) int {
	mysqlConn := common.GetMysqlConn()
	uinfo, err := GetUID(openid)
	if err != nil {
		return -1
	}
	takeUID := uinfo["id"]
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return -1
	}
	var err2 error
	rs, err := tx.Delete("delete from rucksack where uid = ? and id = ? and take = ? and itemid = ?", loseUID, id, 0, 1)
	if rs >= 1 {
		_, err2 = tx.Insert("insert into rucksack set take = ?,uid = ?,itemid = ?,swtype = ?,createtime = ?,updatetime = ?", 1, takeUID, 1, -1, nowTime, nowTime)
	} else {
		tx.Rollback()
		return -1
	}
	if err != nil || err2 != nil {
		log.Println(err, err2)
		tx.Rollback()
		return -2
	}
	tx.Commit()
	return 1
}

// RucksackItemInfo 获取背包物品信息
func RucksackItemInfo(itemid, uid string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select id,uid,take from rucksack where itemid = ? and uid = ? order by take desc limit 1", itemid, uid)
}

// SproutLeaf 生成桑叶
func SproutLeaf(itemid, uid, nowTime, nowDate string, sproutleafs, growthhours int) bool {
	mysqlConn := common.GetMysqlConn()
	commit := true
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	_, err = tx.Update("update user set sproutleafs = ?,sproutleafday = ?,updatetime = ? where id = ?", sproutleafs, nowDate, nowTime, uid)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return false
	}
	for i := 0; i < growthhours; i++ {
		_, err := tx.Insert("insert into rucksack set itemid = ?,uid = ?,itemtype = ?,swtype = ?,updatetime = ?,createtime = ?,take = ?",
			itemid, uid, 1, -1, nowTime, nowTime, 0)
		if err != nil {
			log.Println(err)
			commit = false
			break
		}
	}
	if !commit {
		tx.Rollback()
	}
	tx.Commit()
	return commit
}

// UserInviteAward 发放用户奖励
func UserInviteAward(uid, itemid, num, nowTime string) error {
	mysqlConn := common.GetMysqlConn()
	tx, err := mysqlConn.Begin()
	if err != nil {
		return err
	}
	numI, err := strconv.Atoi(num)
	if err != nil {
		return err
	}
	var itemtype int
	if itemid == "1" {
		itemtype = 0
	} else {
		itemtype = 1
	}
	for i := 0; i < numI; i++ {
		_, err = tx.Insert("insert into rucksack set itemid = ?,uid = ?,itemtype = ?,swtype = ?,updatetime = ?,createtime = ?,take = ?",
			itemid, uid, itemtype, 0, nowTime, nowTime, 1)
		if err != nil {
			break
		}
	}
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return nil
}
