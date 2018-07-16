package main

import (
	"canbaobao/common"
	"canbaobao/route"
	"canbaobao/service"
	// log "canbaobao/service/logs"
)

func main() {
	inits()
}

func inits() {
	// log.SetFlags(log.Lshortfile | log.LstdFlags)
	common.GetConfig()
	common.InitMysql()
	common.GetMysqlConn()
	go service.HourTimer()
	go service.DayTimer()
	route.Web()
	route.API()
}
