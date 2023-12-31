package service

import (
	"test/conf"
	"test/logger"
	"test/repository"
	"test/repository/oracle"
)

type Container struct {
	Store       repository.Store
	Logger      logger.Logger
	JWTService  JWTService
	UserService UserService
}

// create container by intiating all required services
func NewServiceContainer() *Container {
	conf.InitConfig()
	logger := logger.NewLogger()
	store := oracle.SharedStore(logger)
	cacheService := NewCacheService()
	jwtService := NewJWTService(cacheService)
	userService := NewUserService(store, logger, jwtService)

	return &Container{
		JWTService:  jwtService,
		UserService: userService,
		Store:       store,
		Logger:      logger,
	}
}
