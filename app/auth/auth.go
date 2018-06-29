package auth

import (
	"canbaobao/common"
	"canbaobao/db/system"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"regexp"
	"strconv"
	"time"
)

// Login 登陆验证
func Login(c *gin.Context) {
	user := c.PostForm("user")
	pass := c.PostForm("password")
	matchUser, _ := regexp.MatchString("^[0-9a-zA-Z_]{4,12}$", user)
	if !matchUser {
		url := fmt.Sprintf("login?t=%d&c=%d", common.AlertFail, common.AlertLoginFail)
		c.Redirect(302, url)
		return
	}
	if len(pass) < 6 {
		url := fmt.Sprintf("login?t=%d&c=%d", common.AlertFail, common.AlertLoginFail)
		c.Redirect(302, url)
		return
	}
	pass = common.SHA1(pass)
	uinfo, err := system.CheckPass(user, pass)

	if err == 500 {
		url := fmt.Sprintf("login?t=%d&c=%d", common.AlertFail, common.AlertDBFail)
		c.Redirect(302, url)
		return
	} else if err == 401 {
		log.Println(3)
		url := fmt.Sprintf("login?t=%d&c=%d", common.AlertFail, common.AlertLoginFail)
		c.Redirect(302, url)
		return
	}
	id := uinfo["id"]
	vid := uinfo["vendorid"]
	t := time.Now().UnixNano()
	timestamp := strconv.FormatInt(t, 10)
	ip := []byte(c.ClientIP())
	uid := []byte(id)
	token := common.CreateToken(ip, uid, []byte(timestamp))
	cache := common.GetCache()
	tokeninfo := common.GetTokenCache(id, timestamp, user, vid)
	cache.Set(token, tokeninfo, common.TokenTimeOut)
	common.SingleLogin(token)
	common.SetCookie(c, "c", token)
	system.SetMenu(token)
	c.Redirect(302, "/voucher")
}

// Logout 登出
func Logout(c *gin.Context) {
	var token string
	if cookie, err := c.Request.Cookie("c"); err == nil {
		token = cookie.Value
	} else {
		c.Redirect(302, "/login")
		return
	}
	cache := common.GetCache()
	cache.Delete(token)
	system.SetMenu(token)
	c.Redirect(302, "/login")
}
