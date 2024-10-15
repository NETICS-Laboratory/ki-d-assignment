package main

import (
	"ki-d-assignment/common"
	"ki-d-assignment/config"
	"ki-d-assignment/controller"
	"ki-d-assignment/repository"
	"ki-d-assignment/routes"
	"ki-d-assignment/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Terhubung ke Server", err.Error(), common.EmptyObj{})
		(*gin.Context).JSON((&gin.Context{}), http.StatusBadGateway, res)
		return
	}

	var (
		db *gorm.DB = config.SetupDatabaseConnection()

		jwtService service.JWTService = service.NewJWTService()

		userRepository repository.UserRepository = repository.NewUserRepository(db)
		fileRepository repository.FileRepository = repository.NewFileRepository(db)

		userService service.UserService = service.NewUserService(userRepository)
		fileService service.FileService = service.NewFileService(fileRepository, userRepository)

		userController controller.UserController = controller.NewUserController(userService, jwtService)
		fileController controller.FileController = controller.NewFileController(fileService, userService, jwtService)
	)

	server := gin.Default()
	routes.UserRoutes(server, userController, jwtService)
	routes.FileRoutes(server, fileController, jwtService)

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "8090"
	}
	server.Run("127.0.0.1:" + port)
}
