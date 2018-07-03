package route

import (
	"canbaobao/app"
	"canbaobao/app/auth"
	"canbaobao/app/silkworm"
	"canbaobao/app/sys"
	"canbaobao/route/middleware"
	"canbaobao/service"
	"github.com/gin-gonic/gin"
)

var (
	api    *gin.RouterGroup
	router *gin.Engine
)

// Web 路由
func Web() {
	// gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Static("/assets", "./statics/assets")
	router.Static("/upload", "./statics/upload")
	router.StaticFile("/MP_verify_Mf5ZD3XPRzBVSN2v.txt", "./statics/MP_verify_Mf5ZD3XPRzBVSN2v.txt")
	router.LoadHTMLGlob("statics/html/*")
	router.Use(middleware.Cors())

	api = router.Group("/api")
	api = router.Group("/wx")

	// ----------------------- web 路径 ---------------------------
	bms := router.Group("/")
	bms.Use(middleware.TokenAuth())

	// 登录页
	router.GET("/login", app.Login)
	// 验证登录
	router.POST("/checklogin", auth.Login)
	// 登出
	bms.GET("/logout", auth.Logout)
	// 首页
	bms.GET("/", middleware.CheckUserMenu("/"), silkworm.User)
	// 上传图片(wangEditor)
	router.POST("/uploadbywe", service.UploadByWangEditor)

	// 菜单管理页
	bms.GET("/menu", middleware.CheckUserMenu("/menu"), sys.Menu)
	// 菜单删除
	bms.GET("/delmenu", sys.DelMenu)
	// 新增菜单
	bms.POST("/addmenu", sys.AddMenu)
	// 编辑菜单
	bms.POST("/editmenu", sys.EditMenu)
	// 菜单排序
	bms.GET("/menusort", sys.ChangeMenuSort)
	// 获取全部菜单(用于角色权限管理)
	bms.GET("/menulist", sys.ChangeMenuSort)

	// 后台角色管理页
	bms.GET("/role", middleware.CheckUserMenu("/role"), sys.Role)
	// 后台角色删除
	bms.GET("/delrole", sys.DelRole)
	// 新增后台角色
	bms.POST("/addrole", sys.AddRole)
	// 编辑后台角色
	bms.POST("/editrole", sys.EditRole)
	// 获取管理用户岗位
	bms.GET("/rolemenulist", sys.GetRoleMenu)
	// 管理用户岗位
	bms.POST("/rolebindmenu", sys.RoleBindMenu)

	// 后台管理用户管理页
	bms.GET("/admin", middleware.CheckUserMenu("/admin"), sys.AdminUser)
	// 后台管理用户删除
	bms.GET("/deladmin", sys.DelAdminUser)
	// 新增后台管理用户
	bms.POST("/addadmin", sys.AddAdminUser)
	// 编辑后台管理用户
	bms.POST("/editadmin", sys.EditAdminUser)
	// 获取管理用户岗位
	bms.GET("/adminrolelist", sys.GetAdminRole)
	// 管理用户岗位
	bms.POST("/adminbindrole", sys.AdminBindRole)

	// 物品管理页
	bms.GET("/item", middleware.CheckUserMenu("/item"), silkworm.ItemList)
	// 编辑物品
	bms.POST("/edititem", silkworm.EditItem)

	// 蚕宝宝等级管理页
	bms.GET("/level", middleware.CheckUserMenu("/level"), silkworm.Level)
	// 编辑蚕宝宝等级
	bms.POST("/editlevel", silkworm.EditLevel)

	// 桑树等级管理页
	bms.GET("/treelevel", middleware.CheckUserMenu("/treelevel"), silkworm.TreeLevel)
	// 编辑桑树等级
	bms.POST("/edittreelevel", silkworm.EditTreeLevel)

	// 蝴蝶种类管理页
	bms.GET("/butterfly", middleware.CheckUserMenu("/butterfly"), silkworm.Butterfly)
	// 编辑蝴蝶种类
	bms.POST("/editbutterfly", silkworm.EditButterfly)

	// 对话管理页
	bms.GET("/dialog", middleware.CheckUserMenu("/dialog"), silkworm.Dialog)
	// 对话删除
	bms.GET("/deldialog", silkworm.DelDialog)
	// 新增对话
	bms.POST("/adddialog", silkworm.AddDialog)
	// 编辑对话
	bms.POST("/editdialog", silkworm.EditDialog)

	// 店铺管理页
	bms.GET("/vendor", middleware.CheckUserMenu("/vendor"), silkworm.Vendor)
	// 店铺删除
	bms.GET("/delvendor", silkworm.DelVendor)
	// 新增店铺
	bms.POST("/addvendor", silkworm.AddVendor)
	// 编辑店铺
	bms.POST("/editvendor", silkworm.EditVendor)

	// 商品管理页
	bms.GET("/goods", middleware.CheckUserMenu("/goods"), silkworm.Goods)
	// 获取商品内容
	bms.GET("/goodscontent", silkworm.GetGoodsContent)
	// 删除商品
	bms.GET("/delgoods", silkworm.DelGoods)
	// 新增商品
	bms.POST("/addgoods", silkworm.AddGoods)
	// 编辑商品
	bms.POST("/editgoods", silkworm.EditGoods)
	// 商品兑换条件管理
	bms.GET("/goodsredeem", silkworm.GoodsRedeem)
	// 删除商品兑换条件
	bms.GET("/delgoodsredeem", silkworm.DelGoodsRedeem)
	// 新增商品兑换条件
	bms.POST("/addgoodsredeem", silkworm.AddGoodsRedeem)
	// 编辑商品兑换条件
	bms.POST("/editgoodsredeem", silkworm.EditGoodsRedeem)

	// 兑换券管理页
	bms.GET("/voucher", middleware.CheckUserMenu("/voucher"), silkworm.Voucher)
	// 新增兑换券
	bms.POST("/addvoucher", silkworm.AddVoucher)
	// 编辑(使用)兑换券
	bms.GET("/editvoucher", silkworm.EditVoucher)

	// 用户管理
	bms.GET("/user", middleware.CheckUserMenu("/user"), silkworm.User)

	// 问题管理页
	bms.GET("/question", middleware.CheckUserMenu("/question"), silkworm.Question)
	// 删除问题
	bms.GET("/delquestion", silkworm.DelQuestion)
	// 新增问题
	bms.POST("/addquestion", silkworm.AddQuestion)
	// 编辑问题
	bms.POST("/editquestion", silkworm.EditQuestion)
	// 问题选项管理
	bms.GET("/options", silkworm.Options)
	// 删除问题选项
	bms.GET("/deloptions", silkworm.DelOptions)
	// 新增问题选项
	bms.POST("/addoptions", silkworm.AddOptions)
	// 编辑问题选项
	bms.POST("/editoptions", silkworm.EditOptions)

	// 攻略图理页
	bms.GET("/guide", middleware.CheckUserMenu("/guide"), silkworm.Guide)
	// 编辑(上传)攻略图
	bms.POST("/editguide", silkworm.EditGuide)

	// 攻略图理页
	bms.GET("/signed", middleware.CheckUserMenu("/signed"), silkworm.Signed)
	// 编辑(上传)攻略图
	bms.POST("/editsigned", silkworm.EditSigned)
}
