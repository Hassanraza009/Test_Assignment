package rest

import (
	"test/conf"
	"test/logger"
	"test/service"

	"github.com/spf13/viper"
)

// StartServer initate server
func StartServer(container *service.Container) *HttpServer {

	//Inject services instance from ServiceContainer
	userController := NewUserController(container.UserService, container.Logger)
	middleWare := NewMiddleWare(container.JWTService)
	addr := viper.Get(conf.RestServerPort).(string)

	httpServer := NewHttpServer(addr)

	//Inject controller instance to server
	httpServer.middleware = middleWare
	httpServer.userController = userController

	go httpServer.Start()
	logger.Instance().Info("rest server ok")
	return httpServer
}
