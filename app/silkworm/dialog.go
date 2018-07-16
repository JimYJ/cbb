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

// Dialog 对话管理
func Dialog(c *gin.Context) {
	// log.Println(id)
	list, _ := silkworm.GetDialog()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "dialog.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// DelDialog 删除对话
func DelDialog(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		middleware.RedirectErr("dialog", common.AlertError, common.AlertParamsError, c)
		return
	}
	_, err := silkworm.DelDialog(id)
	if err != nil {
		log.Println(err)
		middleware.RedirectErr("dialog", common.AlertFail, common.AlertDelFail, c)
		return
	}
	c.Redirect(302, "/dialog")
}

// AddDialog 新增对话
func AddDialog(c *gin.Context) {
	handelDialog(c, false)
}

// EditDialog 编辑对话
func EditDialog(c *gin.Context) {
	handelDialog(c, true)
}

func handelDialog(c *gin.Context, isEdit bool) {
	content := c.PostForm("content")
	if len([]rune(content)) > 15 {
		middleware.RedirectErr("dialog", common.AlertError, common.AlertDialogLengthError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if !common.CheckInt(id) {
			middleware.RedirectErr("dialog", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditDialog(content, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("dialog", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/dialog")
		return
	}
	_, err := silkworm.AddDialog(content, nowTime)
	if err != nil {
		log.Println("add dialog fail:", err)
		middleware.RedirectErr("dialog", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	c.Redirect(302, "/dialog")
}
