package silkworm

import (
	"canbaobao/common"
	"canbaobao/db"
	"canbaobao/db/silkworm"
	"canbaobao/route/middleware"
	log "canbaobao/service/logs"

	"github.com/gin-gonic/gin"
)

// // UserActive 获取用户动态
// func UserActive(c *gin.Context) {
// 	pageSize := c.PostForm("pageSize")
// 	pageNo := c.PostForm("pageNo")
// 	openid := c.PostForm("openid")
// 	if len(openid) == 0 {
// 		middleware.RespondErr(402, common.Err402Param, c)
// 		return
// 	}
// 	vid, err := silkworm.GetUserVid(openid)
// 	if err != nil || vid == "" {
// 		log.Println("get user vid fail:", err)
// 		middleware.RespondErr(412, common.Err412UserNotBind, c)
// 		return
// 	}
// 	totalCount, err := silkworm.GetUserActiveCount(vid)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	paginaSQL, PageTotal := db.Pagina(pageSize, pageNo, totalCount)
// 	list, err := silkworm.GetUserActive(openid, vid, paginaSQL)
// 	if err != nil {
// 		log.Println(err)
// 		middleware.RespondErr(500, common.Err500DBrequest, c)
// 		return
// 	}
// 	c.JSON(200, gin.H{
// 		"msg":       "success",
// 		"PageTotal": PageTotal,
// 		"pageSize":  pageSize,
// 		"pageNo":    pageNo,
// 		"list":      list,
// 	})
// }

// UserActiveLog 获取重要动态
func UserActiveLog(c *gin.Context) {
	pageSize := c.PostForm("pageSize")
	pageNo := c.PostForm("pageNo")
	totalCount, err := silkworm.GetActiveLogCount()
	if err != nil {
		log.Println(err)
	}
	paginaSQL, PageTotal := db.Pagina(pageSize, pageNo, totalCount)
	list, err := silkworm.GetActiveLog(paginaSQL)
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

// UserActive 获取用户动态
func UserActive(c *gin.Context) {
	pageSize := c.PostForm("pageSize")
	pageNo := c.PostForm("pageNo")
	openid := c.PostForm("openid")
	if len(openid) == 0 {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	uname, err := silkworm.GetUID(openid)
	if err != nil || len(uname) == 0 {
		log.Println("get username  fail:", err)
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	totalCount, err := silkworm.GetMyActiveCount(uname["name"], uname["id"])
	if err != nil {
		log.Println(err)
	}
	paginaSQL, PageTotal := db.Pagina(pageSize, pageNo, totalCount)
	list, err := silkworm.GetMyActive(paginaSQL, uname["id"], uname["name"])
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

// UserSelfActive 获取用户动态
func UserSelfActive(c *gin.Context) {
	pageSize := c.PostForm("pageSize")
	pageNo := c.PostForm("pageNo")
	openid := c.PostForm("openid")
	if len(openid) == 0 {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	uname, err := silkworm.GetUID(openid)
	if err != nil || len(uname) == 0 {
		log.Println("get username  fail:", err)
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	totalCount, err := silkworm.GetUserSelfActiveCount(uname["name"], uname["id"])
	if err != nil {
		log.Println(err)
	}
	paginaSQL, PageTotal := db.Pagina(pageSize, pageNo, totalCount)
	list, err := silkworm.GetUserSelfActive(paginaSQL, uname["id"], uname["name"])
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
