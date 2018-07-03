package silkworm

import (
	"canbaobao/common"
	"fmt"
	"github.com/JimYJ/easysql/mysql"
)

// AddUser 新增未激活用户
func AddUser(avatar, name, loginip, openid, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into user set avatar = ?,name = ?,treelevel = ?,level = ?,loginip = ?,openid = ?,logintime = ?,createtime = ?,updatetime = ?,enabled = ?",
		avatar, name, 0, 0, loginip, openid, nowTime, nowTime, nowTime, 0)
}

// UserBindVendor 绑定店铺，激活用户
func UserBindVendor(id, vid, loginip, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update user set vid = ?,loginip = ?,logintime = ?,updatetime = ?,enabled = ? where id = ?", vid, loginip, nowTime, nowTime, 1, id)
}

// UserUpgrade 用户升级
func UserUpgrade(id, level, loginip, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update user set level = ?,loginip = ?,logintime = ?,updatetime = ? where id = ?", level, loginip, nowTime, nowTime, id)
}

// UserTreeUpgrade 用户桑树升级
func UserTreeUpgrade(id, treelevel, loginip, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update user set treelevel = ?,loginip = ?,logintime = ?,updatetime = ? where id = ?", treelevel, loginip, nowTime, nowTime, id)
}

// GetUser 获取用户
func getUser() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,avatar,vid,treelevel,level,loginip,logintime,createtime from user where enabled = ? ORDER BY id desc", 1)
}

// GetUser 获取用户
func GetUser() ([]map[string]string, error) {
	list, err := getUser()
	for i := 0; i < len(list); i++ {
		if list[i]["vid"] == "0" {
			list[i]["vendor"] = "数据错误"
		} else {
			if !common.CheckInt(list[i]["vendorid"]) {
				list[i]["vendor"] = "绑定错误"
			} else {
				vname, err := GetVendorName(list[i]["vid"])
				if err != nil || vname == "" {
					list[i]["vendor"] = "数据错误"
				} else {
					list[i]["vendor"] = vname
				}
			}
		}
	}
	return list, err
}

// GetSingleUserByID 获取单个用户
func GetSingleUserByID(id string) (map[string]string, error) {
	return getSingleUser(id, false)
}

// GetSingleUserByOpenID 获取单个用户
func GetSingleUserByOpenID(id string) (map[string]string, error) {
	return getSingleUser(id, true)
}

func getSingleUser(id string, isOpenID bool) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	var field string
	if isOpenID {
		field = "openid"
	} else {
		field = "id"
	}
	sql := fmt.Sprintf("select id,name,avatar,vid,treelevel,level,loginip,logintime,createtime from user where %s = ? and enabled = ?", field)
	return mysqlConn.GetRow(mysql.Statement, sql, id, 1)
}

// GetUserName 获取用户名
func GetUserName(id string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select name from user where id = ?", id)
}

// GetFriendList 获取好友列表
func GetFriendList(openid string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	vid, _ := mysqlConn.GetVal(mysql.Statement, "select vid from user where openid = ?", openid)
	return mysqlConn.GetResults(mysql.Statement, "select id,name,avatar,level from user where vid = ? and openid != ?", vid, openid)
}

// GetUserAnswers 获取用户当日回答次数
func GetUserAnswers(openid string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select answers,answerdate from user where openid = ?", openid)
}

// UpdateUserAnswer 更新用户获得的问题数
func UpdateUserAnswer(answers, nowDate, loginip, nowTime, openid string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update user set answers = ?,answerdate = ?,loginip = ?,logintime = ?,updatetime = ? where openid = ?", answers, nowDate, loginip, nowTime, nowTime, openid)
}

// GetUID 获取用户ID
func GetUID(openid string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select id,name from user where openid = ?", openid)
}

// GetWaterFertilizer 获取浇水施肥数量
func GetWaterFertilizer(openid string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select id,name,treewater,treefertilizer,todaywater,todayfertilizer,waterdate,fertilizerdate,treelevel from user where openid = ?", openid)
}

// UpdateWater 更新浇水时间和次数
func UpdateWater(treewater int, waterDate, loginip, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update user set treewater = ?,todaywater = ?,waterdate = ?,loginip = ?,logintime = ?,updatetime = ? where id = ?", treewater, 1, waterDate, loginip, nowTime, nowTime, id)
}

// UpdateFertilizer 更新施肥时间和次数
func UpdateFertilizer(treefertilizer int, fertilizerDate, loginip, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update user set treefertilizer = ?,todayfertilizer = ?,fertilizerdate = ?,loginip = ?,logintime = ?,updatetime = ? where id = ?", treefertilizer, 1, fertilizerDate, loginip, nowTime, nowTime, id)
}

// UpgradeTree 升级桑树
func UpgradeTree(treeLevel int, todayWater, todayFertilizer, waterDate, fertilizerDate, loginip, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update user set treelevel = ?,treewater = ?,treefertilizer = ?,todaywater = ?,todayfertilizer = ?,waterdate = ?,fertilizerdate = ?,loginip = ?,logintime = ?,updatetime = ? where id = ?",
		treeLevel, 0, 0, todayWater, todayFertilizer, waterDate, fertilizerDate, loginip, nowTime, nowTime, id)
}

// GetSignedDays 获取连续签到日期
func GetSignedDays(openid string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select signdate,lastsigndate from user where openid = ?", openid)
}

// Signed 每日签到
func Signed(openid, signdate, lastsigndate, loginip, logintime, updatetime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update user set signdate = ?,lastsigndate = ?,loginip = ?,logintime = ?,updatetime = ? where openid = ?", signdate, lastsigndate, loginip, logintime, updatetime, openid)
}
