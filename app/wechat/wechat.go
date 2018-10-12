package wechat

import (
	sw "canbaobao/app/silkworm"
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/route/middleware"
	"canbaobao/service"
	log "canbaobao/service/logs"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	wechat    = &service.WeChat{}
	returnURL = "http://cbb.naiba168.com/sw/setup.html"
)

// Start 转跳授权
func Start(c *gin.Context) {
	p := c.Query("p")
	c.Redirect(302, wechat.Start(p))
}

// GetUserInfo 获得微信用户信息
func GetUserInfo(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	openid, _, err := wechat.GetOpenID(code)
	if err != nil {
		log.Println(err)
		middleware.RespondErr(common.HTTPExternalErr, common.Err502Wechat, c)
		return
	}
	// cc := common.GetCache()
	// oATKey, oRTKey := common.GetKeyName(openid)
	// cc.Set(oATKey, wechat.AccessToken, time.Minute*100)
	// cc.Set(oRTKey, wechat.RefreshToken, time.Hour*24*25)
	name, avatar, err := wechat.GetUserInfo()
	if err != nil {
		log.Println(err)
		middleware.RespondErr(common.HTTPExternalErr, common.Err502Wechat, c)
		return
	}
	ip := c.ClientIP()
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	rs, _ := silkworm.CheckUserExist(openid)
	var iuid int64
	if rs < 1 {
		iuid, err = silkworm.AddUser(avatar, name, ip, openid, nowTime)
		if err != nil {
			log.Println(err)
			middleware.RespondErr(500, common.Err500DBSave, c)
			return
		}
		go sw.NewUserRuck(openid, nowTime)
		go sw.HandleInviteAward(state, nowTime, iuid)
	}
	url := fmt.Sprintf("%s?openid=%s", returnURL, openid)
	c.Redirect(302, url)
}

// GetAccessToken 获得全局AccessToken
func GetAccessToken(c *gin.Context) {
	cache := common.GetCache()
	acKey := "acKey"
	var AccessToken string
	var err error
	v, found := cache.Get(acKey)
	if found != false {
		AccessToken = v.(string)
	} else {
		AccessToken, err = wechat.GetAccessToken()
		if err != nil {
			log.Println(err)
			middleware.RespondErr(502, common.Err502Wechat, c)
			return
		}
		cache.Set(acKey, AccessToken, time.Minute*100)
	}
	c.JSON(200, gin.H{
		"msg":         "success",
		"accessToken": AccessToken,
	})
}

// GetTicket 获得JS Ticket
func GetTicket(c *gin.Context) {
	cache := common.GetCache()
	acKey := "acKey"
	tcKey := "tcKey"
	var AccessToken, Ticket string
	var err error
	t, found := cache.Get(tcKey)
	if found != false {
		Ticket = t.(string)
	} else {
		v, found := cache.Get(acKey)
		if found != false {
			AccessToken = v.(string)
		} else {
			AccessToken, err = wechat.GetAccessToken()
			if err != nil {
				log.Println(err)
				middleware.RespondErr(502, common.Err502Wechat, c)
				return
			}
			cache.Set(acKey, AccessToken, time.Minute*100)
		}
		Ticket, err = wechat.GetJsapiTicket(AccessToken)
		if err != nil {
			log.Println(err)
			middleware.RespondErr(502, common.Err502Wechat, c)
			return
		}
		cache.Set(tcKey, Ticket, time.Minute*100)
	}
	c.JSON(200, gin.H{
		"msg":    "success",
		"ticket": Ticket,
	})
}
