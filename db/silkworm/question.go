package silkworm

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
	"log"
)

// DelQuestion 删除问题
func DelQuestion(id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete(mysql.Statement, "delete from question where id = ?", id)
}

// GetQuestionName 获得问题内容
func GetQuestionName(id string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select content from question where id = ?", id)
}

// GetQuestionAnswer 获得问题答案ID
func GetQuestionAnswer(id string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select answer from question where id = ?", id)
}

// AddQuestion 新增问题
func AddQuestion(itemid, content, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into question set content = ?,itemid = ?,createtime = ?,updatetime = ?", content, itemid, nowTime, nowTime)
}

// EditQuestion 编辑问题
func EditQuestion(itemid, content, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update question set itemid = ?,content = ?,updatetime = ? where id = ?", itemid, content, nowTime, id)
}

// EditQuestionAnswer 编辑问题答案
func EditQuestionAnswer(optionsid, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update question set optionsid = ?,updatetime = ? where id = ?", optionsid, nowTime, id)
}

// GetQuestion 获取问题
func GetQuestion() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select question.id,name,question.itemid,question.content,question.createtime,question.updatetime from question left join item on itemid = item.id ORDER BY question.id desc")
}

// DelOptions 删除选项
func DelOptions(id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete(mysql.Statement, "delete from options where id = ?", id)
}

// AddOptions 新增选项
func AddOptions(qid, content, nowTime string, answer bool) bool {
	mysqlConn := common.GetMysqlConn()
	var a int
	mysqlConn.TxBegin()
	oldanswer, err4 := mysqlConn.TxGetVal(mysql.Statement, "select answer from question where id = ?", qid)
	if err4 != nil {
		log.Println(err4)
	}
	commit := true
	if answer {
		a = 1
		var err, err2, err3 error
		id, err := mysqlConn.TxInsert(mysql.Statement, "insert into options set qid = ?,content = ?,answer = ?,createtime = ?,updatetime = ?", qid, content, a, nowTime, nowTime)
		_, err2 = mysqlConn.TxUpdate(mysql.Statement, "update question set answer = ?,updatetime = ? where id = ?", id, nowTime, qid)
		if oldanswer != "" {
			_, err3 = mysqlConn.TxUpdate(mysql.Statement, "update options set answer = ?,updatetime = ? where id = ?", 0, nowTime, oldanswer)
		}
		if err != nil || err2 != nil || err3 != nil {
			log.Println(err, err2, err3)
			commit = false
		}
	} else {
		a = 0
		_, err := mysqlConn.TxInsert(mysql.Statement, "insert into options set qid = ?,content = ?,answer = ?,createtime = ?,updatetime = ?", qid, content, a, nowTime, nowTime)
		if err != nil {
			log.Println(err)
			commit = false
		}
	}
	if !commit {
		mysqlConn.TxRollback()
	} else {
		mysqlConn.TxCommit()
	}
	return commit
}

// EditOptions 编辑选项
func EditOptions(content, nowTime, id, qid string, answer bool) bool {
	mysqlConn := common.GetMysqlConn()
	var a int
	var err, err2, err3 error
	mysqlConn.TxBegin()
	oldanswer, _ := mysqlConn.TxGetVal(mysql.Statement, "select answer from question where id = ?", qid)
	if answer {
		a = 1
		_, err = mysqlConn.TxUpdate(mysql.Statement, "update question set answer = ?,updatetime = ? where id = ?", id, nowTime, qid)
	} else {
		a = 0
	}
	_, err2 = mysqlConn.TxUpdate(mysql.Statement, "update options set content = ?,answer = ?,updatetime = ? where id = ?", content, a, nowTime, id)
	if oldanswer != "" {
		_, err3 = mysqlConn.TxUpdate(mysql.Statement, "update options set answer = ?,updatetime = ? where id = ?", 0, nowTime, oldanswer)
	}
	if err != nil || err2 != nil || err3 != nil {
		log.Println(err, err2, err3)
		mysqlConn.TxRollback()
		return false
	}
	mysqlConn.TxCommit()
	return true
}

// GetOptions 获取选项
func GetOptions(qid string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,content,answer,createtime,updatetime from options where qid = ? ORDER BY id", qid)
}

// GetOptionsContent 获取选项内容
func GetOptionsContent(qid string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,content from options where qid = ? ORDER BY id", qid)
}

// GetOptionsCount 获取选项数
func GetOptionsCount(qid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from options where qid = ?", qid)
}

// GetRandomQuestion 随机获取问题
func GetRandomQuestion() []map[string]string {
	mysqlConn := common.GetMysqlConn()
	rs := make([]map[string]string, 3)
	for i := 0; i < 3; i++ {
		rs[i], _ = mysqlConn.GetRow(mysql.Statement, "SELECT t1.id,content,itemid FROM `question` AS t1 JOIN (SELECT ROUND(RAND() * ((SELECT MAX(id) FROM `question`)-(SELECT MIN(id) FROM `question`))+(SELECT MIN(id) FROM `question`)) AS id) AS t2 WHERE t1.id >= t2.id ORDER BY t1.id LIMIT 1")
	}
	return rs
}

// CheckAnswer 检查答案
func CheckAnswer(qid string) (map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetRow(mysql.Statement, "select answer,itemid from question where id = ?", qid)
}
