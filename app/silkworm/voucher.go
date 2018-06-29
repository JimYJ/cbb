package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Voucher 兑换券管理
func Voucher(c *gin.Context) {
	_, _, vid := common.GetUIDByToken(common.GetTokenByCookie(c))
	list, _ := silkworm.GetVoucher(vid)
	vendorlist, _ := silkworm.GetVendor()
	userlist, _ := silkworm.GetUser()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "voucher.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"vendorlist":   vendorlist,
		"userlist":     userlist,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// AddVoucher 新增兑换券
func AddVoucher(c *gin.Context) {
	handelVoucher(c, false)
}

// EditVoucher 编辑兑换券
func EditVoucher(c *gin.Context) {
	handelVoucher(c, true)
}

func handelVoucher(c *gin.Context, isEdit bool) {
	_, _, vid := common.GetUIDByToken(common.GetTokenByCookie(c)) //用户的绑定店铺ID
	content := c.PostForm("content")
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	// 编辑(使用兑换券)
	if isEdit {
		id := c.Query("id")
		if !common.CheckInt(id) {
			middleware.RedirectErr("voucher", common.AlertError, common.AlertParamsError, c)
			return
		}
		rs, err := silkworm.EditVoucher(vid, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("voucher", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		if rs == 0 {
			middleware.RedirectErr("voucher", common.AlertFail, common.AlertVoucherError, c)
			return
		}
		c.Redirect(302, "/voucher")
		return
	}
	// 新增兑换券
	if vid != "0" {
		middleware.RedirectErr("voucher", common.AlertFail, common.AlertVoucherIIError, c)
		return
	}
	vendorid := c.PostForm("vendorid")
	uid := c.PostForm("uid")
	_, err := silkworm.AddVoucher(vendorid, uid, content, nowTime)
	if err != nil {
		log.Println("add voucher fail:", err)
		middleware.RedirectErr("voucher", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	c.Redirect(302, "/voucher")
}
