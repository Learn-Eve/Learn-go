package middleware

import (
	"fast-learn/api"
	"fast-learn/global"
	"fast-learn/global/constants"
	"fast-learn/model"
	"fast-learn/service"
	"fast-learn/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ERR_CODE_INVALID_TOKEN     = 10401 // token无效
	ERR_CODE_TOKEN_PARSE       = 10402 // 解析Token失败
	ERR_CODE_TOKEN_NOT_MATCHED = 10403
	ERR_CODE_TOKEN_EXPIRED     = 10404
	ERR_CODE_TOKEN_RENEW       = 10405
	TOKEN_NAME                 = "Authorization"
	TOKEN_PREFIX               = "Bearer: "
	RENEW_TOKEN_DURATION       = 10 * 60 * time.Second
)

func tokenErr(c *gin.Context, code int) {
	api.Fail(c, api.ResponseJson{
		Status: http.StatusUnauthorized,
		Code:   code,
		Msg:    "Invalid token",
	})
}

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader(TOKEN_NAME)

		// token不存在，直接返回
		if token == "" || !strings.HasPrefix(token, TOKEN_PREFIX) {
			tokenErr(c, ERR_CODE_INVALID_TOKEN)
			return
		}

		// token无法解析，直接返回
		token = token[len(TOKEN_PREFIX):]
		iJwtCustClaims, err := utils.ParseToken(token)
		nUserId := iJwtCustClaims.ID
		if err != nil || nUserId == 0 {
			tokenErr(c, ERR_CODE_TOKEN_PARSE)
			return
		}

		stUserId := strconv.Itoa(int(nUserId))
		stRedisUserIdKey := strings.Replace(constants.LOGIN_USER_TOKEN_REDIS_KEY, "{ID}", stUserId, -1)

		// token与访问者登录对应的token不一致，直接返回
		stRedisToken, err := global.RedisClient.Get(stRedisUserIdKey)
		if err != nil || token != stRedisToken {
			tokenErr(c, ERR_CODE_TOKEN_NOT_MATCHED)
			return
		}

		// token过期，
		nTokenExpireDuration, err := global.RedisClient.GetExpireDuration(stRedisUserIdKey)
		if err != nil || nTokenExpireDuration <= 0 {
			tokenErr(c, ERR_CODE_TOKEN_EXPIRED)
			return
		}

		// token的续期
		if nTokenExpireDuration.Seconds() < RENEW_TOKEN_DURATION.Seconds() {
			stNewToken, err := service.GenerateAndCacheLoginUserToken(nUserId, iJwtCustClaims.Name)
			if err != nil {
				tokenErr(c, ERR_CODE_TOKEN_RENEW)
				return
			}

			c.Header("token", stNewToken)
		}

		//iUser, err := dao.NewUserDao().GetUserByID(nUserId)
		//if err != nil {
		//	tokenErr(c)
		//	return
		//}
		//c.Set(constants.LOGIN_USER, iUser)

		c.Set(constants.LOGIN_USER, model.LoginUser{
			ID:   nUserId,
			Name: iJwtCustClaims.Name,
		})

		c.Next()
	}
}
