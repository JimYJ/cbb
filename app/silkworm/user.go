package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

// User 用户列表
func User(c *gin.Context) {
	list, _ := silkworm.GetUser()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "user.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// UserInfoByID 获取用户信息
func UserInfoByID(c *gin.Context) {
	id := c.PostForm("id")
	var uinfo map[string]string
	if id == "" {
		middleware.RespondErr(common.HTTPParamErr, common.Err402Param, c)
		return
	}
	uinfo, _ = silkworm.GetSingleUserByID(id)
	c.JSON(200, gin.H{
		"msg":      "success",
		"userinfo": uinfo,
	})
}

// UserInfoByOpenID 获取用户信息
func UserInfoByOpenID(c *gin.Context) {
	openid := c.PostForm("openid")
	var uinfo map[string]string
	if openid == "" {
		middleware.RespondErr(common.HTTPParamErr, common.Err402Param, c)
		return
	}
	uinfo, _ = silkworm.GetSingleUserByOpenID(openid)
	c.JSON(200, gin.H{
		"msg":      "success",
		"userinfo": uinfo,
	})
}

// FriendList 获取好友列表
func FriendList(c *gin.Context) {
	openid := c.PostForm("openid")
	if openid == "" {
		middleware.RespondErr(common.HTTPParamErr, common.Err402Param, c)
		return
	}
	list, _ := silkworm.GetFriendList(openid)
	c.JSON(200, gin.H{
		"msg":  "success",
		"list": list,
	})
}

func waterFertilize(c *gin.Context, isWater bool) {
	openid := c.PostForm("openid")
	if openid == "" {
		middleware.RespondErr(common.HTTPParamErr, common.Err402Param, c)
		return
	}
	uInfo, _ := silkworm.GetWaterFertilizer(openid)
	if uInfo == nil {
		middleware.RespondErr(common.HTTPParamErr, common.Err402Param, c)
		return
	}
	treewater, _ := strconv.Atoi(uInfo["treewater"])
	treefertilizer, _ := strconv.Atoi(uInfo["treefertilizer"])
	todayWater, _ := uInfo["todaywater"]
	todayFertilizer, _ := uInfo["todayfertilizer"]
	waterDate := uInfo["waterdate"]
	fertilizerDate := uInfo["fertilizerdate"]
	treeLevel, err := strconv.Atoi(uInfo["treelevel"])
	if err != nil {
		treeLevel = 0
	}
	nowDate := time.Now().Local().Format("2006-01-02")
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	ip := c.ClientIP()
	if isWater {
		check := common.CheckLimit(todayWater, waterDate, nowDate, 999)
		if check == -1 {
			middleware.RespondErr(200, common.Err201Limit, c)
			return
		}
		if treewater >= 15 {
			middleware.RespondErr(202, common.Err202Limit, c)
			return
		}
		if treewater+1 >= 15 && treefertilizer >= 15 {
			rs := upgradeTree(treeLevel, "1", todayFertilizer, nowDate, fertilizerDate, ip, nowTime, uInfo)
			if !rs {
				middleware.RespondErr(200, common.Err500DBSave, c)
				return
			}
			responSuccess(c)
		} else {
			_, err := silkworm.UpdateWater(treewater+1, nowDate, ip, nowTime, uInfo["id"])
			if err != nil {
				log.Println("water for tree fail:", err)
				middleware.RespondErr(200, common.Err500DBSave, c)
				return
			}
			responSuccess(c)
		}
	} else {
		check := common.CheckLimit(todayFertilizer, fertilizerDate, nowDate, 999)
		if check == -1 {
			middleware.RespondErr(200, common.Err201Limit, c)
			return
		}
		if treefertilizer >= 15 {
			middleware.RespondErr(202, common.Err202Limit, c)
			return
		}
		if treefertilizer+1 >= 15 && treewater >= 15 {
			rs := upgradeTree(treeLevel, todayWater, "1", waterDate, nowDate, ip, nowTime, uInfo)
			if !rs {
				middleware.RespondErr(200, common.Err500DBSave, c)
				return
			}
			responSuccess(c)
		} else {
			_, err := silkworm.UpdateFertilizer(treefertilizer+1, nowDate, ip, nowTime, uInfo["id"])
			if err != nil {
				log.Println("water for tree fail:", err)
				middleware.RespondErr(200, common.Err500DBSave, c)
				return
			}
			responSuccess(c)
		}
	}
}

// WaterForTree 浇水
func WaterForTree(c *gin.Context) {
	waterFertilize(c, true)
}

// FertilizerForTree 施肥
func FertilizerForTree(c *gin.Context) {
	waterFertilize(c, false)
}

// UpgradeTree 升级用户桑树
func upgradeTree(treeLevel int, todayWater, todayFertilizer, waterDate, fertilizerDate, loginip, nowTime string, uInfo map[string]string) bool {
	treeLevel++
	_, err := silkworm.UpgradeTree(treeLevel, todayWater, todayFertilizer, waterDate, fertilizerDate, loginip, nowTime, uInfo["id"])
	if err != nil {
		log.Println("Upgrade Tree Fail:", err)
		return false
	}
	_, err = silkworm.SaveUserActive(silkworm.ActiveTreeup, uInfo["name"], uInfo["id"], "", "0", nowTime, strconv.Itoa(treeLevel))
	if err != nil {
		log.Println("Save User Active Fail:", err)
	}
	return true
}

// GetSignedDays 获取本周签到天数
func GetSignedDays(c *gin.Context) {
	lastSignedDate, signDate, err := getUserSignDate(c)
	if err != nil {
		return
	}
	if lastSignedDate == "" || signDate == "" {
		c.JSON(200, gin.H{
			"msg":        "success",
			"signedDays": 0,
		})
		return
	}
	signedDays := calcSignedDays(lastSignedDate, signDate)
	c.JSON(200, gin.H{
		"msg":        "success",
		"signedDays": signedDays,
	})
}

// 计算本周连续签到天数
func calcSignedDays(lastSignedDate, signDate string) time.Duration {
	day, _ := time.ParseDuration("24h")
	nowDate, _ := time.Parse("2006-01-02", common.GetThisWeekMonday())
	lastDate, _ := time.Parse("2006-01-02", lastSignedDate)
	sign, _ := time.Parse("2006-01-02", signDate)
	s := sign.Sub(lastDate)
	var signedDays time.Duration
	if lastDate.Sub(nowDate) < 0 {
		signedDays = 0
	} else if s < 0 {
		signedDays = 0
	} else if s == 0 {
		signedDays = 1
	} else {
		signedDays = s/day + 1
	}
	return signedDays
}

// UserSigned 每日签到
func UserSigned(c *gin.Context) {
	openid := c.PostForm("openid")
	lastSignedDate, signDate, err := getUserSignDate(c)
	if err != nil {
		return
	}
	nowDate := time.Now().Local().Format("2006-01-02")
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	ip := c.ClientIP()
	if lastSignedDate == "" || signDate == "" {
		_, err := silkworm.Signed(openid, nowDate, nowDate, ip, nowTime, nowTime)
		if err != nil {
			middleware.RespondErr(200, common.Err500DBSave, c)
			return
		}
		go handelSignAward(openid, nowTime, true, 1)
		responSignSuccess(c, 1)
		return
	}
	ThisWeekMonday, _ := time.Parse("2006-01-02", common.GetThisWeekMonday())
	lastDate, _ := time.Parse("2006-01-02", lastSignedDate)
	sign, _ := time.Parse("2006-01-02", signDate)
	toDay, _ := time.Parse("2006-01-02", nowDate)
	tSub := toDay.Sub(sign)
	day, _ := time.ParseDuration("24h")
	daySub := tSub / day
	if lastDate.Sub(ThisWeekMonday) < 0 {
		_, err := silkworm.Signed(openid, nowDate, nowDate, ip, nowTime, nowTime)
		if err != nil {
			log.Println(err)
			middleware.RespondErr(500, common.Err500DBSave, c)
			return
		}
		go handelSignAward(openid, nowTime, true, 1)
		responSignSuccess(c, 1)
		return
	} else if daySub == 0 {
		middleware.RespondErr(201, common.Err201Limit, c)
		return
	} else if daySub == 1 {
		_, err := silkworm.Signed(openid, nowDate, lastDate.Format("2006-01-02"), ip, nowTime, nowTime)
		if err != nil {
			log.Println(err)
			middleware.RespondErr(500, common.Err500DBSave, c)
			return
		}
		signedDays := calcSignedDays(lastSignedDate, nowDate)
		go handelSignAward(openid, nowTime, true, int64(signedDays))
		if signedDays >= 7 {
			go handelSignAward(openid, nowTime, false, int64(signedDays))
		}
		responSignSuccess(c, int64(signedDays))
		return
	} else if daySub > 1 {
		_, err := silkworm.Signed(openid, nowDate, nowDate, ip, nowTime, nowTime)
		if err != nil {
			log.Println(err)
			middleware.RespondErr(500, common.Err500DBSave, c)
			return
		}
		go handelSignAward(openid, nowTime, true, 1)
		responSignSuccess(c, 1)
		return
	}
	middleware.RespondErr(500, common.Err500DBrequest, c)
}

// 获得签到记录
func getUserSignDate(c *gin.Context) (string, string, error) {
	openid := c.PostForm("openid")
	if openid == "" {
		middleware.RespondErr(common.HTTPParamErr, common.Err402Param, c)
		return "", "", errors.New(common.Err402Param)
	}
	rs, err := silkworm.GetSignedDays(openid)
	lastSignedDate := rs["lastsigndate"]
	signDate := rs["signdate"]
	if err != nil {
		log.Println(err)
		middleware.RespondErr(common.HTTPParamErr, common.Err402Param, c)
		return "", "", err
	}
	return lastSignedDate, signDate, nil
}

func handelSignAward(openid, nowTime string, isDay bool, signedDays int64) {
	dayitemid, weekitemid := silkworm.SignedItem()
	var itemInfo map[string]string
	var itemid string
	if isDay {
		itemInfo, _ = silkworm.ItemInfo(dayitemid)
		itemid = dayitemid

	} else {
		itemInfo, _ = silkworm.ItemInfo(weekitemid)
		itemid = weekitemid
	}
	moreInfo := fmt.Sprintf("%d天", signedDays)
	AddItemToRucksack(silkworm.ActiveSign, openid, itemid, nowTime, moreInfo, itemInfo)
}