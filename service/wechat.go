package service

var (
	url         = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=&s&state=STATE#wechat_redirect"
	redirectURI = ""
)

// WeChat 微信接口
type WeChat struct {
	appID, appSecret, accessToken, code string
}

// GetCode 转跳获得CODE
func (w *WeChat) GetCode() {

}
