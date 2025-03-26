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

### gorm连接mysql
1. 环境准备
```shell
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```
2. 连接到数据库
```go
db, err := gorm.Open(mysql.New(mysql.Config{
  DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
  DefaultStringSize: 256, // string 类型字段的默认长度
  DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
  DontSupportRenameIndex: true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
  DontSupportRenameColumn: true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
  SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
}), &gorm.Config{})
```
> 然后可以利用链式来写sql语句完成数据库的操作
```go
db.Find(&users, []int{1,2,3})  // SELECT * FROM users WHERE id IN (1,2,3);
db.Where("name = ?", "jinzhu").Delete(&email) // DELETE from emails where id = 10 AND name = "jinzhu";
```

### zap写日志
1. 环境准备
```shell
go get -u go.uber.org/zap
```
2. 使用zap
> zap提供两种log引擎吧，一个是Logger另一个是SugaredLogger，
> 在每一次内存分配都很重要的上下文中是哟ingLogger，只支持强类型结构化日志记录，
> 在性能表现比较好上下文使用SugaredLogger。
> 使用`zap.NewProduction()/zap.NewDevelopment()`或者`zap.Example()`创建一个Logger,
> 通过调用主logger的`.Sugar()`方法来获取一个SugaredLogger
>> 如果想要把日志写入文件等需要定制logger
> `func New(core zapcore.Core, options ...Option) *Logger`
> 
>> zapcore.Core需要三个配置: Encoder如何写入日志，WriterSyncer指定日志写到哪里，Log Level哪种日志被写入
```go
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
    return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./test.log")
	return zapcore.AddSync(file)
}
```
3. 使用Lumberjack进行日志切割归档
+ 环境准备
```shell
go get gopkg.in/natefinch/lumberjack.v2 
```
+ 使用
```go
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log", // 日志文件位置
		MaxSize:    10,  // 日志文件最大大小(MB)
		MaxBackups: 5, // 保留旧文件的最大个数
		MaxAge:     30, // 保留旧文件的最大天数
		Compress:   false, // 是否压缩/归档旧文件
	}
	return zapcore.AddSync(lumberJackLogger)
}
```

### jwt鉴权
> 说来实在惭愧，之前从来没有认真的写过鉴权相关的内容，都是用现有的平台框架直接生成，
> 相当于是在把一些基础的功能都封装了，我只管实现核心的业务功能，所以...
> 鉴权也好，日志也好我其实都没有自己写过...
1. 环境准备
```shell
go get -u github.com/golang-jwt/jwt/v5
```
2. 使用
   
(1) 创建Token（JWT）对象
> 利用`NewWithClaims(method SigningMethod, claims Claims) *Token`来生成token对象
> + method：这是一个 `SigningMethod` 接口参数，用于指定 JWT 的签名算法。常用的签名算法有 `SigningMethodHS256`、`SigningMethodRS256等`
> + claims：这是一个 Claims 接口参数，它表示 JWT 的声明。在 jwt 库中，预定义了一些结构体来实现这个接口，例如 `RegisteredClaims` 和 `MapClaims` 等

(2) 生成JWT字符串
> 利用 `func (t *Token) SignedString(key interface{}) (string, error)`生成token字符串
> + key：该参数是用于签名 token 的密钥。密钥的类型取决于使用的签名算法
> 如果使用 HMAC 算法（如 HS256、HS384 等），key 应该是一个对称密钥（通常是 []byte 类型的密钥）。如果使用 RSA 或 ECDSA 签名算法（如 RS256、ES256），key 应该是一个私钥 *rsa.PrivateKey 或 *ecdsa.PrivateKey
> + 方法返回两个值：一个是成功签名后的 JWT 字符串，另一个是在签名过程中遇到的任何错误
```go
func GenerateToken(id uint, name string) (string, error) {
	iJwtCustClaims := JwtCustClaims{
		ID:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.tokenExpire") * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "Token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, iJwtCustClaims)
	return token.SignedString(stSigningKey)
}
```

