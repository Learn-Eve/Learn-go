package router

import (
	"context"
	_ "fast-learn/docs"
	"fast-learn/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// 负责系统中路由初始化

// rgPublic不需要鉴权 rgAuth需要鉴权
type IFnRegistRoute = func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup)

var (
	gfnRoutes []IFnRegistRoute
)

// 注册路由
func RegistRoute(fn IFnRegistRoute) {
	if fn == nil {
		return
	}

	gfnRoutes = append(gfnRoutes, fn)
}

func InitRouter() {
	// 创建监听ctrl+c，应用退出信号的上下文
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelCtx()

	// 初始化gin框架，并注册相关路由
	r := gin.Default()
	rgPublic := r.Group("/api/v1/public")
	rgAuth := r.Group("/api/v1")

	// 初始基础平台的路由
	InitBasePlatformRoutes()

	// 开始注册系统各模块对应的路由信息
	for _, fnRegistRoute := range gfnRoutes {
		fnRegistRoute(rgPublic, rgAuth)
	}

	// 集成swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 从配置文件中读取并配置web服务配置
	stPort := viper.GetString("server.port")
	if stPort == "" {
		stPort = "8999"
	}

	// 创建web server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", stPort),
		Handler: r,
	}

	// 启动一个goroutine来开启web服务，避免主线程的信号监听被阻塞
	go func() {
		global.Logger.Info(fmt.Sprintf("Start Listen: %s", stPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			//fmt.Println(fmt.Sprintf("Start Server Error: %s", err.Error()))
			global.Logger.Error(fmt.Sprintf("Start Server Error: %s", err.Error()))
			return
		}

		//fmt.Println(fmt.Sprintf("Start Server Listen: %s", server.Addr))
	}()

	<-ctx.Done()

	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctx); err != nil {
		global.Logger.Error(fmt.Sprintf("Shutdown Server Error: %s", err.Error()))
		//fmt.Println(fmt.Sprintf("Shutdown Server Error: %s", err.Error()))
		return
	}

	//fmt.Println("Stop Server Success")
	global.Logger.Info(fmt.Sprintf("Stop Server Success"))
}

// 初始化基础模块路由
func InitBasePlatformRoutes() {
	InitUserRoutes()
}
