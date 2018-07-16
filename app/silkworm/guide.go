package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	"canbaobao/service"
	log "canbaobao/service/logs"
	"github.com/gin-gonic/gin"
	"time"
)

// Guide 获取物品列表
func Guide(c *gin.Context) {
	list, _ := silkworm.GuideList()
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "guide.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// EditGuide 编辑物品
func EditGuide(c *gin.Context) {
	handelGuide(c, true)
}

func handelGuide(c *gin.Context, isEdit bool) {
	imgurl := service.UploadImages(c)
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if !common.CheckInt(id) {
			log.Println("id is error:", id)
			middleware.RedirectErr("guide", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditGuide(imgurl, nowTime, id)
		if err != nil {
			log.Println("edit guide fail:", err)
			middleware.RedirectErr("guide", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/guide")
		return
	}
	c.Redirect(302, "/guide")
}

// GetGuideImg 获得攻略图片接口
func GetGuideImg(c *gin.Context) {
	list, _ := silkworm.GuideList()
	// imgURL := fmt.Sprintf("%s%s", common.AppPath, list[0]["imgurl"])
	c.JSON(200, gin.H{
		"msg": "success",
		"img": list[0]["imgurl"],
	})
}
