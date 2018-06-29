package service

import (
	"canbaobao/common"
	"canbaobao/route/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"path"
	"strings"
	"time"
)

const (
	uploadIamgesPath = "statics\\upload\\img\\" // 图片上传路径
	uploadIamgesURI  = "/upload/img/"
)

var (
	imgExt          = []string{"jpg", "png", "gif", "jpeg"}
	sizeLimit int64 = 3 * 1024 * 1024
)

func checkExt(ext string, extlist []string) bool {
	ext = strings.Replace(ext, ".", "", -1)
	ext = strings.ToLower(ext)
	rs := false
	for i := 0; i < len(extlist); i++ {
		if extlist[i] == ext {
			rs = true
			break
		}
	}
	return rs
}

// UploadByWangEditor 上传图片
func UploadByWangEditor(c *gin.Context) {
	file, err := c.FormFile("data")
	if err != nil {
		log.Println(err)
		middleware.RespondErr(common.HTTPParamErr, common.GetAlertCentent(common.AlertFileEmptyError), c)
		return
	}
	ext := path.Ext(file.Filename)
	if !checkExt(ext, imgExt) {
		log.Println(ext)
		middleware.RespondErr(common.HTTPParamErr, common.GetAlertCentent(common.AlertFileFormatError), c)
		return
	}
	log.Println(file.Size)
	if file.Size > sizeLimit {
		middleware.RespondErr(common.HTTPParamErr, common.GetAlertCentent(common.AlertFileSizeError), c)
		return
	}
	savepath := fmt.Sprintf("%s%s", uploadIamgesPath, file.Filename)
	err = c.SaveUploadedFile(file, savepath)
	if err != nil {
		log.Println(err, savepath)
		middleware.RespondErr(common.HTTPParamErr, common.GetAlertCentent(common.AlertSaveFail), c)
		return
	}
	url := fmt.Sprintf("%s%s", uploadIamgesURI, file.Filename)
	c.JSON(200, gin.H{
		"errno": 0,
		"data":  []string{url},
	})
}

// 上传图片
func uploadImages(c *gin.Context, newFileName, RedirectPath string) string {
	file, err := c.FormFile("data")
	if err != nil {
		log.Println(err)
		middleware.RedirectErr(RedirectPath, common.AlertError, common.AlertFileEmptyError, c)
		return ""
	}
	ext := path.Ext(file.Filename)
	if !checkExt(ext, imgExt) {
		log.Println(ext)
		middleware.RedirectErr(RedirectPath, common.AlertError, common.AlertFileFormatError, c)
		return ""
	}
	log.Println(file.Size)
	if file.Size > sizeLimit {
		middleware.RedirectErr(RedirectPath, common.AlertError, common.AlertFileSizeError, c)
		return ""
	}
	savepath := fmt.Sprintf("%s%s%s", uploadIamgesPath, newFileName, ext)
	err = c.SaveUploadedFile(file, savepath)
	if err != nil {
		log.Println(err, savepath)
		middleware.RedirectErr(RedirectPath, common.AlertError, common.AlertSaveFail, c)
		return ""
	}
	url := fmt.Sprintf("%s%s%s", uploadIamgesURI, newFileName, ext)
	return url
}

// UploadImages 上传游戏攻略图片
func UploadImages(c *gin.Context) string {
	newFileName := "guide"
	RedirectPath := "guide"
	return uploadImages(c, newFileName, RedirectPath)
}

// UploadRenameImages 上传图片并重命名
func UploadRenameImages(c *gin.Context, RedirectPath string) string {
	newFileName := fmt.Sprintf("%s%v", time.Now().Format("20060102150405"), time.Now().UnixNano())
	return uploadImages(c, newFileName, RedirectPath)
}
