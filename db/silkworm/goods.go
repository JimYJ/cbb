package silkworm

import (
	"canbaobao/common"
	"fmt"

	"github.com/JimYJ/easysql/mysql"
)

// DelGoods 删除商品
func DelGoods(id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete(mysql.Statement, "delete from goods where id = ?", id)
}

// GetGoodsName 获得商品名称
func GetGoodsName(id string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select name from goods where id = ?", id)
}

// AddGoods 新增商品
func AddGoods(name, bigimg, content, swcount, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into goods set name = ?,bigimg = ?,swcount = ?,content = ?,createtime = ?,updatetime = ?", name, bigimg, swcount, content, nowTime, nowTime)
}

// EditGoods 编辑商品
func EditGoods(name, bigimg, content, swcount, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	if bigimg == "" {
		return mysqlConn.Insert(mysql.Statement, "update goods set name = ?,swcount = ?,content = ?,updatetime = ? where id = ?", name, swcount, content, nowTime, id)
	}
	return mysqlConn.Insert(mysql.Statement, "update goods set name = ?,bigimg = ?,swcount = ?,content = ?,updatetime = ? where id = ?", name, bigimg, swcount, content, nowTime, id)

}

// GetGoods 获取商品
func GetGoods() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,bigimg,swcount,createtime,updatetime from goods ORDER BY id desc")
}

// GetPaginaGoods 获取分页商品
func GetPaginaGoods(paginaSQL string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	sql := fmt.Sprintf("select id,name,bigimg,swcount,createtime,updatetime from goods ORDER BY id desc %s", paginaSQL)
	return mysqlConn.GetResults(mysql.Statement, sql)
}

// GetGoodsCount 获取商品总数
func GetGoodsCount() (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from goods ORDER BY id desc")
}

// GetGoodsContent 获取商品内容
func GetGoodsContent(id string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select content from goods where id = ?", id)
}

// GetGoodsExchange 获取商品兑换方案2-条件
func GetGoodsExchange(id string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select swcount from goods where id = ?", id)
}

// DelGoodsRedeem 删除商品兑换条件
func DelGoodsRedeem(id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete(mysql.Statement, "delete from goodsredeem where id = ?", id)
}

// AddGoodsRedeem 新增商品兑换条件
func AddGoodsRedeem(butterflyid, numbers, gid, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into goodsredeem set butterflyid = ?,numbers = ?,goodsid = ?,createtime = ?,updatetime = ?", butterflyid, numbers, gid, nowTime, nowTime)
}

// EditGoodsRedeem 编辑商品兑换条件
func EditGoodsRedeem(gid, butterflyid, numbers, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update goodsredeem set butterflyid = ?,numbers = ?,updatetime = ? where id = ?", butterflyid, numbers, nowTime, id)
}

// GetGoodsRedeem 获取商品兑换条件(蝴蝶)表
func GetGoodsRedeem(goodsid string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select goodsredeem.id,name,butterflyid,numbers,goodsredeem.createtime,goodsredeem.updatetime from goodsredeem left join butterfly on butterflyid = butterfly.id where goodsid = ? ORDER BY id", goodsid)
}

// CheckRepeat 检查兑换条件是否重复
func CheckRepeat(gid, butterflyid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from goodsredeem where goodsid = ? and butterflyid = ?", gid, butterflyid)
}

// GetGoodsRedeemList 获取商品兑换条件ID
func GetGoodsRedeemList(goodsid string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select butterflyid,numbers from goodsredeem where goodsid = ? ORDER BY id", goodsid)
}
