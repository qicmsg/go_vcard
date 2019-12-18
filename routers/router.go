package routers

import (
	"github.com/gin-gonic/gin"
	loginController "vcard/app/http/controllers/api/v1/login"
	userController "vcard/app/http/controllers/api/v1/user"
	"vcard/app/http/middleware/logger"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	//engine.Use(gin.Recovery())

	logger.Setup()

	// 日志中间件
	// engine.Use(logger.LoggerToFile())

	apiV1 := engine.Group("/api/v1")

	apiV1.GET("/login", loginController.Login)


	user := apiV1.Group("user")
	user.GET("/info", userController.Info)
	user.GET("/save", userController.Save)
	user.GET("/list", userController.List)
	user.GET("/test", userController.Test)


	return engine
}
