package silkworm

import (
	"canbaobao/common"
	"canbaobao/db"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	"canbaobao/service"
	log "canbaobao/service/logs"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Goods 商品管理
func Goods(c *gin.Context) {
	// log.Println(id)
	list, _ := silkworm.GetGoods()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	var substr string
	for i := 0; i < len(list); i++ {
		sublist, _ := silkworm.GetGoodsRedeem(list[i]["id"])
		substr = ""
		if sublist != nil && len(sublist) > 0 {
			for j := 0; j < len(sublist); j++ {
				substr = fmt.Sprintf("%s%s蝴蝶%s只,", substr, sublist[j]["name"], sublist[j]["numbers"])
			}
		}
		list[i]["redeems"] = substr
	}
	c.HTML(200, "goods.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// GetGoodsContent 获取商品详情
func GetGoodsContent(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		middleware.RespondErr(common.HTTPParamErr, common.Err402Param, c)
		return
	}
	content, err := silkworm.GetGoodsContent(id)
	if err != nil {
		log.Println(err)
		middleware.RespondErr(common.HTTPUnexpectedErr, common.Err406Unexpected, c)
		return
	}
	c.JSON(200, gin.H{
		"content": content,
	})
}

// DelGoods 删除商品
func DelGoods(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		middleware.RedirectErr("goods", common.AlertError, common.AlertParamsError, c)
		return
	}
	_, err := silkworm.DelGoods(id)
	if err != nil {
		log.Println(err)
		middleware.RedirectErr("goods", common.AlertFail, common.AlertDelFail, c)
		return
	}
	c.Redirect(302, "/goods")
}

// AddGoods 新增商品
func AddGoods(c *gin.Context) {
	handelGoods(c, false)
}

// EditGoods 编辑商品
func EditGoods(c *gin.Context) {
	handelGoods(c, true)
}

func handelGoods(c *gin.Context, isEdit bool) {
	name := c.PostForm("names")
	content := c.PostForm("content")
	swcount := c.PostForm("swcount")
	bigimg := service.UploadRenameImages(c, "goods")
	// log.Println(bigimg, content,swcount)
	if len(name) == 0 || len(content) == 0 || len(swcount) == 0 {
		middleware.RedirectErr("goods", common.AlertError, common.AlertParamsError, c)
		return
	}
	if _, err := strconv.Atoi(swcount); err != nil {
		middleware.RedirectErr("goods", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if _, err := strconv.Atoi(id); err != nil {
			log.Println("goods id error:", err)
			middleware.RedirectErr("goods", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditGoods(name, bigimg, content, swcount, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("goods", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/goods")
		return
	}
	_, err := silkworm.AddGoods(name, bigimg, content, swcount, nowTime)
	if err != nil {
		log.Println("add goods fail:", err)
		middleware.RedirectErr("goods", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	c.Redirect(302, "/goods")
}

// GoodsRedeem 商品兑换条件管理
func GoodsRedeem(c *gin.Context) {
	gid := c.Query("gid")
	if gid == "" || !common.CheckInt(gid) {
		log.Println(gid)
		middleware.RedirectErr("goods", common.AlertError, common.AlertParamsError, c)
		return
	}
	list, _ := silkworm.GetGoodsRedeem(gid)
	butterflylist, _ := silkworm.ButterflyList()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "goodsredeem.html", gin.H{
		"menu":          system.GetMenu(),
		"list":          list,
		"butterflylist": butterflylist,
		"gid":           gid,
		"alerttitle":    title,
		"alertcontext":  content,
	})
}

// DelGoodsRedeem 删除商品兑换条件
func DelGoodsRedeem(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	gid := c.Query("gid")
	if gid == "" || !common.CheckInt(gid) {
		middleware.RedirectErr("goods", common.AlertError, common.AlertParamsError, c)
		return
	}
	path := fmt.Sprintf("goodsredeem?gid=%s&", gid)
	if id == "" || !common.CheckInt(gid) {
		middleware.RedirectErr2(path, common.AlertError, common.AlertParamsError, c)
		return
	}
	_, err := silkworm.DelGoodsRedeem(id)
	if err != nil {
		log.Println(err)
		middleware.RedirectErr2(path, common.AlertFail, common.AlertDelFail, c)
		return
	}
	c.Redirect(302, path)
}

// AddGoodsRedeem 新增商品兑换条件
func AddGoodsRedeem(c *gin.Context) {
	handelGoodsRedeem(c, false)
}

// EditGoodsRedeem 编辑商品兑换条件
func EditGoodsRedeem(c *gin.Context) {
	handelGoodsRedeem(c, true)
}

func handelGoodsRedeem(c *gin.Context, isEdit bool) {
	butterflyid := c.PostForm("butterflyid")
	numbers := c.PostForm("numbers")
	gid := c.Query("gid")
	if gid == "" || !common.CheckInt(gid) {
		middleware.RedirectErr("goods", common.AlertError, common.AlertParamsError, c)
		return
	}
	path := fmt.Sprintf("goodsredeem?gid=%s&", gid)
	if butterflyid == "" || numbers == "" {
		middleware.RedirectErr2(path, common.AlertError, common.AlertParamsError, c)
		return
	}
	temp, _ := silkworm.CheckRepeat(gid, butterflyid)
	count, err := strconv.Atoi(temp)
	if err != nil || count > 0 {
		middleware.RedirectErr2(path, common.AlertError, common.AlertRepeatError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.DefaultQuery("id", "")
		if id == "" || !common.CheckInt(id) {
			middleware.RedirectErr2(path, common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditGoodsRedeem(gid, butterflyid, numbers, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr2(path, common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, path)
		return
	}
	_, err = silkworm.AddGoodsRedeem(butterflyid, numbers, gid, nowTime)
	if err != nil {
		log.Println("add goods fail:", err)
		middleware.RedirectErr2(path, common.AlertFail, common.AlertSaveFail, c)
		return
	}
	c.Redirect(302, path)
}

// GoodsList 获取商品列表
func GoodsList(c *gin.Context) {
	pageSize := c.PostForm("pageSize")
	pageNo := c.PostForm("pageNo")
	totalCount, _ := silkworm.GetGoodsCount()
	paginaSQL, PageTotal := db.Pagina(pageSize, pageNo, totalCount)
	list, _ := silkworm.GetPaginaGoods(paginaSQL)
	newList := common.ChangeMapInterface(list)
	for i := 0; i < len(newList); i++ {
		newList[i]["content"] = base64.StdEncoding.EncodeToString([]byte(newList[i]["content"].(string)))
	}
	for i := 0; i < len(newList); i++ {
		sublist, _ := silkworm.GetGoodsRedeem(list[i]["id"])
		newList[i]["redeemlist"] = sublist
	}
	c.JSON(200, gin.H{
		"msg":       "success",
		"list":      newList,
		"pageTotal": PageTotal,
	})
}
