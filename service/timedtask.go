package service

import (
	"canbaobao/common"
	"canbaobao/db/silkworm"
	"canbaobao/db/system"
	log "canbaobao/service/logs"
	"os"
	"strconv"
	"time"
)

var (
	nextHour time.Time
	nextDay  time.Time
	ht       *time.Timer
	dt       *time.Timer
	logFile  *os.File
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
		// nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), nextDay.Hour(), nextDay.Minute(), 0, 0, nextDay.Location())
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 30, 0, 0, nextDay.Location())
		dt = time.NewTimer(nextDay.Sub(time.Now().Local()))
		select {
		case <-dt.C:
			//每天0点执行
			log.Println("=========start exec day task==========")
			checkHealth()
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

// slUpActive 生长桑叶动态
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

// checkHealth 健康值计算
func checkHealth() {
	list, err := silkworm.GetSWlist()
	if err != nil {
		log.Println(err)
		return
	}
	d, _ := time.ParseDuration("-24h")
	nowDate := time.Now().Local().Add(d).Format("2006-01-02")
	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")
	yesterday, _ := time.Parse("2006-01-02", nowDate)
	updateHealth := make(map[string]int)
	var updateListIndex []int
	for i := 0; i < len(list); i++ {
		health, _ := strconv.Atoi(list[i]["health"])
		if health <= 20 {
			continue
		}
		a := handelFeedDate(list[i]["leafday"], yesterday)
		b := handelFeedDate(list[i]["sppday"], yesterday)
		c := handelFeedDate(list[i]["mppday"], yesterday)
		d := handelFeedDate(list[i]["lppday"], yesterday)
		if !a && !b && !c && !d {
			updateHealth[list[i]["id"]] = calcNewHealth(health)
			updateListIndex = append(updateListIndex, i)
		}
	}
	log.Println(updateHealth, updateListIndex)
	silkworm.UpdateHealth(updateHealth)
	time.Sleep(time.Second * 30)
	silkworm.UpdateHealthActive(&updateListIndex, &list, nowTime)
}

func calcNewHealth(health int) int {
	if health-10 >= 20 {
		return health - 10
	}
	return 20
}

func handelFeedDate(feedDate string, yesterday time.Time) bool {
	if len(feedDate) == 0 {
		return false
	}
	leaf, _ := time.Parse("2006-01-02", feedDate)
	if yesterday.Sub(leaf) > 0 {
		return false
	}
	return true
}
