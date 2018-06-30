package route

import (
	"canbaobao/app/silkworm"
)

// API 接口路由
func API() {
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
	// router.RunTLS("127.0.0.1:443", sslcert, sslkey)
	router.Run(":845")
}
