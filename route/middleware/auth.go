package middleware

import (
	"canbaobao/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// TokenAuth 验证token的中间件
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		if cookie, err := c.Request.Cookie("c"); err == nil {
			token = cookie.Value
		} else {
			RedirectErr("login", common.AlertError, common.AlertCheckTokenError, c)
			return
		}
		cache := common.GetCache()
		v, found := cache.Get(token)
		ip := []byte(c.ClientIP())
		if found == false {
			log.Println("check token fail", token)
			RedirectErr("login", common.AlertError, common.AlertCheckTokenError, c)
			return
		}
		uinfo := v.(map[string]string)
		uid := uinfo["uid"]
		timestamp := uinfo["timestamp"]
		reToken := common.CreateToken(ip, []byte(uid), []byte(timestamp))
		if token != reToken {
			RedirectErr("login", common.AlertError, common.AlertCheckTokenError, c)
			return
		}
		c.Next()
	}
}

//RespondErr 错误结果响应
func RespondErr(code int, message string, c *gin.Context) {
	resp := map[string]string{"msg": message}
	c.JSON(code, resp)
	c.Abort()
}

// RedirectErr 错误结果转跳
func RedirectErr(path string, alerttitle, alertcontent int, c *gin.Context) {
	redirectErr("%s?t=%d&c=%d", path, alerttitle, alertcontent, c)
}

// RedirectErr2 错误结果转跳
func RedirectErr2(path string, alerttitle, alertcontent int, c *gin.Context) {
	redirectErr("%st=%d&c=%d", path, alerttitle, alertcontent, c)
}

func redirectErr(format, path string, alerttitle, alertcontent int, c *gin.Context) {
	url := fmt.Sprintf(format, path, alerttitle, alertcontent)
	c.Redirect(302, url)
	c.Abort()
}

