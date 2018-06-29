package sys

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
	"strconv"
	"time"
)

// AdminUser 用户管理
func AdminUser(c *gin.Context) {
	// log.Println(id)
	rolelist, _ := system.GetRole()
	vendorlist, _ := silkworm.GetVendor()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "adminuser.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         system.GetAllAdminUser(),
		"rolelist":     rolelist,
		"vendorlist":   vendorlist,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// DelAdminUser 删除用户
func DelAdminUser(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		middleware.RedirectErr("admin", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	_, err := system.DelAdminUser(id, nowTime)
	if err != nil {
		log.Println(err)
		middleware.RedirectErr("admin", common.AlertFail, common.AlertDelFail, c)
		return
	}
	c.Redirect(302, "/admin")
}

// AddAdminUser 新增用户
func AddAdminUser(c *gin.Context) {
	handelAdminUser(c, false)
}

// EditAdminUser 编辑用户
func EditAdminUser(c *gin.Context) {
	handelAdminUser(c, true)
}

func handelAdminUser(c *gin.Context, isEdit bool) {
	name := c.PostForm("username")
	vid := c.PostForm("vendorid")
	matchUser, _ := regexp.MatchString("^[0-9a-zA-Z_]{4,12}$", name)
	if !matchUser {
		middleware.RedirectErr("admin", common.AlertError, common.AlertUserError, c)
		return
	}
	if len(c.PostForm("password")) < 6 {
		middleware.RedirectErr("admin", common.AlertError, common.AlertPassError, c)
		return
	}
	pass := common.SHA1(c.PostForm("password"))
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if _, err := strconv.Atoi(id); err != nil {
			log.Println("role id error:", err)
			middleware.RedirectErr("admin", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := system.EditAdminUser(name, pass, vid, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("admin", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/admin")
		return
	}
	_, err := system.AddAdminUser(name, pass, vid, nowTime)
	if err != nil {
		log.Println("add adminsuer fail:", err)
		middleware.RedirectErr("admin", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	c.Redirect(302, "/admin")
}

// GetAdminRole 获取用户岗位列表
func GetAdminRole(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" || !common.CheckInt(id) {
		c.JSON(200, gin.H{
			"list": nil,
		})
		return
	}
	list := system.GetAdminRole(id)
	c.JSON(200, gin.H{
		"list": list,
	})
}

// AdminBindRole 管理用户岗位
func AdminBindRole(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	roles := c.PostFormArray("roles")
	if id == "" || !common.CheckInt(id) || roles == nil {
		middleware.RedirectErr("admin", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	err := system.AdminBindRole(id, nowTime, roles)
	if err != nil {
		log.Println(err)
		middleware.RedirectErr("admin", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	cookie, _ := c.Request.Cookie("c")
	token := cookie.Value
	system.SetMenu(token)
	c.Redirect(302, "/admin")
}
