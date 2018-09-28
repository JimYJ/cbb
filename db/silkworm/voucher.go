package silkworm

import (
	"canbaobao/common"
	log "canbaobao/service/logs"
	"fmt"
	"time"

	"github.com/JimYJ/easysql/mysql"
)

// AddVoucher 新增兑换券
func AddVoucher(vendorid, uid, content, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into voucher set vendorid = ?,uid = ?,content = ?,status = ?,createtime = ?,updatetime = ?", vendorid, uid, content, 0, nowTime, nowTime)
}

// EditVoucher 使用兑换券(修改兑换券状态为使用)
func EditVoucher(vid, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	if vid == "0" {
		return mysqlConn.Update(mysql.Statement, "update voucher set status = ?,updatetime = ? where id = ?", 1, nowTime, id)
	}
	return mysqlConn.Update(mysql.Statement, "update voucher set status = ?,updatetime = ? where id = ?", 1, nowTime, id, vid)
}

func getVoucher(vid string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	if vid == "0" {
		return mysqlConn.GetResults(mysql.Statement, "select voucher.id,vendorid,content, uid, voucher.status,startday,endday,vtype,voucher.updatetime,`user`.name,vendor.`name` as vendorname from voucher LEFT JOIN `user` on uid = `user`.id LEFT JOIN vendor on vendorid = vendor.id ORDER BY voucher.id desc")
	}
	return mysqlConn.GetResults(mysql.Statement, "select voucher.id,vendorid,content, uid, status,startday,endday,vtype,voucher.updatetime,`user`.name,vendor.`name` as vendorname from voucher LEFT JOIN `user` on uid = `user`.id LEFT JOIN vendor on vendorid = vendor.id  where vendorid = ? ORDER BY voucher.id desc", vid)
}

// GetVoucher 获取兑换券
func GetVoucher(vid string) ([]map[string]string, error) {
	list, err := getVoucher(vid)
	nowDate, _ := time.Parse("2006-01-02", time.Now().Local().Format("2006-01-02"))
	for i := 0; i < len(list); i++ {
		if list[i]["status"] == "0" {
			list[i]["statustr"] = "未使用"
		} else {
			list[i]["statustr"] = "已使用"
		}
		if list[i]["vtype"] == "1" {
			list[i]["vtypestr"] = "有效"
			list[i]["vtypeint"] = "0"
			startDay, _ := time.Parse("2006-01-02", list[i]["startday"])
			endDay, _ := time.Parse("2006-01-02", list[i]["endday"])
			if nowDate.Sub(startDay) < 0 {
				list[i]["vtypestr"] = "未启用"
				list[i]["vtypeint"] = "1"
				list[i]["status"] = "1"
			} else if nowDate.Sub(endDay) > 0 {
				list[i]["vtypestr"] = "已失效"
				list[i]["vtypeint"] = "2"
				list[i]["status"] = "1"
			}
		} else {
			list[i]["vtypestr"] = "长期有效"
			list[i]["vtypeint"] = "0"
		}
		if list[i]["vendorid"] == "0" {
			list[i]["vendorname"] = "数据错误"
		}
		if list[i]["name"] == "" {
			list[i]["name"] = "用户已被删除"
		}
	}
	return list, err
}

// getVoucherByUser 用户兑换券
func getVoucherByUser(uid, paginaSQL string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, fmt.Sprintf("select voucher.id,content,voucher.createtime,vendor.name,vendorid,startday,endday,vtype from voucher left join vendor on vendorid = vendor.id where voucher.uid = ? and `status` = ? ORDER BY id desc %s", paginaSQL), uid, 0)
}

