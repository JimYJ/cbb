package service

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"log"
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
		nextHour = time.Now().Local().Add(time.Minute * 1)
		nextHour = time.Date(nextHour.Year(), nextHour.Month(), nextHour.Day(), nextHour.Hour(), nextHour.Minute(), 0, 0, nextHour.Location())
		ht = time.NewTimer(nextHour.Sub(time.Now().Local()))
		select {
		case <-ht.C:
			//整小时执行
			log.Println(time.Now().Local().Format("2006-01-02 15:04:05"))
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
			log.Println("exec time:", time.Now().Local().Format("2006-01-02 15:04:05"))
		}
	}
}

// SproutLeaf 长桑叶
func SproutLeaf() {
	userList, _ := silkworm.GetUserForTimer()
	treeLevelList, _ := silkworm.TreeLevelList()
	for i := 0; i < len(userList); i++ {
		userTreeLevel, _ := strconv.Atoi(userList[i]["treelevel"])
		treeLevelInfo := treeLevelList[userTreeLevel-1]
		limitTimes, _ := strconv.Atoi(treeLevelInfo["maxhours"])
		nowDate := time.Now().Local().Format("2006-01-02")
		nowExecTimes := common.CheckLimit(userList[i]["answers"], userList[i]["answerdate"], nowDate, limitTimes)
		if nowExecTimes == -1 {
			continue
		}
	}
}
