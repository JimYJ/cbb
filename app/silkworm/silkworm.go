package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Level 蚕宝宝列表
func Level(c *gin.Context) {
	list, _ := silkworm.LevelList()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "level.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// EditLevel 编辑等级
func EditLevel(c *gin.Context) {
	handelLevel(c, true)
}

func handelLevel(c *gin.Context, isEdit bool) {
	redeemitem := c.PostForm("redeemitem")
	if redeemitem == "" {
		middleware.RedirectErr("level", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if !common.CheckInt(id) {
			middleware.RedirectErr("level", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditLevel(redeemitem, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("level", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/level")
		return
	}
	c.Redirect(302, "/level")
}
