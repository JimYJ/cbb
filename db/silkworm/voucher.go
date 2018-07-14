package silkworm

import (
	"canbaobao/common"
	"fmt"
	"github.com/JimYJ/easysql/mysql"
	"log"
)

// AddVoucher 新增兑换券
func AddVoucher(vendorid, uid, content, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into voucher set vendorid = ?,uid = ?,content = ?,status = ?,createtime = ?,updatetime = ?", vendorid, uid, content, 0, nowTime, nowTime)
}

// EditVoucher 使用兑换券
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
		return mysqlConn.GetResults(mysql.Statement, "select id,vendorid,content, uid, status,createtime,updatetime from voucher ORDER BY id desc")
	}
	return mysqlConn.GetResults(mysql.Statement, "select id,vendorid,content, uid, status,createtime,updatetime from voucher where vendorid = ? ORDER BY id desc", vid)
}

// GetVoucher 获取兑换券
func GetVoucher(vid string) ([]map[string]string, error) {
	list, err := getVoucher(vid)
	for i := 0; i < len(list); i++ {
		if list[i]["status"] == "0" {
			list[i]["statustr"] = "未使用"
		} else {
			list[i]["statustr"] = "已使用"
		}
		if list[i]["vendorid"] == "0" {
			list[i]["vendor"] = "数据错误"
		} else {
			if !common.CheckInt(list[i]["vendorid"]) {
				list[i]["vendor"] = "绑定错误"
			} else {
				vname, err := GetVendorName(list[i]["vendorid"])
				if err != nil || vname == "" {
					list[i]["vendor"] = "数据错误"
				} else {
					list[i]["vendor"] = vname
				}
			}
		}
		if !common.CheckInt(list[i]["uid"]) {
			list[i]["username"] = "UID错误"
		} else {
			name, err := GetUserName(list[i]["uid"])
			if err != nil || name == "" {
				list[i]["username"] = "用户已被删除"
			} else {
				list[i]["username"] = name
			}
		}
	}
	return list, err
}

// GetVoucherByUser 用户兑换券
func GetVoucherByUser(uid, paginaSQL string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, fmt.Sprintf("select voucher.id,content,voucher.createtime,vendor.name,vendorid from voucher left join vendor on vendorid = vendor.id where voucher.uid = ? and `status` = ? ORDER BY id desc %s", paginaSQL), uid, 0)
}

// GetVoucherByUserCount 用户兑换券总数
func GetVoucherByUserCount(uid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from voucher where uid = ? and `status` = ? ORDER BY id desc", uid, 0)
}

// ExchangeGoods 兑换商品
func ExchangeGoods(vendorid, uid, content, nowTime string, idList []string) bool {
	mysqlConn := common.GetMysqlConn()
	rs := true
	mysqlConn.TxBegin()
	_, err := mysqlConn.TxInsert(mysql.Statement, "insert into voucher set vendorid = ?,uid = ?,content = ?,status = ?,createtime = ?,updatetime = ?", vendorid, uid, content, 0, nowTime, nowTime)
	for i := 0; i < len(idList); i++ {
		a, err := mysqlConn.TxDelete(mysql.Statement, "delete from usersw where id = ?", idList[i])
		if err != nil || a == 0 {
			log.Println(err)
			rs = false
			break
		}
	}
	if !rs || err != nil {
		mysqlConn.TxRollback()
		return false
	}
	mysqlConn.TxCommit()
	return rs
}
