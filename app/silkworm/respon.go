package silkworm

import (
	"github.com/gin-gonic/gin"
)

func responSuccess(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "success",
	})
}

func responSignSuccess(c *gin.Context, days int64) {
	c.JSON(200, gin.H{
		"msg":      "success",
		"signdays": days,
	})
}

func responIntInfoSuccess(c *gin.Context, key string, i int) {
	c.JSON(200, gin.H{
		"msg": "success",
		key:   i,
	})
}
