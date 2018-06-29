package system

import (
	"canbaobao/common"
	"github.com/JimYJ/easysql/mysql"
	"log"
)

// GetMenulist 获取多级菜单
func GetMenulist() ([]map[string]interface{}, error) {
	list, err := GetMainMenu("0")
	if err != nil {
		return nil, err
	}
	rs := make([]map[string]interface{}, len(list))
	for i := 0; i < len(list); i++ {
		rs[i] = make(map[string]interface{})
		for k, v := range list[i] {
			rs[i][k] = v
		}
		temp, err := GetMainMenu(list[i]["id"])
		if err != nil {
			log.Println(err)
			rs[i]["list"] = nil
		} else {
			rs[i]["list"] = temp
		}
	}
	return rs, nil
}

// GetMainMenu 获取分类菜单 0父级
func GetMainMenu(parentid string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,path,icon from bms_menu where deleted = ? and parentid = ? ORDER BY sort", 0, parentid)
}

// SetMenu 从数据库加载菜单到缓存
func SetMenu(token string) {
	list, err := GetMenulist()
	if err != nil {
		log.Println(err)
	}
	ulist := GetUserMenuList(token)
	for i := 0; i < len(list); i++ {
		if list[i]["list"] != nil {
			sublist := list[i]["list"].([]map[string]string)
			if sublist != nil {
				templist := make([]map[string]string, 0)
				for j := 0; j < len(sublist); j++ {
					_, ok := ulist[sublist[j]["path"]]
					if ok {
						templist = append(templist, sublist[j])
					}
				}
				list[i]["list"] = templist
			} else {
				list[i]["list"] = nil
			}
		}
	}
	cache := common.GetCache()
	cache.Set(common.Sysmenu, list, -1)
}

// GetMenu 获取菜单用于HTML渲染
func GetMenu() []map[string]interface{} {
	cache := common.GetCache()
	menu, err := cache.Get(common.Sysmenu)
	if !err {
		menu = nil
	}
	if menu == nil {
		return nil
	}
	return menu.([]map[string]interface{})
}

//  获取分类菜单
func getAllMenu(parentid string) ([]map[string]string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetResults(mysql.Statement, "select id,name,path,parentid,icon,createtime,updatetime,deleted from bms_menu where deleted = ? and parentid = ? ORDER BY sort", 0, parentid)
}

//  获取子菜单数量
func getSubMenuCount(parentid string) (string, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.GetVal(mysql.Statement, "select count(*) from bms_menu where deleted = ? and parentid = ? ORDER BY sort", 0, parentid)
}

// GetAllMenu 处理菜单详细列表
func GetAllMenu(parentid string) []map[string]string {
	list, err := getAllMenu(parentid)
	if err != nil {
		return nil
	}
	mysqlConn := common.GetMysqlConn()
	for i := 0; i < len(list); i++ {
		if list[i]["deleted"] == "1" {
			list[i]["delete"] = "是"
		} else {
			list[i]["delete"] = "否"
		}
		if list[i]["parentid"] != "0" {
			name, err := mysqlConn.GetVal(mysql.Statement, "select name from bms_menu where id = ?", list[i]["parentid"])
			if err != nil {
				log.Println(err)
			}
			if name != "" {
				list[i]["parentname"] = name
			} else {
				list[i]["parentname"] = "父级菜单"
				count, err := getSubMenuCount(list[i]["id"])
				if err != nil {
					log.Println(err)
					list[i]["subcount"] = "0"
				} else {
					list[i]["subcount"] = count
				}
			}
		} else {
			list[i]["parentname"] = "父级菜单"
			count, err := getSubMenuCount(list[i]["id"])
			if err != nil {
				log.Println(err)
				list[i]["subcount"] = "0"
			} else {
				list[i]["subcount"] = count
			}
		}
	}
	return list
}

// DelMenu 删除菜单
func DelMenu(id, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Delete(mysql.Statement, "update bms_menu set deleted = ?,updatetime = ? where id = ?", 1, nowTime, id)
}

// AddMenu 新增菜单
func AddMenu(name, path, parentid, icon, nowTime string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "insert into bms_menu set name = ?,path = ?,parentid = ?,icon = ?,createtime = ?,updatetime = ?", name, path, parentid, icon, nowTime, nowTime)
}

// EditMenu 编辑菜单
func EditMenu(name, path, parentid, icon, nowTime, id string) (int64, error) {
	mysqlConn := common.GetMysqlConn()
	return mysqlConn.Insert(mysql.Statement, "update bms_menu set name = ?,path = ?,parentid = ?,icon = ?,updatetime = ? where id = ?", name, path, parentid, icon, nowTime, id)
}

// ChangeSort 更改排序 upordown:true up|false down
func ChangeSort(id, parentid string, upordown bool) bool {
	list, err := GetMainMenu(parentid)
	if err != nil {
		log.Println(err)
		return false
	}
	var j int
	for i := 0; i < len(list); i++ {
		if list[i]["id"] == id {
			j = i
			break
		}
	}
	if upordown {
		if j == 0 {
			return false
		}
		list[j]["id"], list[j-1]["id"] = list[j-1]["id"], list[j]["id"]
	} else {
		if j == len(list)-1 {
			return false
		}
		list[j]["id"], list[j+1]["id"] = list[j+1]["id"], list[j]["id"]
	}
	mysqlConn := common.GetMysqlConn()
	mysqlConn.TxBegin()
	for i := 0; i < len(list); i++ {
		_, err = mysqlConn.TxUpdate(mysql.Statement, "update bms_menu set sort = ? where id = ?", i, list[i]["id"])
		if err != nil {
			break
		}
	}
	if err != nil {
		mysqlConn.TxRollback()
		log.Println("change sort fail:", err)
		return false
	}
	mysqlConn.TxCommit()
	return true
}
