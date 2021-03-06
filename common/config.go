package common

import (
	"fmt"
	"github.com/JimYJ/easysql/mysql"
	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var (
	dbhost, dbname, user, pass string
	port                       int
	once                       sync.Once
	c                          *cache.Cache
	AppPath                    string
	// TokenTimeOut 登陆超时
	TokenTimeOut = 2 * time.Hour
	// CacheTimeOut 缓存超时
	CacheTimeOut = 1 * time.Minute
	// AppID ...
	AppID = ""
	// SecretKey ...
	SecretKey = ""
	sha1salt  = []byte("63d81bc836e86565a5e8668faf1863cbbcd5b392fba28d7d48b39c858b3e4e75")
	//LoginMaxLimit 每分钟登陆请求限制
	LoginMaxLimit = 30
	//LoginGap 计次时间间隔，单位是秒
	LoginGap = 60
	// LogsPath 日志目录
	LogsPath = ""
)

// 错误信息
var (
	// Err201Limit 操作数达到上限
	Err201Limit = "The upper limit has been reached today"
	// Err202Limit 浇水施肥达到本级别上限
	Err202Limit = "The upper limit has been reached this level"
	// Err203Limit 同一时间只可孵化一个
	Err203Limit = "Only one hatch at the same time"
	// Err205Limit 申请的蚕宝宝未孵化成蝴蝶，无法参与配对
	Err205Limit = "the requested silkworm did not hatch into butterflies and could not participate in pairing."
	// Err206Limit 申请的蝴蝶已经申请配对或正在配对中，无法参与配对
	Err206Limit = "The butterfly has already applied for matching or pairing, unable to participate in matching."
	//Err401 认证错误
	Err401 = "Authentication error!"
	//Err401login 登录失败
	Err401login = "User or Pass error!"
	//Err401captcha 验证码错误
	Err401captcha = "captcha is error!"
	//Err401SmsCode 登录失败
	Err401SmsCode = "Sms code is error!"
	//Err402Param 参数不正确
	Err402Param = "param is error!"
	//Err406Unexpected 异常错误
	Err406Unexpected = "request Unexpected!"
	// Err407Limit 桑树已达最大等级
	Err407Limit = "tree has reached the maximum level"
	//Err412UserNotBind 用户未绑定
	Err412UserNotBind = "User Is Not Bind Vendor!"
	//Err413UserItemNoExist 物品不存在背包中
	Err413UserItemNoExist = "item is no exist!"
	//Err414UserIsBind 用户已绑定店铺
	Err414UserIsBind = "User is Bind Vendor!"
	//Err415CannotPairSelf 不允许和自己配对
	Err415CannotPairSelf = "cannot pair self!"
	//Err416NotTaken 未拾取物品不能使用
	Err416NotTaken = "cannot use not taken item!"
	//Err417OtherUserItem 不能使用他人物品
	Err417OtherUserItem = "cannot use others user item!"
	//Err418OtherUserSW 不能喂养他人蚕仔
	Err418OtherUserSW = "cannot feed others user silkworm!"
	//Err419ItemCannotFeed 该物品不能用来喂养
	Err419ItemCannotFeed = "this item cannot use for feed!"
	//Err420CannotFeedButterfly 不能喂养蝴蝶
	Err420CannotFeedButterfly = "feed to silkworm only!"
	//Err421SilkwormHatching 蚕宝宝还在孵化中，无法喂养
	Err421SilkwormHatching = "silkworm still hatching, cannot feed it!"
	//Err422NotEnoughExchange 兑换条件不足
	Err422NotEnoughExchange = "insufficient exchange conditions!"
	//Err423PairingCannotExchange 蝴蝶配对中无法兑换商品
	Err423PairingCannotExchange = "butterfly is pairing,it cannot exchange goods!"
	//Err429Frequent 请求过于频繁
	Err429Frequent = "Request too Frequent!"
	//Err500DBrequest 数据库请求错误
	Err500DBrequest = "Database request error!"
	//Err500DBSave 数据库保存失败
	Err500DBSave = "Database save fail!"
	//Err500CannotGetUID 数据库请求错误
	Err500CannotGetUID = "Cannot get UID by token!"
	//Err502Wechat 微信接口请求错误
	Err502Wechat = "request wechat api fail!"
)

// http状态码
var (
	//HTTPAuthErr 认证错误
	HTTPAuthErr = 401
	//HTTPParamErr 请求参数错误
	HTTPParamErr = 402
	//HTTPForbiddenErr 拒绝请求
	HTTPForbiddenErr = 403
	//HTTPForbiddenErr 请求异常
	HTTPUnexpectedErr = 406
	//HTTPFrequentErr 请求过于频繁
	HTTPFrequentErr = 429
	//HTTPSystemErr 系统内部错误(程序，数据库等)
	HTTPSystemErr = 500
	//HTTPExternalErr 系统外部错误(第三方)
	HTTPExternalErr = 502
)

type config struct {
	Mysql     mysqlconf
	AppID     string
	SecretKey string
	AppPath   string
	LogsPath  string
}

type mysqlconf struct {
	Host   string
	Port   int
	DBname string
	User   string
	Pass   string
}

func (conf *config) getConfig() *config {
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Panicf("yamlFile.Get err   #%v ", err)
		return nil
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Panicf("yamlFile.Get err   #%v ", err)
		return nil
	}
	return conf
}

// GetConfig 获取配置文件
func GetConfig() {
	var conf config
	conf.getConfig()
	dbhost = conf.Mysql.Host
	dbname = conf.Mysql.DBname
	port = conf.Mysql.Port
	user = conf.Mysql.User
	pass = conf.Mysql.Pass
	AppID = conf.AppID
	SecretKey = conf.SecretKey
	AppPath = conf.AppPath
	LogsPath = conf.LogsPath
}

// InitMysql 初始化mysql参数
func InitMysql() {
	mysql.Init(dbhost, port, dbname, user, pass, "utf8mb4", 100, 100)
}

// GetMysqlConn 获取mysql连接
func GetMysqlConn() *mysql.MysqlDB {
	mysqlConn, err := mysql.GetMysqlConn()
	if err != nil {
		log.Panicln("mysql conn error:", err)
		return nil
	}
	return mysqlConn
}

// GetCache 获得缓存对象
func GetCache() *cache.Cache {
	once.Do(func() {
		c = cache.New(TokenTimeOut, 10*time.Minute)
	})
	return c
}

// GetKeyName 获得AccessToken和RefreshToken的缓存KEY
func GetKeyName(openid string) (string, string) {
	return fmt.Sprintf("OpenIDAccessToken:%s", openid), fmt.Sprintf("OpenIDRefreshToken:%s", openid)
}
