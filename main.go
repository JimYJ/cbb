package main

import (
	"canbaobao/common"
	"canbaobao/route"
	"canbaobao/service"
	"log"
)

func main() {
	// inits()
	service.HourTimer()
}

func inits() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	common.GetConfig()
	common.InitMysql()
	common.GetMysqlConn()
	route.Web()
	route.API()
}
