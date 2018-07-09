package service

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	nextHour time.Time
	nextDay  time.Time
	ht       *time.Timer
	dt       *time.Timer
)

// HourTimer 每整点小时定时器
func HourTimer() {
	for {
		nextHour = time.Now().Local().Add(time.Hour * 1)
		nextHour = time.Date(nextHour.Year(), nextHour.Month(), nextHour.Day(), nextHour.Hour(), 0, 0, 0, nextHour.Location())
		ht = time.NewTimer(nextHour.Sub(time.Now().Local()))
		select {
		case <-ht.C:
			//整小时执行
			log.Println("=========start exec hour task==========")
			status()
			sproutLeaf()
			log.Println("=========end exec hour task==========")
		}
	}
}

// DayTimer 每天0点小时定时器
func DayTimer() {
	for {
		nextDay = time.Now().Local().Add(time.Hour * 24)
		nextDay = time.Date(nextHour.Year(), nextHour.Month(), nextHour.Day(), 0, 0, 0, 0, nextHour.Location())
		dt = time.NewTimer(nextHour.Sub(time.Now().Local()))
		select {
		case <-dt.C:
			//每天0点执行
			log.Println("=========start exec day task==========")
			log.Println("=========end exec day task==========")
		}
	}
}

// sproutLeaf 生长桑叶
func sproutLeaf() {
	userList, _ := silkworm.GetUserForTimer()
	treeLevelList, _ := silkworm.TreeLevelList()
	nowDate := time.Now().Local().Format("2006-01-02")
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	for i := 0; i < len(userList); i++ {
		userTreeLevel, _ := strconv.Atoi(userList[i]["treelevel"])
		uname := userList[i]["name"]
		treeLevelInfo := treeLevelList[userTreeLevel-1]
		limitTimes, _ := strconv.Atoi(treeLevelInfo["maxhours"])
		growthhours, _ := strconv.Atoi(treeLevelInfo["growthhours"])
		nowExecTimes := common.CheckLimit(userList[i]["sproutleafs"], userList[i]["sproutleafday"], nowDate, limitTimes)
		if nowExecTimes == -1 {
			continue
		}
		rs := silkworm.SproutLeaf("1", userList[i]["id"], nowTime, nowDate, nowExecTimes, growthhours)
		if !rs {
			log.Println("sprout Leaf for user fail, uid:", userList[i]["id"])
		}
		slUpActive(uname, userList[i]["id"], nowTime, strconv.Itoa(growthhours))
	}
}

// slUpActive 生长桑叶
func slUpActive(uname, uid, nowTime, moreInfo string) {
	_, err := silkworm.SaveUserActive(silkworm.ActiveSproutLeaf, uname, uid, "桑叶", "1", nowTime, moreInfo)
	if err != nil {
		log.Println("Save User Active Fail:", err)
	}
}
func status() {
	rs, err := system.Status()
	if err != nil || rs != "1" {
		os.Exit(0)
	}
}
