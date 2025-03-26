package router

import (
	"fast-learn/api"
	"github.com/gin-gonic/gin"
)

func InitHostRoutes() {
	RegistRoute(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		hostApi := api.NewHostApi()
		rgAuthHost := rgAuth.Group("host")
		{
			rgAuthHost.POST("/shutdown", hostApi.Shutdown)
		}
	})
}
