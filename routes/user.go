package routes

import (
	"ki-d-assignment/controller"
	// "ki-d-assignment/middleware"
	"ki-d-assignment/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, UserController controller.UserController, jwtService service.JWTService) {
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/register", UserController.RegisterUser)
		// userRoutes.GET("getalluser", middleware.Authenticate(jwtService), UserController.GetAllUser)
		// userRoutes.POST("/login", UserController.LoginUser)
		// userRoutes.DELETE("/delete", middleware.Authenticate(jwtService), UserController.DeleteUser)
		// userRoutes.PUT("/update", middleware.Authenticate(jwtService), UserController.UpdateUser)
		// userRoutes.GET("/me", middleware.Authenticate(jwtService), UserController.MeUser)
	}
}
