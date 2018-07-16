package log

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"net/http/httputil"
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
		if logFile.date != time.Now().Minute() {
			logFile.date = time.Now().Minute()
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
			Infof("| %3d | %13v | %15s | %-7s %s", statusCode, clientIP, latency, method, path)
		} else if statusCode >= 400 && statusCode < 500 {
			Warnf("| %3d | %13v | %15s | %-7s %s", statusCode, clientIP, latency, method, path)
		} else {
			Errorf("| %3d | %13v | %15s | %-7s %s", statusCode, clientIP, latency, method, path)
		}
	}
}

// Recovery 异常日志分割
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				Errorf("[Recovery] %s panic recovered:\n%s\n%s\n%s%s", time.Now().Format("2006/01/02 - 15:04:05"), string(httprequest), err, stack, reset)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())

	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// Println 输出日常日志
func Println(args ...interface{}) {
	logFile.info.Println(args...)
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
	logFile.fileName = fmt.Sprintf("%s%s%s", "./logs/", time.Now().Format("2006-01-02_15-04"), ".log")
	logF, err = os.OpenFile(logFile.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	logFile.fileFd = logF
	logFile.date = time.Now().Hour()
	logFile.info = log.New(io.MultiWriter(os.Stdout, logFile.fileFd), "[GIN]", log.Ldate|log.Ltime|log.Lshortfile)
	logFile.warn = log.New(io.MultiWriter(logFile.fileFd), "[GIN]", log.Ldate|log.Ltime|log.Lshortfile)
	logFile.err = log.New(io.MultiWriter(os.Stderr, logFile.fileFd), "[GIN]", log.Ldate|log.Ltime|log.Lshortfile)
	logFile.debug = log.New(io.MultiWriter(logFile.fileFd), "[GIN]", log.Ldate|log.Ltime|log.Lshortfile)
}
