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

// Butterfly 蝴蝶列表
func Butterfly(c *gin.Context) {
	list, _ := silkworm.ButterflyList()
	if list != nil && len(list) > 0 {
		for i := 0; i < len(list); i++ {
			if list[i]["type"] == "0" {
				list[i]["typename"] = "普通"
			} else if list[i]["type"] == "1" {
				list[i]["typename"] = "特殊"
			} else {
				list[i]["typename"] = "未知"
			}
		}
	}
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "butterfly.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// EditButterfly 编辑蝴蝶
func EditButterfly(c *gin.Context) {
	handelButterfly(c, true)
}

func handelButterfly(c *gin.Context, isEdit bool) {
	name := c.PostForm("names")
	if name == "" {
		middleware.RedirectErr("butterfly", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if !common.CheckInt(id) {
			middleware.RedirectErr("butterfly", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditButterfly(name, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("butterfly", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/butterfly")
		return
	}
	c.Redirect(302, "/butterfly")
}
