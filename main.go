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
	"time"

	"github.com/gin-contrib/cors"
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

		userRepository          repository.UserRepository          = repository.NewUserRepository(db)
		fileRepository          repository.FileRepository          = repository.NewFileRepository(db)
		accessRequestRepository repository.AccessRequestRepository = repository.NewAccessRequestRepository(db)

		userService service.UserService = service.NewUserService(userRepository, accessRequestRepository)
		fileService service.FileService = service.NewFileService(fileRepository, userRepository)

		userController controller.UserController = controller.NewUserController(userService, jwtService)
		fileController controller.FileController = controller.NewFileController(fileService, userService, jwtService)
	)

	server := gin.Default()

	// CORS middleware
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Adjust this to your frontend's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.UserRoutes(server, userController, jwtService)
	routes.FileRoutes(server, fileController, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	server.Run("127.0.0.1:" + port)
}
