package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hock1024always/GoEdu/dao"
	"net/http"
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserApi struct {
	Username string `json:"username"`
	Userid   int    `json:"userid"`
	subject  string `json:"subject"`
}

// JWT Claims 结构体
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// JWT 密钥
var jwtKey = []byte("910")

func (User) TableName() string {
	return "user"
}

// 判断用户名是否已经存在
func CheckUserExist(username string) (User, error) {
	var user User
	err := dao.Db.Where("username =?", username).First(&user).Error
	return user, err
}

// 保存用户
func AddUser(username, password string, email string) (UserApi, error) {
	user := User{
		Username: username,
		Password: password,
		Email:    email,
	}
	err := dao.Db.Create(&user).Error
	userapi := UserApi{Username: username, Userid: user.Id}
	return userapi, err
}

// 通过Id来查找用户
func CheckUserById(id int) (User, error) {
	var user User
	err := dao.Db.Where("id =?", id).First(&user).Error
	return user, err
}

// 删除通过用户名用户
func DeleteUserByUsername(username string) error {
	var user User
	err := dao.Db.Where("username =?", username).Delete(&user).Error
	return err
}

// 更新密码
func UpdateUserPassword(username string, password string) (string, error) {
	var user User
	err := dao.Db.Model(&user).Where("username =?", username).UpdateColumn("password", password).Error
	return password, err
}

type JsonStruct struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   interface{} `json:"msg"`
	Count int64       `json:"count"`
}

type JsonErrStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

func ReturnSuccess(c *gin.Context, code int, msg interface{}, data interface{}, count int64) {
	//code：响应码，msg：泛型信息，data：泛型数据，count：信息条数
	json := &JsonStruct{Code: code, Msg: msg, Data: data, Count: count}
	c.JSON(http.StatusOK, json)
}

func ReturnError(c *gin.Context, code int, msg string) {
	//code：响应码，msg：错误信息
	json := &JsonErrStruct{Code: code, Msg: msg}
	c.JSON(http.StatusOK, json)
}

// 定义密钥（用于签名）
var secretKey = []byte("my_secret_key")

// 生成 JWT Token
func GenerateToken(username string) (string, error) {
	// 创建 JWT 负载
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // 过期时间 2 小时
		"iat":      time.Now().Unix(),                    // 签发时间
	}

	// 创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 进行签名
	return token.SignedString(secretKey)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保使用的是正确的签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})
}

// 验证 Token 是否有效
func ValidateToken(tokenString string) string {
	// 解析 Token
	token, _ := ParseToken(tokenString)
	//if err != nil {
	//	// 如果解析失败，打印错误信息并返回空字符串
	//	fmt.Println("Token解析失败")
	//	return ""
	//}
	claims, _ := token.Claims.(jwt.MapClaims)
	return claims["username"].(string)
	// 验证 Token 是否有效
	//if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	//	// 如果 Token 有效，返回用户名
	//	return claims["username"].(string)
	//} else {
	//	// 如果 Token 无效，打印错误信息并返回空字符串
	//	fmt.Println("Token无效")
	//	return ""
	//}
}
