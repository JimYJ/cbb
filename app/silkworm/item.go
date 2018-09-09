package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	log "canbaobao/service/logs"
	"github.com/gin-gonic/gin"
	"time"
)

// ItemList 获取物品列表
func ItemList(c *gin.Context) {
	list, _ := silkworm.ItemList()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "item.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// EditItem 编辑物品
func EditItem(c *gin.Context) {
	handelItem(c, true)
}

func handelItem(c *gin.Context, isEdit bool) {
	exp := c.PostForm("exp")
	limitday := c.PostForm("limitday")
	if exp == "" || limitday == "" {
		middleware.RedirectErr("item", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if !common.CheckInt(id) {
			middleware.RedirectErr("item", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditItem(exp, limitday, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("item", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/item")
		return
	}
	c.Redirect(302, "/item")
}

// AwardItem 编辑奖励物品
func AwardItem(c *gin.Context) {
	list, _ := silkworm.AwardItemList()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	itemlist, _ := silkworm.InviteAwardItem()
	c.HTML(200, "inviteaward.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"itemlist":     itemlist,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// EditAwardItem 编辑物品
func EditAwardItem(c *gin.Context) {
	handelAwardItem(c, true)
}

func handelAwardItem(c *gin.Context, isEdit bool) {
	itemid := c.PostForm("itemid")
	num := c.PostForm("num")
	if len(itemid) == 0 {
		middleware.RedirectErr("inviteaward", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if !common.CheckInt(id) {
			middleware.RedirectErr("inviteaward", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditAwardItem(itemid, num, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("inviteaward", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/inviteaward")
		return
	}
	c.Redirect(302, "/inviteaward")
}
