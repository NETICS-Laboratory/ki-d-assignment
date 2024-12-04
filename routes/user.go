package routes

import (
	"ki-d-assignment/controller"
	"ki-d-assignment/middleware"
	"ki-d-assignment/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, UserController controller.UserController, jwtService service.JWTService) {
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/register", UserController.RegisterUser)
		userRoutes.POST("/login", UserController.LoginUser)
		userRoutes.GET("/me", middleware.Authenticate(jwtService), UserController.MeUser)
		userRoutes.GET("/me-decrypted", middleware.Authenticate(jwtService), UserController.MeUserDecrypted)
		userRoutes.GET("/idcard-decrypted", middleware.Authenticate(jwtService), UserController.DecryptUserIDCard)

		userRoutes.POST("/request-access", middleware.Authenticate(jwtService), UserController.RequestAccess)
		userRoutes.GET("/request-access", middleware.Authenticate(jwtService), UserController.GetAccessRequests)
		userRoutes.PUT("/request-access/:request_id", middleware.Authenticate(jwtService), UserController.UpdateAccessRequestStatus)

		userRoutes.POST("/decrypt-key", middleware.Authenticate(jwtService), UserController.DecryptKeys)
		userRoutes.POST("/get-requested-user-data", middleware.Authenticate(jwtService), UserController.AccessPrivateData)
	}
}
