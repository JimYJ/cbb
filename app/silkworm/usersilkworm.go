package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/route/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
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
		middleware.RespondErr(402, common.Err402UserItemNoExist, c)
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
				list[i]["enabletime"] = "0"
				silkworm.Enable(list[i]["id"])
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
