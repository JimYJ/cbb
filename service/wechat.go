package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	redirectURI = "http://cbb.naiba168.com/wx/getuinfo"
)

// WeChat 微信接口
type WeChat struct {
	appID, appSecret, accessToken, code, openid string
}

var (
	appID     = "wx38c9bfc5478be3e6"
	appSecret = "daee17f625a0377b4c37d7795ef86874"
)

// Start 转跳微信授权页
func (w *WeChat) Start() string {
	w.appID = appID
	w.appSecret = appSecret
	return fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=STATE#wechat_redirect", w.appID, redirectURI, "snsapi_userinfo")
}

//GetOpenID 获取公众号OPENID
func (w *WeChat) GetOpenID(code string) (string, string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", w.appID, w.appSecret, code)
	rs, err := Get(url)
	if err != nil {
		return "", "", err
	}
	results := JSON2Map(rs)
	if v, ok := results["errcode"]; ok {
		errcode := v.(float64)
		str := strconv.Itoa(int(errcode))
		return "", "", errors.New(str)
	}
	if v, ok := results["access_token"]; ok {
		if v != "" {
			w.accessToken = v.(string)
		} else {
			return "", "", errors.New("get access token fail,access_token is empty")
		}
	}
	if v, ok := results["openid"]; ok {
		if v != "" {
			w.openid = v.(string)
		} else {
			return "", "", errors.New("get openid fail,access_token is empty")
		}
	}
	return w.openid, w.accessToken, nil
}

//GetUserInfo 获取用户信息
func (w *WeChat) GetUserInfo() (string, string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", w.accessToken, w.openid)
	rs, err := Get(url)
	if err != nil {
		return "", "", err
	}
	var headimgURL, nickName string
	results := JSON2Map(rs)
	if v, ok := results["errcode"]; ok {
		errcode := v.(float64)
		str := strconv.Itoa(int(errcode))
		return "", "", errors.New(str)
	}
	if v, ok := results["nickname"]; ok {
		if v != "" {
			nickName = v.(string)
		} else {
			return "", "", errors.New("get nickname fail,access_token is empty")
		}
	}
	if v, ok := results["headimgurl"]; ok {
		if v != "" {
			headimgURL = v.(string)
		} else {
			return "", "", errors.New("get headimgurl fail,access_token is empty")
		}
	}
	return nickName, headimgURL, nil
}

//Get http Get请求简单封装
func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

//JSON2Map []byte JSON转map
func JSON2Map(b []byte) map[string]interface{} {
	var f interface{}
	json.Unmarshal(b, &f)
	return f.(map[string]interface{})
}
