# 小白学习Go Web实战
## 技术栈
1. gin
2. viper
3. swagger
4. zap
5. gorm


## 项目开始准备
1. 准备swagger环境
   ```shell
        go get -u github.com/swaggo/swag/cmd/swag
        go install github.com/swaggo/swag/cmd/swag@latest
    ```
   

## 过程防遗忘
> 脑子真是不好使了，断断续续的学，对go也没什么基础，之前仅有的一丁点开发经验也丢了，学了一半发现还是得捋一下
> 对于swagger部分不捋了，感觉用到的概率不大，主要针对学习过程中请求的封装来捋一下

### viper读取配置文件
1. 环境准备
```shell
go get github.com/spf13/viper
```
2. 配置文件读取
```go
viper.SetConfigName("config") // name of config file (without extension)
viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
viper.AddConfigPath(".")               // optionally look for config in the working directory
err :=  viper.ReadInConfig() //Find and read the config file
if err != nil { // Handle errors reading the config file
	panic(fmt.Errorf("fatal error config file: %w", err))
}
```