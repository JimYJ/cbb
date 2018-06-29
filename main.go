package main

import (
	"canbaobao/common"
	"canbaobao/route"
	"log"
)

func main() {
	inits()
}

func inits() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	common.GetConfig()
	common.InitMysql()
	common.GetMysqlConn()
	route.Web()
	route.API()
}
