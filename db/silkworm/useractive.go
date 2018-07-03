package silkworm

import (
	"canbaobao/common"
	"fmt"
	"github.com/JimYJ/easysql/mysql"
)

// UserActive
const (
	ActiveTreeup = iota
	ActiveLevelup
	ActiveBindvendor
	ActiveVoucher
	ActiveShare
	ActiveAnswer
	ActiveHatch
	ActiveSign
	ActivePair
	ActiveTakeLeaf
	ActiveStealLeaf
)

// UserActiveList
var (
	ActiveList    = []string{"treeup", "levelup", "bindvendor", "voucher", "share", "answer", "larva", "sign", "pair", "takeleaf", "stealleaf"}
	ActiveStrList = []string{
		"的桑树升到了 %s 级。",
		"的用户等级升级，获得了物品:",
		"绑定了店铺.",
		"获得了兑换券，可兑换物品:",
		"分享注册成功，获得了物品:",
		"回答正确，获得物品:",
		"成功孵化了蚕仔。",
		"连续签到%s，获得了物品:",
		"配对成功，产下了:",
		"拾取了桑叶。",
		"偷摘了 %s 的桑叶。",
	}
)

// SaveUserActive 保存用户动态
func SaveUserActive(types int, uname, uid, itemname, itemid, nowTime, moreInfo string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	var content, str string

	if types == 2 || types == 6 || types == 0 || types == 9 || types == 10 {
		if types == 0 || types == 10 {
			str = fmt.Sprintf(ActiveStrList[types], moreInfo)
		} else {
			str = ActiveStrList[types]
		}
		content = fmt.Sprintf("%s%s", uname, str)
	} else {
		if types == 7 {
			str = fmt.Sprintf(ActiveStrList[types], moreInfo)
		} else {
			str = ActiveStrList[types]
		}
		content = fmt.Sprintf("%s%s%s", uname, str, itemname)
	}
	return mysqlConn.Insert(mysql.Statement, "insert into useractive set type = ?,uname = ?,uid = ?,itemname = ?,itemid = ?,content = ?,createtime = ?",
		ActiveList[types], uname, uid, itemname, itemid, content, nowTime)
}

// GetUserVid 获取用户VID
func GetUserVid(openid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select vid from user where openid = ?", openid)
}

// GetUserActive 获取好友动态
func GetUserActive(openid, vid, paginaSQL string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	sql := fmt.Sprintf("select useractive.id,user.name as username,user.avatar,useractive.content,useractive.createtime from useractive left join user on uid = user.id where user.vid = ? order by id desc %s", paginaSQL)
	return mysqlConn.GetResults(mysql.Statement, sql, vid)
}

// GetUserActiveCount 获取好友动态总数
func GetUserActiveCount(vid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from useractive left join user on uid = user.id where user.vid = ?", vid)
}
