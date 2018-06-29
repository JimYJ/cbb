package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
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
