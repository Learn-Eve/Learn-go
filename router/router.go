package router

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

	InitBasePlatformRoutes()

	for _, fnRegistRoute := range gfnRoutes {
		fnRegistRoute(rgPublic, rgAuth)
	}

	stPort := viper.GetString("server.port")
	if stPort == "" {
		stPort = "8999"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", stPort),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// TODO: 记录日志
			fmt.Println(fmt.Sprintf("Start Server Error: %s", err.Error()))
			return
		}

		//fmt.Println(fmt.Sprintf("Start Server Listen: %s", server.Addr))
	}()

	<-ctx.Done()

	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctx); err != nil {
		// TODO：记录日志
		fmt.Println(fmt.Sprintf("Shutdown Server Error: %s", err.Error()))
		return
	}

	fmt.Println("Stop Server Success")
}

// 初始化基础模块路由
func InitBasePlatformRoutes() {
	InitUserRoutes()
}
