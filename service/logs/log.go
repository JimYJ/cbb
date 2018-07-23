package log

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	nowMinute int
)

type file struct {
	level    map[string]bool
	fileName string
	date     int
	fileFd   *os.File
	err      *log.Logger
	warn     *log.Logger
	info     *log.Logger
	debug    *log.Logger
	info2    *log.Logger
}

var (
	logFile   file
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
	reset     = string([]byte{27, 91, 48, 109})
	logF      *os.File
)

func init() {
	logFile.level = map[string]bool{"Info": true, "Warn": true, "Debug": true, "Error": true}
	logFile.CreateLogFile()
}

// Logs 日志分割
func Logs() gin.HandlerFunc {
	return func(c *gin.Context) {
		if logFile.date != time.Now().Day() {
			logFile.date = time.Now().Day()
			logFile.CreateLogFile()
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()
		// Stop timer
		latency := time.Since(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		if statusCode >= 200 && statusCode <= 400 {
			Infof("| %3d | %15v | %15s | %-7s %s", statusCode, clientIP, latency, method, path)
		} else if statusCode >= 400 && statusCode < 500 {
			Warnf("| %3d | %15v | %15s | %-7s %s", statusCode, clientIP, latency, method, path)
		} else {
			Errorf("| %3d | %15v | %15s | %-7s %s", statusCode, clientIP, latency, method, path)
		}
	}
}

// Println 输出日常日志
func Println(args ...interface{}) {
	_, filePath, line, _ := runtime.Caller(1)
	str := fmt.Sprintf("%s%s%d", filePath, ":", line)
	a := []interface{}{str}
	a = append(a, args...)
	logFile.info2.Println(a...)
}

// Infof 普通信息输出
func Infof(format string, args ...interface{}) {
	if logFile.level["Info"] {
		logFile.info.Println(fmt.Sprintf(format, args...))
	}
}

// Warnf 警告信息输出
func Warnf(format string, args ...interface{}) {
	if logFile.level["Warn"] {
		logFile.warn.Println(fmt.Sprintf(format, args...))
	}
}

// Debugf 调试信息输出
func Debugf(format string, args ...interface{}) {
	if logFile.level["Debug"] {
		logFile.debug.Println(fmt.Sprintf(format, args...))
	}
}

// Errorf 错误信息输出
func Errorf(format string, args ...interface{}) {
	if logFile.level["Error"] {
		logFile.err.Println(fmt.Sprintf(format, args...))
	}
}

// CreateLogFile 创建日志分割文件
func (m *file) CreateLogFile() {
	if logF != nil {
		logF.Close()
	}
	var err error
	logFile.fileName = fmt.Sprintf("%s%s%s", "./logs/", time.Now().Format("2006-01-02"), ".log")
	logF, err = os.OpenFile(logFile.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	logFile.fileFd = logF
	logFile.date = time.Now().Hour()
	logFile.info2 = log.New(io.MultiWriter(os.Stdout, logFile.fileFd), "[GIN] ", log.Ldate|log.Ltime)
	logFile.info = log.New(io.MultiWriter(os.Stdout, logFile.fileFd), "[GIN] ", log.Ldate|log.Ltime)
	logFile.warn = log.New(io.MultiWriter(os.Stdout, logFile.fileFd), "[GIN] ", log.Ldate|log.Ltime)
	logFile.err = log.New(io.MultiWriter(os.Stderr, logFile.fileFd), "[GIN] ", log.Ldate|log.Ltime)
	logFile.debug = log.New(io.MultiWriter(os.Stderr, logFile.fileFd), "[GIN] ", log.Ldate|log.Ltime)
}
