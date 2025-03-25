package router

import (
	"github.com/gin-contrib/sessions"
	sessions_redis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/hock1024always/GoEdu/config"
	"github.com/hock1024always/GoEdu/controllers"
	"github.com/hock1024always/GoEdu/pkg/logger"
)

// 路由 函数的名字要大写，这样才可以被其他包访问！
func Router() *gin.Engine {
	//创建一个路由的实例
	r := gin.Default()

	//日志中间件
	r.Use(gin.LoggerWithConfig(logger.LoggerToFile()))
	r.Use(logger.Recover)
	//sessions中间件
	store, _ := sessions_redis.NewStore(10, "tcp", config.RedisAddress, "", []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	user := r.Group("/user")
	{
		// 向注册用户发送验证码 username password confirm_password email
		user.POST("/register", controllers.UserController{}.Register)
		// 实现注册用户验证码
		user.POST("/register/verify", controllers.UserController{}.Verify)
		// 登录用户相关的路由 username password
		user.POST("/login", controllers.UserController{}.Login)
		//实现删除用户的路由 username password confirm_sentence
		user.POST("/delete", controllers.UserController{}.UserDelete)
		//实现websocket的路由
		user.GET("/ws", controllers.UserController{}.WsHandler)
		//实现用户视频功能
		user.GET("/video", controllers.UserController{}.LiveHandler)
		//// 实现用户获取自己的聊天记录 token
		//user.POST("/get_chat_records", controllers.UserController{}.GetChatRecords)
		////实现用户修改密码 username password new_password confirm_new_password
		//user.POST("/modify_password", controllers.UserController{}.ModifyPassword)
		//实现用户于机器学习交互 ws
		user.GET("/ws/ai", controllers.UserController{}.AIHandler)
	}

	return r
}
