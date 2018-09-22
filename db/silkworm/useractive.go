package silkworm

import (
	"canbaobao/common"
	log "canbaobao/service/logs"
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
	ActivePairEnd
	ActiveTakeLeaf
	ActiveStealLeaf
	ActivePairApply
	ActivePairApplyed
	ActivePairAllow
	ActivePairReject
	ActivePairEndII
	ActivePairAllowII
	ActivePairRejectII
	ActiveNewUser
	ActiveFirstSWUp
	ActiveSWUp
	ActiveBeButterfly
	ActiveSproutLeaf
	ActiveSubHealth
	ActiveBatchVoucher
)

// UserActiveList
var (
	ActiveList = []string{
		"treeup",
		"levelup",
		"bindvendor",
		"voucher",
		"share",
		"answer",
		"larva",
		"sign",
		"pairend",
		"takeleaf",
		"stealleaf",
		"pairapply",
		"pairapplyed",
		"pairallow",
		"pairreject",
		"pairendii",
		"pairallowii",
		"pairrejectii",
		"newuser",
		"firstswup",
		"swup",
		"bebutterfly",
		"sproutleaf",
		"subhealth",
		"batchvoucher",
	}
	ActiveStrList = []string{
		"的桑树升到了 %s 级。",
		"的用户等级成功升到了 %s 级。",
		"绑定了店铺.",
		"兑换商品成功，获得了兑换券，可兑换物品:",
		"分享注册成功，获得了物品:",
		"回答正确，获得物品:",
		"成功孵化了蚕仔。",
		"连续签到%s，获得了物品:",
		"与 %s 的蝴蝶配对成功，产下了:",
		"拾取了桑叶。",
		"偷摘了 %s 的桑叶。",
		"向 %s 发起了蝴蝶配对申请。",
		"收到了来自 %s 的蝴蝶配对申请。",
		"同意了 %s 的蝴蝶配对申请。",
		"拒绝了 %s 的蝴蝶配对申请。",
		"与 %s 的蝴蝶配对结束。",
		"的蝴蝶配对申请被 %s 通过。",
		"的蝴蝶配对申请被 %s 拒绝。",
		"成为了新用户，获得了物品:",
		"第一次将蚕宝宝升到了 %s 级，获得兑换券，可兑换物品:",
		"的蚕宝宝升到了 %s 级。",
		"成功将蚕仔成长成蝴蝶，用户等级升为 %s 级，并同时获得了新的蚕仔。",
		"的桑树长出了 %s 片桑叶。",
		"用户昨天没有喂蚕宝宝，蚕宝宝健康值下降 10% 。",
		"获得了平台赠送的兑换券，有效期为:%s，可兑换物品:",
	}
)

// SaveUserActive 保存用户动态
func SaveUserActive(types int, uname, uid, itemname, itemid, nowTime, moreInfo string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	var content, str string

	if types == ActiveLevelup || types == ActiveBindvendor || types == 6 || types == ActiveTreeup || types == 9 || types == 10 || types == 11 || types == 12 || types == 13 || types == 14 || types == 15 || types == 16 || types == 17 || types == ActiveSWUp || types == ActiveBeButterfly || types == ActiveSproutLeaf {
		if types == ActiveTreeup || types == 10 || types == 11 || types == 12 || types == 13 || types == 14 || types == 15 || types == 16 || types == 17 || types == ActiveSWUp || types == ActiveBeButterfly || types == ActiveSproutLeaf {
			str = fmt.Sprintf(ActiveStrList[types], moreInfo)
		} else {
			str = ActiveStrList[types]
		}
		content = fmt.Sprintf("%s%s", uname, str)
	} else {
		if types == ActiveSign || types == ActivePairEnd || types == ActiveFirstSWUp || types == ActiveBatchVoucher {
			str = fmt.Sprintf(ActiveStrList[types], moreInfo)
		} else {
			str = ActiveStrList[types]
		}
		content = fmt.Sprintf("%s%s%s", uname, str, itemname)
	}
	return mysqlConn.Insert(mysql.Statement, "insert into useractive set type = ?,uname = ?,uid = ?,itemname = ?,itemid = ?,content = ?,createtime = ?",
		ActiveList[types], uname, uid, itemname, itemid, content, nowTime)
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

// UpdateHealthActive 用户蚕宝宝健康值下降动态
func UpdateHealthActive(updateListIndex *[]int, list *[]map[string]string, nowTime string) bool {
	mysqlConn := common.GetMysqlConn()
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	commit := true
	for i := 0; i < len(*updateListIndex); i++ {
		content := fmt.Sprintf("%s%s", (*list)[(*updateListIndex)[i]]["name"], ActiveStrList[ActiveSubHealth])
		_, err := tx.Insert("insert into useractive set type = ?,uname = ?,uid = ?,itemname = ?,itemid = ?,content = ?,createtime = ?",
			ActiveList[ActiveSubHealth], (*list)[(*updateListIndex)[i]]["name"], (*list)[(*updateListIndex)[i]]["uid"], "", "0", content, nowTime)
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

// GetActiveLog 获取重要滚动动态
func GetActiveLog(paginaSQL, username string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	username = "%" + username + "%"
	sql := fmt.Sprintf("select * from (select useractive.id,user.name as username,user.avatar,useractive.content,useractive.createtime from useractive left join user on uid = user.id where uid = 63 and content LIKE ('%s')UNION select useractive.id,user.name as username,user.avatar,useractive.content,useractive.createtime from useractive left join user on uid = user.id where `type` in ('pairallow','pairend','bebutterfly','batchvoucher','voucher')) as a ORDER BY id DESC %s", username, paginaSQL)
	return mysqlConn.GetResults(mysql.Statement, sql)
}

// GetActiveLogCount 获取重要滚动动态总数
func GetActiveLogCount(username string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	username = "%" + username + "%"
	sql := fmt.Sprintf("select count(*) from ((select t1.id from useractive as t1 left join user as t2 on t1.uid = t2.id where t1.uid = 63 and t1.content LIKE ('%s') )union(select t1.id  from useractive as t1 left join user as t2 on t1.uid = t2.id where t1.`type` in ('pairallow','pairend','bebutterfly','batchvoucher','voucher'))) as t5", username)
	return mysqlConn.GetVal(mysql.Statement, sql)
}
