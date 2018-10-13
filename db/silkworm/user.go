package silkworm

import (
	"canbaobao/common"
	"fmt"
	"github.com/JimYJ/easysql/mysql"
	"strconv"
)

// AddUser 新增未激活用户
func AddUser(avatar, name, loginip, openid, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	tx, _ := mysqlConn.Begin()
	userId, _ := mysqlConn.GetVal(mysql.Statement, "select id from litemall_user where weixin_openid = ?", openid)
	var err error
	if userId == "" {
		result, err := tx.Insert("insert into litemall_user set avatar = ?,username = ?,nickname = ?,last_login_ip = ?,weixin_openid = ?,last_login_time = ?,add_time = ?", avatar, name, name, loginip, openid, nowTime, nowTime)
		if err != nil {
			return result, err
		}
	}
	result, err := tx.Insert("insert into user set avatar = ?,name = ?,treelevel = ?,level = ?,loginip = ?,openid = ?,logintime = ?,createtime = ?,updatetime = ?,enabled = ?",
		avatar, name, 1, 0, loginip, openid, nowTime, nowTime, nowTime, 0)
	if err != nil {
		tx.Rollback()
		return result, err
	}
	tx.Commit()
	return result, err
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
	return mysqlConn.GetResults(mysql.Statement, "select id,name,avatar,vid,treelevel,level,loginip,logintime,createtime from user ORDER BY id desc")
}

// GetUserForTimer 获取用户(定时任务)
func GetUserForTimer() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,sproutleafs,sproutleafday,treelevel from user ORDER BY id desc")
}

// GetUserList 获取用户列表
func GetUserList() (map[string]string, map[string]string, map[string]string) {
	list, _ := GetUser()
	unameList := make(map[string]string)
	ulevelList := make(map[string]string)
	vendorList := make(map[string]string)
	for _, v := range list {
		unameList[v["id"]] = v["name"]
		ulevelList[v["id"]] = v["level"]
		vendorList[v["id"]] = v["vendor"]
	}
	return unameList, ulevelList, vendorList
}

// GetUser 获取用户
func GetUser() ([]map[string]string, error) {
	list, err := getUser()
	for i := 0; i < len(list); i++ {
		if list[i]["vid"] == "" {
			list[i]["vendor"] = "未绑定"
		} else if list[i]["vid"] == "0" {
			list[i]["vendor"] = "数据错误"
		} else {
			if !common.CheckInt(list[i]["vendorid"]) {
				list[i]["vendor"] = "绑定错误"
			} else {
				vname, err := GetVendorName(list[i]["vid"])
				if err != nil || vname == "" {
					list[i]["vendor"] = "未绑定"
				} else {
					list[i]["vendor"] = vname
				}
			}
		}
	}
	return list, err
}

// GetUserVid 获取用户VID
func GetUserVid(openid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select vid from user where openid = ?", openid)
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
	sql := fmt.Sprintf("select id,name,avatar,vid,treelevel,level,loginip,logintime,createtime from user where %s = ?", field)
	return mysqlConn.GetRow(mysql.Statement, sql, id)
}

// GetUserName 获取用户名
func GetUserName(id string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select name from user where id = ?", id)
}

// GetFriendList 获取好友列表
func GetFriendList(openid, vid, paginaSQL string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	sql := fmt.Sprintf("select user.id,user.name,avatar,user.level,exp from user LEFT JOIN usersw ON uid = user.id where vid = ? and openid != ? and hatch = 0 order by level desc,exp desc %s", paginaSQL)
	return mysqlConn.GetResults(mysql.Statement, sql, vid, openid)
}

// GetFriendCount 获取好友总数
func GetFriendCount(openid, vid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from user where vid = ? and openid != ? ", vid, openid)
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
	return mysqlConn.GetRow(mysql.Statement, "select id,name,vid,level from user where openid = ?", openid)
}

