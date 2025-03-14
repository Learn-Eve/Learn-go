package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

func InitConfig() {
	// 设置配置文件名称
	viper.SetConfigName("settings")
	// 设置配置文件类型
	viper.SetConfigType("yml")
	// 读取配置文件的路径
	//fmt.Println(os.Getwd())
	viper.AddConfigPath("./conf/")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Sprintf("Load Config Error: %s \n", err.Error()))
	}
	fmt.Println(viper.GetString("server.port"))

}
