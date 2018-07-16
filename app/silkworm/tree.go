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

// TreeLevel 桑树列表
func TreeLevel(c *gin.Context) {
	list, _ := silkworm.TreeLevelList()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "treelevel.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// EditTreeLevel 编辑桑树等级
func EditTreeLevel(c *gin.Context) {
	handelTreeLevel(c, true)
}

func handelTreeLevel(c *gin.Context, isEdit bool) {
	growthhours := c.PostForm("growthhours")
	maxhours := c.PostForm("maxhours")
	if growthhours == "" || maxhours == "" {
		middleware.RedirectErr("treelevel", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if !common.CheckInt(id) {
			middleware.RedirectErr("treelevel", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditTreeLevel(growthhours, maxhours, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("treelevel", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/treelevel")
		return
	}
	c.Redirect(302, "/treelevel")
}
