package handler

import (
	"github.com/gin-gonic/gin"
)

// Hi 公开的健康检查接口，无需认证
func Hi(c *gin.Context) {
	Success(c, gin.H{
		"message": "hi",
	})
}

// Ping 受保护的健康检查接口，返回用户信息
func Ping(c *gin.Context) {
	userInfo, exists := c.Get("user_info")
	if exists {
		userMap := userInfo.(map[string]interface{})
		Success(c, gin.H{
			"message": "pong",
			"user": gin.H{
				"sub":      userMap["sub"],
				"username": userMap["username"],
				"name":     userMap["name"],
			},
			"full_user_info": userMap,
		})
	} else {
		Success(c, gin.H{
			"message": "pong",
		})
	}
}
