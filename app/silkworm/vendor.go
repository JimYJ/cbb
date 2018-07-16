package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	log "canbaobao/service/logs"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
)

// Vendor 店铺管理
func Vendor(c *gin.Context) {
	// log.Println(id)
	list, _ := silkworm.GetVendor()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	for i := 0; i < len(list); i++ {
		province := gjson.Get(common.CityJSON, fmt.Sprintf("%s.name", list[i]["province"]))
		city := gjson.Get(common.CityJSON, fmt.Sprintf("%s.city.%s.name", list[i]["province"], list[i]["city"]))
		county := gjson.Get(common.CityJSON, fmt.Sprintf("%s.city.%s.districtAndCounty.%s", list[i]["province"], list[i]["city"], list[i]["county"]))
		list[i]["area"] = fmt.Sprintf("%s%s%s", province.String(), city.String(), county.String())
	}
	c.HTML(200, "vendor.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// DelVendor 删除店铺
func DelVendor(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		middleware.RedirectErr("vendor", common.AlertError, common.AlertParamsError, c)
		return
	}
	_, err := silkworm.DelVendor(id)
	if err != nil {
		log.Println(err)
		middleware.RedirectErr("vendor", common.AlertFail, common.AlertDelFail, c)
		return
	}
	c.Redirect(302, "/vendor")
}

// AddVendor 新增店铺
func AddVendor(c *gin.Context) {
	handelVendor(c, false)
}

// EditVendor 编辑店铺
func EditVendor(c *gin.Context) {
	handelVendor(c, true)
}

func handelVendor(c *gin.Context, isEdit bool) {
	name := c.PostForm("names")
	leader := c.PostForm("leader")
	leaderphone := c.PostForm("leaderphone")
	province := c.PostForm("province")
	city := c.PostForm("city")
	county := c.PostForm("county")
	if name == "" || leader == "" || leaderphone == "" || province == "" || city == "" || county == "" {
		middleware.RedirectErr("vendor", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if _, err := strconv.Atoi(id); err != nil {
			log.Println("role id error:", err)
			middleware.RedirectErr("vendor", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditVendor(name, leader, leaderphone, nowTime, id, province, city, county)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("vendor", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/vendor")
		return
	}
	_, err := silkworm.AddVendor(name, leader, leaderphone, nowTime, province, city, county)
	if err != nil {
		log.Println("add vendor fail:", err)
		middleware.RedirectErr("vendor", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	c.Redirect(302, "/vendor")
}

// VendorByArea 区域筛选店铺
func VendorByArea(c *gin.Context) {
	province := c.PostForm("province")
	city := c.PostForm("city")
	county := c.PostForm("county")
	if province == "" || city == "" || county == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	list, err := silkworm.GetVendorByArea(province, city, county)
	if err != nil {
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	c.JSON(200, gin.H{
		"msg":  "success",
		"list": list,
	})
}
