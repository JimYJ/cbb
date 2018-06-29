package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

// Vendor 店铺管理
func Vendor(c *gin.Context) {
	// log.Println(id)
	list, _ := silkworm.GetVendor()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "vendor.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// DelVendor 删除店铺
func DelVendor(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		middleware.RedirectErr("vendor", common.AlertError, common.AlertParamsError, c)
		return
	}
	_, err := silkworm.DelVendor(id)
	if err != nil {
		log.Println(err)
		middleware.RedirectErr("vendor", common.AlertFail, common.AlertDelFail, c)
		return
	}
	c.Redirect(302, "/vendor")
}

// AddVendor 新增店铺
func AddVendor(c *gin.Context) {
	handelVendor(c, false)
}

// EditVendor 编辑店铺
func EditVendor(c *gin.Context) {
	handelVendor(c, true)
}

func handelVendor(c *gin.Context, isEdit bool) {
	name := c.PostForm("names")
	leader := c.PostForm("leader")
	leaderphone := c.PostForm("leaderphone")
	if name == "" || leader == "" || leaderphone == "" {
		middleware.RedirectErr("vendor", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if _, err := strconv.Atoi(id); err != nil {
			log.Println("role id error:", err)
			middleware.RedirectErr("vendor", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditVendor(name, leader, leaderphone, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("vendor", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/vendor")
		return
	}
	_, err := silkworm.AddVendor(name, leader, leaderphone, nowTime)
	if err != nil {
		log.Println("add vendor fail:", err)
		middleware.RedirectErr("vendor", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	c.Redirect(302, "/vendor")
}
