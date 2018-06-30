package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/route/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
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
	uinfo, _ := silkworm.GetUID(openid)
	uid := uinfo["id"]
	uname := uinfo["name"]
	// 奖励物品进背包
	_, err := silkworm.AddItemRucksack(itemid, uid, nowTime, itemtype)
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
	rs, err := silkworm.TakeLeaf(openid, id)
	if err != nil {
		log.Println(err)
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	if rs > 0 {
		responSuccess(c)
	} else {
		middleware.RespondErr(402, common.Err402UserItemNoExist, c)
		return
	}

}
