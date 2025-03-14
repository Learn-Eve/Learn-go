package cmd

import (
	"fast-learn/conf"
	"fast-learn/router"
	"fmt"
)

// 负责系统相关的启动初始化

func Start() {
	// 读取配置文件内容
	conf.InitConfig()
	router.InitRouter()
}

func Clean() {
	fmt.Println("===========Clean===============")
}