(3) 解析JWT字符串
> 利用`func Parse(tokenString string, keyFunc Keyfunc, options ...ParserOption) (*Token, error)`解析token
> + tokenString：要解析的 JWT 字符串
> + keyFunc：这是一个回调函数，返回用于验证 JWT 签名的密钥。
> 
> 还可以利用`func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc, options ...ParserOption) (*Token, error)`来解析
> + tokenString：要解析的 JWT 字符串。 
> + claims：这是一个 Claims 接口参数，用于接收解析 JWT 后的 claims 数据。 
> + keyFunc：与 Parse 函数中的相同，用于提供验证签名所需的密钥。
```go
var stSigningKey = []byte(viper.GetString("jwt.signingKey"))

func ParseToken(tokenStr string) (JwtCustClaims, error) {
    iJwtCustClaims := JwtCustClaims{}
    token, err := jwt.ParseWithClaims(tokenStr, &iJwtCustClaims, func(token *jwt.Token) (interface{}, error) {
        return stSigningKey, nil
    })

    if err == nil && !token.Valid {
        err = errors.New("Invalid Token")
    }

    return iJwtCustClaims, err
}
```

(4) 简单完整示例
```go
func main() {
    jwtKey := make([]byte, 32) // 生成32字节（256位）的密钥
    if _, err := rand.Read(jwtKey); err != nil {
       panic(err)
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
       "iss": "zs", // 发行者
       "sub": "zs.cn", // 主题
       "aud": "Programmer", // 观众
    })
    jwtStr, err := token.SignedString(jwtKey)
    if err != nil {
       panic(err)
    }

    // 解析 jwt
    claims, err := ParseJwtWithClaims(jwtKey, jwtStr)
    if err != nil {
       panic(err)
    }
    fmt.Println(claims) // map[aud:Programmer iss:zs sub:zs.cn]
}
```

### 封装ResponseJson
使用gin框架的最简单的方式是：
```go
import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/", func(ctx *gin.Context) {
        ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
        "message": "hello word",
        "index":   "this is the page",
    })
})
```
但是如果每一个接口都用这种方式来返回的话似乎不够优雅，如果每个接口只需要关注自己的业务逻辑，
然后把要返回的数据直接交给统一的函数来干好像更优雅一点，然后现在来统一一下返回数据的形式:
```go
type ResponseJson struct {
	Status int    `json:"-"`  // 状态码，如200成功，400 bad request等
	Code   int    `json:"code,omitempty"` // 个人理解是用来区分每个状态码下的不同状态,可能更方便接口问题定位
	Msg    string `json:"msg,omitempty"` // 接口返回提示信息
	Data   any    `json:"data,omitempty"` // 接口的数据
}
```
将常见的几种状态返回进行封装
```go
func HttpResponse(ctx *gin.Context, status int, resp ResponseJson) {
    if resp.IsEmpty() {
        ctx.AbortWithStatus(status)
        return
    }
    ctx.AbortWithStatusJSON(status, resp)
}

func buildStatus(resp ResponseJson, nDefaultStatus int) int {
    if resp.Status == 0 {
        return nDefaultStatus
    }
    return resp.Status
}

func Ok(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusOK), resp)
}

func Fail(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusBadRequest), resp)
}

func ServerFail(ctx *gin.Context, resp ResponseJson) {
	HttpResponse(ctx, buildStatus(resp, http.StatusInternalServerError), resp)
}
```
使用的时候就可以直接使用封装好的接口来返回数据就好
```go
func main() {
    r := gin.Default()

    r.GET("/", func(ctx *gin.Context) {
        //ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
        //"message": "hello word",
        // 直接使用下面的方式
        Ok(ctx, ResponseJson{
            Msg: "hello word",
        })
    })
})
```

### Gin参数绑定
Gin框架提供多种绑定方式，包括：Query参数绑定、Form参数帮i的那个、JSON参数绑定
提供了两种方式：Must Bind和Should Bind

| 功能     | Must Bind | Should Bind     |
|--------|-----------|-----------------|
|        | Bind      | ShouldBind      |
| 绑定JSON | BindJSON  | ShouldBindJSON  |
| 绑定XML  | BindXML   | ShouldBindXML   |
| 绑定GET  | BindQuery | ShouldBindQuery |
| 绑定YAML | BindYAML  | ShouldBindYAML  |

1. 使用示例
```go
// 定义待绑定的JSON结构体
type Param struct {
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Likes []string `json:"likes"`
}
// 绑定提交的Json数据
func TestBindJson(engine *gin.Engine) {
	engine.POST("/bindJson", func(context *gin.Context) {
		var jsonParam Param
		var err error
		bindType := context.Query("type")
		fmt.Println(bindType)
		if bindType == "1" {
			err = context.BindJSON(&jsonParam)
		} else {
			err = context.ShouldBindJSON(&jsonParam)
		}
		if err != nil {
			context.JSON(500, gin.H{"error": err})
			return
		}
		context.JSON(200, gin.H{"result": jsonParam})
	})
}
```
2. 自定义参数验证器

