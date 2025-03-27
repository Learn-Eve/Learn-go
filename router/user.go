package router

import (
	"fast-learn/api"
	"github.com/gin-gonic/gin"
)

// 用户管理

func InitUserRoutes() {
	RegistRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		userApi := api.NewUserApi()
		rgPublicUser := rgPublic.Group("user").Use(func() gin.HandlerFunc {
			return func(ctx *gin.Context) {
				//ctx.JSON(299, gin.H{
				//	"msg": "Login MiddleWare",
				//})
			}
		}())
		{
			rgPublicUser.POST("/login", userApi.Login)
		}

		rgAuthUser := rgAuth.Group("user")
		{
			rgAuthUser.POST("add", userApi.AddUser)
			rgAuthUser.GET("/:id", userApi.GetUserById)
			rgAuthUser.POST("/list", userApi.GetUserList)
			rgAuthUser.POST("/update", userApi.UpdateUser)
			rgAuthUser.DELETE("/:id", userApi.DeleteUserById)
		}
	})
}
