package silkworm

import (
	"canbaobao/common"
	"canbaobao/db"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	log "canbaobao/service/logs"
	"github.com/gin-gonic/gin"
	"strconv"
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
	content := c.PostForm("content")
	if len(vendorid) == 0 || len(vendorid) == 0 || common.CheckInt(vendorid) || common.CheckInt(uid) || len(content) == 0 {
		middleware.RedirectErr("voucher", common.AlertError, common.AlertParamsError, c)
		return
	}
	_, err := silkworm.AddVoucher(vendorid, uid, content, nowTime)
	if err != nil {
		log.Println("add voucher fail:", err)
		middleware.RedirectErr("voucher", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	c.Redirect(302, "/voucher")
}

// UserVoucher 用户兑换券列表
func UserVoucher(c *gin.Context) {
	openid := c.PostForm("openid")
	pageSize := c.PostForm("pageSize")
	pageNo := c.PostForm("pageNo")
	if len(openid) == 0 {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	uinfo, err := silkworm.GetUID(openid)
	if err != nil {
		log.Println("get user info fail:", err, "openid", openid)
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	uid := uinfo["id"]
	totalCount, err := silkworm.GetVoucherByUserCount(uid)
	if err != nil {
		log.Println("get voucher total count fail:", err)
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	paginaSQL, PageTotal := db.Pagina(pageSize, pageNo, totalCount)
	list, err := silkworm.GetVoucherByUser(uid, paginaSQL)
	if err != nil {
		log.Println("get voucher list fail:", err)
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	c.JSON(200, gin.H{
		"msg":       "success",
		"list":      list,
		"PageTotal": PageTotal,
		"pageSize":  pageSize,
		"pageNo":    pageNo,
	})
}

// ExchangeGoods 兑换商品(生成兑换券)
func ExchangeGoods(c *gin.Context) {
	openid := c.PostForm("openid")
	goodsid := c.PostForm("goodsid")
	if openid == "" || goodsid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	uinfo, err := silkworm.GetUID(openid)
	if err != nil {
		log.Println(err)
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	uid := uinfo["id"]
	uname := uinfo["name"]
	vid := uinfo["vid"]
	if vid == "" {
		middleware.RespondErr(412, common.Err412UserNotBind, c)
		return
	}
	goodsRedeemList, err := silkworm.GetGoodsRedeemList(goodsid)
	if err != nil || goodsRedeemList == nil {
		log.Println(err)
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	idList := make([]string, 0)
	for i := 0; i < len(goodsRedeemList); i++ {
		numbers, _ := strconv.Atoi(goodsRedeemList[i]["numbers"])
		bfList, err := silkworm.GetUserButterflyList(goodsRedeemList[i]["butterflyid"], uid, goodsRedeemList[i]["numbers"])
		if err != nil || len(bfList) < numbers {
			log.Println(err)
			middleware.RespondErr(422, common.Err422NotEnoughExchange, c)
			return
		}
		for j := 0; j < len(bfList); j++ {
			idList = append(idList, bfList[j]["id"])
		}
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	goodsName, _ := silkworm.GetGoodsName(goodsid)
	rs := silkworm.ExchangeGoods(vid, uid, goodsName, nowTime, idList)
	if !rs {
		middleware.RespondErr(500, common.Err500DBSave, c)
		return
	}
	go exchangeGoodsActive(uname, uid, nowTime, goodsName)
	responSuccess(c)
}

// exchangeGoodsActive 兑换商品动态
func exchangeGoodsActive(uname, uid, nowTime, goodsName string) {
	_, err := silkworm.SaveUserActive(silkworm.ActiveVoucher, uname, uid, goodsName, "0", nowTime, "")
	if err != nil {
		log.Println("Save User Active Fail:", err)
	}
}

// BatchVoucher 批量发放兑换券
func BatchVoucher(c *gin.Context) {
	_, _, vid := common.GetUIDByToken(common.GetTokenByCookie(c)) //用户的绑定店铺ID
	if vid != "0" {
		middleware.RedirectErr("voucher", common.AlertFail, common.AlertVoucherIIError, c)
		return
	}
	province := c.PostForm("province")
	city := c.PostForm("city")
	startDay := c.PostForm("startday")
	endDay := c.PostForm("endday")
	content := c.PostForm("content")
	if len(province) == 0 || len(startDay) == 0 || len(endDay) == 0 || len(content) == 0 {
		middleware.RedirectErr("voucher", common.AlertError, common.AlertParamsError, c)
		return
	}
	s, err := time.Parse("2006-01-02", startDay)
	e, err2 := time.Parse("2006-01-02", endDay)
	t, _ := time.Parse("2006-01-02", time.Now().Local().Format("2006-01-02"))
	if err != nil || err2 != nil {
		middleware.RedirectErr("voucher", common.AlertError, common.AlertParamsError, c)
		return
	}
	if s.Sub(e) > 0 {
		middleware.RedirectErr("voucher", common.AlertError, common.AlertRangeDateError, c)
		return
	}
	if t.Sub(e) > 0 {
		middleware.RedirectErr("voucher", common.AlertError, common.AlertRangeDateErrorII, c)
		return
	}
	ulist, err := silkworm.GetUserForAreaVendor(province, city)
	if err != nil {
		log.Println("Get User For Area Vendor Fail", err)
		middleware.RedirectErr("voucher", common.AlertFail, common.AlertGetDataFail, c)
		return
	}
	if ulist == nil || len(ulist) == 0 {
		middleware.RedirectErr("voucher", common.AlertFail, common.AlertBatchVoucherFail, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	rs := silkworm.BatchVoucher(ulist, nowTime, startDay, endDay, content)
	if !rs {
		middleware.RedirectErr("voucher", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	go silkworm.BatchVoucherActive(ulist, nowTime, startDay, endDay, content)
	c.Redirect(302, "/voucher")
}
