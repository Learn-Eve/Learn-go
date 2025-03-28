package service

import (
	"errors"
	"fast-learn/dao"
	"fast-learn/global"
	"fast-learn/global/constants"
	"fast-learn/model"
	"fast-learn/service/dto"
	"fast-learn/utils"
	"fmt"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

var userService *UserService

type UserService struct {
	BaseService
	Dao *dao.UserDao
}

func NewUserService() *UserService {
	if userService == nil {
		userService = &UserService{
			Dao: dao.NewUserDao(),
		}
	}

	return userService
}

func GenerateAndCacheLoginUserToken(nUserId uint, stUserName string) (string, error) {
	token, err := utils.GenerateToken(nUserId, stUserName)
	if err == nil {
		err = global.RedisClient.Set(strings.Replace(constants.LOGIN_USER_TOKEN_REDIS_KEY, "{ID}", strconv.Itoa(int(nUserId)), -1), token, viper.GetDuration("jwt.tokenExpire")*time.Minute)
	}
	return token, err
}

func (m *UserService) Login(iUserDTO dto.UserLoginDTO) (model.User, string, error) {
	var errResult error
	var token string

	iUser, err := m.Dao.GetUserByName(iUserDTO.Name)
	// 用户名或密码不正确
	if err != nil || !utils.CompareHashAndPassword(iUser.Password, iUserDTO.Password) {
		errResult = errors.New("Invalid UserName or Password")
	} else { // 登录成功，生成token
		token, err = GenerateAndCacheLoginUserToken(iUser.ID, iUser.Name)
		if err != nil {
			errResult = errors.New(fmt.Sprintf("Token Generate Error:%s", err.Error()))
		}
	}

	return iUser, token, errResult
}

func (m *UserService) AddUser(iUserAddDTO *dto.UserAddDTO) error {
	// 判断用户是否存在
	if m.Dao.CheckUserNameExist(iUserAddDTO.Name) {
		return errors.New("User Already Exists")
	}
	return m.Dao.AddUser(iUserAddDTO)
}

func (m *UserService) GetUserById(iCommonIDDTO *dto.CommonIDDTO) (model.User, error) {
	return m.Dao.GetUserByID(iCommonIDDTO.ID)
}

func (m *UserService) GetUserList(iUserListDTO *dto.UserListDTO) ([]model.User, int64, error) {
	return m.Dao.GetUserList(iUserListDTO)
}

func (m *UserService) UpdateUser(iUserUpdateDTO *dto.UserUpdateDTO) error {
	if iUserUpdateDTO.ID <= 0 {
		return errors.New("Invalid ID")
	}

	return m.Dao.UpdateUser(iUserUpdateDTO)
}

func (m *UserService) DeleteUserById(iCommonIDDTO *dto.CommonIDDTO) error {
	return m.Dao.DeleteUserById(iCommonIDDTO.ID)
}
