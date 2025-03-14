package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 用户管理

func InitUserRoutes() {
	RegistRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		// 登录就不需要鉴权
		rgPublic.POST("/login", func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"msg": "Login Success",
			})
		})

		rgAuthUser := rgAuth.Group("user")
		rgAuthUser.GET("", func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"data": []map[string]any{
					{"id": 1, "name": "zs"},
					{"id": 2, "name": "ls"},
				},
			})
		})

		rgAuthUser.GET("/:id", func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"id":   1,
				"name": "zs",
			})
		})
	})
}
