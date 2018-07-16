package system

import (
	"canbaobao/common"
	log "canbaobao/service/logs"
	"github.com/JimYJ/easysql/mysql"
)

// GetRole 获取角色简单列表
func GetRole() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name from bms_role where deleted = ? ORDER BY id", 0)
}

// DelRole 删除后台管理角色
func DelRole(id, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete(mysql.Statement, "update bms_role set deleted = ?,updatetime = ? where id = ?", 1, nowTime, id)
}

// AddRole 新增后台管理角色
func AddRole(name, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into bms_role set name = ?,createtime = ?,updatetime = ?", name, nowTime, nowTime)
}

// EditRole 编辑后台管理角色
func EditRole(name, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update bms_role set name = ?,updatetime = ? where id = ?", name, nowTime, id)
}

//  获取全部角色
func getAllRole() ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,createtime,updatetime,deleted from bms_role where deleted = ? ORDER BY id", 0)
}

// GetAllRole 处理角色详细列表
func GetAllRole() []map[string]string {
	list, err := getAllRole()
	if err != nil {
		return nil
	}
	for i := 0; i < len(list); i++ {
		if list[i]["deleted"] == "1" {
			list[i]["delete"] = "是"
		} else {
			list[i]["delete"] = "否"
		}
	}
	return list
}

// GetRoleMenu 获取角色绑定的权限
func GetRoleMenu(id string) []map[string]string {
	mysqlConn := common.GetMysqlConn()
	list, err := mysqlConn.GetResults(mysql.Statement, "select menuid from bms_permiassion where roleid = ? order by id", id)
	if err != nil {
		return nil
	}
	return list
}

// GetRoleMenuPath 获取角色绑定的可访问菜单路径
func GetRoleMenuPath(id string) []map[string]string {
	mysqlConn := common.GetMysqlConn()
	list, err := mysqlConn.GetResults(mysql.Statement, "select path from bms_permiassion left join bms_menu on menuid = bms_menu.id where roleid = ? order by bms_menu.id", id)
	if err != nil {
		return nil
	}
	return list
}

// RoleBindMenu 绑定权限给角色
func RoleBindMenu(id, nowTime string, list []string) error {
	mysqlConn := common.GetMysqlConn()
	mysqlConn.TxBegin()
	_, err := mysqlConn.TxDelete(mysql.Statement, "delete from bms_permiassion where roleid = ?", id)
	if err != nil {
		mysqlConn.TxRollback()
		return err
	}
	for i := 0; i < len(list); i++ {
		if !common.CheckInt(list[i]) {
			break
		}
		_, err = mysqlConn.TxInsert(mysql.Statement, "insert into bms_permiassion set roleid = ?,menuid = ?,createtime = ?,updatetime = ?", id, list[i], nowTime, nowTime)
		if err != nil {
			log.Println(err)
			break
		}
	}
	if err != nil {
		mysqlConn.TxRollback()
		return err
	}
	mysqlConn.TxCommit()
	return nil
}
