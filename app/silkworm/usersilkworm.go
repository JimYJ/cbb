package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/route/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// 喂养物品类型
const (
	null = iota
	Leaf
	SPowerPack
	LPowerPack
	MPowerPack
)

// HatchForNormal 普通蚕仔孵化
func HatchForNormal(c *gin.Context) {
	hatch(c, false)
}

// HatchForSpecial 特殊蚕仔孵化
func HatchForSpecial(c *gin.Context) {
	hatch(c, true)
}

// 孵化
func hatch(c *gin.Context, isSpecial bool) {
	openid := c.PostForm("openid")
	if openid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	rs, err := silkworm.CheckHatch(openid)
	if err != nil {
		log.Println("CheckHatch Fail", err)
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	rsInt, _ := strconv.Atoi(rs)
	if rsInt > 0 {
		middleware.RespondErr(203, common.Err203Limit, c)
		return
	}
	var swtype int
	if isSpecial {
		swtype = 1
	} else {
		swtype = 0
	}
	nums, rucksackid, err := silkworm.GetUserSWID(openid, swtype)
	// log.Println(nums, rucksackid, err)
	if nums == "0" {
		log.Println("get rucksackid Fail", err)
		middleware.RespondErr(413, common.Err413UserItemNoExist, c)
		return
	}
	uinfo, err := silkworm.GetUID(openid)
	uid := uinfo["id"]
	uname := uinfo["name"]
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	enablehous, _ := time.ParseDuration("24h")
	enabletime := time.Now().Local().Add(enablehous).Unix()
	rsults := silkworm.Hatch(uid, rucksackid, swtype, enabletime)
	if rsults {
		// 记录动态
		_, err = silkworm.SaveUserActive(silkworm.ActiveHatch, uname, uid, "", "-1", nowTime, "")
		if err != nil {
			log.Println("Save User Active Fail", err)
		}
		responSuccess(c)
	} else {
		middleware.RespondErr(500, common.Err500DBSave, c)
	}
}

// 用户的蚕宝宝列表
func userSilkwormList(c *gin.Context, uid string) {
	if uid == "" || !common.CheckInt(uid) {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	list, err := silkworm.GetUserSilkworm(uid)
	if err != nil {
		middleware.RespondErr(402, common.Err500DBrequest, c)
		return
	}
	levelList, _ := silkworm.LevelList()
	butterflyList, _ := silkworm.ButterflyList()
	nowUnix := time.Now().Local().Unix()
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	for i := 0; i < len(list); i++ {
		if list[i]["hatch"] == "0" {
			list[i]["dialog"] = silkworm.GetRandomDialog()
			for j := 0; j < len(levelList); j++ {
				if list[i]["level"] == levelList[j]["level"] {
					list[i]["levelExp"] = levelList[j]["exp"]
					list[i]["img"] = levelList[j]["img"]
					continue
				}
				levelInt, _ := strconv.Atoi(list[i]["level"])
				upLevel := fmt.Sprintf("%d", levelInt+1)
				if upLevel == levelList[j]["level"] {
					list[i]["nextLevelExp"] = levelList[j]["exp"]
					break
				}
			}
			list[i]["expPercent"] = strconv.Itoa(common.CalcExpPercent(list[i]["levelExp"], list[i]["exp"], list[i]["nextLevelExp"]))
		} else {
			for j := 0; j < len(butterflyList); j++ {
				if list[i]["swid"] == butterflyList[j]["id"] {
					list[i]["img"] = butterflyList[j]["img"]
					break
				}
			}
		}
		if list[i]["enable"] == "1" {
			list[i]["enabletime"] = "0"
		} else if list[i]["enable"] == "0" {
			enabletime, _ := strconv.ParseInt(list[i]["enabletime"], 10, 64)
			if enabletime-nowUnix > 0 {
				list[i]["enablestatus"] = "孵化中"
				t, _ := time.ParseDuration(fmt.Sprintf("%ds", enabletime-nowUnix))
				list[i]["enabletime"] = common.FormatTimeGap(t.String())
			} else {
				uinfo, _ := silkworm.GetUinfoByID(uid)
				uname := uinfo["name"]
				level, _ := strconv.Atoi(uinfo["level"])
				vid := uinfo["vid"]
				if vid == "" {
					middleware.RespondErr(412, common.Err412UserNotBind, c)
					return
				}
				list[i]["enabletime"] = "0"
				silkworm.Enable(list[i]["id"])
				if level < 1 {
					linfo, _ := silkworm.GetLevel("1")
					awardItem := linfo["redeemitem"]
					go createVoucher(awardItem, vid, uid, nowTime, uname, "1")
				} else {
					go swUpActive(uname, uid, nowTime, "1")
				}
			}
		}
		if list[i]["pair"] == "0" {
			if list[i]["hatch"] == "1" {
				list[i]["pairstatus"] = "未配对"
				list[i]["pairtime"] = "0"
			}
		} else if list[i]["pair"] == "1" {
			list[i]["pairstatus"] = "配对申请中"
			list[i]["pairtime"] = "0"
		} else if list[i]["pair"] == "2" {
			pairtime, _ := strconv.ParseInt(list[i]["pairtime"], 10, 64)
			if pairtime-nowUnix > 0 {
				t, _ := time.ParseDuration(fmt.Sprintf("%ds", pairtime-nowUnix))
				list[i]["pairtime"] = common.FormatTimeGap(t.String())
				list[i]["pairstatus"] = "配对中"
			} else {
				list[i]["pairtime"] = "0"
				list[i]["pair"] = "0"
				list[i]["pairstatus"] = "未配对"
				list[i]["pairid"] = "0"
				list[i]["pairsrc"] = "0"
				silkworm.EndPair(list[i]["pairid"], list[i]["id"])
				if list[i]["pairsrc"] == "1" {
					go endPairRuck(uid, list[i]["pairuid"], nowTime)
				} else {
					go endPairRuck(list[i]["pairuid"], uid, nowTime)
				}
			}
		}
	}
	c.JSON(200, gin.H{
		"msg":  "success",
		"list": list,
	})
}

// UserSilkwormList 获取用户自己的蚕宝宝列表
func UserSilkwormList(c *gin.Context) {
	openid := c.PostForm("openid")
	if openid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	uinfo, _ := silkworm.GetUID(openid)
	uid := uinfo["id"]
	userSilkwormList(c, uid)
}

// FriendSilkwormList 获取好友的的蚕宝宝列表
func FriendSilkwormList(c *gin.Context) {
	uid := c.PostForm("uid")
	if uid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	userSilkwormList(c, uid)
}

// endPairRuck 结束配对时背包和动态处理
func endPairRuck(uid, pairuid, nowTime string) {
	_, err := silkworm.AddSilkwormRucksack("6", uid, "1", nowTime, 0)
	if err != nil {
		log.Println("Add Item to Rucksack Fail", err)
	}
	uinfo, _ := silkworm.GetSingleUserByID(uid)
	uname := uinfo["name"]
	pairuinfo, _ := silkworm.GetSingleUserByID(pairuid)
	pairuname := pairuinfo["name"]
	itemInfo, _ := silkworm.ItemInfo("6")
	itemName := itemInfo["name"]
	_, err = silkworm.SaveUserActive(silkworm.ActivePairEnd, uname, uid, itemName, "6", nowTime, pairuname)
	if err != nil {
		log.Println("Save User Active Fail", err)
	}
	_, err = silkworm.SaveUserActive(silkworm.ActivePairEndII, pairuname, pairuid, "", "0", nowTime, uname)
	if err != nil {
		log.Println("Save Pair User Active Fail", err)
	}
}

//NewUserRuck 新用户注册背包和动态
func NewUserRuck(openid, nowTime string) {
	uinfo, _ := silkworm.GetUID(openid)
	uid := uinfo["id"]
	_, err := silkworm.AddSilkwormRucksack("5", uid, "0", nowTime, 0)
	if err != nil {
		log.Println("Add Item to Rucksack Fail", err)
	}
	uname := uinfo["name"]
	itemInfo, _ := silkworm.ItemInfo("5")
	itemName := itemInfo["name"]
	_, err = silkworm.SaveUserActive(silkworm.ActiveNewUser, uname, uid, itemName, "5", nowTime, "")
	if err != nil {
		log.Println("Save User Active Fail", err)
	}
}

// ApplyPair 申请配对
func ApplyPair(c *gin.Context) {
	openid := c.PostForm("openid")
	if openid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	uinfo, _ := silkworm.GetUID(openid)
	uid := uinfo["id"]
	pairuid := c.PostForm("pairuid")
	id := c.PostForm("id")
	pairid := c.PostForm("pairid")
	if uid == "" || pairuid == "" || id == "" || pairid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	if pairid == id || uid == pairuid {
		middleware.RespondErr(415, common.Err415CannotPairSelf, c)
		return
	}
	hatch, pair1, _ := silkworm.CheckPairCondition(id)
	pairhatch, pair2, _ := silkworm.CheckPairCondition(pairid)
	if hatch == "" || pair1 == "" || pairhatch == "" || pair2 == "" {
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	if hatch != "1" || pairhatch != "1" {
		middleware.RespondErr(205, common.Err205Limit, c)
		return
	}
	if pair1 != "0" || pair2 != "0" {
		middleware.RespondErr(206, common.Err206Limit, c)
		return
	}
	rs := silkworm.ApplyPair(id, pairid, uid, pairuid)
	if !rs {
		middleware.RespondErr(500, common.Err500DBSave, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	go userActiveForApply(pairuid, uid, nowTime)
	responSuccess(c)
}

// AllowPair 同意配对
func AllowPair(c *gin.Context) {
	handlePair(c, false)
}

// RejectPair 拒绝配对
func RejectPair(c *gin.Context) {
	handlePair(c, true)
}

// 处理配对
func handlePair(c *gin.Context, isReject bool) {
	id := c.PostForm("id")
	pairid := c.PostForm("pairid")
	if id == "" || pairid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	hatch, pair1, uid := silkworm.CheckPairCondition(id)
	pairhatch, pair2, pairuid := silkworm.CheckPairCondition(pairid)
	if hatch == "" || pair1 == "" || pairhatch == "" || pair2 == "" {
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	if hatch != "1" || pairhatch != "1" {
		middleware.RespondErr(205, common.Err205Limit, c)
		return
	}
	if pair1 != "1" || pair2 != "1" {
		middleware.RespondErr(206, common.Err206Limit, c)
		return
	}
	var rs bool
	if isReject {
		rs = silkworm.EndPair(id, pairid)
	} else {
		pairhous, _ := time.ParseDuration("24h")
		pairtime := time.Now().Local().Add(pairhous).Unix()
		rs = silkworm.AllowPair(id, pairid, pairtime)
	}
	if !rs {
		middleware.RespondErr(500, common.Err500DBSave, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isReject {
		go userActiveForReject(pairuid, uid, nowTime)
	} else {
		go userActiveForAllow(pairuid, uid, nowTime)
	}
	responSuccess(c)
}

func userActiveForApply(pairuid, uid, nowTime string) {
	uname, pairuname := getUserInfo(uid, pairuid)
	_, err := silkworm.SaveUserActive(silkworm.ActivePairApply, uname, uid, "", "0", nowTime, pairuname)
	if err != nil {
		log.Println("Save User Active Fail", err)
	}
	_, err = silkworm.SaveUserActive(silkworm.ActivePairApplyed, pairuname, pairuid, "", "0", nowTime, uname)
	if err != nil {
		log.Println("Save Pair User Active Fail", err)
	}
}

func userActiveForAllow(pairuid, uid, nowTime string) {
	uname, pairuname := getUserInfo(uid, pairuid)
	_, err := silkworm.SaveUserActive(silkworm.ActivePairAllowII, uname, uid, "", "0", nowTime, pairuname)
	if err != nil {
		log.Println("Save Pair User Active Fail", err)
	}
	_, err = silkworm.SaveUserActive(silkworm.ActivePairAllow, pairuname, pairuid, "", "0", nowTime, uname)
	if err != nil {
		log.Println("Save User Active Fail", err)
	}
}

func userActiveForReject(pairuid, uid, nowTime string) {
	uname, pairuname := getUserInfo(uid, pairuid)
	_, err := silkworm.SaveUserActive(silkworm.ActivePairRejectII, uname, uid, "", "0", nowTime, pairuname)
	if err != nil {
		log.Println("Save Pair User Active Fail", err)
	}
	_, err = silkworm.SaveUserActive(silkworm.ActivePairReject, pairuname, pairuid, "", "0", nowTime, uname)
	if err != nil {
		log.Println("Save User Active Fail", err)
	}
}

func getUserInfo(uid, pairuid string) (string, string) {
	uinfo, _ := silkworm.GetSingleUserByID(uid)
	uname := uinfo["name"]
	pairuinfo, _ := silkworm.GetSingleUserByID(pairuid)
	pairuname := pairuinfo["name"]
	return uname, pairuname
}

// Feed 喂食
func Feed(c *gin.Context) {
	openid := c.PostForm("openid")
	itemid := c.PostForm("itemid")
	id := c.PostForm("id")
	if openid == "" || itemid == "" || id == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	uinfo, err := silkworm.GetUID(openid)
	uid := uinfo["id"]
	rucksackItemInfo, err3 := silkworm.RucksackItemInfo(itemid, uid)
	userSilkwormInfo, err2 := silkworm.GetSingleUserSWInfo(id)
	if err != nil || err3 != nil || err2 != nil {
		log.Println(err, err2, err3)
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	if uinfo == nil {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	if rucksackItemInfo == nil {
		middleware.RespondErr(413, common.Err413UserItemNoExist, c)
		return
	}
	take := rucksackItemInfo["take"]
	rucksackid := rucksackItemInfo["id"]
	if take == "0" {
		middleware.RespondErr(416, common.Err416NotTaken, c)
		return
	}
	iteminfo, err := silkworm.ItemInfo(itemid)
	if err != nil {
		log.Println(err)
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	if iteminfo == nil {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	if uid != userSilkwormInfo["uid"] {
		middleware.RespondErr(418, common.Err418OtherUserSW, c)
		return
	}
	userLevel, _ := strconv.Atoi(uinfo["level"])
	itemExp, _ := strconv.Atoi(iteminfo["exp"])
	limitday, _ := strconv.Atoi(iteminfo["limitday"])
	itemTypes := iteminfo["types"]
	hatch := userSilkwormInfo["hatch"]
	enable := userSilkwormInfo["enable"]
	level, _ := strconv.Atoi(userSilkwormInfo["level"])
	silkwormExp, _ := strconv.Atoi(userSilkwormInfo["exp"])
	swtype := userSilkwormInfo["swtype"]
	if itemTypes != "1" && itemTypes != "0" {
		middleware.RespondErr(419, common.Err419ItemCannotFeed, c)
		return
	}
	if hatch != "0" {
		middleware.RespondErr(420, common.Err420CannotFeedButterfly, c)
		return
	}
	if enable != "1" {
		middleware.RespondErr(421, common.Err421SilkwormHatching, c)
		return
	}
	feedTimes, keyTimes, keyDate := checkItemDayLimit(itemid, uid, limitday)
	if feedTimes == -1 {
		middleware.RespondErr(201, common.Err201Limit, c)
		return
	}
	if level >= 10 {
		log.Println("usersilkworm level err!", level)
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	silkwormLevel, _ := silkworm.LevelList()
	newlevel := level
	newUserLevel := userLevel
	vid := uinfo["vid"]
	uname := uinfo["name"]
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	nowDate := time.Now().Local().Format("2006-01-02")
	ip := c.ClientIP()
	for i := 0; i < 11; i++ {
		nextLevelExp, _ := strconv.Atoi(silkwormLevel[newlevel]["exp"])
		if itemExp+silkwormExp < nextLevelExp {
			// 增加经验值
			rs := silkworm.UpExp(strconv.Itoa(itemExp+silkwormExp), strconv.Itoa(newlevel), silkwormLevel[newlevel-1]["name"], id, uid, rucksackid, keyTimes, keyDate, ip, nowTime, nowDate, feedTimes)
			if !rs {
				middleware.RespondErr(500, common.Err500DBSave, c)
				return
			}
			break
		} else {
			newlevel++
			if userLevel < 1 {
				// 发送兑换券
				if vid == "" {
					middleware.RespondErr(412, common.Err412UserNotBind, c)
					return
				}
				awardItem := silkwormLevel[newlevel-1]["redeemitem"]
				go createVoucher(awardItem, vid, uid, nowTime, uname, strconv.Itoa(newlevel))
			} else {
				go swUpActive(uname, uid, nowTime, strconv.Itoa(newlevel))
			}
			if newlevel >= 10 {
				//成为蝴蝶,用户等级+1
				newUserLevel++
				var randBFid int
				rand.Seed(time.Now().UnixNano())
				if swtype == "0" {
					randBFid = rand.Intn(5) + 1
				} else {
					randBFid = rand.Intn(5) + 6
				}
				bfname, _ := silkworm.ButterflyName(randBFid)
				rs := silkworm.BeButterfly(strconv.Itoa(itemExp+silkwormExp), bfname, strconv.Itoa(randBFid), id, uid, rucksackid, strconv.Itoa(newUserLevel), ip, nowTime, nowDate, keyTimes, keyDate, feedTimes)
				if !rs {
					middleware.RespondErr(500, common.Err500DBSave, c)
					return
				}
				go beButterflyActive(uname, uid, nowTime, strconv.Itoa(newUserLevel))
				break
			}
		}
	}
	if newlevel > level && newUserLevel > userLevel {
		c.JSON(200, gin.H{
			"msg":          "success",
			"newlevel":     newlevel,
			"newUserLevel": newUserLevel,
		})
	} else if newlevel > level {
		c.JSON(200, gin.H{
			"msg":      "success",
			"newlevel": newlevel,
		})
	} else {
		c.JSON(200, gin.H{
			"msg": "success",
		})
	}
}

func checkItemDayLimit(itemid, uid string, limitday int) (int, string, string) {
	itemID, err := strconv.Atoi(itemid)
	if err != nil {
		log.Println(err)
		return -1, "", ""
	}
	feedInfo, err := silkworm.GetUserFeeds(uid)
	if err != nil {
		log.Println(err)
		return -1, "", ""
	}
	var feedTimes, feedDate, keyTimes, keyDate string
	if itemID == Leaf {
		feedTimes = feedInfo["leafusetoday"]
		feedDate = feedInfo["leafday"]
		keyTimes = "leafusetoday"
		keyDate = "leafday"
	} else if itemID == SPowerPack {
		feedTimes = feedInfo["sppusetoday"]
		feedDate = feedInfo["sppday"]
		keyTimes = "sppusetoday"
		keyDate = "sppday"
	} else if itemID == LPowerPack {
		feedTimes = feedInfo["mppusetoday"]
		feedDate = feedInfo["mppday"]
		keyTimes = "mppusetoday"
		keyDate = "mppday"
	} else if itemID == MPowerPack {
		feedTimes = feedInfo["lppusetoday"]
		feedDate = feedInfo["lppday"]
		keyTimes = "lppusetoday"
		keyDate = "lppday"
	}
	nowDate := time.Now().Local().Format("2006-01-02")
	return common.CheckLimit(feedTimes, feedDate, nowDate, limitday), keyTimes, keyDate
}

// createVoucher 生成兑换券及动态
func createVoucher(awardItem, vid, uid, nowTime, uname, moreInfo string) {
	_, err := silkworm.AddVoucher(vid, uid, awardItem, nowTime)
	if err != nil {
		log.Println("Create Voucher Fail:", err)
	}
	_, err = silkworm.SaveUserActive(silkworm.ActiveFirstSWUp, uname, uid, awardItem, "0", nowTime, moreInfo)
	if err != nil {
		log.Println("Save User Active Fail:", err)
	}
}

// swUpActive 蚕宝宝升级动态
func swUpActive(uname, uid, nowTime, moreInfo string) {
	_, err := silkworm.SaveUserActive(silkworm.ActiveSWUp, uname, uid, "", "0", nowTime, moreInfo)
	if err != nil {
		log.Println("Save User Active Fail:", err)
	}
}

// beButterflyActive 化蝶动态
func beButterflyActive(uname, uid, nowTime, moreInfo string) {
	_, err := silkworm.SaveUserActive(silkworm.ActiveBeButterfly, uname, uid, "", "0", nowTime, moreInfo)
	if err != nil {
		log.Println("Save User Active Fail:", err)
	}
}
