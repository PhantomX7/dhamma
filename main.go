package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PhantomX7/dhamma/command"
	"github.com/PhantomX7/dhamma/config"
	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/controller"
	"github.com/PhantomX7/dhamma/middleware"
	"github.com/PhantomX7/dhamma/repository"
	"github.com/PhantomX7/dhamma/routes"
	"github.com/PhantomX7/dhamma/service"
	"github.com/PhantomX7/dhamma/utils/validators"
	"github.com/go-playground/validator/v10"

	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func args(db *gorm.DB) bool {
	if len(os.Args) > 1 {
		flag := command.Commands(db)
		if !flag {
			return false
		}
	}

	return true
}

func main() {
	// load environment
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db := config.SetUpDatabaseConnection()
	// defer config.CloseDatabaseConnection(db)

	if !args(db) {
		return
	}

	// var (
	// 	jwtService service.JWTService = service.NewJWTService()

	// 	// Implementation Dependency Injection
	// 	// Repository
	// 	userRepository repository.UserRepository = repository.NewUserRepository(db)

	// 	// Service
	// 	userService service.UserService = service.NewUserService(userRepository, jwtService)

	// 	// Controller
	// 	userController controller.UserController = controller.NewUserController(userService)
	// )

	app := fx.New(
		fx.Provide(
			config.SetUpDatabaseConnection,
			setUpServer,
			// initLibs,
		),
		repository.Module,
		service.Module,
		controller.Module,
		routes.Module,
		fx.Invoke(
			validators.NewValidator,
			//startCron,
			//startQueue,
			startServer,
		),
	)
	app.Run()

	// server := gin.Default()
	// server.Use(middleware.CORSMiddleware())

	// routes
	// routes.User(server, userController, jwtService)

}

func startServer(lc fx.Lifecycle, server *gin.Engine, db *gorm.DB) {
	myFigure := figure.NewColorFigure("Phantom", "", "green", true)
	myFigure.Print()

	// register all custom validator here
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("unique",
			validators.CustomValidator.Unique())
		if err != nil {
			log.Println("error when applying unique validator")
		}
		err = v.RegisterValidation("exist", validators.CustomValidator.Exist())
		if err != nil {
			log.Println("error when applying exist validator")
		}
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// start server

			return nil
		},
		OnStop: func(context.Context) error {
			config.CloseDatabaseConnection(db)

			return nil
		},
	})

	if err := server.Run(fmt.Sprintf(":%s", os.Getenv("GOLANG_PORT"))); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

func setUpServer() *gin.Engine {
	if os.Getenv("APP_ENV") == constants.ENUM_RUN_PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	server.Static("/assets", "./assets")

	return server
}
