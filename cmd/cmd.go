package cmd

import (
	"fast-learn/conf"
	"fast-learn/global"
	"fast-learn/router"
	"fmt"
)

// 负责系统相关的启动初始化

func Start() {
	// 读取配置文件内容
	conf.InitConfig()
	
	// 初始化日志组件
	global.Logger = conf.InitLogger()

	// 初始化系统路由
	router.InitRouter()
}

func Clean() {
	fmt.Println("===========Clean===============")
}
