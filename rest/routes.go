package rest

import "github.com/gin-gonic/gin"

func AttachRoutes(r *gin.Engine, server *HttpServer) {
	public := r.Group("/")
	{
		public.POST("/api/user/signup", server.userController.CreateUserAccount)
		public.POST("/api/user/login", server.userController.Login)
	}
	validated := r.Group("/", server.middleware.ValidateToken())
	{
		validated.GET("/api/user/get", server.userController.GetUserByEmail)
	}

}
