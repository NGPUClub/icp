package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SendError(c *gin.Context, codeId int, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"codeId": codeId,
		"msg":    msg,
	})
}

func SendSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"codeId": 200,
		"msg":    "success",
	})
}

func ClientIP(c *gin.Context) string {
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

func TaskID() string {
	return fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), fmt.Sprintf("%06d", time.Now().Nanosecond()/1e3))
}
