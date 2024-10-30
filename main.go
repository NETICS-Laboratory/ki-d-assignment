package main

import (
	"ki-d-assignment/cmd"
	"ki-d-assignment/common"
	"ki-d-assignment/config"
	"ki-d-assignment/controller"
	"ki-d-assignment/database"
	"ki-d-assignment/repository"
	"ki-d-assignment/routes"
	"ki-d-assignment/service"
	"ki-d-assignment/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		res := common.BuildErrorResponse("Gagal Terhubung ke Server", err.Error(), common.EmptyObj{})
		(*gin.Context).JSON((&gin.Context{}), http.StatusBadGateway, res)
		return
	}

	db := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	if len(os.Args) > 1 {
		cmd.Commands(db)
		return
	}

	var (
		jwtService service.JWTService = service.NewJWTService()

		userRepository          repository.UserRepository          = repository.NewUserRepository(db)
		fileRepository          repository.FileRepository          = repository.NewFileRepository(db)
		accessRequestRepository repository.AccessRequestRepository = repository.NewAccessRequestRepository(db)

		userService service.UserService = service.NewUserService(userRepository, accessRequestRepository)
		fileService service.FileService = service.NewFileService(fileRepository, userRepository)

		userController controller.UserController = controller.NewUserController(userService, jwtService)
		fileController controller.FileController = controller.NewFileController(fileService, userService, jwtService)
	)

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

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

	port := utils.MustGetenv("PORT")
	if port == "" {
		port = "8090"
	}
	server.Run("127.0.0.1:" + port)
}
