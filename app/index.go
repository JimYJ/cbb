package app

import (
	"canbaobao/common"
	"canbaobao/db/system"
	"github.com/gin-gonic/gin"
)

// Index 首页
func Index(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"menu": system.GetMenu(),
	})
}

// Login 登录页
func Login(c *gin.Context) {
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "login.html", gin.H{
		"alerttitle":   title,
		"alertcontext": content,
	})
}
