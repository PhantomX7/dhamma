package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PhantomX7/dhamma/constants"
	"github.com/PhantomX7/dhamma/modules"
	"github.com/PhantomX7/dhamma/routes"
	"github.com/PhantomX7/dhamma/utility"
	"github.com/PhantomX7/go-core/utility/validators"

	"github.com/common-nighthawk/go-figure"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// func args(db *gorm.DB) bool {
// 	if len(os.Args) > 1 {
// 		flag := command.Commands(db)
// 		if !flag {
// 			return false
// 		}
// 	}

// 	return true
// }

func main() {
	// load environment
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	// db := config.SetUpDatabaseConnection()
	// defer config.CloseDatabaseConnection(db)

	// if !args(db) {
	// 	return
	// }

	app := fx.New(
		fx.Provide(
			setupDatabase,
			setUpServer,
			// initLibs,
		),
		modules.RepositoryModule,
		modules.ServiceModule,
		modules.ControllerModule,
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
			// config.CloseDatabaseConnection(db)

			return nil
		},
	})

	if err := server.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

func setUpServer() *gin.Engine {
	if os.Getenv("APP_ENV") == constants.ENUM_RUN_PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()
	// server.Use(middleware.CORSMiddleware())

	server.Static("/assets", "./assets")

	return server
}

func setupDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}

	if err = utility.RunMigration(db); err != nil {
		log.Println(err)
		panic(err)
	}
	return db
}
