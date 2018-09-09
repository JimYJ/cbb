package common

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/speps/go-hashids"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	loginlist = "LoginList"
	secure    = false // cookie 是否只在HTTPS中使用
	// Sysmenu 菜单缓存KEY
	Sysmenu        = "SysMenu"
	hashIDSalt     = "64becc3c23843942b1040ffd4743d1368d988ddf046d17d448a6e199c02c3044b425a680112b399d4dbe9b35b7ccc989"
	hashIDAlphabet = "abcdefghijklmnopqrstuvwxyz"
	hashIDLen      = 10
)

// SetCookie 设置COOKIE
func SetCookie(c *gin.Context, name, value string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  value,
		Secure: secure,
	}
	http.SetCookie(c.Writer, cookie)
}

// SHA1 sha1加盐加密
func SHA1(data string) string {
	t := sha1.New()
	data = fmt.Sprintf("%s%s", data, sha1salt)
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//CreateToken 生成token
func CreateToken(ip, uid, timestamp []byte) string {
	unToken := fmt.Sprintf("%s%s%s",
		base64.URLEncoding.EncodeToString(ip),
		base64.URLEncoding.EncodeToString(uid),
		base64.URLEncoding.EncodeToString(timestamp))
	h := hmac.New(sha256.New, []byte(AppID))
	h.Write([]byte(unToken))
	token := base64.URLEncoding.EncodeToString(h.Sum([]byte(SecretKey)))
	return token
}

//JSON2Map []byte JSON转map
func JSON2Map(b []byte) map[string]interface{} {
	var f interface{}
	json.Unmarshal(b, &f)
	return f.(map[string]interface{})
}

//GetUIDByToken 根据TOKEN获取UID
func GetUIDByToken(token string) (string, string, string) {
	cache := GetCache()
	v, found := cache.Get(token)
	if found {
		uinfo := v.(map[string]string)
		uid := uinfo["uid"]
		user := uinfo["user"]
		vid := uinfo["vid"]
		// log.Println(uid)
		return uid, user, vid
	}
	return "", "", ""
}

// GetTokenCache 填入处理token缓存缓存内容
func GetTokenCache(uid, timestamp, user, vid string) map[string]string {
	a := make(map[string]string)
	a["uid"] = uid
	a["user"] = user
	a["timestamp"] = timestamp
	a["vid"] = vid
	return a
}

// SingleLogin 单用户登陆，禁止用户多用户在线，原则是新登陆覆盖旧登陆
func SingleLogin(token string) {
	uid, _, _ := GetUIDByToken(token)
	cache := GetCache()
	newlist := make(map[string]string)
	loginList, found := cache.Get(loginlist)
	if found {
		newlist = loginList.(map[string]string)
		oldToken, ok := newlist[uid]
		if ok {
			cache.Delete(oldToken)
		}
	}
	newlist[uid] = token
	cache.Set(loginlist, newlist, -1)
	// log.Println(newlist)
}

// CheckInt 检查是否整数
func CheckInt(i string) bool {
	match, _ := regexp.MatchString("^[0-9]*$", i)
	return match
}

// GetTokenByCookie 从COOKIE获取TOKEN
func GetTokenByCookie(c *gin.Context) string {
	cookie, _ := c.Request.Cookie("c")
	return cookie.Value
}

// Remove 删除[]string函数
func Remove(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
}

// GetThisWeekMonday 获得本周周一日期
func GetThisWeekMonday() string {
	t := time.Now().Local()
	nowday, _ := strconv.Atoi(fmt.Sprintf("%d", t.Weekday()))
	if nowday == 0 { // 周日默认为0
		nowday = 7
	}
	subNumber := nowday - 1
	if subNumber == 0 {
		return t.Format("2006-01-02")
	}
	return t.AddDate(0, 0, 0-subNumber).Format("2006-01-02")
}

// ChangeMapInterface 转换[]map[string]string 为[]map[string]interface{}
func ChangeMapInterface(list []map[string]string) []map[string]interface{} {
	newList := make([]map[string]interface{}, len(list))
	for i := 0; i < len(list); i++ {
		temp := make(map[string]interface{})
		for k, v := range list[i] {
			temp[k] = v
		}
		newList[i] = temp
	}
	return newList
}

// CheckLimit 检测有没有超过每日限额
func CheckLimit(execTimes, execLastData, nowDate string, limitTimes int) int {
	nowExecTimes := 1
	if execLastData != "" {
		lastData, _ := time.Parse("2006-01-02", execLastData)
		today, _ := time.Parse("2006-01-02", nowDate)
		if lastData.Sub(today) == 0 {
			if CheckInt(execTimes) {
				a, _ := strconv.Atoi(execTimes)
				if a >= limitTimes {
					return -1
				}
				nowExecTimes += a
			}
		}
	}
	return nowExecTimes
}

// ChangeMapInt 转换map[string]string 为map[string]int
func ChangeMapInt(m map[string]string) map[string]int {
	temp := make(map[string]int)
	for k, v := range m {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			log.Println(err)
			return nil
		}
		temp[k] = vInt
	}
	return temp
}

// CalcExpPercent 计算经验值百分比
func CalcExpPercent(levelExp, nowExp, nextLevelExp string) int {
	a, err := strconv.Atoi(levelExp)
	b, err2 := strconv.Atoi(nowExp)
	c, err3 := strconv.Atoi(nextLevelExp)
	if err != nil || err2 != nil || err3 != nil {
		log.Println(err, err2, err3)
		return 0
	}
	return int(math.Floor((float64(b-a) / float64(c-a)) * 100))
}

// FormatTimeGap 格式化时间间隔-
func FormatTimeGap(s string) string {
	if s == "" {
		return ""
	}
	s = strings.Replace(s, "s", "", -1)
	s = strings.Replace(s, "h", ":", -1)
	s = strings.Replace(s, "m", ":", -1)
	return s
}

// GetHashID 生成HASHID
func GetHashID(id int64) string {
	hd := hashids.NewData()
	hd.Salt = hashIDSalt
	hd.Alphabet = hashIDAlphabet
	hd.MinLength = hashIDLen
	h, _ := hashids.NewWithData(hd)
	e, _ := h.EncodeInt64([]int64{id})
	return e

}

// GetIDByHashID 根据HASHID获得ID
func GetIDByHashID(hashid string) ([]int64, error) {
	hd := hashids.NewData()
	hd.Salt = hashIDSalt
	hd.Alphabet = hashIDAlphabet
	hd.MinLength = hashIDLen
	h, _ := hashids.NewWithData(hd)
	return h.DecodeInt64WithError(hashid)
}
