package service

import (
	"canbaobao/common"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// WeChat 微信接口
type WeChat struct {
	appID, appSecret, AccessToken, code, openid, RefreshToken string
}

var (
	appID     = "wx38c9bfc5478be3e6"
	appSecret = "daee17f625a0377b4c37d7795ef86874"
)

// Start 转跳微信授权页
func (w *WeChat) Start() string {
	w.appID = appID
	w.appSecret = appSecret
	redirectURI := fmt.Sprintf("%s/wx/getuinfo", common.AppPath)
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
			w.AccessToken = v.(string)
		} else {
			return "", "", errors.New("get access token fail,access_token is empty")
		}
	}
	if v, ok := results["refresh_token"]; ok {
		if len(v.(string)) > 0 {
			w.RefreshToken = v.(string)
		} else {
			return "", "", errors.New("get refresh token fail,refresh_token is empty")
		}
	}
	if v, ok := results["openid"]; ok {
		if v != "" {
			w.openid = v.(string)
		} else {
			return "", "", errors.New("get openid fail,access_token is empty")
		}
	}
	return w.openid, w.AccessToken, nil
}

//GetUserInfo 获取用户信息
func (w *WeChat) GetUserInfo() (string, string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", w.AccessToken, w.openid)
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

// RefreshAccessToken 刷新授权TOKEN
func (w *WeChat) RefreshAccessToken(refreshToken string) error {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&grant_type=refresh_token&refresh_token=%s", w.appID, refreshToken)
	rs, err := Get(url)
	if err != nil {
		return err
	}
	results := JSON2Map(rs)
	if v, ok := results["errcode"]; ok {
		errcode := v.(float64)
		str := strconv.Itoa(int(errcode))
		return errors.New(str)
	}
	if v, ok := results["access_token"]; ok {
		if v != "" {
			w.AccessToken = v.(string)
		} else {
			return errors.New("get access token fail,access_token is empty")
		}
	}
	if v, ok := results["refresh_token"]; ok {
		if len(v.(string)) > 0 {
			w.RefreshToken = v.(string)
		} else {
			return errors.New("get refresh token fail,refresh_token is empty")
		}
	}
	return nil
}

//GetAccessToken 获取全局 AccessToken
func (w *WeChat) GetAccessToken() (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", w.appID, w.appSecret)
	rs, err := Get(url)
	if err != nil {
		return "", err
	}
	results := JSON2Map(rs)
	if _, ok := results["errcode"]; ok {
		// errcode := v.(float64)
		// str := strconv.Itoa(int(errcode))
		return "", errors.New(results["errmsg"].(string))
	}
	if v, ok := results["access_token"]; ok {
		if v != "" {
			return v.(string), nil
		}
		return "", errors.New("get access token fail,access_token is empty")
	}
	return "", errors.New("get access token fail,respon key access_token is null")
}

//GetJsapiTicket 获取公众号Jsapi Ticket
func (w *WeChat) GetJsapiTicket(accessToken string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken)
	rs, err := Get(url)
	if err != nil {
		return "", err
	}
	results := JSON2Map(rs)
	if v, ok := results["ticket"]; ok {
		if v != "" {
			return v.(string), nil
		}
		return "", errors.New("get ticket fail,ticket is empty")
	}
	return "", errors.New("get ticket fail,respon key ticket is null")
}
