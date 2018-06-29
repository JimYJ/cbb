package silkworm

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"canbaobao/route/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

// Question 问题管理
func Question(c *gin.Context) {
	// log.Println(id)
	list, _ := silkworm.GetQuestion()
	itemlist, _ := silkworm.ItemTypeList(false)
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	for i := 0; i < len(list); i++ {
		counts, _ := silkworm.GetOptionsCount(list[i]["id"])
		list[i]["counts"] = counts
	}
	c.HTML(200, "question.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"itemlist":     itemlist,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// DelQuestion 删除问题
func DelQuestion(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		middleware.RedirectErr("question", common.AlertError, common.AlertParamsError, c)
		return
	}
	_, err := silkworm.DelQuestion(id)
	if err != nil {
		log.Println(err)
		middleware.RedirectErr("question", common.AlertFail, common.AlertDelFail, c)
		return
	}
	c.Redirect(302, "/question")
}

// AddQuestion 新增问题
func AddQuestion(c *gin.Context) {
	handelQuestion(c, false)
}

// EditQuestion 编辑问题
func EditQuestion(c *gin.Context) {
	handelQuestion(c, true)
}

func handelQuestion(c *gin.Context, isEdit bool) {
	itemid := c.PostForm("itemid")
	content := c.PostForm("content")
	if itemid == "" || content == "" {
		middleware.RedirectErr("question", common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.Query("id")
		if _, err := strconv.Atoi(id); err != nil {
			log.Println("Question id error:", err)
			middleware.RedirectErr("question", common.AlertError, common.AlertParamsError, c)
			return
		}
		_, err := silkworm.EditQuestion(itemid, content, nowTime, id)
		if err != nil {
			log.Println(err)
			middleware.RedirectErr("question", common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, "/question")
		return
	}
	_, err := silkworm.AddQuestion(itemid, content, nowTime)
	if err != nil {
		log.Println("add Question fail:", err)
		middleware.RedirectErr("question", common.AlertFail, common.AlertSaveFail, c)
		return
	}
	c.Redirect(302, "/question")
}

// Options 问题选项管理
func Options(c *gin.Context) {
	qid := c.Query("qid")
	if qid == "" || !common.CheckInt(qid) {
		log.Println(qid)
		middleware.RedirectErr("question", common.AlertError, common.AlertParamsError, c)
		return
	}
	list, _ := silkworm.GetOptions(qid)
	for i := 0; i < len(list); i++ {
		if list[i]["answer"] == "0" {
			list[i]["answerstr"] = "否"
		} else {
			list[i]["answerstr"] = "是"
		}
	}
	question, _ := silkworm.GetQuestionName(qid)
	title, content := common.GetAlertMsg(c.Query("t"), c.Query("c"))
	c.HTML(200, "options.html", gin.H{
		"menu":         system.GetMenu(),
		"list":         list,
		"qid":          qid,
		"question":     question,
		"alerttitle":   title,
		"alertcontext": content,
	})
}

// DelOptions 删除问题选项
func DelOptions(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	qid := c.Query("qid")
	if qid == "" || !common.CheckInt(qid) {
		middleware.RedirectErr("question", common.AlertError, common.AlertParamsError, c)
		return
	}
	path := fmt.Sprintf("options?qid=%s&", qid)
	if id == "" || !common.CheckInt(qid) {
		middleware.RedirectErr2(path, common.AlertError, common.AlertParamsError, c)
		return
	}
	_, err := silkworm.DelOptions(id)
	if err != nil {
		log.Println(err)
		middleware.RedirectErr2(path, common.AlertFail, common.AlertDelFail, c)
		return
	}
	c.Redirect(302, path)
}

// AddOptions 新增问题选项
func AddOptions(c *gin.Context) {
	handelOptions(c, false)
}

// EditOptions 编辑问题选项
func EditOptions(c *gin.Context) {
	handelOptions(c, true)
}

func handelOptions(c *gin.Context, isEdit bool) {
	content := c.PostForm("content")
	qid := c.Query("qid")
	var answer bool
	if c.PostForm("answer") == "1" {
		answer = true
	} else {
		answer = false
	}
	if qid == "" || !common.CheckInt(qid) {
		middleware.RedirectErr("question", common.AlertError, common.AlertParamsError, c)
		return
	}
	path := fmt.Sprintf("options?qid=%s&", qid)
	if content == "" {
		middleware.RedirectErr2(path, common.AlertError, common.AlertParamsError, c)
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	if isEdit {
		id := c.DefaultQuery("id", "")
		if id == "" || !common.CheckInt(id) {
			middleware.RedirectErr2(path, common.AlertError, common.AlertParamsError, c)
			return
		}
		rs := silkworm.EditOptions(content, nowTime, id, qid, answer)
		if !rs {
			middleware.RedirectErr2(path, common.AlertFail, common.AlertSaveFail, c)
			return
		}
		c.Redirect(302, path)
		return
	}
	rs := silkworm.AddOptions(qid, content, nowTime, answer)
	if !rs {
		log.Println("add option fail")
		middleware.RedirectErr2(path, common.AlertFail, common.AlertSaveFail, c)
		return
	}
	c.Redirect(302, path)
}

// UserQuestionList 用户每日题库
func UserQuestionList(c *gin.Context) {
	// log.Println(id)
	openid := c.PostForm("openid")
	if openid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	checkAnswer, _ := silkworm.GetUserAnswers(openid)
	nowDate := time.Now().Local().Format("2006-01-02")
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	nowAnswers := common.CheckLimit(checkAnswer["answers"], checkAnswer["answerdate"], nowDate, 999)
	if nowAnswers == -1 {
		middleware.RespondErr(201, common.Err201Limit, c)
		return
	}
	ip := c.ClientIP()
	_, err := silkworm.UpdateUserAnswer(strconv.Itoa(nowAnswers), nowDate, ip, nowTime, openid)
	if err != nil {
		log.Println(err)
	}
	list := silkworm.GetRandomQuestion()
	newList := common.ChangeMapInterface(list)
	for i := 0; i < len(newList); i++ {
		qid := newList[i]["id"].(string)
		options, _ := silkworm.GetOptionsContent(qid)
		itemname, _ := silkworm.ItemInfo(newList[i]["itemid"].(string))
		newList[i]["options"] = options
		newList[i]["itemname"] = itemname["name"]
	}
	c.JSON(200, gin.H{
		"msg":  "success",
		"list": newList,
	})
}

// CheckAnswer 检查答案
func CheckAnswer(c *gin.Context) {
	openid := c.PostForm("openid")
	qid := c.PostForm("questionid")
	optionsid := c.PostForm("optionid")
	if openid == "" || qid == "" || optionsid == "" {
		middleware.RespondErr(402, common.Err402Param, c)
		return
	}
	rs, _ := silkworm.CheckAnswer(qid)
	rightAnswer := rs["answer"]
	itemid := rs["itemid"]
	if optionsid != rightAnswer {
		c.JSON(200, gin.H{
			"msg":         "success",
			"results":     false,
			"rightAnswer": rightAnswer,
		})
		return
	}
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	itemInfo, _ := silkworm.ItemInfo(itemid)
	AddItemToRucksack(silkworm.ActiveAnswer, openid, itemid, nowTime, "", itemInfo)
	c.JSON(200, gin.H{
		"msg":      "success",
		"results":  true,
		"itemname": itemInfo["name"],
	})
}
