package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/route/middleware"
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
	log.Println(nums, rucksackid, err)
	if nums == "0" {
		log.Println("get rucksackid Fail", err)
		middleware.RespondErr(402, common.Err402UserItemNoExist, c)
		return
	}
	uinfo, err := silkworm.GetUID(openid)
	uid := uinfo["id"]
	uname := uinfo["name"]
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	rsults := silkworm.Hatch(uid, rucksackid, swtype)
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
