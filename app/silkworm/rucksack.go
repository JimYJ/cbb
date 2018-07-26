package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/route/middleware"
	log "canbaobao/service/logs"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// UserRucksack 获取用户背包列表
func UserRucksack(c *gin.Context) {
	openid := c.PostForm("openid")
	if openid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	list, err := silkworm.UserRucksack(openid, true)
	for i := 0; i < len(list); i++ {
		iteminfo, _ := silkworm.ItemInfo(list[i]["itemid"])
		list[i]["itemname"] = iteminfo["name"]
		list[i]["itemimg"] = iteminfo["img"]
		list[i]["itemexp"] = iteminfo["exp"]
	}
	if err != nil {
		log.Println(err)
	}
	c.JSON(200, gin.H{
		"msg":  "success",
		"list": list,
	})
}

// AddItemToRucksack 普通物品进背包/记录动态
func AddItemToRucksack(activeType int, openid, itemid, nowTime, moreInfo string, itemInfo map[string]string) {
	itemname := itemInfo["name"]
	itemtype, _ := strconv.Atoi(itemInfo["types"])
	uinfo, err := silkworm.GetUID(openid)
	if err != nil {
		log.Println("get user info fail:", err, "openid is:", openid)
		return
	}
	uid := uinfo["id"]
	uname := uinfo["name"]
	if uinfo == nil || len(uid) == 0 || len(uname) == 0 {
		log.Println("get user info fail:", err, "openid is:", openid)
		return
	}
	// 奖励物品进背包
	_, err = silkworm.AddItemRucksack(itemid, uid, nowTime, itemtype)
	if err != nil {
		log.Println("Add Item to Rucksack Fail", err)
	}
	// 记录动态
	_, err = silkworm.SaveUserActive(activeType, uname, uid, itemname, itemid, nowTime, moreInfo)
	if err != nil {
		log.Println("Save User Active Fail", err)
	}
}

// GetUntakeLeaf 获取未拾取桑叶列表
func GetUntakeLeaf(c *gin.Context) {
	openid := c.PostForm("openid")
	if openid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	list, err := silkworm.GetUserLeafUntake(openid)
	if err != nil {
		log.Println(err)
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	itemInfo, _ := silkworm.ItemInfo("1")
	img := itemInfo["img"]
	for i := 0; i < len(list); i++ {
		list[i]["img"] = img
	}
	c.JSON(200, gin.H{
		"msg":  "success",
		"list": list,
	})
}

// TakeLeaf 收取桑叶
func TakeLeaf(c *gin.Context) {
	openid := c.PostForm("openid")
	id := c.PostForm("id")
	if openid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	rs, err := silkworm.TakeLeaf(openid, id, nowTime)
	if err != nil {
		log.Println(err)
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	log.Println(rs)
	if rs > 0 {
		uinfo, _ := silkworm.GetUID(openid)
		uid := uinfo["id"]
		uname := uinfo["name"]
		// 记录动态
		_, err = silkworm.SaveUserActive(silkworm.ActiveTakeLeaf, uname, uid, "", "0", nowTime, "")
		if err != nil {
			log.Println("Save User Active Fail", err)
		}
		responSuccess(c)
	} else {
		middleware.RespondErr(413, common.Err413UserItemNoExist, c)
		return
	}
}

// GetFriendUntakeLeaf 获取未拾取桑叶列表
func GetFriendUntakeLeaf(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	list, err := silkworm.GetUserLeafUntakeByID(id)
	if err != nil {
		log.Println(err)
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	itemInfo, _ := silkworm.ItemInfo("1")
	img := itemInfo["img"]
	for i := 0; i < len(list); i++ {
		list[i]["img"] = img
	}
	c.JSON(200, gin.H{
		"msg":  "success",
		"list": list,
	})
}

// TakeFriendLeaf 偷桑叶
func TakeFriendLeaf(c *gin.Context) {
	openid := c.PostForm("openid")
	loseUID := c.PostForm("loseuid") // 被偷用户
	id := c.PostForm("id")
	if openid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	rs := silkworm.TakeLeafByID(openid, loseUID, id, nowTime)
	if rs == -1 {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	} else if rs == -2 {
		middleware.RespondErr(500, common.Err500DBSave, c)
		return
	}
	uinfo, _ := silkworm.GetUID(openid)
	loseUserName, _ := silkworm.GetUserName(loseUID)
	uid := uinfo["id"]
	uname := uinfo["name"]
	// 记录动态
	_, err := silkworm.SaveUserActive(silkworm.ActiveStealLeaf, uname, uid, "", "0", nowTime, loseUserName)
	if err != nil {
		log.Println("Save User Active Fail", err)
	}
	responSuccess(c)
}

//Up2ButterflyRuck 进化成蝴蝶后，给用户新增一个蚕仔
func Up2ButterflyRuck(uid, nowTime string) {
	_, err := silkworm.AddSilkwormRucksack("5", uid, "0", nowTime, 0)
	if err != nil {
		log.Println("Add Item to Rucksack Fail", err, uid)
	}
}
