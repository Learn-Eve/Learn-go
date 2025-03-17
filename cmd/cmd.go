package cmd

import (
	"fast-learn/conf"
	"fast-learn/global"
	"fast-learn/router"
	"fast-learn/utils"
	"fmt"
)

// 负责系统相关的启动初始化

func Start() {
	var initErr error
	// 读取配置文件内容
	conf.InitConfig()

	// 初始化日志组件
	global.Logger = conf.InitLogger()

	// 初始化数据库连接
	db, err := conf.InitDB()
	global.DB = db
	if err != nil {
		initErr = utils.AppendError(initErr, err)
	}

	if initErr != nil {
		if global.Logger != nil {
			global.Logger.Error(initErr.Error())
		}

		panic(initErr.Error())
	}

	// 初始化系统路由
	router.InitRouter()
}

func Clean() {
	fmt.Println("===========Clean===============")
}