// LoginRequestLimit 限制接口请求次数的中间件-限制次数多接口共用
func LoginRequestLimit(maxLimit int, post bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			RespondErr(common.HTTPUnexpectedErr, common.Err406Unexpected, c)
			return
		}
		ipKey := fmt.Sprintf("Limited:%s", ip)
		timestamp := time.Now().Unix()
		var user string
		if post {
			user = c.PostForm("phone")
		} else {
			user = c.Query("phone")
		}
		userKey := fmt.Sprintf("Limit:%s", user)
		cache := common.GetCache()
		v1, found := cache.Get(ipKey)
		if found {
			t := v1.([]int64)[0]
			freq := v1.([]int64)[1]
			if t > timestamp {
				if int(freq) >= maxLimit {
					RespondErr(common.HTTPFrequentErr, common.Err429Frequent, c)
					return
				}
				cache.Set(ipKey, []int64{t, freq + 1}, common.CacheTimeOut)
			} else {
				cache.Set(ipKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
			}
		} else {
			cache.Set(ipKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
		}
		v2, found := cache.Get(userKey)
		if found {
			t := v2.([]int64)[0]
			freq := v2.([]int64)[1]
			if t > timestamp {
				if int(freq) >= maxLimit {
					RespondErr(common.HTTPFrequentErr, common.Err429Frequent, c)
					return
				}
				cache.Set(userKey, []int64{t, freq + 1}, common.CacheTimeOut)
			} else {
				cache.Set(userKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
			}
		} else {
			cache.Set(userKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
		}
		c.Next()
	}
}

// UserRequestLimit 登陆用户-限制接口请求次数的中间件-限制次数多接口共用
func UserRequestLimit(maxLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			RespondErr(common.HTTPUnexpectedErr, common.Err406Unexpected, c)
			return
		}
		timestamp := time.Now().Unix()
		token := c.GetHeader("token")
		tokenKey := fmt.Sprintf("Limit:%s", token)
		cache := common.GetCache()
		v1, found := cache.Get(ip)
		if found {
			t := v1.([]int64)[0]
			freq := v1.([]int64)[1]
			if t > timestamp {
				if int(freq) >= maxLimit {
					RespondErr(common.HTTPFrequentErr, common.Err429Frequent, c)
					return
				}
				cache.Set(ip, []int64{t, freq + 1}, common.CacheTimeOut)
			} else {
				cache.Set(ip, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
			}
		} else {
			cache.Set(ip, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
		}
		v2, found := cache.Get(tokenKey)
		if found {
			t := v2.([]int64)[0]
			freq := v2.([]int64)[1]
			if t > timestamp {
				if int(freq) >= maxLimit {
					RespondErr(common.HTTPFrequentErr, common.Err429Frequent, c)
					return
				}
				cache.Set(tokenKey, []int64{t, freq + 1}, common.CacheTimeOut)
			} else {
				cache.Set(tokenKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
			}
		} else {
			cache.Set(tokenKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
		}
		c.Next()
	}
}

// RequestLimit 普通限制接口请求次数的中间件-限制次数多接口共用
func RequestLimit(maxLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			RespondErr(common.HTTPUnexpectedErr, common.Err406Unexpected, c)
			return
		}
		timestamp := time.Now().Unix()
		cache := common.GetCache()
		ipKey := fmt.Sprintf("Limit:%s", ip)
		v1, found := cache.Get(ipKey)
		if found {
			t := v1.([]int64)[0]
			freq := v1.([]int64)[1]
			if t > timestamp {
				if int(freq) >= maxLimit {
					RespondErr(common.HTTPFrequentErr, common.Err429Frequent, c)
					return
				}
				cache.Set(ipKey, []int64{t, freq + 1}, common.CacheTimeOut)
			} else {
				cache.Set(ipKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
			}
		} else {
			cache.Set(ipKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
		}
		c.Next()
	}
}

// SingleRequestLimit 每个接口独立限制请求次数的中间件
func SingleRequestLimit(maxLimit int, post bool, name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			RespondErr(common.HTTPUnexpectedErr, common.Err406Unexpected, c)
			return
		}
		timestamp := time.Now().Unix()
		var user string
		if post {
			user = c.PostForm("phone")
		} else {
			user = c.Query("phone")
		}
		cache := common.GetCache()
		ipKey := fmt.Sprintf("%s:%s", name, ip)
		v1, found := cache.Get(ipKey)
		if found {
			t := v1.([]int64)[0]
			freq := v1.([]int64)[1]
			if t > timestamp {
				if int(freq) >= maxLimit {
					RespondErr(common.HTTPFrequentErr, common.Err429Frequent, c)
					return
				}
				cache.Set(ipKey, []int64{t, freq + 1}, common.CacheTimeOut)
			} else {
				cache.Set(ipKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
			}
		} else {
			cache.Set(ipKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
		}
		userKey := fmt.Sprintf("%s:%s", name, user)
		v2, found := cache.Get(userKey)
		if found {
			t := v2.([]int64)[0]
			freq := v2.([]int64)[1]
			if t > timestamp {
				if int(freq) >= maxLimit {
					RespondErr(common.HTTPFrequentErr, common.Err429Frequent, c)
					return
				}
				cache.Set(userKey, []int64{t, freq + 1}, common.CacheTimeOut)
			} else {
				cache.Set(userKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
			}
		} else {
			cache.Set(userKey, []int64{timestamp + int64(common.LoginGap), 1}, common.CacheTimeOut)
		}
		c.Next()
	}
}

// CheckUserMenu 检查用户是否有访问路径的权限
func CheckUserMenu(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		if cookie, err := c.Request.Cookie("c"); err == nil {
			token = cookie.Value
		} else {
			RedirectErr("login", common.AlertError, common.AlertCheckTokenError, c)
			return
		}
		uid, _, _ := common.GetUIDByToken(token)
		k := fmt.Sprintf("UserMenu:%s", uid)
		cache := common.GetCache()
		p, found := cache.Get(k)
		if found == false {
			log.Println("check path visit error", token)
			RedirectErr("login", common.AlertError, common.AlertPathVisitError, c)
			return
		}
		pathlist := p.(map[string]bool)
		if _, ok := pathlist[path]; !ok {
			log.Println("check path visit error", path)
			RedirectErr("login", common.AlertError, common.AlertPathVisitError, c)
			return
		}
		c.Next()
	}
}

//Cors 允许跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		// origin := c.Request.Header.Get("origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Headers", "Access-Token")
		c.Next()
	}
}
