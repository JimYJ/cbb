package wechat

import (
	sw "canbaobao/app/silkworm"
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/route/middleware"
	"canbaobao/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

var (
	wechat    = &service.WeChat{}
	returnURL = "http://cbb.naiba168.com/sw/index.html"
)

// Start 转跳授权
func Start(c *gin.Context) {
	c.Redirect(302, wechat.Start())
}

// GetUserInfo 获得微信用户信息
func GetUserInfo(c *gin.Context) {
	code := c.Query("code")
	openid, _, err := wechat.GetOpenID(code)
	if err != nil {
		log.Println(err)
		middleware.RespondErr(common.HTTPExternalErr, common.Err502Wechat, c)
		return
	}
	name, avatar, err := wechat.GetUserInfo()
	if err != nil {
		log.Println(err)
		middleware.RespondErr(common.HTTPExternalErr, common.Err502Wechat, c)
		return
	}
	ip := c.ClientIP()
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	rs, _ := silkworm.CheckUserExist(openid)
	if rs < 1 {
		_, err = silkworm.AddUser(avatar, name, ip, openid, nowTime)
		if err != nil {
			log.Println(err)
			middleware.RespondErr(500, common.Err500DBSave, c)
			return
		}
		go sw.NewUserRuck(openid, nowTime)
	}
	url := fmt.Sprintf("%s?openid=%s", returnURL, openid)
	c.Redirect(302, url)
}