// GetVoucherByUser 用户兑换券
func GetVoucherByUser(uid, paginaSQL string) ([]map[string]string, error) {
	list, err := getVoucherByUser(uid, paginaSQL)
	nowDate, _ := time.Parse("2006-01-02", time.Now().Local().Format("2006-01-02"))
	for i := 0; i < len(list); i++ {
		if list[i]["vtype"] == "1" {
			list[i]["vtypestr"] = "有效"
			startDay, _ := time.Parse("2006-01-02", list[i]["startday"])
			endDay, _ := time.Parse("2006-01-02", list[i]["endday"])
			if nowDate.Sub(startDay) < 0 {
				list[i]["vtypestr"] = "未启用"
			} else if nowDate.Sub(endDay) > 0 {
				list[i]["vtypestr"] = "已失效"
			}
		} else {
			list[i]["vtypestr"] = "长期有效"
		}
	}
	return list, err
}

// GetVoucherByUserCount 用户兑换券总数
func GetVoucherByUserCount(uid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from voucher where uid = ? and `status` = ? ORDER BY id desc", uid, 0)
}

// ExchangeGoods 兑换商品(生成兑换券)
func ExchangeGoods(vendorid, uid, content, nowTime string, idList []string) bool {
	mysqlConn := common.GetMysqlConn()
	rs := true
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	_, err = tx.Insert("insert into voucher set vendorid = ?,uid = ?,content = ?,status = ?,createtime = ?,updatetime = ?", vendorid, uid, content, 0, nowTime, nowTime)
	for i := 0; i < len(idList); i++ {
		a, err := tx.Delete("delete from usersw where id = ?", idList[i])
		if err != nil || a == 0 {
			log.Println(err)
			rs = false
			break
		}
	}
	if !rs || err != nil {
		tx.Rollback()
		return false
	}
	tx.Commit()
	return rs
}

// ExchangeGoodsII 兑换商品-方案2(生成兑换券)
func ExchangeGoodsII(vendorid, uid, content, nowTime string, idList []string) bool {
	mysqlConn := common.GetMysqlConn()
	rs := true
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	_, err = tx.Insert("insert into voucher set vendorid = ?,uid = ?,content = ?,status = ?,createtime = ?,updatetime = ?", vendorid, uid, content, 0, nowTime, nowTime)
	for i := 0; i < len(idList); i++ {
		a, err := tx.Delete("delete from usersw where id = ?", idList[i])
		if err != nil || a == 0 {
			log.Println(err)
			rs = false
			break
		}
	}
	if !rs || err != nil {
		tx.Rollback()
		return false
	}
	tx.Commit()
	return rs
}

// BatchVoucher 批量发放兑换券
func BatchVoucher(ulist []map[string]string, nowTime, startDay, endDay, content string) bool {
	mysqlConn := common.GetMysqlConn()
	rs := true
	tx, err := mysqlConn.Begin()
	if err != nil {
		log.Println("begin tx fail", err)
		return false
	}
	var a int64
	for i := 0; i < len(ulist); i++ {
		a, err = tx.Insert("insert into voucher set vendorid = ?,uid = ?,content = ?,status = ?,startday = ?,endday = ?,vtype = ?,createtime = ?,updatetime = ?",
			ulist[i]["vid"], ulist[i]["id"], content, 0, startDay, endDay, 1, nowTime, nowTime)
		if err != nil || a == 0 {
			log.Println(err)
			rs = false
			break
		}
		a = 0
	}
	if !rs || err != nil {
		tx.Rollback()
		return false
	}
	tx.Commit()
	return rs
}

// BatchVoucherActive 批量发放兑换券的用户动态
func BatchVoucherActive(ulist []map[string]string, nowTime, startDay, endDay, content string) {
	dateStr := fmt.Sprintf("%s - %s", startDay, endDay)
	for i := 0; i < len(ulist); i++ {
		_, err := SaveUserActive(ActiveBatchVoucher, ulist[i]["name"], ulist[i]["id"], content, "0", nowTime, dateStr)
		if err != nil {
			log.Println("Save User Active Fail:", err, "user info:", ulist[i]["name"], ulist[i]["id"])
		}
		// time.Sleep(time.Millisecond)
	}
}
