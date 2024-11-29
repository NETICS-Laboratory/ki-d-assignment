package routes

import (
	"ki-d-assignment/controller"
	"ki-d-assignment/middleware"
	"ki-d-assignment/service"

	"github.com/gin-gonic/gin"
)

func FileRoutes(router *gin.Engine, fileController controller.FileController, jwtService service.JWTService) {
	fileRoutes := router.Group("/api/files")
	{
		fileRoutes.POST("/upload", middleware.Authenticate(jwtService), fileController.UploadFile)
		fileRoutes.GET("/get-files", middleware.Authenticate(jwtService), fileController.GetUserFiles)
		fileRoutes.POST("/get-file-decrypted", middleware.Authenticate(jwtService), fileController.GetUserFileDecrypted)
		fileRoutes.POST("/verify-digital-signature", middleware.Authenticate(jwtService), fileController.VerifyFileSignature)
	}
}
