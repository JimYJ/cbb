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

// Signed 获取签到规则
func Signed(c *gin.Context) {
	list, _ := silkworm.SignedList()
	itemlist, _ := silkworm.ItemList()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	for i := 0; i < len(list); i++ {
		for j := 0; j < len(itemlist); j++ {
			if list[i]["dayitemid"] == itemlist[j]["id"] {
				list[i]["dayitemname"] = itemlist[j]["name"]
				break
			}
		}
		for j := 0; j < len(itemlist); j++ {
			if list[i]["weekitemid"] == itemlist[j]["id"] {
				list[i]["weekitemname"] = itemlist[j]["name"]
				break
			}
		}
	}
	c.HTML(200, "signed.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"itemlist":     itemlist,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// EditSigned 编辑签到规则
func EditSigned(c *gin.Context) {
	handelSigned(c, true)
}

func handelSigned(c *gin.Context, isEdit bool) {
	intro := c.PostForm("intro")
	dayitemid := c.PostForm("dayitemid")
	weekitemid := c.PostForm("weekitemid")
	if intro == "" || dayitemid == "" || weekitemid == "" {
		middleware.RedirectErr("signed", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if !common.CheckInt(id) {
			middleware.RedirectErr("signed", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditSigned(intro, dayitemid, weekitemid, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("signed", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/signed")
		return
	}
	c.Redirect(302, "/signed")
}
