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
	enablehous, _ := time.ParseDuration("12h")
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
	for i := 0; i < len(list); i++ {
		if list[i]["hatch"] == "0" {
			list[i]["dialog"] = silkworm.GetRandomDialog()
			for j := 0; j < len(levelList); j++ {
				if list[i]["level"] == levelList[j]["level"] {
					list[i]["levelExp"] = levelList[j]["exp"]
					list[i]["img"] = levelList[j]["img"]
				}
				levelInt, _ := strconv.Atoi(list[i]["level"])
				upLevel := fmt.Sprintf("%d", levelInt+1)
				if upLevel == levelList[j]["level"] {
					list[i]["nextLevelExp"] = levelList[j]["exp"]
				}
				list[i]["expPercent"] = strconv.Itoa(common.CalcExpPercent(list[i]["levelExp"], list[i]["exp"], list[i]["nextLevelExp"]))
			}
		} else {
			for j := 0; j < len(butterflyList); j++ {
				if list[i]["swid"] == butterflyList[j]["id"] {
					list[i]["img"] = butterflyList[j]["img"]
				}
			}
		}
		if list[i]["enable"] == "1" {
			list[i]["enabletime"] = "0"
		} else if list[i]["enable"] == "0" {
			enabletime, _ := strconv.ParseInt(list[i]["enabletime"], 10, 64)
			if enabletime-nowUnix > 0 {
				t, _ := time.ParseDuration(fmt.Sprintf("%ds", enabletime-nowUnix))
				list[i]["enabletime"] = common.FormatTimeGap(t.String())
			} else {
				list[i]["enabletime"] = "0"
				silkworm.Enable(list[i]["id"])
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
