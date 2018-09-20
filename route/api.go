package route

import (
	"canbaobao/app/silkworm"
	"canbaobao/app/wechat"
)

// API 接口路由
func API() {
	// ----------------------- 微信相关路径 ---------------------------
	wx.GET("/start", wechat.Start)
	wx.GET("/getuinfo", wechat.GetUserInfo)
	wx.GET("/getaccesstoken", wechat.GetAccessToken)
	wx.GET("/getticket", wechat.GetTicket)
	// ----------------------- api 接口路径 ---------------------------
	// api.GET("/login", app.Login)
	// 获取游戏攻略图
	api.GET("/guide", silkworm.GetGuideImg)
	// 获取签到天数
	api.POST("/getsigndays", silkworm.GetSignedDays)
	// 获取用户信息(ID)
	api.POST("/userinfobyid", silkworm.UserInfoByID)
	// 获取用户信息(OpenID)
	api.POST("/userinfobyopenid", silkworm.UserInfoByOpenID)
	// 获取好友列表
	api.POST("/friendlist", silkworm.FriendList)
	// 获取问题列表
	api.POST("/questionlist", silkworm.UserQuestionList)
	// 检查问题答案
	api.POST("/checkanswer", silkworm.CheckAnswer)
	// 获得用户背包
	api.POST("/userrucksack", silkworm.UserRucksack)
	// 获得用户动态
	api.POST("/useractive", silkworm.UserActive)
	// 获得重要动态
	api.POST("/impoactive", silkworm.UserActiveLog)
	// 获得兑换商品列表
	api.POST("/goodslist", silkworm.GoodsList)
	// 浇水
	api.POST("/waterfortree", silkworm.WaterForTree)
	// 施肥
	api.POST("/fertilizerfortree", silkworm.FertilizerForTree)
	// 施肥
	api.POST("/usersigned", silkworm.UserSigned)
	// 孵化普通蚕仔
	api.POST("/hatchfornormal", silkworm.HatchForNormal)
	// 孵化特殊蚕仔
	api.POST("/hatchforspecial", silkworm.HatchForSpecial)
	// 获得自己的蚕宝宝列表
	api.POST("/usersilkwormlist", silkworm.UserSilkwormList)
	// 获得好友的蚕宝宝列表
	api.POST("/friendsilkwormlist", silkworm.FriendSilkwormList)
	// 获得自己未拾取桑叶列表
	api.POST("/untakeleaf", silkworm.GetUntakeLeaf)
	// 拾取桑叶
	api.POST("/takeleaf", silkworm.TakeLeaf)
	// 获得好友未拾取桑叶列表
	api.POST("/friendleaf", silkworm.GetFriendUntakeLeaf)
	// 偷桑叶
	api.POST("/stealleaf", silkworm.TakeFriendLeaf)
	// 申请蝴蝶配对
	api.POST("/applypair", silkworm.ApplyPair)
	// 同意蝴蝶配对申请
	api.POST("/allowpair", silkworm.AllowPair)
	// 拒绝蝴蝶配对申请
	api.POST("/rejectpair", silkworm.RejectPair)
	// 排行榜
	api.GET("/billboard", silkworm.BillBoard)
	// 店铺列表
	api.POST("/vendor", silkworm.VendorByArea)
	// 绑定店铺
	api.POST("/bindvendor", silkworm.BindVendor)
	// 检测用户是否看过引导页
	api.POST("/intropage", silkworm.CheckUserIntroPage)
	// 喂食
	api.POST("/feed", silkworm.Feed)
	// 用户兑换券
	api.POST("/voucher", silkworm.UserVoucher)
	// 用户兑换券
	api.POST("/exchangegoods", silkworm.ExchangeGoods)
	// 获取邀请链接
	api.POST("/invitelink", silkworm.GetUserInviteLink)
	// 邀请记录
	api.POST("/invitelog", silkworm.GetUserAwardLog)
	// router.RunTLS("127.0.0.1:443", sslcert, sslkey)
	router.Run(":845")
}