// GetUinfoByID 获取用户信息
func GetUinfoByID(id string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select id,name,vid,level from user where id = ?", id)
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

// BillBoard 全国排行榜
func BillBoard(paginaSQL string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	sql := fmt.Sprintf("select user.id,user.name,avatar,treelevel,user.level,exp from user LEFT JOIN usersw ON uid = user.id AND usersw.hatch = 0 ORDER BY level desc,exp desc %s", paginaSQL)
	return mysqlConn.GetResults(mysql.Statement, sql)
}

// GetUserCount 获取用户总数
func GetUserCount() (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from user")
}

// GetUserFeeds 获取用户当日喂养次数
func GetUserFeeds(id string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select leafusetoday,sppusetoday,mppusetoday,lppusetoday,leafday,sppday,mppday,lppday from user where id = ?", id)
}

// GetIntroPageRecord 获取用户是否看过引导页
func GetIntroPageRecord(openid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select intropage from user where openid = ?", openid)
}

// UpdateIntroPageRecord 更新用户是否看过引导页
func UpdateIntroPageRecord(openid string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Update(mysql.Statement, "update user set intropage = ? where openid = ?", 1, openid)
}

// CheckUserExist 判断用户是否存在
func CheckUserExist(openid string) (int, error) {
	mysqlConn := common.GetMysqlConn()
	rs, _ := mysqlConn.GetVal(mysql.Statement, "select count(*) from user where openid = ?", openid)
	count, err := strconv.Atoi(rs)
	return count, err
}

// GetUserForAreaVendor 获取指定地区的店铺的绑定用户
func GetUserForAreaVendor(province, city string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	if len(city) == 0 {
		return mysqlConn.GetResults(mysql.Statement, "select `user`.id,`user`.`name`,`user`.vid FROM litemall_vendor vendor LEFT JOIN `user` on vendor.id = `user`.vid WHERE province = ? and `user`.id is not NULL", province)
	}
	return mysqlConn.GetResults(mysql.Statement, "select `user`.id,`user`.`name`,`user`.vid FROM litemall_vendor vendor LEFT JOIN `user` on vendor.id = `user`.vid WHERE province = ? and city = ? and `user`.id is not NULL", province, city)
}

// UserInviteAwardLog 发放邀请用户奖励记录
func UserInviteAwardLog(uid, iuid, itemid, num, nowTime string) error {
	mysqlConn := common.GetMysqlConn()
	_, err := mysqlConn.Insert(mysql.Statement, "INSERT INTO invitelogs SET uid = ?,iuid = ?,itemid = ?,num = ?,createtime = ?", uid, iuid, itemid, num, nowTime)
	if err != nil {
		return err
	}
	return nil
}

// GetUserAwardLog 获取奖励记录
func GetUserAwardLog(uid string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select invitelogs.id,user.name as uname,item.name as itemname,num,invitelogs.createtime from invitelogs left join user on iuid = user.id left join item on itemid = item.id where uid = ? and `read` = ? order by invitelogs.id desc", uid, 0)
}

// ReadAwardLog 标记记录已读
func ReadAwardLog(list []map[string]string) error {
	mysqlConn := common.GetMysqlConn()
	tx, _ := mysqlConn.Begin()
	var err error
	for _, v := range list {
		_, err = tx.Update("update invitelogs set `read` = ? where id = ?", 1, v["id"])
		if err != nil {
			break
		}
	}
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// UserInviteList 用户邀请统计表
func UserInviteList(startDay, endDay string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	if len(startDay) == 0 || len(endDay) == 0 {
		return mysqlConn.GetResults(mysql.Statement, "SELECT count(*) as num,uid FROM invitelogs GROUP BY uid ORDER BY num desc")
	}
	return mysqlConn.GetResults(mysql.Statement, "SELECT count(*) as num,uid FROM invitelogs where createtime < ? and createtime > ? GROUP BY uid ORDER BY num desc", endDay, startDay)
}