(1)创建自定义验证器函数
```go
var testValidatorFunc Validator.Func = func(fl validator.FieldLevel) bool {
    if value, ok := fl.Field().Interface().(string); ok {
        if value != "" && 0 == strings.Index(value, "a") {
            return true
        }
    }
	
    return false
}
```
+ `validator.Func` 是一个自定义验证器函数类型，它需要包含自定义验证逻辑的函数体。
+ `validator`.FieldLevel 是一个接口类型，提供了用于获取要验证字段的信息的方法，例如 `FieldName()`、`Field()` 和 `Param()`。
+ `fl.Field()` 返回一个反射值，允许我们获取要验证的字段的值。
+ 使用 `Interface().(string)` 将反射值储存在一个空的 interface{} 变量中，并将其转换为 `string` 类型的变量。

(2) 注册验证器
```go
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
   v.RegisterValidation("first_is_a", testValidatorFunc)
}
```
(3) 将标签添加到对应的字段上，然后就可以用了
```go
type UserLoginDTO struct {
	Name     string `json:"name" binding:"required,first_is_a"`
	Password string `json:"password" binding:"required"`
}
```
(4) 自定义错误提示

`validator` 校验返回的结果只有 3 种情况：
+ `nil`：没有错误；
+ `InvalidValidationError`：输入参数错误；
+ `ValidationErrors`：字段违反约束。
```go
func formatValidationErrors(err error, target interface{}) error {
    // 检查是否是验证错误
    validationErrs, ok := err.(validator.ValidationErrors)
    if !ok { // 如果不是就返回原始错误信息
        return err
    }

    // 获取结构体类型信息
    structType := reflect.TypeOf(target).Elem()

    var errMsgs []string
    for _, fieldErr := range validationErrs {
        // 获取字段定义
        field, _ := structType.FieldByName(fieldErr.Field())

        // 尝试获取特定错误消息（如first_is_a_err）
        customErrTag := fieldErr.Tag() + "_err"
        errMsg := field.Tag.Get(customErrTag)

        // 如果没有特定消息，使用默认格式
        if errMsg == "" {
            errMsg = fmt.Sprintf("%s验证失败(%s)", fieldErr.Field(), fieldErr.Tag())
        }

        errMsgs = append(errMsgs, errMsg)
    }

    return fmt.Errorf(strings.Join(errMsgs, "; "))
}
```
需要注意的是设置了两个验证规则：`required` 和 `first_is_a`，Go的validator会按顺序执行验证规则，先检查`required`规则，再检查`first_is_a`规则

### Cors跨域
1. 环境准备
```shell
go get github.com/gin-contrib/cors
```
2. 使用
```go
func main() {
    // 创建一个默认的 Gin 实例
    server := gin.Default()

    // 使用 CORS 中间件处理跨域问题，配置 CORS 参数
    server.Use(cors.New(cors.Config{
        // 允许的源地址（CORS中的Access-Control-Allow-Origin）
        // AllowOrigins: []string{"https://foo.com"},
        // 允许的 HTTP 方法（CORS中的Access-Control-Allow-Methods）
        AllowMethods: []string{"PUT", "PATCH"},
        // 允许的 HTTP 头部（CORS中的Access-Control-Allow-Headers）
        AllowHeaders: []string{"Origin"},
        // 暴露的 HTTP 头部（CORS中的Access-Control-Expose-Headers）
        ExposeHeaders: []string{"Content-Length"},
        // 是否允许携带身份凭证（CORS中的Access-Control-Allow-Credentials）
        AllowCredentials: true,
        // 允许源的自定义判断函数，返回true表示允许，false表示不允许
        AllowOriginFunc: func(origin string) bool {
            if strings.HasPrefix(origin, "http://localhost") {
                // 允许你的开发环境
                return true
            }
            // 允许包含 "yourcompany.com" 的源
            return strings.Contains(origin, "yourcompany.com")
        },
        // 用于缓存预检请求结果的最大时间（CORS中的Access-Control-Max-Age）
        MaxAge: 12 * time.Hour,
    }))

    // 启动 Gin 服务器，监听在 0.0.0.0:8080 上
    server.Run(":8080")
}
```
