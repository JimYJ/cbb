package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"log"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ResponBodyLog 记录响应报文日志
func ResponBodyLog(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	statusCode := c.Writer.Status()
	if statusCode > 200 {
		log.Println("Response body:", blw.body.String())
	}
}
