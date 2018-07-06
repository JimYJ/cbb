package silkworm

import (
	"canbaobao/common"
	"canbaobao/db"
	"canbaobao/db/silkworm"
	"canbaobao/route/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

// UserActive 获取用户动态
func UserActive(c *gin.Context) {
	pageSize := c.PostForm("pageSize")
	pageNo := c.PostForm("pageNo")
	openid := c.PostForm("openid")
	if openid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	vid, err := silkworm.GetUserVid(openid)
	if err != nil || vid == "" {
		log.Println("get user vid fail:", err)
		middleware.RespondErr(412, common.Err412UserNotBind, c)
		return
	}
	totalCount, err := silkworm.GetUserActiveCount(vid)
	if err != nil {
		log.Println(err)
	}
	paginaSQL, PageTotal := db.Pagina(pageSize, pageNo, totalCount)
	list, err := silkworm.GetUserActive(openid, vid, paginaSQL)
	if err != nil {
		log.Println(err)
		middleware.RespondErr(500, common.Err500DBrequest, c)
		return
	}
	c.JSON(200, gin.H{
		"msg":       "success",
		"PageTotal": PageTotal,
		"pageSize":  pageSize,
		"pageNo":    pageNo,
		"list":      list,
	})
}
